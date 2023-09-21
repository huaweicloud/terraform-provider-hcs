package eps

import (
	"context"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/eps/v1/enterpriseprojects"
)

func ResourceEnterpriseProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnterpriseProjectCreate,
		ReadContext:   resourceEnterpriseProjectRead,
		UpdateContext: resourceEnterpriseProjectUpdate,
		DeleteContext: resourceEnterpriseProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		//request and response parameters
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa5a-zA-Z0-9_-]{1,64}$"),
						"The name consists of 1 to 64 characters, and only contains letters, digits, "+
							"underscores (_), and hyphens (-)."),
					validation.StringDoesNotMatch(regexp.MustCompile("(?i)default"),
						"The name cannot include any form of the word 'default'"),
				),
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile("^[a-z0-9-]{1,36}$"),
						"Resource set. The value can contain 1 to 36 characters, "+
							"including only lowercase letters, digits, and hyphens (-)."),
				),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceEnterpriseProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := config.GetHcsConfig(meta)
	epsClient, err := config.EnterpriseProjectClient(config.GetRegion(d))

	if err != nil {
		return fmtp.DiagErrorf("Unable to create HuaweiCloudStack EPS client : %s", err)
	}

	createOpts := enterpriseprojects.CreateOpts{
		Name:        d.Get("name").(string),
		ProjectId:   d.Get("project_id").(string),
		Description: d.Get("description").(string),
	}

	project, err := enterpriseprojects.Create(epsClient, createOpts).Extract()

	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack Enterprise Project: %s", err)
	}

	d.SetId(project.ID)

	return resourceEnterpriseProjectRead(ctx, d, meta)
}

func resourceEnterpriseProjectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := config.GetHcsConfig(meta)
	epsClient, err := config.EnterpriseProjectClient(config.GetRegion(d))

	if err != nil {
		return fmtp.DiagErrorf("Unable to create HuaweiCloudStack EPS client : %s", err)
	}

	project, err := enterpriseprojects.Get(epsClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error retrieving HuaweiCloudStack Enterprise Project")
	}

	mErr := multierror.Append(nil,
		d.Set("name", project.Name),
		d.Set("description", project.Description),
		d.Set("project_id", project.ProjectId),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting HuaweiCloudStack enterprise project fields: %w", err)
	}

	return nil
}

func resourceEnterpriseProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := config.GetHcsConfig(meta)
	epsClient, err := config.EnterpriseProjectClient(config.GetRegion(d))

	if err != nil {
		return fmtp.DiagErrorf("Unable to create HuaweiCloudStack EPS client : %s", err)
	}

	if d.HasChanges("name", "description") {
		updateOpts := enterpriseprojects.CreateOpts{
			Name:        d.Get("name").(string),
			ProjectId:   d.Get("project_id").(string),
			Description: d.Get("description").(string),
		}

		_, err = enterpriseprojects.Update(epsClient, updateOpts, d.Id()).Extract()

		if err != nil {
			return fmtp.DiagErrorf("Error updating HuaweiCloudStack Enterprise Project: %s", err)
		}
	}

	return resourceEnterpriseProjectRead(ctx, d, meta)
}

func resourceEnterpriseProjectDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
