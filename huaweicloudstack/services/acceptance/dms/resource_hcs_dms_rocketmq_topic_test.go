package dms

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func getDmsRocketMQTopicResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getRocketmqTopic: query DMS rocketmq topic
	var (
		getRocketmqTopicHttpUrl = "v2/{project_id}/instances/{instance_id}/topics/{topic}"
		getRocketmqTopicProduct = "dmsv2"
	)
	getRocketmqTopicClient, err := cfg.NewServiceClient(getRocketmqTopicProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DmsRocketMQTopic Client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<topic>")
	}
	instanceID := parts[0]
	topic := parts[1]
	getRocketmqTopicPath := getRocketmqTopicClient.Endpoint + getRocketmqTopicHttpUrl
	getRocketmqTopicPath = strings.ReplaceAll(getRocketmqTopicPath, "{project_id}", getRocketmqTopicClient.ProjectID)
	getRocketmqTopicPath = strings.ReplaceAll(getRocketmqTopicPath, "{instance_id}", instanceID)
	getRocketmqTopicPath = strings.ReplaceAll(getRocketmqTopicPath, "{topic}", topic)

	getRocketmqTopicOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRocketmqTopicResp, err := getRocketmqTopicClient.Request("GET", getRocketmqTopicPath, &getRocketmqTopicOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DmsRocketMQTopic: %s", err)
	}
	return utils.FlattenResponse(getRocketmqTopicResp)
}

func TestAccDmsRocketMQTopic_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_dms_rocketmq_topic.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDmsRocketMQTopicResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQTopic_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "queue_num", "3"),
					resource.TestCheckResourceAttr(rName, "permission", "all"),
					resource.TestCheckResourceAttr(rName, "message_type", "NORMAL"),
					resource.TestCheckResourceAttr(rName, "total_read_queue_num", "4"),
					resource.TestCheckResourceAttr(rName, "total_write_queue_num", "4"),
				),
			},
			{
				Config: testDmsRocketMQTopic_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "permission", "sub"),
					resource.TestCheckResourceAttr(rName, "queue_num", "4"),
					resource.TestCheckResourceAttr(rName, "message_type", "NORMAL"),
					resource.TestCheckResourceAttr(rName, "total_read_queue_num", "4"),
					resource.TestCheckResourceAttr(rName, "total_write_queue_num", "4"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "queue_num"},
			},
		},
	})
}

func testAccDmsRocketmqTopic_Base(rName string) string {
	return fmt.Sprintf(`
%s

data "hcs_availability_zones" "test" {}

resource "hcs_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = "4.8.0"
  storage_space     = 600
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  availability_zones = [
    data.hcs_availability_zones.test.names[0]
  ]

  flavor_id         = "c6.4u8g.cluster"
  storage_spec_code = "dms.physical.storage.high.v2"
  broker_num        = 1
}
`, common.TestBaseNetwork(rName), rName)
}

func testDmsRocketMQTopic_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dms_rocketmq_topic" "test" {
  instance_id  = hcs_dms_rocketmq_instance.test.id
  name         = "%s"
  queue_num    = 3
  permission   = "all"
  message_type = "NORMAL"

  brokers {
    name      = "broker-0"
    queue_num = 3
  }
}
`, testAccDmsRocketmqTopic_Base(name), name)
}

func testDmsRocketMQTopic_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dms_rocketmq_topic" "test" {
  instance_id  = hcs_dms_rocketmq_instance.test.id
  name         = "%s"
  queue_num    = 4
  permission   = "sub"
  message_type = "NORMAL"

  brokers {
    name      = "broker-0"
    queue_num = 3
  }
}
`, testAccDmsRocketmqTopic_Base(name), name)
}
