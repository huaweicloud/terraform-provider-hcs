package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/dcs/v2/instances"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getDcsResourceFunc(c *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DcsV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DCS client(V2): %s", err)
	}
	return instances.Get(client, state.Primary.ID)
}

func TestAccDcsInstances_basic(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = acceptance.RandomAccResourceName()
	var pwd = acceptance.RandomPassword()
	resourceName := "hcs_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_basic(instanceName, pwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "6.0"),
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "0.125"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "18:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.id", "1"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "timeout"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "100"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
				),
			},
			{
				Config: testAccDcsV1Instance_updated(instanceName, pwd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6388"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "02:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.begin_at", "01:00-02:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.save_days", "2"),
					resource.TestCheckResourceAttr(resourceName, "backup_policy.0.backup_at.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.id", "10"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "latency-monitor-threshold"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "120"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password", "auto_renew", "period", "period_unit", "rename_commands",
					"internal_version", "save_days", "backup_type", "begin_at", "period_type", "backup_at", "parameters", "bandwidth_info"},
			},
		},
	})
}

func TestAccDcsInstances_whitelists(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	var pwd = acceptance.RandomPassword()
	resourceName := "hcs_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_whitelists(instanceName, pwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "6.0"),
					resource.TestCheckResourceAttr(resourceName, "whitelist_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.0.group_name", "testgroup1"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.0.ip_address.0", "192.168.10.100"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.0.ip_address.1", "192.168.0.0/24"),
				),
			},
			{
				Config: testAccDcsV1Instance_whitelists_update(instanceName, pwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "6.0"),
					resource.TestCheckResourceAttr(resourceName, "whitelist_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.0.group_name", "testgroup2"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.0.ip_address.0", "172.16.10.100"),
					resource.TestCheckResourceAttr(resourceName, "whitelists.0.ip_address.1", "172.16.0.0/24"),
				),
			},
		},
	})
}

func TestAccDcsInstances_single(t *testing.T) {
	var instance instances.DcsInstance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	var pwd = acceptance.RandomPassword()
	resourceName := "hcs_dcs_instance.instance_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDcsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_single(instanceName, pwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
				),
			},
		},
	})
}

func testAccDcsV1Instance_basic(instanceName, pwd string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" test {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "hcs_vpc_subnet" "test" {
  name        = "%[1]s"
  cidr        = "192.168.1.0/24"
  gateway_ip  = "192.168.1.1"
  vpc_id      = hcs_vpc.test.id
  description = "created by acc test"
}

resource "hcs_dcs_instance" "instance_1" {
  name               = "%[1]s"
  engine_version     = "6.0"
  password           = "%[2]s"
  engine             = "Redis"
  port               = 6388
  capacity           = 0.125
  vpc_id             = resource.hcs_vpc.test.id
  subnet_id          = resource.hcs_vpc_subnet.test.id
  availability_zones = ["az0.dc0"]
  flavor             = "redis.ha.xu1.tiny.r2.128"
  maintain_begin     = "18:00:00"
  maintain_end       = "22:00:00"

  backup_policy {
    backup_type = "auto"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [4]
    save_days   = 1
  }

  rename_commands = {
    command  = "command001"
    keys     = "keys001"
    flushall = "flushall001"
    flushdb  = "flushdb001"
    hgetall  = "hgetall001"
  }

  parameters {
    id    = "1"
    name  = "timeout"
    value = "100"
  }
}`, instanceName, pwd)
}

func testAccDcsV1Instance_updated(instanceName, pwd string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" test {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "hcs_vpc_subnet" "test" {
  name        = "%[1]s"
  cidr        = "192.168.1.0/24"
  gateway_ip  = "192.168.1.1"
  vpc_id      = hcs_vpc.test.id
  description = "created by acc test"
}

resource "hcs_dcs_instance" "instance_1" {
  name               = "%[1]s"
  engine_version     = "5.0"
  password           = "Nsno1278das321"
  engine             = "Redis"
  port               = 6388
  capacity           = 4
  vpc_id             = resource.hcs_vpc.test.id
  subnet_id          = resource.hcs_vpc_subnet.test.id
  availability_zones = ["az0.dc0"]
  flavor             = "redis.ha.xu1.large.r2.4"
  maintain_begin     = "02:00:00"
  maintain_end       = "06:00:00"

  backup_policy {
    backup_type = "auto"
    begin_at    = "01:00-02:00"
    period_type = "weekly"
    backup_at   = [1, 2, 4]
    save_days   = 2
  }

  rename_commands = {
    command  = "command001"
    keys     = "keys001"
    flushall = "flushall001"
    flushdb  = "flushdb001"
    hgetall  = "hgetall001"
  }

  parameters {
    id    = "10"
    name  = "latency-monitor-threshold"
    value = "120"
  }
}`, instanceName, pwd)
}

func testAccDcsV1Instance_whitelists(instanceName, pwd string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" test {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "hcs_vpc_subnet" "test" {
  name        = "%[1]s"
  cidr        = "192.168.1.0/24"
  gateway_ip  = "192.168.1.1"
  vpc_id      = hcs_vpc.test.id
  description = "created by acc test"
}

resource "hcs_dcs_instance" "instance_1" {
  name               = "%[1]s"
  engine_version     = "6.0"
  password           = "Nsno1278das321"
  engine             = "Redis"
  port               = 6388
  capacity           = 0.125
  vpc_id             = resource.hcs_vpc.test.id
  subnet_id          = resource.hcs_vpc_subnet.test.id
  availability_zones = ["az0.dc0"]
  flavor             = "redis.ha.xu1.tiny.r2.128"

  whitelists {
    group_name = "testgroup1"
    ip_address = ["192.168.10.100", "192.168.0.0/24"]
  }
}`, instanceName, pwd)
}

func testAccDcsV1Instance_whitelists_update(instanceName, pwd string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" test {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "hcs_vpc_subnet" "test" {
  name        = "%[1]s"
  cidr        = "192.168.1.0/24"
  gateway_ip  = "192.168.1.1"
  vpc_id      = hcs_vpc.test.id
  description = "created by acc test"
}

resource "hcs_dcs_instance" "instance_1" {
  name               = "%[1]s"
  engine_version     = "6.0"
  password           = "Nsno1278das321"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = resource.hcs_vpc.test.id
  subnet_id          = resource.hcs_vpc_subnet.test.id
  availability_zones = ["az0.dc0"]
  flavor             = "redis.ha.xu1.tiny.r2.128"

  whitelists {
    group_name = "testgroup2"
    ip_address = ["172.16.10.100", "172.16.0.0/24"]
  }
}`, instanceName, pwd)
}

func testAccDcsV1Instance_single(instanceName, pwd string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" test {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "hcs_vpc_subnet" "test" {
  name        = "%[1]s"
  cidr        = "192.168.1.0/24"
  gateway_ip  = "192.168.1.1"
  vpc_id      = hcs_vpc.test.id
  description = "created by acc test"
}

resource "hcs_dcs_instance" "instance_1" {
  name               = "%[1]s"
  engine_version     = "5.0"
  password           = "Nsno1278das321"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = resource.hcs_vpc.test.id
  subnet_id          = resource.hcs_vpc_subnet.test.id
  availability_zones = ["az0.dc0"]
  flavor             = "redis.single.xu1.tiny.128"
}`, instanceName, pwd)
}
