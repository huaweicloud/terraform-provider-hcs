package elb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func getSecurityPoliciesV3ResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getSecurityPolicy: Query the ELB security policy
	var (
		getSecurityPolicyHttpUrl = "v3/{project_id}/elb/security-policies/{security_policy_id}"
		getSecurityPolicyProduct = "elb"
	)
	getSecurityPolicyClient, err := cfg.NewServiceClient(getSecurityPolicyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecurityPolicies Client: %s", err)
	}

	getSecurityPolicyPath := getSecurityPolicyClient.Endpoint + getSecurityPolicyHttpUrl
	getSecurityPolicyPath = strings.ReplaceAll(getSecurityPolicyPath, "{project_id}", getSecurityPolicyClient.ProjectID)
	getSecurityPolicyPath = strings.ReplaceAll(getSecurityPolicyPath, "{security_policy_id}", fmt.Sprintf("%v", state.Primary.ID))

	getSecurityPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSecurityPolicyResp, err := getSecurityPolicyClient.Request("GET", getSecurityPolicyPath, &getSecurityPolicyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SecurityPolicies: %s", err)
	}
	return utils.FlattenResponse(getSecurityPolicyResp)
}

func TestAccSecurityPoliciesV3_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_elb_security_policy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSecurityPoliciesV3ResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSecurityPoliciesV3_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "protocols.0", "TLSv1.1"),
					resource.TestCheckResourceAttr(rName, "protocols.1", "TLSv1.2"),
					resource.TestCheckResourceAttr(rName, "ciphers.0", "ECDHE-ECDSA-AES128-SHA"),
					resource.TestCheckResourceAttr(rName, "ciphers.1", "ECDHE-RSA-AES256-SHA"),
				),
			},
			{
				Config: testSecurityPoliciesV3_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "protocols.0", "TLSv1.2"),
					resource.TestCheckResourceAttr(rName, "ciphers.0", "ECDHE-ECDSA-AES128-SHA"),
					resource.TestCheckResourceAttr(rName, "name", name),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testSecurityPoliciesV3_basic(name string) string {
	return fmt.Sprintf(`
resource "hcs_elb_security_policy" "test" {
  protocols = [
    "TLSv1.1",
    "TLSv1.2"
  ]
  ciphers = [
    "ECDHE-ECDSA-AES128-SHA",
    "ECDHE-RSA-AES256-SHA"
  ]
  name = "%s"
}
`, name)
}

func testSecurityPoliciesV3_basic_update(name string) string {
	return fmt.Sprintf(`
resource "hcs_elb_security_policy" "test" {
  protocols = [
    "TLSv1.2"
  ]
  ciphers = [
    "ECDHE-ECDSA-AES128-SHA"
  ]
  name = "%s"
}
`, name)
}
