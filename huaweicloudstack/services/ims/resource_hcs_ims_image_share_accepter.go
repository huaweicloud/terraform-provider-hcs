// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product IMS
// ---------------------------------------------------------------

package ims

import (
	"context"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func ResourceImsImageShareAccepter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImsImageShareAccepterCreate,
		ReadContext:   resourceImsImageShareAccepterRead,
		DeleteContext: resourceImsImageShareAccepterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the image.`,
			},
		},
	}
}

func resourceImsImageShareAccepterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)

	// createImageShareAccepter: create IMS image share accepter
	var (
		createImageShareAccepterHttpUrl = "v1/cloudimages/members"
		createImageShareAccepterProduct = "ims"
	)
	createImageShareAccepterClient, err := cfg.NewServiceClient(createImageShareAccepterProduct, region)
	if err != nil {
		return diag.Errorf("error creating IMS Client: %s", err)
	}

	createImageShareAccepterPath := createImageShareAccepterClient.Endpoint + createImageShareAccepterHttpUrl

	createImageShareAccepterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createImageShareAccepterOpt.JSONBody = utils.RemoveNil(buildCreateImageShareAccepterBodyParams(d,
		createImageShareAccepterClient.ProjectID))

	createImageShareAccepterResp, err := createImageShareAccepterClient.Request("PUT",
		createImageShareAccepterPath, &createImageShareAccepterOpt)
	if err != nil {
		return diag.Errorf("error creating IMS image share accepter: %s", err)
	}

	createImageShareAccepterRespBody, err := utils.FlattenResponse(createImageShareAccepterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId, err := jmespath.Search("job_id", createImageShareAccepterRespBody)
	if err != nil {
		return diag.Errorf("error creating IMS image share accepter: job_id is not found in API response")
	}

	err = waitForJobSuccess(ctx, d, createImageShareAccepterClient, jobId.(string), schema.TimeoutCreate)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId)

	return resourceImsImageShareAccepterRead(ctx, d, meta)
}

func buildCreateImageShareAccepterBodyParams(d *schema.ResourceData, projectId string) map[string]interface{} {
	imagesParams := []interface{}{
		utils.ValueIgnoreEmpty(d.Get("image_id")),
	}
	bodyParams := map[string]interface{}{
		"images":     imagesParams,
		"project_id": projectId,
		"status":     "accepted",
	}
	return bodyParams
}

func resourceImsImageShareAccepterRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceImsImageShareAccepterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)

	// deleteImageShareAccepter: delete IMS image share accepter
	var (
		deleteImageShareAccepterHttpUrl = "v1/cloudimages/members"
		deleteImageShareAccepterProduct = "ims"
	)
	deleteImageShareAccepterClient, err := cfg.NewServiceClient(deleteImageShareAccepterProduct, region)
	if err != nil {
		return diag.Errorf("error creating IMS Client: %s", err)
	}

	deleteImageShareAccepterPath := deleteImageShareAccepterClient.Endpoint + deleteImageShareAccepterHttpUrl

	deleteImageShareAccepterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	deleteImageShareAccepterOpt.JSONBody = utils.RemoveNil(buildDeleteImageShareAccepterBodyParams(d,
		deleteImageShareAccepterClient.ProjectID))
	deleteImageShareAccepterResp, err := deleteImageShareAccepterClient.Request("PUT",
		deleteImageShareAccepterPath, &deleteImageShareAccepterOpt)
	if err != nil {
		return diag.Errorf("error deleting IMS image share accepter: %s", err)
	}

	deleteImageShareAccepterRespBody, err := utils.FlattenResponse(deleteImageShareAccepterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId, err := jmespath.Search("job_id", deleteImageShareAccepterRespBody)
	if err != nil {
		return diag.Errorf("error deleting IMS image share accepter: job_id is not found in API response")
	}

	err = waitForJobSuccess(ctx, d, deleteImageShareAccepterClient, jobId.(string), schema.TimeoutDelete)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildDeleteImageShareAccepterBodyParams(d *schema.ResourceData, projectId string) map[string]interface{} {
	imagesParams := []interface{}{
		utils.ValueIgnoreEmpty(d.Get("image_id")),
	}
	bodyParams := map[string]interface{}{
		"images":     imagesParams,
		"project_id": projectId,
		"status":     "rejected",
	}
	return bodyParams
}
