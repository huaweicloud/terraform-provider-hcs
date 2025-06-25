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

func TestAccOpenGaussInstance_dorado(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "hcs_gaussdb_opengauss_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = acctest.RandString(10)
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
			acceptance.TestAccPreCheckOpengaussDoradoPool(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstance_dorado(rName, password, 3),
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
					resource.TestCheckResourceAttr(resourceName, "dorado_storage_pool_id", acceptance.DORADO_STORAGE_POOL_ID),
				),
			},
			{
				Config: testAccOpenGaussInstance_dorado_update(rName, password, 3),
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
					resource.TestCheckResourceAttr(resourceName, "dorado_storage_pool_id", acceptance.DORADO_STORAGE_POOL_ID),
				),
			},
		},
	})
}

func TestAccOpenGaussInstance_float_ip(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "hcs_gaussdb_opengauss_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = acctest.RandString(10)
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
			acceptance.TestAccPreCheckOpengaussDoradoPool(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstance_float_ip(rName, password, 3),
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
					resource.TestCheckResourceAttr(resourceName, "enable_single_float_ip", "true"),
				),
			},
			{
				Config: testAccOpenGaussInstance_float_ip_update(rName, password, 3),
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
					resource.TestCheckResourceAttr(resourceName, "enable_single_float_ip", "false"),
				),
			},
		},
	})
}

func TestAccOpenGaussInstance_ha(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "hcs_gaussdb_opengauss_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = acctest.RandString(10)
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
			acceptance.TestAccPreCheckOpengaussDoradoPool(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstance_ha(rName, password, 3),
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
					resource.TestCheckResourceAttr(resourceName, "ha.0.consistency_protocol", "quorum"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "replica_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
				),
			},
			{
				Config: testAccOpenGaussInstance_ha_update1(rName, password, 3),
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
					resource.TestCheckResourceAttr(resourceName, "ha.0.consistency_protocol", "syncStorage"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "replica_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
				),
			},
			{
				Config: testAccOpenGaussInstance_ha_update2(rName, password, 3),
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
					resource.TestCheckResourceAttr(resourceName, "ha.0.consistency_protocol", "paxos"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "replica_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
				),
			},
		},
	})
}

func TestAccOpenGaussInstance_volume(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "hcs_gaussdb_opengauss_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = acctest.RandString(10)
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
			acceptance.TestAccPreCheckOpengaussDoradoPool(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstance_volume(rName, password, 3),
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
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "LOCALSSD"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "480"),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "replica_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "enable_single_float_ip", "true"),
				),
			},
			{
				Config: testAccOpenGaussInstance_volume_update(rName, password, 3),
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
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "DORADO"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "480"),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "replica_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "enable_single_float_ip", "false"),
				),
			},
		},
	})
}

func TestAccOpenGaussInstance_kms(t *testing.T) {
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
			acceptance.TestAccPreCheckOpengaussKmsProjectName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstance_kms(rName, password, 3),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "hcs_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id", "hcs_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"hcs_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "kms_tde_key_id", "hcs_kms_key.test1.id"),
					resource.TestCheckResourceAttr(resourceName, "kms_project_name", acceptance.HCS_KMS_KEY_ID),
					resource.TestCheckResourceAttr(resourceName, "ha.0.mode", "enterprise"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.replication_mode", "sync"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.consistency", "strong"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "replica_num", "3"),
				),
			},
			{
				Config: testAccOpenGaussInstance_kms_update(rName, newPassword, 3),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "password", newPassword),
					resource.TestCheckResourceAttr(resourceName, "kms_tde_key_id", "hcs_kms_key.test2.id"),
					resource.TestCheckResourceAttr(resourceName, "kms_project_name", acceptance.HCS_KMS_KEY_ID),
					resource.TestCheckResourceAttr(resourceName, "kms_tde_status", "on"),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
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

func testAccOpenGaussInstance_dorado(rName, password string, replicaNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_gaussdb_opengauss_instance" "test" {
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  dorado_storage_pool_id = "%[5]s"
  enable_single_float_ip = true

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
`, testAccOpenGaussInstance_base(rName), rName, password, replicaNum, acceptance.DORADO_STORAGE_POOL_ID)
}

func testAccOpenGaussInstance_dorado_update(rName, password string, replicaNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_gaussdb_opengauss_instance" "test" {
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  dorado_storage_pool_id = "%[5]s"
  enable_single_float_ip = false

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
`, testAccOpenGaussInstance_base(rName), rName, password, replicaNum, acceptance.DORADO_STORAGE_POOL_ID)
}

func testAccOpenGaussInstance_ha(rName, password string, replicaNum int) string {
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
    mode                 = "enterprise"
    replication_mode     = "sync"
    consistency          = "strong"
	consistency_protocol = "quorum"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, testAccOpenGaussInstance_base(rName), rName, password, replicaNum)
}

func testAccOpenGaussInstance_ha_update1(rName, password string, replicaNum int) string {
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
    mode                 = "enterprise"
    replication_mode     = "sync"
    consistency          = "strong"
	consistency_protocol = "syncStorage"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, testAccOpenGaussInstance_base(rName), rName, password, replicaNum)
}

func testAccOpenGaussInstance_ha_update2(rName, password string, replicaNum int) string {
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

  solution = "logger"

  ha {
    mode                 = "enterprise"
    replication_mode     = "sync"
    consistency          = "strong"
	consistency_protocol = "quorum"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, testAccOpenGaussInstance_base(rName), rName, password, replicaNum)
}

func testAccOpenGaussInstance_float_ip(rName, password string, replicaNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_gaussdb_opengauss_instance" "test" {
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  enable_single_float_ip = true

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

func testAccOpenGaussInstance_float_ip_update(rName, password string, replicaNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_gaussdb_opengauss_instance" "test" {
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  enable_single_float_ip = false

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

func testAccOpenGaussInstance_volume(rName, password string, replicaNum int) string {
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
    type = "LOCALSSD"
    size = 480
  }
}
`, testAccOpenGaussInstance_base(rName), rName, password, replicaNum)
}

func testAccOpenGaussInstance_volume_update(rName, password string, replicaNum int) string {
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
    type = "DORADO"
    size = 480
  }
}
`, testAccOpenGaussInstance_base(rName), rName, password, replicaNum)
}

func testAccOpenGaussInstance_kms(rName, password string, replicaNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_kms_key" "test1" {
  key_alias    = "tf-test1"
  pending_days = "7"
}

resource "hcs_kms_key" "test2" {
  key_alias    = "tf-test2"
  pending_days = "7"
}

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

  kms_tde_key_id   = hcs_kms_key.test1.id
  kms_project_name = "%[5]s"

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
`, testAccOpenGaussInstance_base(rName), rName, password, replicaNum, acceptance.OPENGAUSS_KMS_PROJECT_NAME)
}

func testAccOpenGaussInstance_kms_update(rName, password string, replicaNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_kms_key" "test1" {
  key_alias    = "tf-test1"
  pending_days = "7"
}

resource "hcs_kms_key" "test2" {
  key_alias    = "tf-test2"
  pending_days = "7"
}

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

  kms_tde_key_id   = hcs_kms_key.test2.id
  kms_project_name = "%[5]s"
  kms_tde_status   = "on"

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
`, testAccOpenGaussInstance_base(rName), rName, password, replicaNum, acceptance.OPENGAUSS_KMS_PROJECT_NAME)
}
