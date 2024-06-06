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

func getLtsTransferResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getTransfer: Query the log transfer task.
	var (
		getTransferHttpUrl = "v2/{project_id}/transfers"
		getTransferProduct = "lts"
	)
	getTransferClient, err := cfg.NewServiceClient(getTransferProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}

	getTransferPath := getTransferClient.Endpoint + getTransferHttpUrl
	getTransferPath = strings.ReplaceAll(getTransferPath, "{project_id}", getTransferClient.ProjectID)

	getTransferOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getTransferResp, err := getTransferClient.Request("GET", getTransferPath, &getTransferOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving LTS transfer: %s", err)
	}

	getTransferRespBody, err := utils.FlattenResponse(getTransferResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving LTS transfer: %s", err)
	}

	jsonPath := fmt.Sprintf("log_transfers[?log_transfer_id =='%s']|[0]", state.Primary.ID)
	getTransferRespBody = utils.PathSearch(jsonPath, getTransferRespBody, nil)
	if getTransferRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getTransferRespBody, nil
}

func TestAccLtsTransfer_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	obsName := acceptance.HCS_OBS_BUCKET_NAME
	rName := "hcs_lts_transfer.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLtsTransferResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLtsTransfer_basic(name, obsName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "hcs_lts_group.group", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_streams.0.log_stream_id",
						"hcs_lts_stream.stream", "id"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_type", "OBS"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_mode", "cycle"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_storage_format", "RAW"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_status", "ENABLE"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_period", "3"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_period_unit", "hour"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_dir_prefix_name", "lts_transfer_obs_"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_prefix_name", "obs_"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_time_zone", "UTC"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_time_zone_id", "Etc/GMT"),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
				),
			},
			{
				Config: testLtsTransfer_basic_update(name, obsName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "hcs_lts_group.group", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_streams.0.log_stream_id",
						"hcs_lts_stream.stream", "id"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_type", "OBS"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_mode", "cycle"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_storage_format", "RAW"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_status", "DISABLE"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_period", "2"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_period_unit", "min"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_dir_prefix_name", "lts_transfer_obs_2_"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_prefix_name", "obs_2_"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_time_zone", "UTC-02:00"),
					resource.TestCheckResourceAttr(rName, "log_transfer_info.0.log_transfer_detail.0.obs_time_zone_id", "Etc/GMT+2"),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testLtsTransfer_basic(name, obsName string) string {
	return fmt.Sprintf(`
resource "hcs_lts_group" "group" {
  group_name  = "%[1]s"
  ttl_in_days = 1
}
resource "hcs_lts_stream" "stream" {
  group_id    = hcs_lts_group.group.id
  stream_name = "%[1]s"
}

resource "hcs_lts_transfer" "test" {
  log_group_id = hcs_lts_group.group.id

  log_streams {
    log_stream_id = hcs_lts_stream.stream.id
  }

  log_transfer_info {
    log_transfer_type   = "OBS"
    log_transfer_mode   = "cycle"
    log_storage_format  = "RAW"
    log_transfer_status = "ENABLE"

    log_transfer_detail {
      obs_period          = 3
      obs_period_unit     = "hour"
      obs_bucket_name     = "%[2]s"
      obs_dir_prefix_name = "lts_transfer_obs_"
      obs_prefix_name     = "obs_"
      obs_time_zone       = "UTC"
      obs_time_zone_id    = "Etc/GMT"
    }
  }
}
`, name, obsName)
}

func testLtsTransfer_basic_update(name, obsName string) string {
	return fmt.Sprintf(`
resource "hcs_lts_group" "group" {
  group_name  = "%[1]s"
  ttl_in_days = 1
}
resource "hcs_lts_stream" "stream" {
  group_id    = hcs_lts_group.group.id
  stream_name = "%[1]s"
}

resource "hcs_lts_transfer" "test" {
  log_group_id = hcs_lts_group.group.id

  log_streams {
    log_stream_id = hcs_lts_stream.stream.id
  }

  log_transfer_info {
    log_transfer_type   = "OBS"
    log_transfer_mode   = "cycle"
    log_storage_format  = "RAW"
    log_transfer_status = "DISABLE"

    log_transfer_detail {
      obs_period          = 2
      obs_period_unit     = "min"
      obs_bucket_name     = "%[2]s"
      obs_dir_prefix_name = "lts_transfer_obs_2_"
      obs_prefix_name     = "obs_2_"
      obs_time_zone       = "UTC-02:00"
      obs_time_zone_id    = "Etc/GMT+2"
    }
  }
}
`, name, obsName)
}
