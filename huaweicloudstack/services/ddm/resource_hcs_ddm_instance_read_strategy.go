package ddm

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

var ddmInstanceReadStrategyNonUpdatableParams = []string{
	"instance_id",

}

// @API DDM PUT /v2/{project_id}/instances/{instance_id}/action/read-write-strategy
func ResourceDdmInstanceReadStrategy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDdmReadStrategyCreateOrUpdate,
		ReadContext:   resourceDdmReadStrategyRead,
		UpdateContext: resourceDdmReadStrategyCreateOrUpdate,
		DeleteContext: resourceDdmReadStrategyDelete,

		CustomizeDiff: config.FlexibleForceNew(ddmInstanceReadStrategyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of the DDM instance.",
			},
			"read_weights": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the ID of the DB instance associated with the DDM schema.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Specifies read weight of the DB instance associated with the DDM schema.",
						},
					},
				},
				Required:    true,
				Description: `Specifies the list of read weights of the primary DB instance and its read replicas.`,
			},

			// Internal
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceDdmReadStrategyCreateOrUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)

	var (
		dbReadStrategyHttpUrl = "v2/{project_id}/instances/{instance_id}/action/read-write-strategy"
		dbReadStrategyProduct = "ddm"
	)
	dbReadStrategyClient, err := cfg.NewServiceClient(dbReadStrategyProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	dbReadStrategyPath := dbReadStrategyClient.Endpoint + dbReadStrategyHttpUrl
	dbReadStrategyPath = strings.ReplaceAll(dbReadStrategyPath, "{project_id}", dbReadStrategyClient.ProjectID)
	dbReadStrategyPath = strings.ReplaceAll(dbReadStrategyPath, "{instance_id}", fmt.Sprintf("%v", instanceID))

	dbReadStrategyOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	dbReadStrategyOpt.JSONBody = buildDbReadStrategyBodyParams(d)
	_, err = dbReadStrategyClient.Request("PUT", dbReadStrategyPath, &dbReadStrategyOpt)
	if err != nil {
		return diag.Errorf("error setting read strategy of the DDM instance: %s", err)
	}
	d.SetId(instanceID)

	return nil
}

func buildDbReadStrategyBodyParams(d *schema.ResourceData) map[string]interface{} {
	readWeightMap := make(map[string]interface{})
	readWeightList := d.Get("read_weights").(*schema.Set).List()
	for _, readWeight := range readWeightList {
		variable := readWeight.(map[string]interface{})
		readWeightMap[variable["db_id"].(string)] = variable["weight"]
	}

	bodyParams := map[string]interface{}{
		"read_weight": readWeightMap,
	}
	return bodyParams
}

func resourceDdmReadStrategyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDdmReadStrategyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to read ddm instance strategy. Deleting this resource
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
