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

func getDmsRocketMQConsumerGroupResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getRocketmqConsumerGroup: query DMS rocketmq consumer group
	var (
		getRocketmqConsumerGroupHttpUrl = "v2/{project_id}/instances/{instance_id}/groups/{group}"
		getRocketmqConsumerGroupProduct = "dmsv2"
	)
	getRocketmqConsumerGroupClient, err := cfg.NewServiceClient(getRocketmqConsumerGroupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DmsRocketMQConsumerGroup Client: %s", err)
	}

	// Split instance_id and group from resource id
	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<consumerGroup>")
	}
	instanceID := parts[0]
	name := parts[1]
	getRocketmqConsumerGroupPath := getRocketmqConsumerGroupClient.Endpoint + getRocketmqConsumerGroupHttpUrl
	getRocketmqConsumerGroupPath = strings.ReplaceAll(getRocketmqConsumerGroupPath, "{project_id}",
		getRocketmqConsumerGroupClient.ProjectID)
	getRocketmqConsumerGroupPath = strings.ReplaceAll(getRocketmqConsumerGroupPath, "{instance_id}", instanceID)
	getRocketmqConsumerGroupPath = strings.ReplaceAll(getRocketmqConsumerGroupPath, "{group}", name)

	getRocketmqConsumerGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRocketmqConsumerGroupResp, err := getRocketmqConsumerGroupClient.Request("GET", getRocketmqConsumerGroupPath,
		&getRocketmqConsumerGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DmsRocketMQConsumerGroup: %s", err)
	}
	return utils.FlattenResponse(getRocketmqConsumerGroupResp)
}

func TestAccDmsRocketMQConsumerGroup_basic(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_dms_rocketmq_consumer_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsRocketMQConsumerGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQConsumerGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "broadcast", "true"),
					resource.TestCheckResourceAttr(resourceName, "retry_max_times", "3"),
					resource.TestCheckResourceAttr(resourceName, "description", "add description."),
				),
			},
			{
				Config: testDmsRocketMQConsumerGroup_basic_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "broadcast", "false"),
					resource.TestCheckResourceAttr(resourceName, "retry_max_times", "5"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
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

func TestAccDmsRocketMQConsumerGroup_ebable(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_dms_rocketmq_consumer_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsRocketMQConsumerGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQConsumerGroup_enable(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "broadcast", "true"),
					resource.TestCheckResourceAttr(resourceName, "retry_max_times", "3"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "add description."),
				),
			},
			{
				Config: testDmsRocketMQConsumerGroup_enable_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "broadcast", "false"),
					resource.TestCheckResourceAttr(resourceName, "retry_max_times", "5"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "update description"),
				),
			},
		},
	})
}

func TestAccDmsRocketMQConsumerGroup_consume(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_dms_rocketmq_consumer_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsRocketMQConsumerGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQConsumerGroup_consume(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "broadcast", "true"),
					resource.TestCheckResourceAttr(resourceName, "consume_orderly", "false"),
					resource.TestCheckResourceAttr(resourceName, "retry_max_times", "3"),
					resource.TestCheckResourceAttr(resourceName, "description", "add description."),
				),
			},
			{
				Config: testDmsRocketMQConsumerGroup_consume_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "broadcast", "false"),
					resource.TestCheckResourceAttr(resourceName, "consume_orderly", "true"),
					resource.TestCheckResourceAttr(resourceName, "retry_max_times", "5"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
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

func testAccDmsRocketmqConsumerGroup_Base(rName string) string {
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

func testDmsRocketMQConsumerGroup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dms_rocketmq_consumer_group" "test" {
  instance_id = hcs_dms_rocketmq_instance.test.id
  broadcast   = true

  name            = "%s"
  retry_max_times = "3"
  description     = "add description."
}
`, testAccDmsRocketmqConsumerGroup_Base(name), name)
}

func testDmsRocketMQConsumerGroup_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dms_rocketmq_consumer_group" "test" {
  instance_id = hcs_dms_rocketmq_instance.test.id
  broadcast   = false

  name            = "%s"
  retry_max_times = "5"
  description     = ""
}
`, testAccDmsRocketmqConsumerGroup_Base(name), name)
}

func testDmsRocketMQConsumerGroup_enable(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dms_rocketmq_consumer_group" "test" {
  instance_id = hcs_dms_rocketmq_instance.test.id
  broadcast   = true

  name            = "%s"
  retry_max_times = "3"
  description     = "add description."

  enable = true
}
`, testAccDmsRocketmqConsumerGroup_Base(name), name)
}

func testDmsRocketMQConsumerGroup_enable_update(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dms_rocketmq_consumer_group" "test" {
  instance_id = hcs_dms_rocketmq_instance.test.id
  broadcast   = false

  name            = "%s"
  retry_max_times = "5"
  description     = "update description"

  enable = false
}
`, testAccDmsRocketmqConsumerGroup_Base(name), name)
}

func testDmsRocketMQConsumerGroup_consume(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dms_rocketmq_consumer_group" "test" {
  instance_id = hcs_dms_rocketmq_instance.test.id
  broadcast   = true

  name            = "%s"
  retry_max_times = "3"
  description     = "add description."

  consume_orderly = false
}
`, testAccDmsRocketmqConsumerGroup_Base(name), name)
}

func testDmsRocketMQConsumerGroup_consume_update(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dms_rocketmq_consumer_group" "test" {
  instance_id = hcs_dms_rocketmq_instance.test.id
  broadcast   = false

  name            = "%s"
  retry_max_times = "5"
  description     = ""

  consume_orderly = true
}
`, testAccDmsRocketmqConsumerGroup_Base(name), name)
}
