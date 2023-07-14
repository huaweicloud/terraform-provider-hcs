#!/usr/bin/env bash

function usage() {
    echo "Usage: ./scripts/codecheck.sh {package}"
    echo "Example1: ./scripts/codecheck.sh ./huaweicloud/services/vpc"
    echo "Example2: ./scripts/codecheck.sh ./huaweicloud/services/..."
    echo ""
}

function checkImporter() {
    dir=$1
    for f in $(ls $dir); do
        if [[ $f =~ "resource_huaweicloud_" ]]; then
            hasImporter=$(grep -w "Importer:" $dir/$f)
            if [ "X$hasImporter" == "X" ]; then
                echo -e "\033[31m  -> the resource in $f should can be imported\n\033[0m"
            fi
        fi
    done
}

function checkCheckDeleted() {
    dir=$1
    for f in $(ls $dir); do
        if [[ $f =~ "resource_huaweicloud_" ]]; then
            checkDeleted=$(grep "CheckDeleted" $dir/$f)
             if [ "X$checkDeleted" == "X" ]; then
                checkDeleted=$(grep "\"Resource not found\"" $dir/$f)
            fi

            if [ "X$checkDeleted" == "X" ]; then
                echo -e "\033[31m  -> $f: please use common.CheckDeletedDiag in ReadContext\n\033[0m"
            fi
        fi
    done
}

function checkMultierror() {
    dir=$1
    for f in $(ls $dir); do
        if [[ $f =~ "_huaweicloud_" ]]; then
            hasMultierror=$(grep -w "go-multierror" $dir/$f)
            if [ "X$hasMultierror" == "X" ]; then
                echo -e "\033[31m  -> please use go-multierror package in $f\n\033[0m"
            fi
        fi
    done
}

function checkHuaweiCloudKey() {
    dir=$1
    key_words=("fmt.Errorf" "diag.Errorf" "log.Printf")
    for key in ${key_words[@]}; do
        result=$(grep -rn $key $dir | grep -i " huaweicloud")
        if [ "X$result" != "X" ]; then
            echo -e "\033[31m  -> the following $key statements contain 'HuaweiCloud' key:\033[0m"
            echo -e "$result\n"
        fi
    done
}

# Check parameters
package=$1
if [ "X$package" == "X" ]; then
    echo -e "error: package is missing!\n"
    usage
    exit 1
fi
# trim right "/" if necessary
package=${package%/}
packageDir=${package%...}
service=${package##*/}

# Check working directory
workDir=`pwd`
thisDir=${workDir##*/}
if [ "X$thisDir" != "Xterraform-provider-huaweicloud" ]; then
    echo -e "error: the working directory must be terraform-provider-hcs!\n"
    usage
    exit 1
fi

git status >/dev/null
if [ $? -ne 0 ]; then
    echo -e "error: the working directory is not a git repository!\n"
    exit 2
fi

# Check running environment
echo -e "\n==> Checking for running environment..."
LINT=$(which golangci-lint)
SCC=$(which scc)
MISSPELL=$(which misspell)
CYCLO=$(which gocyclo)

if [ "X$LINT" == "X" ] || [ "X$SCC" == "X" ] || [ "X$MISSPELL" == "X" ] || [ "X$CYCLO" == "X" ]; then
    echo "    ==> Checking PATH..."
    GOBIN=$(go env GOPATH)/bin
    added=$(echo $PATH | grep -w $GOBIN)
    if [ "X$added" == "X" ]; then
        echo -e "error: the GOBIN is not in PATH, please add it manually!\n"
        exit 2
    fi
fi

if [ "X$LINT" == "X" ]; then
    echo "    ==> Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

if [ "X$SCC" == "X" ]; then
    echo "    ==> Installing boyter/scc..."
    go install github.com/boyter/scc/v3@latest
fi

if [ "X$MISSPELL" == "X" ]; then
    echo "    ==> Installing misspell..."
    go install github.com/client9/misspell/cmd/misspell@latest
fi

if [ "X$CYCLO" == "X" ]; then
    echo "    ==> Installing gocyclo..."
    go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
fi

# Apply patch
echo -e "\n==> Applying patch..."
git apply --check ./scripts/0001-deprecate-fmtp-and-logp.patch
if [ $? -ne 0 ]; then
    echo -e "warning: cannot apply patch\n"
else
    git apply ./scripts/0001-deprecate-fmtp-and-logp.patch
    applied=TRUE
fi

# Check Code Complexity
echo -e "\n==> Checking for code complexity..."
scc --by-file -s complexity --no-cocomo -w $packageDir | grep -v "/deprecated/"
if [ $? -ne 0 ]; then
    exit 1
fi

echo "the TOP10 most complex functions:"
gocyclo -top 10 -avg $packageDir

# Check golangci-lint
echo -e "\n==> Checking for golangci-lint..."
golangci-lint run $package

# Nolint Directiving
echo -e "\n==> Checking for Nolint directives..."
grep -rn "nolint:" $packageDir | grep -v "/deprecated/"
grep -rn "lintignore:" $packageDir | grep -v "/deprecated/"

if [ "X$service" != "X..." ] && [[ $package == ./huaweicloudstack/services/* ]] && [[ $package != ./huaweicloudstack/services/acceptance/* ]]; then
    grep -rn "markdownlint" ./docs | grep "/${service}_"

    echo -e "\n==> Checking for TF features in $service..."
    checkImporter $packageDir
    checkCheckDeleted $packageDir
    checkMultierror $packageDir
    checkHuaweiCloudKey $packageDir

    echo -e "\n==> Checking for misspell in $service..."
    misspell ./docs | grep "/${service}_"
    misspell ./examples | grep -w "${service}"

    # update path to "./huaweicloudstack/services/acceptance/xxx"
    testpackage=${package/"services"/"services/acceptance"}
    if [ ! -d $testpackage ]; then
        echo -e "error: the acceptance directory is not exist!\n"
        exit 1
    fi
    
    echo -e "\n==> Checking for code complexity in $testpackage..."
    scc --by-file -s complexity --no-cocomo -w $testpackage

    echo "the TOP5 most complex functions:"
    gocyclo -top 5 -avg $testpackage

    echo -e "\n==> Checking for golangci-lint in $testpackage..."
    golangci-lint run $testpackage

    echo -e "\n==> Checking for Nolint directives in $testpackage..."
    grep -rn "nolint:" $testpackage
    grep -rn "lintignore:" $testpackage
fi

# cleanup
if [ "X$applied" == "XTRUE" ]; then
    echo -e "\n==> Cleanup patch..."
    git checkout -- huaweicloudstack/utils/fmtp/errors.go
    git checkout -- huaweicloudstack/utils/logp/log.go
fi

echo -e "\nCheck Completed!\n"
exit 0
