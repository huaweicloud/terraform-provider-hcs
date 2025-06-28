package vdc

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	sdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/project"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
	"testing"
)

func TestAccVdcProjectResourceCreate(t *testing.T) {
	var project sdk.QueryProjectDetailV31
	resourceName := "hcs_vdc_project.test"

	rName := acceptance.RandomAccResourceNameWithPrefix(acceptance.HCS_REGION_NAME)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVdcProjectDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVdcProjectCreate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcProjectExists(resourceName, &project),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func TestAccVdcProjectResourceUpdate(t *testing.T) {
	var project sdk.QueryProjectDetailV31
	resourceName := "hcs_vdc_project.test"

	rName := acceptance.RandomAccResourceNameWithPrefix(acceptance.HCS_REGION_NAME)

	rNameUpdate := acceptance.RandomAccResourceNameWithPrefix(acceptance.HCS_REGION_NAME)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVdcProjectDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVdcProjectCreate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcProjectExists(resourceName, &project),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				Config: testAccResourceVdcProjectUpdate(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcProjectExists(resourceName, &project),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
				),
			},
		},
	})
}

func TestAccVdcProjectResourceDelete(t *testing.T) {
	var project sdk.QueryProjectDetailV31
	resourceName := "hcs_vdc_project.test"
	rName := acceptance.RandomAccResourceNameWithPrefix(acceptance.HCS_REGION_NAME)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVdcProjectDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVdcProjectCreate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcProjectExists(resourceName, &project),
				),
			},
			{
				Config:  " ",
				Destroy: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcProjectDestroy(resourceName),
				),
			},
		},
	})
}

func testAccResourceVdcProjectCreate(rName string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_project" "test" {
		  vdc_id = "%s"
		  name = "%s"
		}`,
		acceptance.HCS_VDC_PROJECT_VDC_ID,
		rName)
}

func testAccResourceVdcProjectUpdate(rName string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_project" "test" {
		  vdc_id = "%s"
		  name = "%s"
		}`,
		acceptance.HCS_VDC_PROJECT_VDC_ID,
		rName)
}
func testAccCheckVdcProjectExists(n string, project *sdk.QueryProjectDetailV31) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		hcsConfig := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		vdcClient, err := hcsConfig.VdcClient(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating hcs vdc client: %s", err)
		}

		found, err := sdk.Get(vdcClient, rs.Primary.ID, sdk.GetOpts{}).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmtp.Errorf("vdc not found")
		}

		project = found

		return nil
	}
}

func testAccCheckVdcProjectDestroy(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		hcsConfig := config.GetHcsConfig(acceptance.TestAccProvider.Meta())

		vdcClient, err := hcsConfig.VdcClient(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating hcs vdc client: %s", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceName {
				continue
			}
			if rs.Primary.ID == "" {
				return fmt.Errorf("resource ID is empty")
			}
			project, err := sdk.Get(vdcClient, rs.Primary.ID, sdk.GetOpts{}).Extract()
			if err == nil {
				return fmtp.Errorf("Vdc still exists")
			}

			if project.Id != rs.Primary.ID {
				return fmtp.Errorf("resource ID does not match: expected=%s, got=%s", rs.Primary.ID, project.Id)
			}
		}

		return nil
	}
}
