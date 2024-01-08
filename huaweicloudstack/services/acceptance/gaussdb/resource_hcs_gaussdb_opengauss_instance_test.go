package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/opengauss/v3/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	hwacceptance "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func getOpenGaussInstanceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.OpenGaussV3Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating HuaweiCloudStack GaussDB client: %s", err)
	}
	return instances.GetInstanceByID(client, state.Primary.ID)
}

func TestAccOpenGaussInstance_basic(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "hcs_gaussdb_opengauss_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = acctest.RandString(10)
		newPassword  = acctest.RandString(10)
	)

	rc := hwacceptance.InitResourceCheck(
		resourceName,
		&instance,
		getOpenGaussInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstance_basic(rName, password, 3),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "hcs_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id", "hcs_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"hcs_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "ha.0.mode", "enterprise"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.replication_mode", "sync"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.consistency", "strong"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "replica_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
				),
			},
			{
				Config: testAccOpenGaussInstance_update(rName, newPassword, 3),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "password", newPassword),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "replica_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "80"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "08:00-09:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "8"),
				),
			},
		},
	})
}

func TestAccOpenGaussInstance_replicaNumTwo(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "hcs_gaussdb_opengauss_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = acctest.RandString(10)
		newPassword  = acctest.RandString(10)
	)

	rc := hwacceptance.InitResourceCheck(
		resourceName,
		&instance,
		getOpenGaussInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstance_basic(rName, password, 2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "hcs_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id", "hcs_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"hcs_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "ha.0.mode", "enterprise"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.replication_mode", "sync"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.consistency", "strong"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "replica_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
				),
			},
			{
				Config: testAccOpenGaussInstance_update(rName, newPassword, 2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "password", newPassword),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "replica_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "80"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "08:00-09:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "8"),
				),
			},
		},
	})
}

func TestAccOpenGaussInstance_haModeCentralized(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "hcs_gaussdb_opengauss_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = acctest.RandString(10)
		newPassword  = acctest.RandString(10)
	)

	rc := hwacceptance.InitResourceCheck(
		resourceName,
		&instance,
		getOpenGaussInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstance_haModeCentralized(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "hcs_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id", "hcs_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"hcs_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "gaussdb.opengauss.ee.m6.2xlarge.x868.ha"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "ha.0.mode", "centralization_standard"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.replication_mode", "sync"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.consistency", "strong"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "replica_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
				),
			},
			{
				Config: testAccOpenGaussInstance_haModeCentralizedUpdate(rName, newPassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "password", newPassword),
					resource.TestCheckResourceAttr(resourceName, "replica_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "80"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "08:00-09:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "8"),
				),
			},
		},
	})
}

func testAccOpenGaussInstance_base(rName string) string {
	return fmt.Sprintf(`
%s

data "hcs_availability_zones" "test" {}

// opengauss requires more sg ports open
resource "hcs_networking_secgroup_rule" "in_v4_tcp_opengauss" {
  security_group_id = hcs_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "hcs_networking_secgroup_rule" "in_v4_tcp_opengauss_egress" {
  security_group_id = hcs_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "egress"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}
`, common.TestBaseNetwork(rName))
}

func testAccOpenGaussInstance_basic(rName, password string, replicaNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_gaussdb_opengauss_instance" "test" {
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  flavor            = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  name              = "%[2]s"
  password          = "%[3]s"
  sharding_num      = 1
  coordinator_num   = 2
  replica_num       = %[4]d
  availability_zone = "${data.hcs_availability_zones.test.names[0]},${data.hcs_availability_zones.test.names[0]},${data.hcs_availability_zones.test.names[0]}"

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "strong"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, testAccOpenGaussInstance_base(rName), rName, password, replicaNum)
}

func testAccOpenGaussInstance_update(rName, password string, replicaNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_gaussdb_opengauss_instance" "test" {
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  flavor            = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  name              = "%[2]s-update"
  password          = "%[3]s"
  sharding_num      = 2
  coordinator_num   = 3
  replica_num       = %[4]d
  availability_zone = "${data.hcs_availability_zones.test.names[0]},${data.hcs_availability_zones.test.names[0]},${data.hcs_availability_zones.test.names[0]}"

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "strong"
  }

  volume {
    type = "ULTRAHIGH"
    size = 80
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 8
  }
}
`, testAccOpenGaussInstance_base(rName), rName, password, replicaNum)
}

func testAccOpenGaussInstance_haModeCentralized(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_gaussdb_opengauss_instance" "test" {
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  flavor            = "gaussdb.opengauss.ee.m6.2xlarge.x868.ha"
  name              = "%[2]s"
  password          = "%[3]s"
  replica_num       = 3
  availability_zone = "${data.hcs_availability_zones.test.names[0]},${data.hcs_availability_zones.test.names[0]},${data.hcs_availability_zones.test.names[0]}"

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "strong"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, testAccOpenGaussInstance_base(rName), rName, password)
}

func testAccOpenGaussInstance_haModeCentralizedUpdate(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_gaussdb_opengauss_instance" "test" {
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  flavor            = "gaussdb.opengauss.ee.m6.2xlarge.x868.ha"
  name              = "%[2]s-update"
  password          = "%[3]s"
  replica_num       = 3
  availability_zone = "${data.hcs_availability_zones.test.names[0]},${data.hcs_availability_zones.test.names[0]},${data.hcs_availability_zones.test.names[0]}"

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "strong"
  }

  volume {
    type = "ULTRAHIGH"
    size = 80
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 8
  }
}
`, testAccOpenGaussInstance_base(rName), rName, password)
}
