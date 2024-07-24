package dew

import (
	"fmt"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dew"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/kms/v1/keys"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getKmsKeyResourceFunc(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.KmsKeyV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating kms client: %s", err)
	}
	key, err := keys.Get(client, state.Primary.ID).ExtractKeyInfo()

	if err == nil && key.KeyState == dew.PendingDeletionState {
		return nil, golangsdk.ErrDefault404{}
	}
	return key, err
}

// Keystore_id scenario testing is currently not supported.
func TestAccKmsKey_Basic(t *testing.T) {
	var keyAlias = acceptance.RandomAccResourceName()
	var keyAliasUpdate = acceptance.RandomAccResourceName()
	var resourceName = "hcs_kms_key.key_1"

	var key keys.Key

	rc := acceptance.InitResourceCheck(
		resourceName,
		&key,
		getKmsKeyResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckKms(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_Basic(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(resourceName, "rotation_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "region", acceptance.HCS_REGION_NAME),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"pending_days",
				},
			},
			{
				Config: testAccKmsKeyUpdate(keyAliasUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAliasUpdate),
					resource.TestCheckResourceAttr(resourceName, "key_description", "key update description"),
					resource.TestCheckResourceAttr(resourceName, "region", acceptance.HCS_REGION_NAME),
				),
			},
		},
	})
}

func TestAccKmsKey_Enable(t *testing.T) {
	var rName = acceptance.RandomAccResourceName()
	var resourceName = "hcs_kms_key.key_1"

	var key keys.Key
	rc := acceptance.InitResourceCheck(
		resourceName,
		&key,
		getKmsKeyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckKms(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_enabled(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
				),
			},
			{
				Config: testAccKmsKey_disabled(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
				),
			},
			{
				Config: testAccKmsKey_enabled(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
				),
			},
		},
	})
}

func TestAccKmsKey_WithTags(t *testing.T) {
	var keyAlias = acceptance.RandomAccResourceName()
	var resourceName = "hcs_kms_key.key_1"

	var key keys.Key
	rc := acceptance.InitResourceCheck(
		resourceName,
		&key,
		getKmsKeyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckKms(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_WithTags(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
		},
	})
}

func TestAccKmsKey_WithEpsId(t *testing.T) {
	var keyAlias = acceptance.RandomAccResourceName()
	var resourceName = "hcs_kms_key.key_1"

	var key keys.Key
	rc := acceptance.InitResourceCheck(
		resourceName,
		&key,
		getKmsKeyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckKms(t); acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_epsId(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						acceptance.HCS_ENTERPRISE_PROJECT_ID),
				),
			},
		},
	})
}

func TestAccKmsKey_rotation(t *testing.T) {
	var keyAlias = acceptance.RandomAccResourceName()
	var resourceName = "hcs_kms_key.key_1"

	var key keys.Key
	rc := acceptance.InitResourceCheck(
		resourceName,
		&key,
		getKmsKeyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckKms(t); acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_Basic(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(resourceName, "rotation_enabled", "false"),
				),
			},
			{
				Config: testAccKmsKey_rotation(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(resourceName, "rotation_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rotation_interval", "365"),
				),
			},
			{
				Config: testAccKmsKey_rotation_interval(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(resourceName, "rotation_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rotation_interval", "200"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"pending_days",
				},
			},
		},
	})
}

func testAccKmsKey_Basic(keyAlias string) string {
	return fmt.Sprintf(`
resource "hcs_kms_key" "key_1" {
  key_alias    = "%s"
  pending_days = "7"
  region       = "%s"
}
`, keyAlias, acceptance.HCS_REGION_NAME)
}

func testAccKmsKey_WithTags(keyAlias string) string {
	return fmt.Sprintf(`
resource "hcs_kms_key" "key_1" {
  key_alias    = "%s"
  pending_days = "7"
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, keyAlias)
}

func testAccKmsKey_epsId(keyAlias string) string {
	return fmt.Sprintf(`
resource "hcs_kms_key" "key_1" {
  key_alias    = "%s"
  pending_days = "7"

  enterprise_project_id = "%s"
}
`, keyAlias, acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccKmsKeyUpdate(keyAliasUpdate string) string {
	return fmt.Sprintf(`
resource "hcs_kms_key" "key_1" {
  key_alias       = "%s"
  key_description = "key update description"
  pending_days    = "7"
}
`, keyAliasUpdate)
}

func testAccKmsKey_enabled(rName string) string {
	return fmt.Sprintf(`
resource "hcs_kms_key" "key_1" {
  key_description = "Terraform acc test is_enabled %s"
  pending_days    = "7"
  key_alias       = "%s"
}`, rName, rName)
}

func testAccKmsKey_disabled(rName string) string {
	return fmt.Sprintf(`
resource "hcs_kms_key" "key_1" {
  key_description = "Terraform acc test is_enabled %s"
  pending_days    = "7"
  key_alias       = "%s"
  is_enabled      = false
}`, rName, rName)
}

func testAccKmsKey_rotation(rName string) string {
	return fmt.Sprintf(`
resource "hcs_kms_key" "key_1" {
  key_alias        = "%s"
  pending_days     = "7"
  rotation_enabled = true
}`, rName)
}

func testAccKmsKey_rotation_interval(rName string) string {
	return fmt.Sprintf(`
resource "hcs_kms_key" "key_1" {
  key_alias         = "%s"
  pending_days      = "7"
  rotation_enabled  = true
  rotation_interval = 200
}`, rName)
}
