package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccPgDatabasePrivilege_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resourceName := "hcs_rds_pg_database_privilege.test"
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsPgDatabasePrivilege_basic(name, password),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test_description"),
				),
			},
			{
				Config: testAccRdsPgDatabasePrivilege_basic_update(name, password),
			},
		},
	})
}

func testAccRdsPgDatabasePrivilege_base(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_rds_pg_database" "test" {
  depends_on = [hcs_rds_pg_account.test]

  instance_id   = hcs_rds_instance.test.id
  name          = "%[2]s"
  owner         = hcs_rds_pg_account.test.name
  character_set = "UTF8"
}
`, testPgAccount_basic(name, password), name)
}

func testAccRdsPgDatabasePrivilege_basic(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_rds_pg_database_privilege" "test" {
  depends_on = [
    hcs_rds_pg_database.test,
    hcs_rds_pg_account.test,
  ]

  instance_id = hcs_rds_instance.test.id
  db_name     = hcs_rds_pg_database.test.name

  users {
    name        = hcs_rds_pg_account.member.name
    readonly    = true
    schema_name = "public"
  }
}
`, testAccRdsPgDatabasePrivilege_base(name, password))
}

func testAccRdsPgDatabasePrivilege_basic_update(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_rds_pg_database_privilege" "test" {
  depends_on = [
    hcs_rds_pg_database.test,
    hcs_rds_pg_account.test,
  ]

  instance_id = hcs_rds_instance.test.id
  db_name     = hcs_rds_pg_database.test.name

  users {
    name        = hcs_rds_pg_account.member.name
    readonly    = false
    schema_name = "public"
  }
  users {
    name        = hcs_rds_pg_account.test.name
    readonly    = true
    schema_name = "public"
  }

  enable_force_new = "true"
}
`, testAccRdsPgDatabasePrivilege_base(name, password))
}
