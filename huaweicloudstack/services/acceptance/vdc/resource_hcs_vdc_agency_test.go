package vdc

import (
	"fmt"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/agency"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceVdcAgency_basic(t *testing.T) {

	var agencyDetail agency.AgencyDetail

	resourceName := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVdcAgencyDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccVdcAgencyBasic(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcAgencyExists(resourceName, &agencyDetail),
					resource.TestCheckResourceAttr(resourceName, "name", resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "delegated_domain_name"),
					resource.TestCheckNoResourceAttr(resourceName, "project_role"),
					resource.TestCheckNoResourceAttr(resourceName, "domain_roles"),
					resource.TestCheckNoResourceAttr(resourceName, "all_resources_roles"),
				),
			},
			{
				Config: testAccVdcAgencyUpdate1(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcAgencyExists(resourceName, &agencyDetail),
					resource.TestCheckResourceAttr(resourceName, "name", resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "delegated_domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "project_role"),
					resource.TestCheckResourceAttr(resourceName, "project_role.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.1", "VDC Readonly"),
					resource.TestCheckNoResourceAttr(resourceName, "domain_roles"),
					resource.TestCheckNoResourceAttr(resourceName, "all_resources_roles"),
				),
			},
			{
				Config: testAccVdcAgencyUpdate2(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcAgencyExists(resourceName, &agencyDetail),
					resource.TestCheckResourceAttr(resourceName, "name", resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "delegated_domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "project_role"),
					resource.TestCheckResourceAttr(resourceName, "project_role.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.1", "VDC Readonly"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_roles"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.1", "Application Admin"),
					resource.TestCheckNoResourceAttr(resourceName, "all_resources_roles"),
				),
			},
			{
				Config: testAccVdcAgencyUpdate3(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcAgencyExists(resourceName, &agencyDetail),
					resource.TestCheckResourceAttr(resourceName, "name", resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "delegated_domain_name"),
					resource.TestCheckNoResourceAttr(resourceName, "project_role"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_roles"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.1", "Application Admin"),
					resource.TestCheckNoResourceAttr(resourceName, "all_resources_roles"),
				),
			},
			{
				Config: testAccVdcAgencyUpdate4(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcAgencyExists(resourceName, &agencyDetail),
					resource.TestCheckResourceAttr(resourceName, "name", resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "delegated_domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "project_role"),
					resource.TestCheckResourceAttr(resourceName, "project_role.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.1", "VDC Readonly"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_roles"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.1", "Application Admin"),
					resource.TestCheckNoResourceAttr(resourceName, "all_resources_roles"),
				),
			},
			{
				Config: testAccVdcAgencyUpdate5(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcAgencyExists(resourceName, &agencyDetail),
					resource.TestCheckResourceAttr(resourceName, "name", resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "delegated_domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "project_role"),
					resource.TestCheckResourceAttr(resourceName, "project_role.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.1", "VDC Readonly"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_roles"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.1", "Application Admin"),
					resource.TestCheckNoResourceAttr(resourceName, "all_resources_roles"),
				),
			},
			{
				Config: testAccVdcAgencyUpdate6(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcAgencyExists(resourceName, &agencyDetail),
					resource.TestCheckResourceAttr(resourceName, "name", resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "delegated_domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "project_role"),
					resource.TestCheckResourceAttr(resourceName, "project_role.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.1", "VDC Readonly"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_roles"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.1", "Application Admin"),
					resource.TestCheckNoResourceAttr(resourceName, "all_resources_roles"),
				),
			},
			{
				Config: testAccVdcAgencyUpdate7(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcAgencyExists(resourceName, &agencyDetail),
					resource.TestCheckResourceAttr(resourceName, "name", resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "delegated_domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "project_role"),
					resource.TestCheckResourceAttr(resourceName, "project_role.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.1", "VDC Readonly"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_roles"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.1", "Application Admin"),
					resource.TestCheckNoResourceAttr(resourceName, "all_resources_roles"),
				),
			},
			{
				Config: testAccVdcAgencyUpdate8(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcAgencyExists(resourceName, &agencyDetail),
					resource.TestCheckResourceAttr(resourceName, "name", resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "delegated_domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "project_role"),
					resource.TestCheckResourceAttr(resourceName, "project_role.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.1", "VDC Readonly"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_roles"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.1", "Application Admin"),
					resource.TestCheckNoResourceAttr(resourceName, "all_resources_roles"),
				),
			},
			{
				Config: testAccVdcAgencyUpdate9(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcAgencyExists(resourceName, &agencyDetail),
					resource.TestCheckResourceAttr(resourceName, "name", resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "delegated_domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "project_role"),
					resource.TestCheckResourceAttr(resourceName, "project_role.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.1", "VDC Readonly"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_roles"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.1", "Application Admin"),
					resource.TestCheckNoResourceAttr(resourceName, "all_resources_roles"),
				),
			},
			{
				Config: testAccVdcAgencyUpdate10(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcAgencyExists(resourceName, &agencyDetail),
					resource.TestCheckResourceAttr(resourceName, "name", resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "delegated_domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "project_role"),
					resource.TestCheckResourceAttr(resourceName, "project_role.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.1", "VDC Readonly"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_roles"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.1", "Application Admin"),
					resource.TestCheckNoResourceAttr(resourceName, "all_resources_roles"),
				),
			},
			{
				Config: testAccVdcAgencyUpdate11(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcAgencyExists(resourceName, &agencyDetail),
					resource.TestCheckResourceAttr(resourceName, "name", resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "delegated_domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "project_role"),
					resource.TestCheckResourceAttr(resourceName, "project_role.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.1", "VDC Readonly"),
					resource.TestCheckResourceAttr(resourceName, "project_role.2.role.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "project_role.2.role.1", "VDC Admin"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_roles"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.1", "Application Admin"),
					resource.TestCheckNoResourceAttr(resourceName, "all_resources_roles"),
				),
			},
			{
				Config: testAccVdcAgencyUpdate12(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcAgencyExists(resourceName, &agencyDetail),
					resource.TestCheckResourceAttr(resourceName, "name", resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "delegated_domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "project_role"),
					resource.TestCheckResourceAttr(resourceName, "project_role.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.1", "VDC Readonly"),
					resource.TestCheckResourceAttr(resourceName, "project_role.2.role.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "project_role.2.role.1", "VDC Admin"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_roles"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.1", "Application Admin"),
					resource.TestCheckNoResourceAttr(resourceName, "all_resources_roles"),
				),
			},
			{
				Config: testAccVdcAgencyUpdate13(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcAgencyExists(resourceName, &agencyDetail),
					resource.TestCheckResourceAttr(resourceName, "name", resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "delegated_domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "project_role"),
					resource.TestCheckResourceAttr(resourceName, "project_role.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.1", "VDC Readonly"),
					resource.TestCheckResourceAttr(resourceName, "project_role.2.role.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "project_role.2.role.1", "VDC Admin"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_roles"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.1", "Application Admin"),
					resource.TestCheckNoResourceAttr(resourceName, "all_resources_roles"),
				),
			},
			{
				Config: testAccVdcAgencyUpdate14(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcAgencyExists(resourceName, &agencyDetail),
					resource.TestCheckResourceAttr(resourceName, "name", resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "delegated_domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "project_role"),
					resource.TestCheckResourceAttr(resourceName, "project_role.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "project_role.1.role.1", "VDC Readonly"),
					resource.TestCheckResourceAttr(resourceName, "project_role.2.role.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "project_role.2.role.1", "VDC Admin"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_roles"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.1", "Application Admin"),
					resource.TestCheckNoResourceAttr(resourceName, "all_resources_roles"),
				),
			},
		},
	})
}

func testAccVdcAgencyBasic(name string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_agency" "agency" {
		  name                  = "%s"
		  description           = "aa"
		  delegated_domain_name = "%s"
		}
	`, name, acceptance.HCS_AGENCY_DOMAIN_NAME)
}

func testAccVdcAgencyUpdate1(name string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_agency" "agency" {
		  name                  = "%s"
		  description           = "aa"
		  delegated_domain_name = "%s"
		  project_role {
			project = "%s"
			roles   = ["VDC Readonly"]
		  }
		}
	`, name, acceptance.HCS_AGENCY_DOMAIN_NAME, acceptance.HCS_PROJECT_NAME)
}

func testAccVdcAgencyUpdate2(name string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_agency" "agency" {
		  name                  = "%s"
		  description           = "aa"
		  delegated_domain_name = "%s"
		  project_role {
			project = "%s"
			roles   = ["VDC Readonly"]
		  }
		  domain_roles        = ["Application Admin", "VDC Admin"]
		}
	`, name, acceptance.HCS_AGENCY_DOMAIN_NAME, acceptance.HCS_PROJECT_NAME)
}

func testAccVdcAgencyUpdate3(name string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_agency" "agency" {
		  name                  = "%s"
		  description           = "aa"
		  delegated_domain_name = "%s"
		  domain_roles        = ["Application Admin", "VDC Admin"]
		}
	`, name, acceptance.HCS_AGENCY_DOMAIN_NAME)
}

func testAccVdcAgencyUpdate4(name string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_agency" "agency" {
		  name                  = "%s"
		  description           = "aa"
		  delegated_domain_name = "%s"
		  project_role {
			project = "%s"
			roles   = ["VDC Readonly", "TagAdmin"]
		  }
		  domain_roles        = ["Application Admin", "VDC Admin"]
		}
	`, name, acceptance.HCS_AGENCY_DOMAIN_NAME, acceptance.HCS_PROJECT_NAME)
}

func testAccVdcAgencyUpdate5(name string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_agency" "agency" {
		  name                  = "%s"
		  description           = "aa"
		  delegated_domain_name = "%s"
		  project_role {
			project = "%s"
			roles   = ["VDC Readonly", "TagAdmin"]
		  }
		  domain_roles        = ["Application Admin", "VDC Admin", "Tenant Guest"]
		}
	`, name, acceptance.HCS_AGENCY_DOMAIN_NAME, acceptance.HCS_PROJECT_NAME)
}

func testAccVdcAgencyUpdate6(name string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_agency" "agency" {
		  name                  = "%s"
		  description           = "aa"
		  delegated_domain_name = "%s"
		  project_role {
			project = "%s"
			roles   = ["VDC Readonly", "TagAdmin"]
		  }
		  domain_roles        = ["Application Admin", "VDC Admin"]
		  all_resources_roles = []
		}
	`, name, acceptance.HCS_AGENCY_DOMAIN_NAME, acceptance.HCS_PROJECT_NAME)
}

func testAccVdcAgencyUpdate7(name string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_agency" "agency" {
		  name                  = "%s"
		  description           = "aa"
		  delegated_domain_name = "%s"
		  project_role {
			project = "%s"
			roles   = ["VDC Readonly", "TagAdmin", "Tenant Guest"]
		  }
		  domain_roles        = ["Application Admin", "VDC Admin"]
		  all_resources_roles = []
		}
	`, name, acceptance.HCS_AGENCY_DOMAIN_NAME, acceptance.HCS_PROJECT_NAME)
}

func testAccVdcAgencyUpdate8(name string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_agency" "agency" {
		  name                  = "%s"
		  description           = "aa"
		  delegated_domain_name = "%s"
		  project_role {
			project = "%s"
			roles   = ["VDC Readonly", "TagAdmin"]
		  }
		  domain_roles        = ["Application Admin", "VDC Admin"]
		  all_resources_roles = []
		}
	`, name, acceptance.HCS_AGENCY_DOMAIN_NAME, acceptance.HCS_PROJECT_NAME)
}

func testAccVdcAgencyUpdate9(name string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_agency" "agency" {
		  name                  = "%s"
		  description           = "aa"
		  delegated_domain_name = "%s"
		  project_role {
			project = "%s"
			roles   = ["VDC Readonly", "TagAdmin", "Tenant Guest"]
		  }
		  domain_roles        = ["Application Admin", "VDC Admin", "Tenant Guest"]
		  all_resources_roles = []
		}
	`, name, acceptance.HCS_AGENCY_DOMAIN_NAME, acceptance.HCS_PROJECT_NAME)
}

func testAccVdcAgencyUpdate10(name string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_agency" "agency" {
		  name                  = "%s"
		  description           = "aa"
		  delegated_domain_name = "%s"
		  project_role {
			project = "%s"
			roles   = ["VDC Readonly", "TagAdmin"]
		  }
		  domain_roles        = ["Application Admin", "VDC Admin", "Tenant Guest"]
		  all_resources_roles = []
		}
	`, name, acceptance.HCS_AGENCY_DOMAIN_NAME, acceptance.HCS_PROJECT_NAME)
}

func testAccVdcAgencyUpdate11(name string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_agency" "agency" {
		  name                  = "%s"
		  description           = "aa"
		  delegated_domain_name = "%s"
		  project_role {
			project = "%s"
			roles   = ["VDC Readonly", "TagAdmin"]
		  }
		  project_role {
			project = "%s"
			roles   = ["VDC Readonly", "VDC Admin"]
		  }
		  domain_roles        = ["Application Admin", "VDC Admin", "Tenant Guest"]
		  all_resources_roles = []
		}
	`, name, acceptance.HCS_AGENCY_DOMAIN_NAME, acceptance.HCS_PROJECT_NAME, acceptance.HCS_PROJECT_NAME)
}

func testAccVdcAgencyUpdate12(name string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_agency" "agency" {
		  name                  = "%s"
		  description           = "aa"
		  delegated_domain_name = "%s"
		  project_role {
			project = "%s"
			roles   = ["VDC Readonly", "TagAdmin"]
		  }
		  project_role {
			project = "%s"
			roles   = ["VDC Readonly", "VDC Admin", "Tenant Guest"]
		  }
		  domain_roles        = ["Application Admin", "VDC Admin", "Tenant Guest"]
		  all_resources_roles = []
		}
	`, name, acceptance.HCS_AGENCY_DOMAIN_NAME, acceptance.HCS_PROJECT_NAME, acceptance.HCS_PROJECT_NAME)
}

func testAccVdcAgencyUpdate13(name string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_agency" "agency" {
		  name                  = "%s"
		  description           = "aa"
		  delegated_domain_name = "%s"
		  project_role {
			project = "%s"
			roles   = ["VDC Readonly", "TagAdmin"]
		  }
		  project_role {
			project = "%s"
			roles   = ["VDC Readonly", "VDC Admin"]
		  }
		  domain_roles        = ["Application Admin", "VDC Admin", "Tenant Guest"]
		  all_resources_roles = []
		}
	`, name, acceptance.HCS_AGENCY_DOMAIN_NAME, acceptance.HCS_PROJECT_NAME, acceptance.HCS_PROJECT_NAME)
}

func testAccVdcAgencyUpdate14(name string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_agency" "agency" {
		  name                  = "%s"
		  description           = "aa"
		  delegated_domain_name = "%s"
		  project_role {
			project = "%s"
			roles   = ["VDC Readonly"]
		  }
		  project_role {
			project = "%s"
			roles   = ["VDC Readonly", "VDC Admin", "TagAdmin"]
		  }
		  domain_roles        = ["Application Admin", "VDC Admin", "Tenant Guest"]
		  all_resources_roles = []
		}
	`, name, acceptance.HCS_AGENCY_DOMAIN_NAME, acceptance.HCS_PROJECT_NAME, acceptance.HCS_PROJECT_NAME)
}

func testAccCheckVdcAgencyExists(n string, model *agency.AgencyDetail) resource.TestCheckFunc {
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

		found, err := agency.GetAgency(vdcClient, agency.GetAgencyOpts{AgencyName: n})
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Vdc agency not found")
		}

		*model = *found

		return nil
	}
}

func testAccCheckVdcAgencyDestroy(n string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		hcsConfig := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		vdcClient, err := hcsConfig.VdcClient(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating hcs vdc user client: %s", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hcs_vdc_agency" {
				continue
			}

			_, err := agency.GetAgency(vdcClient, agency.GetAgencyOpts{AgencyName: n})
			if err != nil {
				return fmtp.Errorf("Vdc agency still exists")
			}
		}

		return nil
	}
}
