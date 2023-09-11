package smn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccDataTopics_basic(t *testing.T) {
	dataSourceName := "data.hcs_smn_topics.test"
	resourcerName := "hcs_smn_topic.topic_1"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTopicsConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttrPair(dataSourceName, "topics.0.id", resourcerName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "topics.0.topic_urn", resourcerName, "topic_urn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "topics.0.display_name", resourcerName, "display_name"),
				),
			},
		},
	})
}

func testAccDataTopicsConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "hcs_smn_topics" "test" {
  name = "%s"

  depends_on = [
    hcs_smn_topic.topic_1
  ]
}
`, testAccSMNV2TopicConfig_basic(rName), rName)
}
