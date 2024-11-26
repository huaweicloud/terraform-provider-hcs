//nolint:revive
package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack"
)

var (
	HCS_SKIP_UNSUPPORTED_TEST = os.Getenv("HCS_SKIP_UNSUPPORTED_TEST")

	HCS_REGION_NAME                        = os.Getenv("HCS_REGION_NAME")
	HCS_CLOUD                              = os.Getenv("HCS_CLOUD")
	HCS_CUSTOM_REGION_NAME                 = os.Getenv("HCS_CUSTOM_REGION_NAME")
	HCS_AVAILABILITY_ZONE                  = os.Getenv("HCS_AVAILABILITY_ZONE")
	HCS_ACCESS_KEY                         = os.Getenv("HCS_ACCESS_KEY")
	HCS_SECRET_KEY                         = os.Getenv("HCS_SECRET_KEY")
	HCS_USER_ID                            = os.Getenv("HCS_USER_ID")
	HCS_USER_NAME                          = os.Getenv("HCS_USER_NAME")
	HCS_PROJECT_ID                         = os.Getenv("HCS_PROJECT_ID")
	HCS_DOMAIN_ID                          = os.Getenv("HCS_DOMAIN_ID")
	HCS_DOMAIN_NAME                        = os.Getenv("HCS_DOMAIN_NAME")
	HCS_ENTERPRISE_PROJECT_ID_TEST         = os.Getenv("HCS_ENTERPRISE_PROJECT_ID_TEST")
	HCS_ENTERPRISE_MIGRATE_PROJECT_ID_TEST = os.Getenv("HCS_ENTERPRISE_MIGRATE_PROJECT_ID_TEST")

	HCS_FLAVOR_ID                 = os.Getenv("HCS_FLAVOR_ID")
	HCS_FLAVOR_NAME               = os.Getenv("HCS_FLAVOR_NAME")
	HCS_IMAGE_ID                  = os.Getenv("HCS_IMAGE_ID")
	HCS_IMAGE_NAME                = os.Getenv("HCS_IMAGE_NAME")
	HCS_VPC_ID                    = os.Getenv("HCS_VPC_ID")
	HCS_NETWORK_ID                = os.Getenv("HCS_NETWORK_ID")
	HCS_SUBNET_ID                 = os.Getenv("HCS_SUBNET_ID")
	HCS_ENTERPRISE_PROJECT_ID     = os.Getenv("HCS_ENTERPRISE_PROJECT_ID")
	HCS_ADMIN                     = os.Getenv("HCS_ADMIN")
	HCS_KEYPAIR_NAME              = os.Getenv("HCS_KEYPAIR_NAME")
	HCS_SERVER_GROUP_ID           = os.Getenv("HCS_SERVER_GROUP_ID")
	HCS_ECS_INSTANCE_ID           = os.Getenv("HCS_ECS_INSTANCE_ID")
	HCS_EIP_ID                    = os.Getenv("HCS_EIP_ID")
	HCS_EIP_NAME                  = os.Getenv("HCS_EIP_NAME")
	HCS_EIP_ADDRESS               = os.Getenv("HCS_EIP_ADDRESS")
	HCS_EIP_EXTERNAL_NETWORK_NAME = os.Getenv("HCS_EIP_EXTERNAL_NETWORK_NAME")

	HCS_MAPREDUCE_CUSTOM           = os.Getenv("HCS_MAPREDUCE_CUSTOM")
	HCS_MAPREDUCE_BOOTSTRAP_SCRIPT = os.Getenv("HCS_MAPREDUCE_BOOTSTRAP_SCRIPT")

	HCS_OBS_BUCKET_NAME        = os.Getenv("HCS_OBS_BUCKET_NAME")
	HCS_OBS_DESTINATION_BUCKET = os.Getenv("HCS_OBS_DESTINATION_BUCKET")

	HCS_OMS_ENABLE_FLAG = os.Getenv("HCS_OMS_ENABLE_FLAG")

	HCS_DEPRECATED_ENVIRONMENT = os.Getenv("HCS_DEPRECATED_ENVIRONMENT")
	HCS_INTERNAL_USED          = os.Getenv("HCS_INTERNAL_USED")

	HCS_WAF_ENABLE_FLAG = os.Getenv("HCS_WAF_ENABLE_FLAG")

	HCS_DEST_REGION          = os.Getenv("HCS_DEST_REGION")
	HCS_DEST_PROJECT_ID      = os.Getenv("HCS_DEST_PROJECT_ID")
	HCS_DEST_PROJECT_ID_TEST = os.Getenv("HCS_DEST_PROJECT_ID_TEST")
	HCS_CHARGING_MODE        = os.Getenv("HCS_CHARGING_MODE")
	HCS_HIGH_COST_ALLOW      = os.Getenv("HCS_HIGH_COST_ALLOW")
	HCS_SWR_SHARING_ACCOUNT  = os.Getenv("HCS_SWR_SHARING_ACCOUNT")

	HCS_RAM_SHARE_ACCOUNT_ID   = os.Getenv("HCS_RAM_SHARE_ACCOUNT_ID")
	HCS_RAM_SHARE_RESOURCE_URN = os.Getenv("HCS_RAM_SHARE_RESOURCE_URN")

	HCS_CERTIFICATE_KEY_PATH         = os.Getenv("HCS_CERTIFICATE_KEY_PATH")
	HCS_CERTIFICATE_CHAIN_PATH       = os.Getenv("HCS_CERTIFICATE_CHAIN_PATH")
	HCS_CERTIFICATE_PRIVATE_KEY_PATH = os.Getenv("HCS_CERTIFICATE_PRIVATE_KEY_PATH")
	HCS_CERTIFICATE_SERVICE          = os.Getenv("HCS_CERTIFICATE_SERVICE")
	HCS_CERTIFICATE_PROJECT          = os.Getenv("HCS_CERTIFICATE_PROJECT")
	HCS_CERTIFICATE_PROJECT_UPDATED  = os.Getenv("HCS_CERTIFICATE_PROJECT_UPDATED")
	HCS_CERTIFICATE_NAME             = os.Getenv("HCS_CERTIFICATE_NAME")
	HCS_DMS_ENVIRONMENT              = os.Getenv("HCS_DMS_ENVIRONMENT")
	HCS_SMS_SOURCE_SERVER            = os.Getenv("HCS_SMS_SOURCE_SERVER")

	HCS_DLI_FLINK_JAR_OBS_PATH           = os.Getenv("HCS_DLI_FLINK_JAR_OBS_PATH")
	HCS_DLI_DS_AUTH_CSS_OBS_PATH         = os.Getenv("HCS_DLI_DS_AUTH_CSS_OBS_PATH")
	HCS_DLI_DS_AUTH_KAFKA_TRUST_OBS_PATH = os.Getenv("HCS_DLI_DS_AUTH_KAFKA_TRUST_OBS_PATH")
	HCS_DLI_DS_AUTH_KAFKA_KEY_OBS_PATH   = os.Getenv("HCS_DLI_DS_AUTH_KAFKA_KEY_OBS_PATH")
	HCS_DLI_DS_AUTH_KRB_CONF_OBS_PATH    = os.Getenv("HCS_DLI_DS_AUTH_KRB_CONF_OBS_PATH")
	HCS_DLI_DS_AUTH_KRB_TAB_OBS_PATH     = os.Getenv("HCS_DLI_DS_AUTH_KRB_TAB_OBS_PATH")
	HCS_DLI_AGENCY_FLAG                  = os.Getenv("HCS_DLI_AGENCY_FLAG")

	HCS_GITHUB_REPO_HOST        = os.Getenv("HCS_GITHUB_REPO_HOST")        // Repository host (Github, Gitlab, Gitee)
	HCS_GITHUB_PERSONAL_TOKEN   = os.Getenv("HCS_GITHUB_PERSONAL_TOKEN")   // Personal access token (Github, Gitlab, Gitee)
	HCS_GITHUB_REPO_PWD         = os.Getenv("HCS_GITHUB_REPO_PWD")         // Repository password (DevCloud, BitBucket)
	HCS_GITHUB_REPO_URL         = os.Getenv("HCS_GITHUB_REPO_URL")         // Repository URL (Github, Gitlab, Gitee)
	HCS_OBS_STORAGE_URL         = os.Getenv("HCS_OBS_STORAGE_URL")         // OBS storage URL where ZIP file is located
	HCS_BUILD_IMAGE_URL         = os.Getenv("HCS_BUILD_IMAGE_URL")         // SWR Image URL for component deployment
	HCS_BUILD_IMAGE_URL_UPDATED = os.Getenv("HCS_BUILD_IMAGE_URL_UPDATED") // SWR Image URL for component deployment update

	HCS_VOD_WATERMARK_FILE   = os.Getenv("HCS_VOD_WATERMARK_FILE")
	HCS_VOD_MEDIA_ASSET_FILE = os.Getenv("HCS_VOD_MEDIA_ASSET_FILE")

	HCS_CHAIR_EMAIL              = os.Getenv("HCS_CHAIR_EMAIL")
	HCS_GUEST_EMAIL              = os.Getenv("HCS_GUEST_EMAIL")
	HCS_MEETING_ACCOUNT_NAME     = os.Getenv("HCS_MEETING_ACCOUNT_NAME")
	HCS_MEETING_ACCOUNT_PASSWORD = os.Getenv("HCS_MEETING_ACCOUNT_PASSWORD")
	HCS_MEETING_APP_ID           = os.Getenv("HCS_MEETING_APP_ID")
	HCS_MEETING_APP_KEY          = os.Getenv("HCS_MEETING_APP_KEY")
	HCS_MEETING_USER_ID          = os.Getenv("HCS_MEETING_USER_ID")
	HCS_MEETING_ROOM_ID          = os.Getenv("HCS_MEETING_ROOM_ID")

	HCS_AAD_INSTANCE_ID = os.Getenv("HCS_AAD_INSTANCE_ID")
	HCS_AAD_IP_ADDRESS  = os.Getenv("HCS_AAD_IP_ADDRESS")

	HCS_WORKSPACE_AD_DOMAIN_NAME = os.Getenv("HCS_WORKSPACE_AD_DOMAIN_NAME") // Domain name, e.g. "example.com".
	HCS_WORKSPACE_AD_SERVER_PWD  = os.Getenv("HCS_WORKSPACE_AD_SERVER_PWD")  // The password of AD server.
	HCS_WORKSPACE_AD_DOMAIN_IP   = os.Getenv("HCS_WORKSPACE_AD_DOMAIN_IP")   // Active domain IP, e.g. "192.168.196.3".
	HCS_WORKSPACE_AD_VPC_ID      = os.Getenv("HCS_WORKSPACE_AD_VPC_ID")      // The VPC ID to which the AD server and desktops belongs.
	HCS_WORKSPACE_AD_NETWORK_ID  = os.Getenv("HCS_WORKSPACE_AD_NETWORK_ID")  // The network ID to which the AD server belongs.

	HCS_FGS_TRIGGER_LTS_AGENCY = os.Getenv("HCS_FGS_TRIGGER_LTS_AGENCY")

	HCS_KMS_ENVIRONMENT    = os.Getenv("HCS_KMS_ENVIRONMENT")
	HCS_KMS_KEY_ID         = os.Getenv("HCS_KMS_KEY_ID")
	HCS_KMS_HSM_CLUSTER_ID = os.Getenv("HCS_KMS_HSM_CLUSTER_ID")

	HCS_ORGANIZATIONS_ENVIRONMENT            = os.Getenv("HCS_ORGANIZATIONS_ENVIRONMENT")
	HCS_ORGANIZATIONS_INVITE_ACCOUNT_ID      = os.Getenv("HCS_ORGANIZATIONS_INVITE_ACCOUNT_ID")
	HCS_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID = os.Getenv("HCS_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID")
	HCS_ORGANIZATIONS_INVITATION_ID          = os.Getenv("HCS_ORGANIZATIONS_INVITATION_ID")

	HCS_ER_TEST_ON = os.Getenv("HCS_ER_TEST_ON") // Whether to run the ER related tests.

	// The OBS address where the HCL/JSON template archive (No variables) is located.
	HCS_RF_TEMPLATE_ARCHIVE_NO_VARS_URI = os.Getenv("HCS_RF_TEMPLATE_ARCHIVE_NO_VARS_URI")
	// The OBS address where the HCL/JSON template archive is located.
	HCS_RF_TEMPLATE_ARCHIVE_URI = os.Getenv("HCS_RF_TEMPLATE_ARCHIVE_URI")
	// The OBS address where the variable archive corresponding to the HCL/JSON template is located.
	HCS_RF_VARIABLES_ARCHIVE_URI = os.Getenv("HCS_RF_VARIABLES_ARCHIVE_URI")

	// The direct connection ID (provider does not support direct connection resource).
	HCS_DC_DIRECT_CONNECT_ID = os.Getenv("HCS_DC_DIRECT_CONNECT_ID")

	// The CFW instance ID
	HCS_CFW_INSTANCE_ID = os.Getenv("HCS_CFW_INSTANCE_ID")

	// The cluster ID of the CCE
	HCS_CCE_CLUSTER_ID = os.Getenv("HCS_CCE_CLUSTER_ID")
	// The partition az of the CCE
	HCS_CCE_PARTITION_AZ = os.Getenv("HCS_CCE_PARTITION_AZ")
	// The namespace of the workload is located
	HCS_WORKLOAD_NAMESPACE = os.Getenv("HCS_WORKLOAD_NAMESPACE")
	// The workload type deployed in CCE/CCI
	HCS_WORKLOAD_TYPE = os.Getenv("HCS_WORKLOAD_TYPE")
	// The workload name deployed in CCE/CCI
	HCS_WORKLOAD_NAME = os.Getenv("HCS_WORKLOAD_NAME")
	// The target region of SWR image auto sync
	HCS_SWR_TARGET_REGION = os.Getenv("HCS_SWR_TARGET_REGION")
	// The target organization of SWR image auto sync
	HCS_SWR_TARGET_ORGANIZATION = os.Getenv("HCS_SWR_TARGET_ORGANIZATION")

	// The ID of the CBR backup
	HCS_IMS_BACKUP_ID = os.Getenv("HCS_IMS_BACKUP_ID")

	// The SecMaster workspace ID
	HCS_SECMASTER_WORKSPACE_ID = os.Getenv("HCS_SECMASTER_WORKSPACE_ID")

	// Deprecated
	HCS_SRC_ACCESS_KEY = os.Getenv("HCS_SRC_ACCESS_KEY")
	HCS_SRC_SECRET_KEY = os.Getenv("HCS_SRC_SECRET_KEY")
	HCS_EXTGW_ID       = os.Getenv("HCS_EXTGW_ID")
	HCS_POOL_NAME      = os.Getenv("HCS_POOL_NAME")

	HCS_IMAGE_SHARE_SOURCE_IMAGE_ID = os.Getenv("HCS_IMAGE_SHARE_SOURCE_IMAGE_ID")

	HCS_LTS_ENABLE_FLAG                 = os.Getenv("HCS_LTS_ENABLE_FLAG")
	HCS_LTS_STRUCT_CONFIG_TEMPLATE_ID   = os.Getenv("HCS_LTS_STRUCT_CONFIG_TEMPLATE_ID")
	HCS_LTS_STRUCT_CONFIG_TEMPLATE_NAME = os.Getenv("HCS_LTS_STRUCT_CONFIG_TEMPLATE_NAME")

	HCS_RDS_INSTANCE_ID = os.Getenv("HCS_RDS_INSTANCE_ID")
	HCS_RDS_BACKUP_ID   = os.Getenv("HCS_RDS_BACKUP_ID")
)

// TestAccProviders is a static map containing only the main provider instance.
//
// Deprecated: Terraform Plugin SDK version 2 uses TestCase.ProviderFactories
// but supports this value in TestCase.Providers for backwards compatibility.
// In the future Providers: TestAccProviders will be changed to
// ProviderFactories: TestAccProviderFactories
var TestAccProviders map[string]*schema.Provider

// TestAccProviderFactories is a static map containing only the main provider instance
var TestAccProviderFactories map[string]func() (*schema.Provider, error)

// TestAccProvider is the "main" provider instance
var TestAccProvider *schema.Provider

func init() {
	TestAccProvider = huaweicloudstack.Provider()

	TestAccProviders = map[string]*schema.Provider{
		"hcs": TestAccProvider,
	}

	TestAccProviderFactories = map[string]func() (*schema.Provider, error){
		"hcs": func() (*schema.Provider, error) {
			return TestAccProvider, nil
		},
	}
}

func preCheckRequiredEnvVars(t *testing.T) {
	if HCS_REGION_NAME == "" {
		t.Fatal("HCS_REGION_NAME must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckSwrOrganization(t *testing.T) {
	if HCS_CLOUD == "" {
		t.Fatal("HCS_CLOUD must be set for acceptance tests")
	}
}

// TestAccSkipUnsupportedTest is a method for skipping acceptance tests for which all related features are not yet
// supported.
// lintignore:AT003
func TestAccSkipUnsupportedTest(t *testing.T) {
	if HCS_SKIP_UNSUPPORTED_TEST == "" {
		t.Skip("Skip this test case because some functions involved in this test case are not supported.")
	}
}

// lintignore:AT003
func TestAccPreCheckOrganizations(t *testing.T) {
	if HCS_ORGANIZATIONS_ENVIRONMENT == "" {
		t.Skip("This environment does not support Organizations tests")
	}
}

// lintignore:AT003
func TestAccPreCheckOrganizationsInviteAccountId(t *testing.T) {
	if HCS_ORGANIZATIONS_INVITE_ACCOUNT_ID == "" {
		t.Skip("HCS_ORGANIZATIONS_INVITE_ACCOUNT_ID must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckOrganizationsInvitationId(t *testing.T) {
	if HCS_ORGANIZATIONS_INVITATION_ID == "" {
		t.Skip("HCS_ORGANIZATIONS_INVITATION_ID must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckOrganizationsOrganizationalUnitId(t *testing.T) {
	if HCS_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID == "" {
		t.Skip("HCS_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheck(t *testing.T) {
	// Do not run the test if this is a deprecated testing environment.
	if HCS_DEPRECATED_ENVIRONMENT != "" {
		t.Skip("This environment only runs deprecated tests")
	}

	preCheckRequiredEnvVars(t)
}

// lintignore:AT003
func TestAccPrecheckDomainId(t *testing.T) {
	if HCS_DOMAIN_ID == "" {
		t.Skip("HCS_DOMAIN_ID must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPrecheckCustomRegion(t *testing.T) {
	if HCS_CUSTOM_REGION_NAME == "" {
		t.Skip("HCS_CUSTOM_REGION_NAME must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckDeprecated(t *testing.T) {
	if HCS_DEPRECATED_ENVIRONMENT == "" {
		t.Skip("This environment does not support deprecated tests")
	}

	preCheckRequiredEnvVars(t)
}

// lintignore:AT003
func TestAccPreCheckInternal(t *testing.T) {
	if HCS_INTERNAL_USED == "" {
		t.Skip("HCS_INTERNAL_USED must be set for internal acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckEpsID(t *testing.T) {
	// The environment variables in tests take HCS_ENTERPRISE_PROJECT_ID_TEST instead of HCS_ENTERPRISE_PROJECT_ID to
	// ensure that other data-resources that support enterprise projects query the default project without being
	// affected by this variable.
	if HCS_ENTERPRISE_PROJECT_ID_TEST == "" {
		t.Skip("The environment variables does not support Enterprise Project ID for acc tests")
	}
}

// lintignore:AT003
func TestAccPreCheckMigrateEpsID(t *testing.T) {
	if HCS_ENTERPRISE_PROJECT_ID_TEST == "" {
		t.Skip("The environment variables does not support Migrate Enterprise Project ID for acc tests")
	}
}

// lintignore:AT003
func TestAccPreCheckUserId(t *testing.T) {
	if HCS_USER_ID == "" {
		t.Skip("The environment variables does not support the user ID (HCS_USER_ID) for acc tests")
	}
}

// lintignore:AT003
func TestAccPreCheckSms(t *testing.T) {
	if HCS_SMS_SOURCE_SERVER == "" {
		t.Skip("HCS_SMS_SOURCE_SERVER must be set for SMS acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckMrsCustom(t *testing.T) {
	if HCS_MAPREDUCE_CUSTOM == "" {
		t.Skip("HCS_MAPREDUCE_CUSTOM must be set for acceptance tests:custom type cluster of map reduce")
	}
}

// lintignore:AT003
func TestAccPreCheckFgsTrigger(t *testing.T) {
	if HCS_FGS_TRIGGER_LTS_AGENCY == "" {
		t.Skip("HCS_FGS_TRIGGER_LTS_AGENCY must be set for FGS trigger acceptance tests")
	}
}

// Deprecated
// lintignore:AT003
func TestAccPreCheckMaas(t *testing.T) {
	if HCS_ACCESS_KEY == "" || HCS_SECRET_KEY == "" || HCS_SRC_ACCESS_KEY == "" || HCS_SRC_SECRET_KEY == "" {
		t.Skip("HCS_ACCESS_KEY, HCS_SECRET_KEY, HCS_SRC_ACCESS_KEY, and HCS_SRC_SECRET_KEY  must be set for MAAS acceptance tests")
	}
}

func RandomAccResourceName() string {
	return fmt.Sprintf("tf_test_%s", acctest.RandString(5))
}

func RandomAccResourceNameWithDash() string {
	return fmt.Sprintf("tf-test-%s", acctest.RandString(5))
}

func RandomCidr() string {
	return fmt.Sprintf("172.16.%d.0/24", acctest.RandIntRange(0, 255))
}

func RandomCidrAndGatewayIp() (string, string) {
	seed := acctest.RandIntRange(0, 255)
	return fmt.Sprintf("172.16.%d.0/24", seed), fmt.Sprintf("172.16.%d.1", seed)
}

func RandomPassword() string {
	return fmt.Sprintf("%s%s%s%d", acctest.RandStringFromCharSet(2, "ABCDEFGHIJKLMNOPQRSTUVWXZY"),
		acctest.RandString(3), acctest.RandStringFromCharSet(2, "~!@#%^*-_=+?"), acctest.RandIntRange(1000, 9999))
}

// lintignore:AT003
func TestAccPrecheckWafInstance(t *testing.T) {
	if HCS_WAF_ENABLE_FLAG == "" {
		t.Skip("Jump the WAF acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckOmsInstance(t *testing.T) {
	if HCS_OMS_ENABLE_FLAG == "" {
		t.Skip("Jump the OMS acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckAdminOnly(t *testing.T) {
	if HCS_ADMIN == "" {
		t.Skip("Skipping test because it requires the admin privileges")
	}
}

// lintignore:AT003
func TestAccPreCheckEipId(t *testing.T) {
	if HCS_EIP_ID == "" {
		t.Skip("Skipping test because it requires the public IP ID (HCS_EIP_ID)")
	}
}

// lintignore:AT003
func TestAccPreCheckEipAddress(t *testing.T) {
	if HCS_EIP_ADDRESS == "" {
		t.Skip("Skipping test because it requires the public IP address (HCS_EIP_ADDRESS)")
	}
}

// lintignore:AT003
func TestAccPreCheckReplication(t *testing.T) {
	if HCS_DEST_REGION == "" || HCS_DEST_PROJECT_ID == "" {
		t.Skip("Jump the replication policy acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckProjectId(t *testing.T) {
	if HCS_DEST_PROJECT_ID_TEST == "" {
		t.Skip("Skipping test because it requires the test project id.")
	}
}

// lintignore:AT003
func TestAccPreCheckProject(t *testing.T) {
	if HCS_ENTERPRISE_PROJECT_ID_TEST != "" {
		t.Skip("This environment does not support project tests")
	}
}

// lintignore:AT003
func TestAccPreCheckOBS(t *testing.T) {
	if HCS_ACCESS_KEY == "" || HCS_SECRET_KEY == "" {
		t.Skip("HCS_ACCESS_KEY and HCS_SECRET_KEY must be set for OBS acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckOBSBucket(t *testing.T) {
	if HCS_OBS_BUCKET_NAME == "" {
		t.Skip("HCS_OBS_BUCKET_NAME must be set for OBS object acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckOBSDestinationBucket(t *testing.T) {
	if HCS_OBS_DESTINATION_BUCKET == "" {
		t.Skip("HCS_OBS_DESTINATION_BUCKET must be set for OBS destination tests")
	}
}

// lintignore:AT003
func TestAccPreCheckChargingMode(t *testing.T) {
	if HCS_CHARGING_MODE != "prePaid" {
		t.Skip("This environment does not support prepaid tests")
	}
}

func TestAccPreCheckEip(t *testing.T) {
	if HCS_EIP_ID == "" || HCS_EIP_ADDRESS == "" {
		t.Skip("HCS_EIP_ID and HCS_EIP_ADDRESS must be set for EIP corresponding acceptance test")
	}
}

// lintignore:AT003
func TestAccPreCheckHighCostAllow(t *testing.T) {
	if HCS_HIGH_COST_ALLOW == "" {
		t.Skip("Do not allow expensive testing")
	}
}

// lintignore:AT003
func TestAccPreCheckScm(t *testing.T) {
	if HCS_CERTIFICATE_KEY_PATH == "" || HCS_CERTIFICATE_CHAIN_PATH == "" ||
		HCS_CERTIFICATE_PRIVATE_KEY_PATH == "" || HCS_CERTIFICATE_SERVICE == "" ||
		HCS_CERTIFICATE_PROJECT == "" || HCS_CERTIFICATE_PROJECT_UPDATED == "" {
		t.Skip("HCS_CERTIFICATE_KEY_PATH, HCS_CERTIFICATE_CHAIN_PATH, HCS_CERTIFICATE_PRIVATE_KEY_PATH, " +
			"HCS_CERTIFICATE_SERVICE, HCS_CERTIFICATE_PROJECT and HCS_CERTIFICATE_TARGET_UPDATED " +
			"can not be empty for SCM certificate tests")
	}
}

// lintignore:AT003
func TestAccPreCheckSWRDomian(t *testing.T) {
	if HCS_SWR_SHARING_ACCOUNT == "" {
		t.Skip("HCS_SWR_SHARING_ACCOUNT must be set for swr domian tests, " +
			"the value of HCS_SWR_SHARING_ACCOUNT should be another IAM user name")
	}
}

// lintignore:AT003
func TestAccPreCheckRAM(t *testing.T) {
	if HCS_RAM_SHARE_ACCOUNT_ID == "" || HCS_RAM_SHARE_RESOURCE_URN == "" {
		t.Skip("HCS_RAM_SHARE_ACCOUNT_ID and HCS_RAM_SHARE_RESOURCE_URN must be set for ram tests, " +
			"the value of HCS_RAM_SHARE_ACCOUNT_ID should be another account id")
	}
}

// lintignore:AT003
func TestAccPreCheckDms(t *testing.T) {
	if HCS_DMS_ENVIRONMENT == "" {
		t.Skip("This environment does not support DMS tests")
	}
}

// lintignore:AT003
func TestAccPreCheckDliJarPath(t *testing.T) {
	if HCS_DLI_FLINK_JAR_OBS_PATH == "" {
		t.Skip("HCS_DLI_FLINK_JAR_OBS_PATH must be set for DLI Flink Jar job acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckDliDsAuthCss(t *testing.T) {
	if HCS_DLI_DS_AUTH_CSS_OBS_PATH == "" {
		t.Skip("HCS_DLI_DS_AUTH_CSS_OBS_PATH must be set for DLI datasource CSS Auth acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckDliDsAuthKafka(t *testing.T) {
	if HCS_DLI_DS_AUTH_KAFKA_TRUST_OBS_PATH == "" || HCS_DLI_DS_AUTH_KAFKA_KEY_OBS_PATH == "" {
		t.Skip("HCS_DLI_DS_AUTH_KAFKA_TRUST_OBS_PATH,HCS_DLI_DS_AUTH_KAFKA_KEY_OBS_PATH must be set for DLI datasource Kafka Auth acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckDliDsAuthKrb(t *testing.T) {
	if HCS_DLI_DS_AUTH_KRB_CONF_OBS_PATH == "" || HCS_DLI_DS_AUTH_KRB_TAB_OBS_PATH == "" {
		t.Skip("HCS_DLI_DS_AUTH_KRB_CONF_OBS_PATH,HCS_DLI_DS_AUTH_KRB_TAB_OBS_PATH must be set for DLI datasource Kafka Auth acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckDliAgency(t *testing.T) {
	if HCS_DLI_AGENCY_FLAG == "" {
		t.Skip("HCS_DLI_AGENCY_FLAG must be set for DLI datasource DLI agency acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckRepoTokenAuth(t *testing.T) {
	if HCS_GITHUB_REPO_HOST == "" || HCS_GITHUB_PERSONAL_TOKEN == "" {
		t.Skip("Repository configurations are not completed for acceptance test of personal access token authorization.")
	}
}

// lintignore:AT003
func TestAccPreCheckRepoPwdAuth(t *testing.T) {
	if HCS_DOMAIN_NAME == "" || HCS_USER_NAME == "" || HCS_GITHUB_REPO_PWD == "" {
		t.Skip("Repository configurations are not completed for acceptance test of password authorization.")
	}
}

// lintignore:AT003
func TestAccPreCheckComponent(t *testing.T) {
	if HCS_DOMAIN_NAME == "" || HCS_GITHUB_REPO_URL == "" || HCS_OBS_STORAGE_URL == "" {
		t.Skip("Repository (package) configurations are not completed for acceptance test of component.")
	}
}

// lintignore:AT003
func TestAccPreCheckComponentDeployment(t *testing.T) {
	if HCS_BUILD_IMAGE_URL == "" {
		t.Skip("SWR image URL configuration is not completed for acceptance test of component deployment.")
	}
}

// lintignore:AT003
func TestAccPreCheckImageUrlUpdated(t *testing.T) {
	if HCS_BUILD_IMAGE_URL_UPDATED == "" {
		t.Skip("SWR image update URL configuration is not completed for acceptance test of component deployment.")
	}
}

// lintignore:AT003
func TestAccPreCheckVODWatermark(t *testing.T) {
	if HCS_VOD_WATERMARK_FILE == "" {
		t.Skip("HCS_VOD_WATERMARK_FILE must be set for VOD watermark template acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckVODMediaAsset(t *testing.T) {
	if HCS_VOD_MEDIA_ASSET_FILE == "" {
		t.Skip("HCS_VOD_MEDIA_ASSET_FILE must be set for VOD media asset acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckPwdAuth(t *testing.T) {
	if HCS_MEETING_ACCOUNT_NAME == "" || HCS_MEETING_ACCOUNT_PASSWORD == "" {
		t.Skip("The account name (HCS_MEETING_ACCOUNT_NAME) or password (HCS_MEETING_ACCOUNT_PASSWORD) is not " +
			"completed for acceptance test of conference.")
	}
}

// lintignore:AT003
func TestAccPreCheckAppAuth(t *testing.T) {
	if HCS_MEETING_APP_ID == "" || HCS_MEETING_APP_KEY == "" || HCS_MEETING_USER_ID == "" {
		t.Skip("The app ID (HCS_MEETING_APP_ID), app KEY (HCS_MEETING_APP_KEY) or user ID (HCS_MEETING_USER_ID) is not " +
			"completed for acceptance test of conference.")
	}
}

// lintignore:AT003
func TestAccPreCheckMeetingRoom(t *testing.T) {
	if HCS_MEETING_ROOM_ID == "" {
		t.Skip("The vmr ID (HCS_MEETING_ROOM_ID) is not completed for acceptance test of conference.")
	}
}

// lintignore:AT003
func TestAccPreCheckParticipants(t *testing.T) {
	if HCS_CHAIR_EMAIL == "" || HCS_GUEST_EMAIL == "" {
		t.Skip("The chair (HCS_CHAIR_EMAIL) or guest (HCS_GUEST_EMAIL) mailbox is not completed for acceptance test of " +
			"conference.")
	}
}

// lintignore:AT003
func TestAccPreCheckAadForwardRule(t *testing.T) {
	if HCS_AAD_INSTANCE_ID == "" || HCS_AAD_IP_ADDRESS == "" {
		t.Skip("The instance information is not completed for AAD rule acceptance test.")
	}
}

// lintignore:AT003
func TestAccPreCheckScmCertificateName(t *testing.T) {
	if HCS_CERTIFICATE_NAME == "" {
		t.Skip("HCS_CERTIFICATE_NAME must be set for SCM acceptance tests.")
	}
}

// TestAccPreCheckKms is a pre-check method that used to control whether KMS resource-related tests are performed.
// KMS resources cannot be deleted immediately (at least 7 days) and the quota is small.
// To save resources and ensure normal pipeline testing, automatic verification is canceled.
// Manual verification is used to ensure the availability of this function (KMS encryption).
// lintignore:AT003
func TestAccPreCheckKms(t *testing.T) {
	if HCS_KMS_ENVIRONMENT == "" {
		t.Skip("This environment does not support KMS tests")
	}
}

// lintignore:AT003
func TestAccPreCheckKmsKey(t *testing.T) {
	if HCS_KMS_KEY_ID == "" {
		t.Skip("HCS_KMS_KEY_ID must be set for encryption acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckProjectID(t *testing.T) {
	if HCS_PROJECT_ID == "" {
		t.Skip("HCS_PROJECT_ID must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckWorkspaceAD(t *testing.T) {
	if HCS_WORKSPACE_AD_DOMAIN_NAME == "" || HCS_WORKSPACE_AD_SERVER_PWD == "" || HCS_WORKSPACE_AD_DOMAIN_IP == "" ||
		HCS_WORKSPACE_AD_VPC_ID == "" || HCS_WORKSPACE_AD_NETWORK_ID == "" {
		t.Skip("The configuration of AD server is not completed for Workspace service acceptance test.")
	}
}

// lintignore:AT003
func TestAccPreCheckER(t *testing.T) {
	if HCS_ER_TEST_ON == "" {
		t.Skip("Skip all ER acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckRfArchives(t *testing.T) {
	if HCS_RF_TEMPLATE_ARCHIVE_NO_VARS_URI == "" || HCS_RF_TEMPLATE_ARCHIVE_URI == "" ||
		HCS_RF_VARIABLES_ARCHIVE_URI == "" {
		t.Skip("Skip the archive URI parameters acceptance test for RF resource stack.")
	}
}

// lintignore:AT003
func TestAccPreCheckDcDirectConnection(t *testing.T) {
	if HCS_DC_DIRECT_CONNECT_ID == "" {
		t.Skip("Skip the interface acceptance test because of the direct connection ID is missing.")
	}
}

// lintignore:AT003
func TestAccPreCheckCfw(t *testing.T) {
	if HCS_CFW_INSTANCE_ID == "" {
		t.Skip("HCS_CFW_INSTANCE_ID must be set for CFW acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckWorkloadType(t *testing.T) {
	if HCS_WORKLOAD_TYPE == "" {
		t.Skip("HCS_WORKLOAD_TYPE must be set for SWR image trigger acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckWorkloadName(t *testing.T) {
	if HCS_WORKLOAD_NAME == "" {
		t.Skip("HCS_WORKLOAD_NAME must be set for SWR image trigger acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckCceClusterId(t *testing.T) {
	if HCS_CCE_CLUSTER_ID == "" {
		t.Skip("HCS_CCE_CLUSTER_ID must be set for SWR image trigger acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckWorkloadNameSpace(t *testing.T) {
	if HCS_WORKLOAD_NAMESPACE == "" {
		t.Skip("HCS_WORKLOAD_NAMESPACE must be set for SWR image trigger acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckSwrTargetRegion(t *testing.T) {
	if HCS_SWR_TARGET_REGION == "" {
		t.Skip("HCS_SWR_TARGET_REGION must be set for SWR image auto sync tests")
	}
}

// lintignore:AT003
func TestAccPreCheckSwrTargetOrigination(t *testing.T) {
	if HCS_SWR_TARGET_ORGANIZATION == "" {
		t.Skip("HCS_SWR_TARGET_ORGANIZATION must be set for SWR image auto sync tests")
	}
}

// lintignore:AT003
func TestAccPreCheckImsBackupId(t *testing.T) {
	if HCS_IMS_BACKUP_ID == "" {
		t.Skip("HCS_IMS_BACKUP_ID must be set for IMS whole image with CBR backup id")
	}
}

// lintignore:AT003
func TestAccPreCheckSourceImage(t *testing.T) {
	if HCS_IMAGE_SHARE_SOURCE_IMAGE_ID == "" {
		t.Skip("Skip the interface acceptance test because of the source image ID is missing.")
	}
}

// lintignore:AT003
func TestAccPreCheckSecMaster(t *testing.T) {
	if HCS_SECMASTER_WORKSPACE_ID == "" {
		t.Skip("HCS_SECMASTER_WORKSPACE_ID must be set for SecMaster acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckCcePartitionAz(t *testing.T) {
	if HCS_CCE_PARTITION_AZ == "" {
		t.Skip("Skip the interface acceptance test because of the cce partition az is missing.")
	}
}

// lintignore:AT003
func TestAccPreCheckCnEast3(t *testing.T) {
	if HCS_REGION_NAME != "cn-east-3" {
		t.Skip("HCS_REGION_NAME must be cn-east-3 for this test.")
	}
}

// lintignore:AT003
func TestAccPreCheckBms(t *testing.T) {
	if HCS_USER_ID == "" {
		t.Skip("HW_USER_ID must be set for BMS acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckMrsBootstrapScript(t *testing.T) {
	if HCS_MAPREDUCE_BOOTSTRAP_SCRIPT == "" {
		t.Skip("HCS_MAPREDUCE_BOOTSTRAP_SCRIPT must be set for acceptance tests: cluster of map reduce with bootstrap")
	}
}

// lintignore:AT003
func TestAccPreCheckLtsEnableFlag(t *testing.T) {
	if HCS_LTS_ENABLE_FLAG == "" {
		t.Skip("HCS_LTS_ENABLE_FLAG must be set for acceptance tests. Skip the LTS acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckLtsStructConfigCustom(t *testing.T) {
	if HCS_LTS_STRUCT_CONFIG_TEMPLATE_ID == "" || HCS_LTS_STRUCT_CONFIG_TEMPLATE_NAME == "" {
		t.Skip("HCS_LTS_STRUCT_CONFIG_TEMPLATE_ID and HCS_LTS_STRUCT_CONFIG_TEMPLATE_NAME must be" +
			" set for LTS struct config custom acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckKmsHsmClusterId(t *testing.T) {
	if HCS_KMS_HSM_CLUSTER_ID == "" {
		t.Skip("HCS_KMS_HSM_CLUSTER_ID must be set for KMS dedicated keystore acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckRdsInstance(t *testing.T) {
	if HCS_RDS_INSTANCE_ID == "" || HCS_RDS_BACKUP_ID == "" {
		t.Skip("HCS_RDS_INSTANCE_ID and HCS_RDS_BACKUP_ID must be set for RDS PostgreSql Restore acceptance tests.")
	}
}
