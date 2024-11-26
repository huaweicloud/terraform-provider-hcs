package rds

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/rds/v3/instances"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func TestAccRdsInstance_basic(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "hcs_rds_instance"
	resourceName := "hcs_rds_instance.test"
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test_description"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.large.2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8634"),
				),
			},
			{
				Config: testAccRdsInstance_update(name, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", name)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "2"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.large.2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_updated"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip", "192.168.0.230"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8636"),
					resource.TestCheckResourceAttrSet(resourceName, "db.0.password"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"db",
					"status",
					"availability_zone",
				},
			},
		},
	})
}

func TestAccRdsInstance_ha(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "hcs_rds_instance"
	resourceName := "hcs_rds_instance.test"
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_ha(name, password, "async"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.c2.large.ha"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "ha_replication_mode", "async"),
				),
			},
			{
				Config: testAccRdsInstance_ha(name, password, "sync"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.c2.large.ha"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "ha_replication_mode", "sync"),
				),
			},
		},
	})
}

func TestAccRdsInstance_restore_pg(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "hcs_rds_instance"
	resourceName := "hcs_rds_instance.test_backup"
	pwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_restore_pg(name, pwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "CLOUDSSD"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8732"),
				),
			},
		},
	})
}

func testAccCheckRdsInstanceDestroy(rsType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acceptance.TestAccProvider.Meta().(*config.HcsConfig)
		client, err := config.RdsV3Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating rds client: %s", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != rsType {
				continue
			}

			id := rs.Primary.ID
			instance, err := GetRdsInstanceByID(client, id)
			if err != nil {
				return err
			}
			if instance.Id != "" {
				return fmt.Errorf("%s (%s) still exists", rsType, id)
			}
		}
		return nil
	}
}

func GetRdsInstanceByID(client *golangsdk.ServiceClient, instanceID string) (*instances.RdsInstanceResponse, error) {
	listOpts := instances.ListOpts{
		Id: instanceID,
	}
	pages, err := instances.List(client, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("An error occurred while querying rds instance %s: %s", instanceID, err)
	}

	resp, err := instances.ExtractRdsInstances(pages)
	if err != nil {
		return nil, err
	}

	instanceList := resp.Instances
	if len(instanceList) == 0 {
		// return an empty rds instance
		log.Printf("[WARN] can not find the specified rds instance %s", instanceID)
		instance := new(instances.RdsInstanceResponse)
		return instance, nil
	}

	if len(instanceList) > 1 {
		return nil, fmt.Errorf("retrieving more than one rds instance by %s", instanceID)
	}
	if instanceList[0].Id != instanceID {
		return nil, fmt.Errorf("the id of rds instance was expected %s, but got %s",
			instanceID, instanceList[0].Id)
	}

	return &instanceList[0], nil
}

func testAccCheckRdsInstanceExists(name string, instance *instances.RdsInstanceResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.HcsConfig)
		client, err := config.RdsV3Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating rds client: %s", err)
		}

		found, err := GetRdsInstanceByID(client, id)
		if err != nil {
			return fmt.Errorf("error checking %s exist, err=%s", name, err)
		}
		if found.Id == "" {
			return fmt.Errorf("resource %s does not exist", name)
		}

		instance = found
		return nil
	}
}

func testAccRdsInstance_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_rds_instance" "test" {
  name              = "%[2]s"
  description       = "test_description"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = ["az1.dc1"]
  security_group_id = hcs_networking_secgroup.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  vpc_id            = hcs_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    type    = "PostgreSQL"
    version = "12"
    port    = 8634
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, common.TestBaseNetwork(name), name)
}

// name, volume.size, backup_strategy, flavor, tags and password will be updated
func testAccRdsInstance_update(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_rds_instance" "test" {
  name              = "%[2]s-update"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = ["az1.dc1"]
  security_group_id = hcs_networking_secgroup.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  vpc_id            = hcs_vpc.test.id
  time_zone         = "UTC+08:00"
  fixed_ip          = "192.168.0.230"

  db {
    password = "%[3]s"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8636
  }

  volume {
    type = "CLOUDSSD"
    size = 100
  }

  backup_strategy {
    start_time = "09:00-10:00"
    keep_days  = 2
  }

  tags = {
    key1 = "value"
    foo  = "bar_updated"
  }
}
`, common.TestBaseNetwork(name), name, password)
}

func testAccRdsInstance_ha(name, password, replicationMode string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_rds_instance" "test" {
  name                = "%[2]s"
  flavor              = "rds.pg.c2.large.ha"
  security_group_id   = hcs_networking_secgroup.test.id
  subnet_id           = hcs_vpc_subnet.test.id
  vpc_id              = hcs_vpc.test.id
  time_zone           = "UTC+08:00"
  ha_replication_mode = "%[4]s"
  availability_zone   = [
    "az1.dc1",
    "az1.dc1"
  ]

  db {
    password = "%[3]s"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8634
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, common.TestBaseNetwork(name), name, password, replicationMode)
}

func testAccRdsInstance_restore_pg(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_rds_instance" "test_backup" {
  name              = "%[2]s"
  flavor            = "rds.pg.c2.large.ha"
  security_group_id = hcs_networking_secgroup.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  vpc_id            = hcs_vpc.test.id
  availability_zone = ["az1.dc1","az1.dc1"]

  restore {
    instance_id = "%[3]s"
    backup_id   = "%[4]s"
  }

  db {
    password = "%[5]s"
    type     = "PostgreSQL"
    version  = "14"
    port     = 8732
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}
`, common.TestBaseNetwork(name), name, acceptance.HCS_RDS_INSTANCE_ID, acceptance.HCS_RDS_BACKUP_ID, password)
}
