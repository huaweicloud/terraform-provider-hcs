package bms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/bms/v1/baremetalservers"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func TestAccBmsInstance_basic(t *testing.T) {
	var instance baremetalservers.CloudServer

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "hcs_bms_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckBms(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBmsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBmsInstance_basic(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBmsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				Config: testAccBmsInstance_basic(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBmsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
		},
	})
}

func testAccCheckBmsInstanceDestroy(s *terraform.State) error {
	cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	bmsClient, err := cfg.BmsV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating HuaweiCloudStack bms client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_bms_instance" {
			continue
		}
		opts := &baremetalservers.ListOpts{
			Tags:         "__type_baremetal",
			ExpectFields: "id,tags",
		}
		server, err := baremetalservers.Get(bmsClient, rs.Primary.ID, opts).Extract()
		if err == nil {
			if server.Status != "DELETED" {
				return fmt.Errorf("instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckBmsInstanceExists(n string, instance *baremetalservers.CloudServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		bmsClient, err := cfg.BmsV1Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloudStack bms client: %s", err)
		}

		opts := &baremetalservers.ListOpts{
			Tags:         "__type_baremetal",
			ExpectFields: "id,tags",
		}
		found, err := baremetalservers.Get(bmsClient, rs.Primary.ID, opts).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Instance not found")
		}

		*instance = *found

		return nil
	}
}

func testAccBmsInstance_base(rName string) string {
	return fmt.Sprintf(`
%s

data "hcs_bms_flavors" "test" {
  availability_zone = try(element(data.hcs_availability_zones.test.names, 0), "")
}

resource "hcs_kps_keypair" "test" {
  name = "%s"
}`, common.TestBaseNetwork(rName), rName)
}

func testAccBmsInstance_basic(rName string, isAutoRenew bool) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_vpc_eip" "myeip" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%[2]s"
    size        = 8
    share_type  = "PER"
  }
}

resource "hcs_bms_instance" "test" {
  security_groups   = [hcs_networking_secgroup.test.id]
  availability_zone = data.hcs_availability_zones.test.names[0]
  vpc_id            = hcs_vpc.test.id
  flavor_id         = ""
  key_pair          = ""
  image_id          = "519ea918-1fea-4ebc-911a-593739b1a3bc" # CentOS 7.4 64bit for BareMetal

  name                  = "%[2]s"
  user_id               = "%[3]s"
  enterprise_project_id = "%[4]s"

  nics {
    subnet_id = hcs_vpc_subnet.test.id
  }
}
`, testAccBmsInstance_base(rName), rName, acceptance.HCS_USER_ID, acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST, isAutoRenew)
}
