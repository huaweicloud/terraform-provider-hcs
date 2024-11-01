package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	v1peerings "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/peerings"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func ResourceVpcPeering() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcPeeringCreate,
		ReadContext:   resourceVpcPeeringRead,
		UpdateContext: resourceVpcPeeringUpdate,
		DeleteContext: resourceVpcPeeringDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{ // request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.ValidateString64WithChinese,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"peer_vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"peer_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"peer_region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVpcPeeringCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	region := hcsConfig.GetRegion(d)
	vpcPeeringClient, err := hcsConfig.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC Peering client: %s", err)
	}

	createOpts := v1peerings.CreateOpts{
		Name:          d.Get("name").(string),
		LocalVpcId:    d.Get("vpc_id").(string),
		PeerVpcId:     d.Get("peer_vpc_id").(string),
		PeerRegion:    d.Get("peer_region").(string),
		PeerProjectId: d.Get("peer_project_id").(string),
	}

	n, err := v1peerings.Create(vpcPeeringClient, createOpts).ExtractCreate()
	if err != nil {
		return diag.Errorf("error creating VPC Peering: %s", err)
	}

	d.SetId(n.ID)
	log.Printf("[DEBUG] Vpc Peering ID: %s", n.ID)

	return resourceVpcPeeringRead(ctx, d, meta)
}

func resourceVpcPeeringRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	region := hcsConfig.GetRegion(d)
	vpcPeeringClient, err := hcsConfig.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC Peering client: %s", err)
	}

	n, err := v1peerings.Get(vpcPeeringClient, d.Id()).ExtractList()
	if err != nil {
		return diag.Errorf("error geting VPC Peering client: %s", err)
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
		d.Set("name", n[0].Name),
		d.Set("vpc_id", n[0].RequesterVpcInfo.VpcId),
		d.Set("peer_vpc_id", n[0].AccepterVpcInfo.VpcId),
		d.Set("peer_project_id", n[0].AccepterVpcInfo.TenantId),
		d.Set("status", n[0].Status),
		d.Set("region", hcsConfig.GetRegion(d)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceVpcPeeringUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	region := hcsConfig.GetRegion(d)
	vpcPeeringClient, err := hcsConfig.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC Peering client: %s", err)
	}

	peerId := d.Id()
	if d.HasChanges("name") {
		updateOpts := v1peerings.UpdateOpts{
			Name: d.Get("name").(string),
		}
		if err := v1peerings.Update(vpcPeeringClient, peerId, updateOpts).Err; err != nil {
			return diag.Errorf("error updating VPC Peering: %s", err)
		}
	}

	return resourceVpcPeeringRead(ctx, d, meta)
}

func resourceVpcPeeringDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	region := hcsConfig.GetRegion(d)
	vpcPeeringClient, err := hcsConfig.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC Peering client: %s", err)
	}

	err = v1peerings.Delete(vpcPeeringClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting VPC Peering %s: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}
