// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product IMS
// ---------------------------------------------------------------

package ims

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func ResourceImsImageShare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImsImageShareCreate,
		UpdateContext: resourceImsImageShareUpdate,
		ReadContext:   resourceImsImageShareRead,
		DeleteContext: resourceImsImageShareDelete,

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
			"source_image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the source image.`,
			},
			"target_project_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: `Specifies the IDs of the target projects.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceImsImageShareCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	sourceImageId := d.Get("source_image_id")
	projectIds := d.Get("target_project_ids")
	err := dealImageMembers(ctx, d, cfg, "POST", schema.TimeoutCreate, sourceImageId, projectIds.(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(sourceImageId.(string))

	return resourceImsImageShareRead(ctx, d, meta)
}

func resourceImsImageShareUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)

	if d.HasChange("target_project_ids") {
		oProjectIdsRaw, nProjectIdsRaw := d.GetChange("target_project_ids")
		shareProjectIds := nProjectIdsRaw.(*schema.Set).Difference(oProjectIdsRaw.(*schema.Set))
		unShareProjectIds := oProjectIdsRaw.(*schema.Set).Difference(nProjectIdsRaw.(*schema.Set))
		if shareProjectIds.Len() > 0 {
			err := dealImageMembers(ctx, d, cfg, "POST", schema.TimeoutCreate, nil, shareProjectIds.List())
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if unShareProjectIds.Len() > 0 {
			err := dealImageMembers(ctx, d, cfg, "DELETE", schema.TimeoutDelete, nil, unShareProjectIds.List())
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceImsImageShareRead(ctx, d, meta)
}

func resourceImsImageShareRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceImsImageShareDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)

	projectIds := d.Get("target_project_ids")
	err := dealImageMembers(ctx, d, cfg, "DELETE", schema.TimeoutDelete, nil, projectIds.(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func dealImageMembers(ctx context.Context, d *schema.ResourceData, cfg *config.HcsConfig, requestMethod,
	timeout string, imageId interface{}, projectIds []interface{}) error {
	region := cfg.GetRegion(d)
	var (
		imageMemberHttpUrl = "v1/cloudimages/members"
		imageMemberProduct = "ims"
	)

	imageMemberClient, err := cfg.NewServiceClient(imageMemberProduct, region)
	if err != nil {
		return fmt.Errorf("error creating IMS Client: %s", err)
	}

	imageMemberPath := imageMemberClient.Endpoint + imageMemberHttpUrl

	imageMemberOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	imageMemberOpt.JSONBody = utils.RemoveNil(buildImageMemberBodyParams(d, imageId, projectIds))
	imageMemberResp, err := imageMemberClient.Request(requestMethod, imageMemberPath, &imageMemberOpt)
	operateMethod := "creating"
	if requestMethod == "DELETE" {
		operateMethod = "deleting"
	}
	if err != nil {
		return fmt.Errorf("error %s IMS image share: %s", operateMethod, err)
	}

	imageMemberRespBody, err := utils.FlattenResponse(imageMemberResp)
	if err != nil {
		return err
	}

	jobId, err := jmespath.Search("job_id", imageMemberRespBody)
	if err != nil {
		return fmt.Errorf("error %s IMS image share: job_id is not found in API response", operateMethod)
	}

	err = waitForJobSuccess(ctx, d, imageMemberClient, jobId.(string), timeout)
	if err != nil {
		return err
	}
	return nil
}

func buildImageMemberBodyParams(d *schema.ResourceData, imageId interface{}, projectIds []interface{}) map[string]interface{} {
	imagesParams := []interface{}{
		utils.ValueIgnoreEmpty(d.Id()),
	}
	if len(imagesParams) == 0 || imagesParams[0] == nil {
		imageSlice := make([]interface{}, 0)
		imagesParams = append(imageSlice, imageId)
	}
	bodyParams := map[string]interface{}{
		"images":   imagesParams,
		"projects": projectIds,
	}
	return bodyParams
}

func waitForJobSuccess(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	jobId, timeout string) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INIT", "RUNNING"},
		Target:     []string{"SUCCESS"},
		Refresh:    imsJobStatusRefreshFunc(jobId, client),
		Timeout:    d.Timeout(timeout),
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for job (%s) success: %s", jobId, err)
	}
	return nil
}

func imsJobStatusRefreshFunc(jobId string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			getJobStatusHttpUrl = "v1/{project_id}/jobs/{job_id}"
		)

		getJobStatusPath := client.Endpoint + getJobStatusHttpUrl
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{project_id}", client.ProjectID)
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{job_id}", fmt.Sprintf("%v", jobId))

		getJobStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		getJobStatusResp, err := client.Request("GET", getJobStatusPath, &getJobStatusOpt)
		if err != nil {
			return getJobStatusResp, "FAIL", nil
		}

		getJobStatusRespBody, err := utils.FlattenResponse(getJobStatusResp)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch("status", getJobStatusRespBody, "")
		return getJobStatusRespBody, status.(string), nil
	}
}
