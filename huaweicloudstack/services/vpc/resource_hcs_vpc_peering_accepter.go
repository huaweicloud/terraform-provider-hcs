package vpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	v1peerings "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/peerings"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceVpcPeeringAccepter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePeeringAccepterCreate,
		ReadContext:   resourcePeeringAccepterRead,
		UpdateContext: resourcePeeringAccepterUpdate,
		DeleteContext: resourcePeeringAccepterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_peering_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"accept": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePeeringAccepterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	peeringClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC Peering Connection client: %s", err)
	}

	id := d.Get("vpc_peering_id").(string)
	n, err := v1peerings.Get(peeringClient, id).ExtractList()
	if err != nil {
		return diag.Errorf("error retrieving Vpc Peering Connection: %s", err)
	}
	if len(n) == 0 {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("the resource %s is gone and will be removed in Terraform state.", d.Id()),
			},
		}
	}

	if n[0].Status != "accepting" {
		return diag.Errorf("VPC peering action not permitted: Can not accept/reject peering request not in PENDING_ACCEPTANCE state.")
	}

	var expectedStatus string

	if _, ok := d.GetOk("accept"); ok {
		expectedStatus = "active"

		err := v1peerings.Accept(peeringClient, id).ExtractErr()
		if err != nil {
			return diag.Errorf("unable to accept VPC Peering Connection: %s", err)
		}
	} else {
		expectedStatus = "rejected"

		err := v1peerings.Reject(peeringClient, id).ExtractErr()
		if err != nil {
			return diag.Errorf("unable to reject VPC Peering Connection: %s", err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"accepting"},
		Target:     []string{expectedStatus},
		Refresh:    waitForPeeringConnStatus(peeringClient, n[0].ID, expectedStatus),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the VPC Peering Connection: %s", err)
	}

	d.SetId(n[0].ID)
	return resourcePeeringAccepterRead(ctx, d, meta)
}

func resourcePeeringAccepterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	peeringClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC Peering Connection client: %s", err)
	}

	n, err := v1peerings.Get(peeringClient, d.Id()).ExtractList()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VPC Peering Connection")
	}
	if len(n) == 0 {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("the resource %s is gone and will be removed in Terraform state.", d.Id()),
			},
		}
	}
	if len(n) == 0 {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("the resource %s is gone and will be removed in Terraform state.", d.Id()),
			},
		}
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", n[0].Name),
		d.Set("status", n[0].Status),
		d.Set("vpc_id", n[0].RequesterVpcInfo.VpcId),
		d.Set("peer_vpc_id", n[0].AccepterVpcInfo.VpcId),
		d.Set("peer_project_id", n[0].AccepterVpcInfo.TenantId),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting VPC Peering Connection fields: %s", err)
	}

	return nil
}

func resourcePeeringAccepterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("accept") {
		return diag.Errorf("VPC peering action not permitted: Can not accept/reject peering request not in pending_acceptance state.")
	}

	return resourceVpcPeeringAccepterRead(ctx, d, meta)
}

func resourcePeeringAccepterDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	log.Printf("[WARN] Will not delete VPC peering connection. Terraform will remove this resource from the state file, resources may remain.")
	d.SetId("")
	return nil
}

func waitForPeeringConnStatus(peeringClient *golangsdk.ServiceClient, peeringId, expectedStatus string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := v1peerings.Get(peeringClient, peeringId).ExtractList()
		if err != nil {
			return nil, "", err
		}
		if len(n) == 0 {
			return nil, "", errors.New(fmt.Sprintf("the resource %s is gone and will be removed in Terraform state.", peeringId))
		}

		if n[0].Status == expectedStatus {
			return n[0], expectedStatus, nil
		}

		return n, "accepting", nil
	}
}
