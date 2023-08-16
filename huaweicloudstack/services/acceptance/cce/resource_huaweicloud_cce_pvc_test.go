package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cce"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/cce/v1/persistentvolumeclaims"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
)

func getPvcResourceFunc(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.Config.CceV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloudStack CCE v1 client: %s", err)
	}
	resp, err := cce.GetCcePvcInfoById(c, state.Primary.Attributes["cluster_id"],
		state.Primary.Attributes["namespace"], state.Primary.ID)
	if resp == nil && err == nil {
		return resp, fmt.Errorf("unable to find the persistent volume claim (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccCcePersistentVolumeClaimsV1_basic(t *testing.T) {
	var pvc persistentvolumeclaims.PersistentVolumeClaim
	resourceName := "hcs_cce_pvc.test"
	randName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pvc,
		getPvcResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCcePersistentVolumeClaimsV1_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", acceptance.HCS_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttr(resourceName, "namespace", "default"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "storage_class_name", "csi-disk"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCEPVCImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"annotations",
				},
			},
		},
	})
}

func TestAccCcePersistentVolumeClaimsV1_obs(t *testing.T) {
	var pvc persistentvolumeclaims.PersistentVolumeClaim
	resourceName := "hcs_cce_pvc.test"
	randName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pvc,
		getPvcResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCcePersistentVolumeClaimsV1_obs(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", acceptance.HCS_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttr(resourceName, "namespace", "default"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "storage_class_name", "csi-obs"),
				),
			},
		},
	})
}

func TestAccCcePersistentVolumeClaimsV1_sfs(t *testing.T) {
	var pvc persistentvolumeclaims.PersistentVolumeClaim
	resourceName := "hcs_cce_pvc.test"
	randName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pvc,
		getPvcResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCcePersistentVolumeClaimsV1_sfs(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", acceptance.HCS_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttr(resourceName, "namespace", "default"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "storage_class_name", "csi-nas"),
				),
			},
		},
	})
}

func testAccCCEPVCImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmtp.Errorf("unable to find the PVC in the terraform state storage")
		}
		if rs.Primary.Attributes["cluster_id"] == "" || rs.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/default/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.ID), nil
	}
}

func testAccCcePersistentVolumeClaimsV1_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcs_cce_pvc" "test" {
  cluster_id  = "%[1]s"
  namespace   = "default"
  name        = "%[2]s"
  annotations = {
    "everest.io/disk-volume-type" = "SSD"
  }
  storage_class_name = "csi-disk"
  access_modes = ["ReadWriteOnce"]
  storage = "10Gi"
}
`, acceptance.HCS_CCE_CLUSTER_ID, rName)
}

func testAccCcePersistentVolumeClaimsV1_obs(rName string) string {
	return fmt.Sprintf(`
resource "hcs_cce_pvc" "test" {
  cluster_id  = "%[1]s"
  namespace   = "default"
  name        = "%[2]s"
  annotations = {
    "everest.io/obs-volume-type" = "STANDARD"
    "csi.storage.k8s.io/fstype" =  "obsfs"
  }
  storage_class_name = "csi-obs"
  access_modes = ["ReadWriteMany"]
  storage = "1Gi"
}
`, acceptance.HCS_CCE_CLUSTER_ID, rName)
}

func testAccCcePersistentVolumeClaimsV1_sfs(rName string) string {
	return fmt.Sprintf(`
resource "hcs_cce_pvc" "test" {
  cluster_id  = "%[1]s"
  namespace   = "default"
  name        = "%[2]s"
  storage_class_name = "csi-nas"
  access_modes = ["ReadWriteMany"]
  storage = "10Gi"
}
`, acceptance.HCS_CCE_CLUSTER_ID, rName)
}
