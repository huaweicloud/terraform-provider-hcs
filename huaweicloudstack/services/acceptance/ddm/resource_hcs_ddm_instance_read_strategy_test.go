package ddm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func TestAccDdmInstanceReadStrategy_basic(t *testing.T) {
	name := acceptance.RandomAccResourceNameWithDash()
	schemaName := acceptance.RandomAccResourceName()
	rName := "hcs_ddm_instance_read_strategy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDdmInstanceReadStrategy_basic(name, schemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"hcs_ddm_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "read_weights.#", "1"),
				),
			},
			{
				Config: testDdmInstanceReadStrategy_update(name, schemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"hcs_ddm_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "read_weights.#", "2"),
				),
			},
		},
	})
}

func testDdmInstanceReadStrategy_basic(name, schemaName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ddm_instance_read_strategy" "test" {
  depends_on  = [hcs_ddm_schema.test]
  instance_id = hcs_ddm_instance.test.id

  read_weights {
    db_id  = hcs_rds_instance.test.id
    weight = 100
  }
}
`, testDdmInstanceReadStrategyBase(name, schemaName))
}

func testDdmInstanceReadStrategy_update(name, schemaName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ddm_instance_read_strategy" "test" {
  depends_on  = [hcs_ddm_schema.test]
  instance_id = hcs_ddm_instance.test.id

  read_weights {
    db_id  = hcs_rds_instance.test.id
    weight = 60
  }

  read_weights {
    db_id  = hcs_rds_read_replica_instance.test.id
    weight = 40
  }
}
`, testDdmInstanceReadStrategyBase(name, schemaName))
}

func testDdmInstanceReadStrategyBase(name, schemaName string) string {
	return fmt.Sprintf(`
%[1]s

data "hcs_vpc" "test" {
  name = "vpc-default"
}

data "hcs_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "hcs_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 3306
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = hcs_networking_secgroup.test.id
}

resource "hcs_networking_secgroup_rule" "egress" {
  direction         = "egress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = hcs_networking_secgroup.test.id
}

data "hcs_availability_zones" "test" {}

data "hcs_ddm_engines" test {
  version = "3.0.8.5"
}

data "hcs_ddm_flavors" test {
  engine_id = data.hcs_ddm_engines.test.engines[0].id
  cpu_arch  = "X86"
}

data "hcs_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
}
  
resource "hcs_rds_instance" "test" {
  depends_on = [
    hcs_networking_secgroup_rule.ingress,
    hcs_networking_secgroup_rule.egress,
  ]

  name              = "%[2]s"
  flavor            = data.hcs_rds_flavors.test.flavors[0].name
  vpc_id            = data.hcs_vpc.test.id
  subnet_id         = data.hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  availability_zone = [
    data.hcs_availability_zones.test.names[0]
  ]

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
  
data "hcs_rds_flavors" "replica" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "replica"
  group_type    = "dedicated"
  memory        = 4
  vcpus         = 2
}
	
resource "hcs_rds_read_replica_instance" "test" {
  name                = "%[2]s-read-replica"
  flavor              = data.hcs_rds_flavors.replica.flavors[0].name
  primary_instance_id = hcs_rds_instance.test.id
  availability_zone   = data.hcs_availability_zones.test.names[0]

  volume {
    type              = "CLOUDSSD"
    size              = 50
    limit_size        = 400
    trigger_threshold = 10
  }
}

resource "hcs_ddm_instance" "test" {
  depends_on = [
    hcs_rds_read_replica_instance.test,
    hcs_networking_secgroup_rule.ingress,
    hcs_networking_secgroup_rule.egress
  ]

  name              = "%[2]s"
  flavor_id         = data.hcs_ddm_flavors.test.flavors[0].id
  node_num          = 2
  engine_id         = data.hcs_ddm_engines.test.engines[0].id
  vpc_id            = data.hcs_vpc.test.id
  subnet_id         = data.hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  availability_zones = [
    data.hcs_availability_zones.test.names[0]
  ]
}

resource "hcs_ddm_schema" "test" {
  instance_id  = hcs_ddm_instance.test.id
  name         = "%[3]s"
  shard_mode   = "single"
  shard_number = "1"
  
  data_nodes {
    id             = hcs_rds_instance.test.id
    admin_user     = "root"
    admin_password = "Huangwei!120521"
  }
  
  delete_rds_data = "true"
  
  lifecycle {
    ignore_changes = [
      data_nodes,
    ]
  }
}
`, common.TestSecGroup(name), name, schemaName)
}
