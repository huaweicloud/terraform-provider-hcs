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

func getDmsRocketMQInstanceResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getRocketmqInstance: Query DMS rocketmq instance
	var (
		getRocketmqInstanceHttpUrl = "v2/{project_id}/instances/{instance_id}"
		getRocketmqInstanceProduct = "dmsv2"
	)
	getRocketmqInstanceClient, err := cfg.NewServiceClient(getRocketmqInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DmsRocketMQInstance Client: %s", err)
	}

	getRocketmqInstancePath := getRocketmqInstanceClient.Endpoint + getRocketmqInstanceHttpUrl
	getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{project_id}", getRocketmqInstanceClient.ProjectID)
	getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{instance_id}", fmt.Sprintf("%v", state.Primary.ID))

	getRocketmqInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRocketmqInstanceResp, err := getRocketmqInstanceClient.Request("GET", getRocketmqInstancePath, &getRocketmqInstanceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DmsRocketMQInstance: %s", err)
	}
	return utils.FlattenResponse(getRocketmqInstanceResp)
}

func TestAccDmsRocketMQInstance_basic(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_dms_rocketmq_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsRocketMQInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.x"),
					resource.TestCheckResourceAttr(resourceName, "enable_acl", "true"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				Config: testDmsRocketMQInstance_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.x"),
					resource.TestCheckResourceAttr(resourceName, "enable_acl", "false"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
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

func TestAccDmsRocketMQInstance_broker_publicip(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_dms_rocketmq_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsRocketMQInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQInstance_broker_publicip(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enable_acl", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_publicip", "true"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "broker_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_space", "300"),
					resource.TestCheckResourceAttrSet(resourceName, "publicip_id"),
					resource.TestCheckResourceAttrSet(resourceName, "publicip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_broker_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_namesrv_address"),
				),
			},
			{
				Config: testDmsRocketMQInstance_broker_publicip_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "enable_acl", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_publicip", "true"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "broker_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "storage_space", "800"),
					resource.TestCheckResourceAttrSet(resourceName, "publicip_id"),
					resource.TestCheckResourceAttrSet(resourceName, "publicip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_broker_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_namesrv_address"),
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

func TestAccDmsRocketMQInstance_publicip(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_dms_rocketmq_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsRocketMQInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQInstance_publicip(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enable_publicip", "false"),
				),
			},
			{
				Config: testDmsRocketMQInstance_publicip_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "enable_publicip", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "publicip_id"),
					resource.TestCheckResourceAttrSet(resourceName, "publicip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_broker_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_namesrv_address"),
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

func TestAccDmsRocketMQInstance_updateWithEpsId(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_dms_rocketmq_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsRocketMQInstanceResourceFunc,
	)
	srcEPS := acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST
	destEPS := acceptance.HCS_ENTERPRISE_MIGRATE_PROJECT_ID_TEST

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQInstance_withEpsId(rName, srcEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testDmsRocketMQInstance_withEpsId(rName, destEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", destEPS),
				),
			},
		},
	})
}

func TestAccDmsRocketMQInstance_spec_code(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_dms_rocketmq_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsRocketMQInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQInstance_spec_code(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.x"),
					resource.TestCheckResourceAttr(resourceName, "enable_acl", "true"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "storage_spec_code", "dms.physical.storage.high.v2"),
				),
			},
			{
				Config: testDmsRocketMQInstance_spec_code_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.x"),
					resource.TestCheckResourceAttr(resourceName, "enable_acl", "false"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "storage_spec_code", "dms.physical.storage.ultra.v2"),
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

func testDmsRocketMQInstance_basic(name string) string {
	return fmt.Sprintf(`
%s

data "hcs_availability_zones" "test" {}

resource "hcs_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = "5.x"
  storage_space     = 300
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  availability_zones = [
    data.hcs_availability_zones.test.names[0],
    data.hcs_availability_zones.test.names[1],
    data.hcs_availability_zones.test.names[2],
  ]

  flavor_id         = "c6.4u8g.cluster"
  storage_spec_code = "dms.physical.storage.high.v2"
  broker_num        = 1
  enable_acl        = true
}
`, common.TestBaseNetwork(name), name)
}

func testDmsRocketMQInstance_update(name string) string {
	return fmt.Sprintf(`
%s

data "hcs_availability_zones" "test" {}

resource "hcs_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = "5.x"
  storage_space     = 300
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  availability_zones = [
    data.hcs_availability_zones.test.names[2],
    data.hcs_availability_zones.test.names[0],
    data.hcs_availability_zones.test.names[1],
  ]

  flavor_id         = "c6.4u8g.cluster"
  storage_spec_code = "dms.physical.storage.high.v2"
  broker_num        = 1
  enable_acl        = false
}
`, common.TestBaseNetwork(name), name)
}

func testDmsRocketMQInstance_broker_publicip(name string) string {
	return fmt.Sprintf(`
%s

data "hcs_availability_zones" "test" {}

resource "hcs_vpc_eip" "test_eip" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "test_eip_${count.index}"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }

  count = 3
}

locals {
  publicip_id = "${hcs_vpc_eip.test_eip[0].id},${hcs_vpc_eip.test_eip[1].id},${hcs_vpc_eip.test_eip[2].id}"
}

resource "hcs_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = "5.x"
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  availability_zones = [
    data.hcs_availability_zones.test.names[0],
    data.hcs_availability_zones.test.names[1],
    data.hcs_availability_zones.test.names[2],
  ]

  storage_spec_code = "dms.physical.storage.high.v2"
  enable_acl        = true
  flavor_id         = "c6.2u8g.single.x86"
  storage_space     = 300  
  broker_num        = 1
  enable_publicip   = true
  publicip_id       = local.publicip_id
}
`, common.TestBaseNetwork(name), name)
}

func testDmsRocketMQInstance_broker_publicip_update(name string) string {
	return fmt.Sprintf(`
%s

data "hcs_availability_zones" "test" {}

resource "hcs_vpc_eip" "test_eip" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "test_eip_${count.index}"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }

  count = 6
}

locals {
  publicip_id = format("%%s,%%s,%%s,%%s,%%s,%%s", hcs_vpc_eip.test_eip[0].id,hcs_vpc_eip.test_eip[1].id,
  hcs_vpc_eip.test_eip[2].id,hcs_vpc_eip.test_eip[3].id,hcs_vpc_eip.test_eip[4].id,hcs_vpc_eip.test_eip[5].id)
}

resource "hcs_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = "5.x"
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  availability_zones = [
    data.hcs_availability_zones.test.names[0],
    data.hcs_availability_zones.test.names[1],
    data.hcs_availability_zones.test.names[2],
  ]

  storage_spec_code = "dms.physical.storage.high.v2"
  enable_acl        = false
  flavor_id         = "c6.2u8g.single.x86"
  storage_space     = 800
  broker_num        = 2  
  enable_publicip   = true
  publicip_id       = local.publicip_id
}
`, common.TestBaseNetwork(name), name)
}

func testDmsRocketMQInstance_publicip(name string) string {
	return fmt.Sprintf(`
%[1]s

data "hcs_availability_zones" "test" {}

resource "hcs_dms_rocketmq_instance" "test" {
  name              = "%[2]s"
  engine_version    = "5.x"
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  availability_zones = [
    data.hcs_availability_zones.test.names[0],
  ]

  storage_spec_code = "dms.physical.storage.high.v2"
  flavor_id         = "c6.2u8g.single.x86"
  storage_space     = 300
  enable_publicip   = false
}
`, common.TestBaseNetwork(name), name)
}

func testDmsRocketMQInstance_publicip_update(name string) string {
	return fmt.Sprintf(`
%[1]s

data "hcs_availability_zones" "test" {}

resource "hcs_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "hcs_dms_rocketmq_instance" "test" {
  name              = "%[2]s"
  engine_version    = "5.x"
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  availability_zones = [
    data.hcs_availability_zones.test.names[0],
  ]

  storage_spec_code = "dms.physical.storage.high.v2"
  flavor_id         = "c6.2u8g.single.x86"
  storage_space     = 300
  enable_publicip   = true
  publicip_id       = hcs_vpc_eip.test.id
}
`, common.TestBaseNetwork(name), name)
}

func testDmsRocketMQInstance_withEpsId(name, epsId string) string {
	return fmt.Sprintf(`
%s

data "hcs_availability_zones" "test" {}

resource "hcs_dms_rocketmq_instance" "test" {
  name                  = "%s"
  engine_version        = "5.x"
  storage_space         = 300
  vpc_id                = hcs_vpc.test.id
  subnet_id             = hcs_vpc_subnet.test.id
  security_group_id     = hcs_networking_secgroup.test.id
  enterprise_project_id = "%s"

  availability_zones = [
    data.hcs_availability_zones.test.names[0],
    data.hcs_availability_zones.test.names[1],
    data.hcs_availability_zones.test.names[2],
  ]
  
  flavor_id         = "c6.4u8g.cluster"
  storage_spec_code = "dms.physical.storage.high.v2"
  broker_num        = 1
  enable_acl        = true

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}`, common.TestBaseNetwork(name), name, epsId)
}

func testDmsRocketMQInstance_spec_code(name string) string {
	return fmt.Sprintf(`
%s

data "hcs_availability_zones" "test" {}

resource "hcs_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = "5.x"
  storage_space     = 300
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  availability_zones = [
    data.hcs_availability_zones.test.names[0],
    data.hcs_availability_zones.test.names[1],
    data.hcs_availability_zones.test.names[2],
  ]

  flavor_id         = "c6.4u8g.cluster"
  storage_spec_code = "dms.physical.storage.high.v2"
  broker_num        = 1
  enable_acl        = true
}
`, common.TestBaseNetwork(name), name)
}

func testDmsRocketMQInstance_spec_code_update(name string) string {
	return fmt.Sprintf(`
%s

data "hcs_availability_zones" "test" {}

resource "hcs_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = "5.x"
  storage_space     = 300
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  availability_zones = [
    data.hcs_availability_zones.test.names[0],
    data.hcs_availability_zones.test.names[1],
    data.hcs_availability_zones.test.names[2],
  ]

  flavor_id         = "c6.4u8g.cluster"
  storage_spec_code = "dms.physical.storage.ultra.v2"
  broker_num        = 1
  enable_acl        = true
}
`, common.TestBaseNetwork(name), name)
}
