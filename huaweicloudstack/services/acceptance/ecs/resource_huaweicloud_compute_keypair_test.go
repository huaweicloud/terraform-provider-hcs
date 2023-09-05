package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/compute/v2/extensions/keypairs"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccComputeV2Keypair_basic(t *testing.T) {
	var keypair keypairs.KeyPair
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "hcs_ecs_compute_keypair.test"
	publicKey, _, _ := acctest.RandSSHKeyPair("Generated-by-AccTest")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2KeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Keypair_basic(rName, publicKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2KeypairExists(resourceName, &keypair),
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

func TestAccComputeV2Keypair_privateKey(t *testing.T) {
	var keypair keypairs.KeyPair
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "hcs_ecs_compute_keypair.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2KeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Keypair_privateKey(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2KeypairExists(resourceName, &keypair),
					resource.TestCheckResourceAttrSet(resourceName, "key_file"),
				),
			},
		},
	})
}

func testAccCheckComputeV2KeypairDestroy(s *terraform.State) error {
	cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	computeClient, err := cfg.ComputeV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating hcs compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_ecs_compute_keypair" {
			continue
		}

		_, err := keypairs.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Keypair still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2KeypairExists(n string, kp *keypairs.KeyPair) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		computeClient, err := cfg.ComputeV2Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating hcs compute client: %s", err)
		}

		found, err := keypairs.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Keypair not found")
		}

		*kp = *found

		return nil
	}
}

func testAccComputeV2Keypair_basic(rName, keypair string) string {
	return fmt.Sprintf(`
resource "hcs_ecs_compute_keypair" "test" {
  name       = "%s"
  public_key = "%s"
}
`, rName, keypair)
}

func testAccComputeV2Keypair_privateKey(rName string) string {
	return fmt.Sprintf(`
resource "hcs_ecs_compute_keypair" "test" {
  name = "%s"
}
`, rName)
}
