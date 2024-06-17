package lts

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getHostAccessConfigResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	ltsClient, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}

	listHostAccessConfigHttpUrl := "v3/{project_id}/lts/access-config-list"
	listHostAccessConfigPath := ltsClient.Endpoint + listHostAccessConfigHttpUrl
	listHostAccessConfigPath = strings.ReplaceAll(listHostAccessConfigPath, "{project_id}", ltsClient.ProjectID)

	listHostAccessConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	name := state.Primary.Attributes["name"]
	listHostAccessConfigOpt.JSONBody = map[string]interface{}{
		"access_config_name_list": []string{name},
	}

	listHostAccessConfigResp, err := ltsClient.Request("POST", listHostAccessConfigPath, &listHostAccessConfigOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving host access config: %s", err)
	}

	listHostAccessConfigRespBody, err := utils.FlattenResponse(listHostAccessConfigResp)
	if err != nil {
		return nil, fmt.Errorf("error flatten host access config response: %s", err)
	}

	jsonPath := fmt.Sprintf("result[?access_config_name=='%s']|[0]", name)
	listHostAccessConfigRespBody = utils.PathSearch(jsonPath, listHostAccessConfigRespBody, nil)
	if listHostAccessConfigRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return listHostAccessConfigRespBody, nil
}

func TestAccHostAccessConfig_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_lts_host_access.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getHostAccessConfigResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testHostAccessConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "access_config.0.paths.0", "/var/log/*"),
					resource.TestCheckResourceAttr(rName, "host_group_ids.#", "0"),
					resource.TestCheckResourceAttr(rName, "access_type", "AGENT"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "hcs_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id", "hcs_lts_stream.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
					resource.TestCheckResourceAttrSet(rName, "log_stream_name"),
				),
			},
			{
				Config: testHostAccessConfig_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "access_config.0.paths.0", "/var/log/*/*.log"),
					resource.TestCheckResourceAttr(rName, "access_config.0.black_paths.0", "/var/log/*/a.log"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value-updated"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(rName, "host_group_ids.#", "1"),
					resource.TestCheckResourceAttrPair(rName, "host_group_ids.0", "hcs_lts_host_group.test", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccHostAccessConfigImportStateFunc(rName),
			},
		},
	})
}

func TestAccHostAccessConfig_windows(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_lts_host_access.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getHostAccessConfigResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testHostAccessConfig_windows_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "access_config.0.paths.0", "D:\\data\\log\\*"),
					resource.TestCheckResourceAttr(rName, "access_config.0.black_paths.0", "D:\\data\\log\\a.log"),
					resource.TestCheckResourceAttr(rName, "access_config.0.windows_log_info.0.time_offset", "7"),
					resource.TestCheckResourceAttr(rName, "access_config.0.windows_log_info.0.time_offset_unit", "day"),
					resource.TestCheckResourceAttr(rName, "host_group_ids.#", "1"),
					resource.TestCheckResourceAttr(rName, "access_type", "AGENT"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "hcs_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id", "hcs_lts_stream.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
					resource.TestCheckResourceAttrSet(rName, "log_stream_name"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccHostAccessConfigImportStateFunc(rName),
			},
		},
	})
}

func testAccHostAccessConfigImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			return "", fmt.Errorf("can not find the 'name' parameter in %s", rName)
		}
		return name, nil
	}
}

func testHostAccessConfig_base(name string) string {
	return fmt.Sprintf(`
resource "hcs_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 7
}

resource "hcs_lts_stream" "test" {
  group_id    = hcs_lts_group.test.id
  stream_name = "%[1]s"
}
`, name)
}

func testHostAccessConfig_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# A host group without hosts
resource "hcs_lts_host_group" "test" {
  name = "%[2]s"
  type = "linux"
}

resource "hcs_lts_host_access" "test" {
  name          = "%[2]s"
  log_group_id  = hcs_lts_group.test.id
  log_stream_id = hcs_lts_stream.test.id

  access_config {
    paths = ["/var/log/*"]

    single_log_format {
      mode = "system"
    }
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testHostAccessConfig_base(name), name)
}

func testHostAccessConfig_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

# A host group without hosts
resource "hcs_lts_host_group" "test" {
  name = "%[2]s"
  type = "linux"
}

resource "hcs_lts_host_access" "test" {
  name           = "%[2]s"
  log_group_id   = hcs_lts_group.test.id
  log_stream_id  = hcs_lts_stream.test.id
  host_group_ids = [hcs_lts_host_group.test.id]

  access_config {
    paths       = ["/var/log/*/*.log"]
    black_paths = ["/var/log/*/a.log"]

    multi_log_format {
      mode  = "time"
      value = "YYYY-MM-DD hh:mm:ss"
    }
  }

  tags = {
    key   = "value-updated"
    owner = "terraform"
  }
}
`, testHostAccessConfig_base(name), name)
}

func testHostAccessConfig_windows_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# A host group without hosts
resource "hcs_lts_host_group" "test" {
  name = "%[2]s"
  type = "windows"
}

resource "hcs_lts_host_access" "test" {
  name           = "%[2]s"
  log_group_id   = hcs_lts_group.test.id
  log_stream_id  = hcs_lts_stream.test.id
  host_group_ids = [hcs_lts_host_group.test.id]

  access_config {
    paths       = ["D:\\data\\log\\*"]
    black_paths = ["D:\\data\\log\\a.log"]

    windows_log_info {
      categorys        = ["System", "Application"]
      event_level      = ["warning", "error"]
      time_offset_unit = "day"
      time_offset      = 7
    }

    single_log_format {
      mode = "system"
    }
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testHostAccessConfig_base(name), name)
}
