package dew

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/csms/v1/secrets"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func getCsmsSecretFunc(c *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.KmsV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS client: %s", err)
	}
	name := state.Primary.Attributes["name"]
	return secrets.Get(client, name)
}

func TestAccDewCsmsSecret_basic(t *testing.T) {
	var (
		secret       secrets.Secret
		name         = acceptance.RandomAccResourceName()
		resourceName = "hcs_csms_secret.test"
		secretText   = utils.HashAndHexEncode("this is a password")
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secret,
		getCsmsSecretFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDewCsmsSecret_basic(name, secretText),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "secret_text", secretText),
					resource.TestCheckResourceAttr(resourceName, "description", "csms secret test"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDewCsmsSecret_kms_key(t *testing.T) {
	var (
		secret       secrets.Secret
		name         = acceptance.RandomAccResourceName()
		resourceName = "hcs_csms_secret.test"
		secretText   = utils.HashAndHexEncode("this is a password")
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secret,
		getCsmsSecretFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDewCsmsSecret_kms_key(name, secretText),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "secret_text", secretText),
					resource.TestCheckResourceAttr(resourceName, "description", "csms secret test"),
					resource.TestCheckResourceAttrPair(resourceName, "kms_key_id", "hcs_kms_key.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
				),
			},
			{
				Config: testAccDewCsmsSecret_kms_key_update(name, secretText),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "secret_text", secretText),
					resource.TestCheckResourceAttr(resourceName, "description", "csms secret test update"),
					resource.TestCheckResourceAttrPair(resourceName, "kms_key_id", "hcs_kms_key.test_second", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value2"),
				),
			},
		},
	})
}

func testAccDewCsmsSecret_base(name string) string {
	return fmt.Sprintf(`
resource "hcs_kms_key" "test" {
  key_alias    = "%[1]s"
  pending_days = "7"
}

resource "hcs_kms_key" "test_second" {
  key_alias    = "%[1]s_second"
  pending_days = "7"
}
`, name)
}

func testAccDewCsmsSecret_basic(name, secretText string) string {
	return fmt.Sprintf(`
resource "hcs_csms_secret" "test" {
  name        = "%[1]s"
  secret_text = "%[2]s"
  description = "csms secret test"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, secretText)
}

func testAccDewCsmsSecret_kms_key(name, secretText string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_csms_secret" "test" {
  name        = "%[1]s"
  secret_text = "%[2]s"
  description = "csms secret test"
  kms_key_id  = hcs_kms_key.test.id

  tags = {
    foo = "bar1"
    key = "value1"
  }
}
`, testAccDewCsmsSecret_base(name), name, secretText)
}

func testAccDewCsmsSecret_kms_key_update(name, secretText string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_csms_secret" "test" {
  name        = "%[1]s"
  secret_text = "%[2]s"
  description = "csms secret test update"
  kms_key_id  = hcs_kms_key.test_second.id

  tags = {
    foo = "bar2"
    key = "value2"
  }
}
`, testAccDewCsmsSecret_base(name), name, secretText)
}
