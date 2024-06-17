package dms

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/dms/v2/kafka/instances"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func getKafkaInstanceFunc(c *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DmsV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloudStack DMS client(V2): %s", err)
	}
	return instances.Get(client, state.Primary.ID).Extract()
}

func TestAccKafkaInstance_basic(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := rName + "update"
	resourceName := "hcs_dms_kafka_instance.test"
	password := acceptance.RandomPassword()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getKafkaInstanceFunc,
	)

	// DMS instances use the tenant-level shared lock, the instances cannot be created or modified in parallel.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstance_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "security_protocol", "SASL_PLAINTEXT"),
					resource.TestCheckResourceAttr(resourceName, "enabled_mechanisms.0", "SCRAM-SHA-512"),
					resource.TestMatchResourceAttr(resourceName, "cross_vpc_accesses.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
			{
				Config: testAccKafkaInstance_update(rName, updateName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "kafka test update"),
					resource.TestCheckResourceAttr(resourceName, "enable_auto_topic", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"manager_password",
					"used_storage_space",
					"cross_vpc_accesses",
					"security_protocol",
					"enabled_mechanisms",
				},
			},
		},
	})
}

func TestAccKafkaInstance_compatible(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_dms_kafka_instance.test"
	password := acceptance.RandomPassword()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getKafkaInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstance_compatible(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
		},
	})
}

func testAccKafkaInstance_basic(rName, password string) string {
	return fmt.Sprintf(`
%s

data "hcs_dms_kafka_flavors" "test" {
  type               = "cluster"
  flavor_id          = 0
  availability_zones = ["az0.dc0"]
  storage_spec_code  = "dms.physical.storage.ultra"
}

resource "hcs_dms_kafka_instance" "test" {
  name               = "%s"
  description        = "kafka test"
  access_user        = "user"
  password           = "%[3]s"
  vpc_id             = hcs_vpc.test.id
  network_id         = hcs_vpc_subnet.test.id
  security_group_id  = hcs_networking_secgroup.test.id
  availability_zones = [
    "az0.dc0"
  ]
  product_id        = data.hcs_dms_kafka_flavors.test.id
  engine_version    = "2.7"
  storage_spec_code = "dms.physical.storage.ultra"

  manager_user       = "kafka-user"
  manager_password   = "%[3]s"
  security_protocol  = "SASL_PLAINTEXT"
  enabled_mechanisms = ["SCRAM-SHA-512"]
}
`, common.TestBaseNetwork(rName), rName, password)
}

func testAccKafkaInstance_update(rName, updateName, password string) string {
	return fmt.Sprintf(`
%s

data "hcs_dms_kafka_flavors" "test" {
  type               = "cluster"
  flavor_id          = 0
  availability_zones = ["az0.dc0"]
  storage_spec_code  = "dms.physical.storage.ultra"
}

resource "hcs_dms_kafka_instance" "test" {
  name               = "%s"
  description        = "kafka test update"
  access_user        = "user"
  password           = "%[3]s"
  vpc_id             = hcs_vpc.test.id
  network_id         = hcs_vpc_subnet.test.id
  security_group_id  = hcs_networking_secgroup.test.id
  availability_zones = [
    "az0.dc0"
  ]
  product_id        = data.hcs_dms_kafka_flavors.test.id
  engine_version    = "2.7"
  storage_spec_code = "dms.physical.storage.ultra"

  manager_user       = "kafka-user"
  manager_password   = "%[3]s"
  security_protocol  = "SASL_PLAINTEXT"
  enabled_mechanisms = ["SCRAM-SHA-512"]
  enable_auto_topic  = true
}
`, common.TestBaseNetwork(rName), updateName, password)
}

func testAccKafkaInstance_compatible(rName, password string) string {
	return fmt.Sprintf(`
%s

data "hcs_dms_kafka_flavors" "test" {
  type               = "cluster"
  flavor_id          = 0
  availability_zones = ["az0.dc0"]
  storage_spec_code  = "dms.physical.storage.ultra"
}

resource "hcs_dms_kafka_instance" "test" {
  name        = "%s"
  description = "kafka test"

  availability_zones = [
    "az0.dc0"
  ]

  vpc_id            = hcs_vpc.test.id
  network_id        = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  product_id        = data.hcs_dms_kafka_flavors.test.id
  engine_version    = "2.7"
  storage_spec_code = "dms.physical.storage.ultra"
  storage_space     = "600"
  # use deprecated argument "bandwidth"

  access_user      = "user"
  password         = "%[3]s"
  manager_user     = "kafka-user"
  manager_password = "%[3]s"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, common.TestBaseNetwork(rName), rName, password)
}
