package common

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/eps/v1/enterpriseprojects"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// MigrateEnterpriseProjectWithoutWait is a method that used to a migrate resource from an enterprise project to
// another.
// NOTE: Please read the following contents carefully before using this method.
//   - This method only sends an asynchronous request and does not guarantee the result.
func MigrateEnterpriseProjectWithoutWait(cfg *config.Config, d *schema.ResourceData,
	opts enterpriseprojects.MigrateResourceOpts) error {
	targetEpsId := cfg.GetEnterpriseProjectID(d)
	if targetEpsId == "" {
		targetEpsId = "0"
	}

	client, err := cfg.EnterpriseProjectClient(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating EPS client: %s", err)
	}
	_, err = enterpriseprojects.Migrate(client, opts, targetEpsId).Extract()
	if err != nil {
		return fmt.Errorf("failed to migrate resource (%s) to the enterprise project (%s): %s",
			opts.ResourceId, targetEpsId, err)
	}
	return nil
}
