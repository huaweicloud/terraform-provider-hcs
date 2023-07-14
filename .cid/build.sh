#!/bin/bash

PROJECT_ROOT=$(cd `dirname $0/`/..;pwd)

cd $PROJECT_ROOT/

go get

CGO_CFLAGS="-fstack-protector-strong -D_FORTIFY_SOURCE=2 -O2" go build -buildmode=pie --ldflags "-s -linkmode 'external' -extldflags '-Wl,-z,now'" main.go

