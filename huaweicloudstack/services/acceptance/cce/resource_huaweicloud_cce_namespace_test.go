package cce

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/cce/v1/namespaces"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
)

func getNamespaceResourceFunc(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CceV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE v1 client: %s", err)
	}
	resp, err := namespaces.Get(c, state.Primary.Attributes["cluster_id"],
		state.Primary.Attributes["name"]).Extract()
	if resp == nil && err == nil {
		return resp, fmt.Errorf("unable to find the namespace (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccCCENamespaceV1_basic(t *testing.T) {
	var namespace namespaces.Namespace
	resourceName := "hcs_cce_namespace.test"
	randName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&namespace,
		getNamespaceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCCENamespaceV1_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cluster_id",
						acceptance.HCS_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCENamespaceImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccCCENamespaceV1_generateName(t *testing.T) {
	var namespace namespaces.Namespace
	resourceName := "hcs_cce_namespace.test"
	randName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&namespace,
		getNamespaceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCCENamespaceV1_generateName(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cluster_id",
						acceptance.HCS_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttr(resourceName, "prefix", randName),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(fmt.Sprintf(`^%s[a-z0-9-]*`, randName))),
				),
			},
		},
	})
}

func testAccCCENamespaceImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmtp.Errorf("unable to find the namespace in the terraform state storage")
		}
		if rs.Primary.Attributes["cluster_id"] == "" || rs.Primary.Attributes["name"] == "" {
			return "", fmtp.Errorf("resource not found: %s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.Attributes["name"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.Attributes["name"]), nil

	}
}

func testAccCCENamespaceV1_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcs_cce_namespace" "test" {
  cluster_id = "%[1]s"
  name       = "%[2]s"
}
`, acceptance.HCS_CCE_CLUSTER_ID, rName)
}

func testAccCCENamespaceV1_generateName(rName string) string {
	return fmt.Sprintf(`
resource "hcs_cce_namespace" "test" {
  cluster_id = "%[1]s"
  prefix     = "%[2]s"
}
`, acceptance.HCS_CCE_CLUSTER_ID, rName)
}
