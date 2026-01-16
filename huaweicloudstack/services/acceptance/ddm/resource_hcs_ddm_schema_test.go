package ddm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func getDdmSchemaResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getSchema: Query DDM schema
	var (
		getSchemaHttpUrl = "v1/{project_id}/instances/{instance_id}/databases/{ddm_dbname}"
		getSchemaProduct = "ddm"
	)
	getSchemaClient, err := cfg.NewServiceClient(getSchemaProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DDM client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<db_name>")
	}
	instanceID := parts[0]
	schemaName := parts[1]
	getSchemaPath := getSchemaClient.Endpoint + getSchemaHttpUrl
	getSchemaPath = strings.ReplaceAll(getSchemaPath, "{project_id}", getSchemaClient.ProjectID)
	getSchemaPath = strings.ReplaceAll(getSchemaPath, "{instance_id}", fmt.Sprintf("%v", instanceID))
	getSchemaPath = strings.ReplaceAll(getSchemaPath, "{ddm_dbname}", fmt.Sprintf("%v", schemaName))

	getSchemaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSchemaResp, err := getSchemaClient.Request("GET", getSchemaPath, &getSchemaOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DdmSchema: %s", err)
	}
	getSchemaRespBody, err := utils.FlattenResponse(getSchemaResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("database", getSchemaRespBody, nil), nil
}

func TestAccDdmSchema_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	instanceName := strings.ReplaceAll(name, "_", "-")
	rName := "hcs_ddm_schema.test"
	dbPwd := "test_1234"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdmSchemaResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdmSchema_basic(instanceName, name, dbPwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "shard_mode", "single"),
					resource.TestCheckResourceAttr(rName, "shard_number", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"instance_id", "data_nodes.0.admin_user",
					"data_nodes.0.admin_password", "delete_rds_data"},
			},
		},
	})
}

func testDdmSchema_base(name, dbPwd string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_rds_instance" "instance" {
  name              = "%[2]s"
  flavor            = "rds.mysql.xlarge.arm4.single"
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id
  availability_zone = ["az0.dc0"]
  description       = "test_description"  

  db {
    type     = "MySQL"
    version  = "5.7"
    password = "%[3]s"
  }

  volume {
    type = "ULTRAHIGH"
    size = 100
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, common.TestVpc(name), name, dbPwd)
}

func testDdmSchema_basic(instanceName, name, dbPwd string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_ddm_schema" "test" {
  instance_id  = hcs_ddm_instance.test.id
  name         = "%[2]s"
  shard_mode   = "single"
  shard_number = "1"

  data_nodes {
    id             = hcs_rds_instance.test.id
    admin_user     = "root"
    admin_password = "%[3]s"
  }

  delete_rds_data = "true"

  lifecycle {
    ignore_changes = [
      data_nodes,
    ]
  }
}
`, testDdmInstance_basic(name), testDdmSchema_base(instanceName, dbPwd), name, dbPwd)
}
