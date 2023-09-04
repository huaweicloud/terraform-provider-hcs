package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/evs/v2/snapshots"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
)

func TestAccEvsSnapshotV2_basic(t *testing.T) {
	var snapshot snapshots.Snapshot

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "hcs_evs_snapshot.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEvsSnapshotV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsSnapshotV2_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsSnapshotV2Exists(resourceName, &snapshot),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Daily backup"),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
				),
			},
		},
	})
}

func testAccCheckEvsSnapshotV2Destroy(s *terraform.State) error {
	cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	evsClient, err := cfg.BlockStorageV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweicloudStack EVS storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_evs_snapshot" {
			continue
		}

		_, err := snapshots.Get(evsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("EVS snapshot still exists")
		}
	}

	return nil
}

func testAccCheckEvsSnapshotV2Exists(n string, sp *snapshots.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		evsClient, err := cfg.BlockStorageV2Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweicloudStack EVS storage client: %s", err)
		}

		found, err := snapshots.Get(evsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("EVS snapshot not found")
		}

		*sp = *found

		return nil
	}
}

func testAccEvsSnapshotV2_basic(rName string) string {
	return fmt.Sprintf(`
data "hcs_availability_zones" "test" {}

resource "hcs_evs_volume" "test" {
  name              = "%s"
  description       = "Created by acc test"
  availability_zone = data.hcs_availability_zones.test.names[0]
  volume_type       = "business_type_01"
  size              = 12
}

resource "hcs_evs_snapshot" "test" {
  volume_id   = hcs_evs_volume.test.id
  name        = "%s"
  description = "Daily backup"
}
`, rName, rName)
}
