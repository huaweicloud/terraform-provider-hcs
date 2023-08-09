package huaweicloudstack

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/helper/pathorcontents"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
)

//nolint:revive
var (
	HCS_AVAILABILITY_ZONE          = os.Getenv("HCS_AVAILABILITY_ZONE")
	HCS_DEPRECATED_ENVIRONMENT     = os.Getenv("HCS_DEPRECATED_ENVIRONMENT")
	HCS_EXTGW_ID                   = os.Getenv("HCS_EXTGW_ID")
	HCS_FLAVOR_ID                  = os.Getenv("HCS_FLAVOR_ID")
	HCS_FLAVOR_NAME                = os.Getenv("HCS_FLAVOR_NAME")
	HCS_IMAGE_ID                   = os.Getenv("HCS_IMAGE_ID")
	HCS_IMAGE_NAME                 = os.Getenv("HCS_IMAGE_NAME")
	HCS_NETWORK_ID                 = os.Getenv("HCS_NETWORK_ID")
	HCS_SUBNET_ID                  = os.Getenv("HCS_SUBNET_ID")
	HCS_POOL_NAME                  = os.Getenv("HCS_POOL_NAME")
	HCS_REGION_NAME                = os.Getenv("HCS_REGION_NAME")
	HCS_ACCESS_KEY                 = os.Getenv("HCS_ACCESS_KEY")
	HCS_SECRET_KEY                 = os.Getenv("HCS_SECRET_KEY")
	HCS_VPC_ID                     = os.Getenv("HCS_VPC_ID")
	HCS_CCI_NAMESPACE              = os.Getenv("HCS_CCI_NAMESPACE")
	HCS_PROJECT_ID                 = os.Getenv("HCS_PROJECT_ID")
	HCS_DOMAIN_ID                  = os.Getenv("HCS_DOMAIN_ID")
	HCS_DOMAIN_NAME                = os.Getenv("HCS_DOMAIN_NAME")
	HCS_MRS_ENVIRONMENT            = os.Getenv("HCS_MRS_ENVIRONMENT")
	HCS_KMS_ENVIRONMENT            = os.Getenv("HCS_KMS_ENVIRONMENT")
	HCS_CCI_ENVIRONMENT            = os.Getenv("HCS_CCI_ENVIRONMENT")
	HCS_CDN_DOMAIN_NAME            = os.Getenv("HCS_CDN_DOMAIN_NAME")
	HCS_CDN_CERT_PATH              = os.Getenv("HCS_CDN_CERT_PATH")
	HCS_CDN_PRIVATE_KEY_PATH       = os.Getenv("HCS_CDN_PRIVATE_KEY_PATH")
	HCS_ENTERPRISE_PROJECT_ID_TEST = os.Getenv("HCS_ENTERPRISE_PROJECT_ID_TEST")
	HCS_USER_ID                    = os.Getenv("HCS_USER_ID")
	HCS_CHARGING_MODE              = os.Getenv("HCS_CHARGING_MODE")
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"huaweicloudstack": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	// Do not run the test if this is a deprecated testing environment.
	if HCS_DEPRECATED_ENVIRONMENT != "" {
		t.Skip("This environment only runs deprecated tests")
	}
}

func testAccPreCheckDeprecated(t *testing.T) {
	if HCS_DEPRECATED_ENVIRONMENT == "" {
		t.Skip("This environment does not support deprecated tests")
	}
}

func testAccPreCheckMrs(t *testing.T) {
	if HCS_MRS_ENVIRONMENT == "" {
		t.Skip("This environment does not support MRS tests")
	}
}

func testAccPreCheckKms(t *testing.T) {
	if HCS_KMS_ENVIRONMENT == "" {
		t.Skip("This environment does not support KMS tests")
	}
}

func testAccPreCheckCDN(t *testing.T) {
	if HCS_CDN_DOMAIN_NAME == "" {
		t.Skip("This environment does not support CDN tests")
	}
}

func testAccPreCheckCERT(t *testing.T) {
	if HCS_CDN_CERT_PATH == "" || HCS_CDN_PRIVATE_KEY_PATH == "" {
		t.Skip("This environment does not support CDN certificate tests")
	}
}

func testAccPreCheckCCINamespace(t *testing.T) {
	if HCS_CCI_NAMESPACE == "" {
		t.Skip("This environment does not support CCI Namespace tests")
	}
}

func testAccPreCheckCCI(t *testing.T) {
	if HCS_CCI_ENVIRONMENT == "" {
		t.Skip("This environment does not support CCI tests")
	}
}

func testAccPreCheckEpsID(t *testing.T) {
	if HCS_ENTERPRISE_PROJECT_ID_TEST == "" {
		t.Skip("This environment does not support Enterprise Project ID tests")
	}
}

func testAccPreCheckChargingMode(t *testing.T) {
	if HCS_CHARGING_MODE != "prePaid" {
		t.Skip("This environment does not support prepaid tests")
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// Steps for configuring HuaweiCloudStack with SSL validation are here:
// https://github.com/hashicorp/terraform/pull/6279#issuecomment-219020144
func TestAccProvider_caCertFile(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping HuaweiCloudStack SSL test.")
	}
	if os.Getenv("OS_CACERT") == "" {
		t.Skip("OS_CACERT is not set; skipping HuaweiCloudStack CA test.")
	}

	p := Provider()

	caFile, err := envVarFile("OS_CACERT")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(caFile)

	raw := map[string]interface{}{
		"cacert_file": caFile,
	}

	diags := p.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected err when specifying HuaweiCloudStack CA by file: %s", diags[0].Summary)
	}
}

func TestAccProvider_caCertString(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping HuaweiCloudStack SSL test.")
	}
	if os.Getenv("OS_CACERT") == "" {
		t.Skip("OS_CACERT is not set; skipping HuaweiCloudStack CA test.")
	}

	p := Provider()

	caContents, err := envVarContents("OS_CACERT")
	if err != nil {
		t.Fatal(err)
	}
	raw := map[string]interface{}{
		"cacert_file": caContents,
	}

	diags := p.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected err when specifying HuaweiCloudStack CA by string: %s", diags[0].Summary)
	}
}

func TestAccProvider_clientCertFile(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping HuaweiCloudStack SSL test.")
	}
	if os.Getenv("OS_CERT") == "" || os.Getenv("OS_KEY") == "" {
		t.Skip("OS_CERT or OS_KEY is not set; skipping HuaweiCloudStack client SSL auth test.")
	}

	p := Provider()

	certFile, err := envVarFile("OS_CERT")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(certFile)
	keyFile, err := envVarFile("OS_KEY")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(keyFile)

	raw := map[string]interface{}{
		"cert": certFile,
		"key":  keyFile,
	}

	diags := p.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected err when specifying HuaweiCloudStack Client keypair by file: %s", diags[0].Summary)
	}
}

func TestAccProvider_clientCertString(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping HuaweiCloudStack SSL test.")
	}
	if os.Getenv("OS_CERT") == "" || os.Getenv("OS_KEY") == "" {
		t.Skip("OS_CERT or OS_KEY is not set; skipping HuaweiCloudStack client SSL auth test.")
	}

	p := Provider()

	certContents, err := envVarContents("OS_CERT")
	if err != nil {
		t.Fatal(err)
	}
	keyContents, err := envVarContents("OS_KEY")
	if err != nil {
		t.Fatal(err)
	}

	raw := map[string]interface{}{
		"cert": certContents,
		"key":  keyContents,
	}

	diags := p.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected err when specifying HuaweiCloudStack Client keypair by contents: %s", diags[0].Summary)
	}
}

func envVarContents(varName string) (string, error) {
	contents, _, err := pathorcontents.Read(os.Getenv(varName))
	if err != nil {
		return "", fmtp.Errorf("Error reading %s: %s", varName, err)
	}
	return contents, nil
}

func envVarFile(varName string) (string, error) {
	contents, err := envVarContents(varName)
	if err != nil {
		return "", err
	}

	tmpFile, err := os.CreateTemp("", varName)
	if err != nil {
		return "", fmtp.Errorf("Error creating temp file: %s", err)
	}
	if _, err := tmpFile.Write([]byte(contents)); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmtp.Errorf("Error writing temp file: %s", err)
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmtp.Errorf("Error closing temp file: %s", err)
	}
	return tmpFile.Name(), nil
}
