package smn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/smn/v2/topics"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getResourceSMNTopic(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	smnClient, err := conf.SmnV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SMN client: %s", err)
	}

	return topics.Get(smnClient, state.Primary.ID).Extract()
}

func TestAccSMNV2Topic_basic(t *testing.T) {
	var topic topics.TopicGet
	resourceName := "hcs_smn_topic.topic_1"
	rName := acceptance.RandomAccResourceNameWithDash()
	displayName := fmt.Sprintf("The display name of %s", rName)
	update_displayName := fmt.Sprintf("The update display name of %s", rName)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&topic,
		getResourceSMNTopic,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSMNV2TopicConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", displayName),
				),
			},
			{
				Config: testAccSMNV2TopicConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", update_displayName),
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

func TestAccSMNV2Topic_withEpsID(t *testing.T) {
	var topic topics.TopicGet
	resourceName := "hcs_smn_topic.topic_1"
	rName := acceptance.RandomAccResourceNameWithDash()
	displayName := fmt.Sprintf("The display name of %s", rName)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&topic,
		getResourceSMNTopic,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSMNV2TopicConfig_withEpsID(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(
						resourceName, "enterprise_project_id", acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccSMNV2TopicConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcs_smn_topic" "topic_1" {
  name         = "%s"
  display_name = "The display name of %s"
}
`, rName, rName)
}

func testAccSMNV2TopicConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "hcs_smn_topic" "topic_1" {
  name         = "%s"
  display_name = "The update display name of %s"
}
`, rName, rName)
}

func testAccSMNV2TopicConfig_withEpsID(rName string) string {
	return fmt.Sprintf(`
resource "hcs_smn_topic" "topic_1" {
  name                  = "%s"
  display_name          = "The display name of %s"
  enterprise_project_id = "%s"
}
`, rName, rName, acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST)
}
