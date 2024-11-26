package rds

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/pagination"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func getPgAccountResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getPgAccount: query RDS PostgreSQL account
	var (
		getPgAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db_user/detail?page=1&limit=100"
		getPgAccountProduct = "rds"
	)
	getPgAccountClient, err := cfg.NewServiceClient(getPgAccountProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and user from resource id
	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format, must be <instance_id>/<name>")
	}
	instanceId := parts[0]
	accountName := parts[1]

	getPgAccountPath := getPgAccountClient.Endpoint + getPgAccountHttpUrl
	getPgAccountPath = strings.ReplaceAll(getPgAccountPath, "{project_id}", getPgAccountClient.ProjectID)
	getPgAccountPath = strings.ReplaceAll(getPgAccountPath, "{instance_id}", instanceId)

	getPgAccountResp, err := pagination.ListAllItems(
		getPgAccountClient,
		"page",
		getPgAccountPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL account: %s", err)
	}

	getPgAccountRespJson, err := json.Marshal(getPgAccountResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL account: %s", err)
	}

	var getPgAccountRespBody interface{}
	err = json.Unmarshal(getPgAccountRespJson, &getPgAccountRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL account: %s", err)
	}

	account := utils.PathSearch(fmt.Sprintf("users[?name=='%s']|[0]", accountName), getPgAccountRespBody, nil)

	if account != nil {
		return account, nil
	}

	return nil, fmt.Errorf("error retrieving RDS PostgreSQL account by instanceID %s and account %s", instanceId,
		accountName)
}

func TestAccPgAccount_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_rds_pg_account.test"
	password := acceptance.RandomPassword()

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPgAccountResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPgAccount_basic(name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"hcs_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test_description"),
					resource.TestCheckResourceAttrPair(rName, "memberof.0",
						"hcs_rds_pg_account.member", "name"),
					resource.TestCheckResourceAttr(rName, "attributes.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.rol_super"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.rol_inherit"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.rol_create_role"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.rol_create_db"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.rol_can_login"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.rol_conn_limit"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.rol_replication"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.rol_bypass_rls"),
				),
			},
			{
				Config: testPgAccount_update(name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"hcs_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "test_description_update"),
					resource.TestCheckResourceAttrPair(rName, "memberof.0",
						"hcs_rds_pg_account.member_update", "name"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testPgAccount_base(name string) string {
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
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}
`, common.TestBaseNetwork(name), name)
}

func testPgAccount_basic(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_rds_pg_account" "member" {
  instance_id = hcs_rds_instance.test.id
  name        = "%[2]s_member"
  password    = "%[3]s"
}

resource "hcs_rds_pg_account" "test" {
  depends_on  = [hcs_rds_pg_account.member]
  instance_id = hcs_rds_instance.test.id
  name        = "%[2]s"
  password    = "%[3]s"
  description = "test_description"
  memberof    = [hcs_rds_pg_account.member.name]
}
`, testPgAccount_base(name), name, password)
}

func testPgAccount_update(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_rds_pg_account" "member" {
  instance_id = hcs_rds_instance.test.id
  name        = "%[2]s_member"
  password    = "%[3]s"
}

resource "hcs_rds_pg_account" "member_update" {
  instance_id = hcs_rds_instance.test.id
  name        = "%[2]s_member_update"
  password    = "%[3]s"
}

resource "hcs_rds_pg_account" "test" {
  depends_on  = [hcs_rds_pg_account.member_update]
  instance_id = hcs_rds_instance.test.id
  name        = "%[2]s"
  password    = "%[3]s"
  description = "test_description_update"
  memberof    = [hcs_rds_pg_account.member_update.name]
}
`, testPgAccount_base(name), name, password)
}
