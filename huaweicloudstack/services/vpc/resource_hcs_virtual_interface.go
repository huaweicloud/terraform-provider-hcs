package vpc

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v2/dc/dc_endpoint_groups"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v2/dc/virtual_interfaces"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func ResourceVirtualInterface() *schema.Resource {
	linkInfo := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"interface_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hosting_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"local_gateway_v4_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
			},
			"local_gateway_v6_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
			},
			"remote_gateway_v4_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
			},
			"remote_gateway_v6_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
			},
			"vlan": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 4063),
			},
			"bgp_asn": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"bgp_asn_dot": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}

	return &schema.Resource{
		CreateContext: resourceVirtualInterfaceCreate,
		ReadContext:   resourceVirtualInterfaceRead,
		UpdateContext: resourceVirtualInterfaceUpdate,
		DeleteContext: resourceVirtualInterfaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.ValidateString64WithChinese,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"direct_connect_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vgw_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"remote_ep_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remote_ep_group": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
				Set: schema.HashString,
			},
			"link_infos": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     linkInfo,
			},
		},
	}
}

func resourceVirtualInterfaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	vifClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Direct Connect client: %s", err)
	}

	name := d.Get("name").(string)
	remoteEpGroup, err := createRemoteEpGroup(name, d.Get("remote_ep_group").(*schema.Set).List(), vifClient)
	if err != nil {
		return diag.Errorf("error creating Dc Endpoint Group of VIF: %s", err)
	}

	defer func() {
		// delete remote ep group when encounter errors
		if err != nil && remoteEpGroup != nil {
			deleteErr := dc_endpoint_groups.Delete(vifClient, remoteEpGroup.ID).Err
			if deleteErr != nil {
				log.Printf("[ERROR] Error deleting Dc Endpoint Group VIF: %s", deleteErr)
			}
		}
	}()

	createOpts := virtual_interfaces.CreateOpts{
		Name:            name,
		DirectConnectId: d.Get("direct_connect_id").(string),
		VgwId:           d.Get("vgw_id").(string),
		RemoteEpGroupId: remoteEpGroup.ID,
		Description:     d.Get("description").(string),
		LinkInfos:       buildLinkInfos(d),
	}

	n, err := virtual_interfaces.Create(vifClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating Direct Connect: %s", err)
	}

	d.SetId(n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE", "ERROR"},
		Refresh:    refreshVifStatus(vifClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return diag.Errorf(
			"error waiting for Virtual Interface (%s) to become ACTIVE or ERROR: %s",
			n.ID, stateErr)
	}

	log.Printf("[DEBUG] Virtual Gateway ID: %s", n.ID)
	return resourceVirtualInterfaceRead(ctx, d, meta)
}

func resourceVirtualInterfaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	vifClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Virtual Interface client: %s", err)
	}

	n, err := virtual_interfaces.Get(vifClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Virtual Interface Connection")
	}

	endpointGroup, err := dc_endpoint_groups.Get(vifClient, n.RemoteEpGroupId).Extract()
	if err != nil {
		return diag.Errorf("error retrieving Dc Endpoint Group of Virtual Interface: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", n.Name),
		d.Set("status", n.Status),
		d.Set("description", n.Description),
		d.Set("direct_connect_id", n.DirectConnectId),
		d.Set("vgw_id", n.VgwId),
		d.Set("remote_ep_group_id", n.RemoteEpGroupId),
		d.Set("remote_ep_group", endpointGroup.Endpoints),
		d.Set("link_infos", getLinkInfos(n)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Virtual Interface fields: %s", err)
	}

	return nil
}

func resourceVirtualInterfaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	vifClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Virtual Interface client: %s", err)
	}

	name := d.Get("name").(string)

	updateOpts := virtual_interfaces.UpdateOpts{
		Name:        name,
		Description: d.Get("description").(string),
	}

	if d.HasChange("remote_ep_group") {
		_, n := d.GetChange("remote_ep_group")
		remoteEpGroup, err := createRemoteEpGroup(name, n.(*schema.Set).List(), vifClient)
		if err != nil {
			return diag.Errorf("error creating Dc Endpoint Group of Virtual Interface: %s", err)
		}
		updateOpts.RemoteEpGroupId = remoteEpGroup.ID
		defer func() {
			// delete new remote ep group when encounter errors
			if err != nil && remoteEpGroup != nil {
				deleteErr := dc_endpoint_groups.Delete(vifClient, remoteEpGroup.ID).Err
				if deleteErr != nil {
					log.Printf("[ERROR] Error deleting Dc Endpoint Group of Virtual Interface: %s", deleteErr)
				}
			}
		}()
	}

	_, err = virtual_interfaces.Update(vifClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating Virtual Interface: %s", err)
	}

	if d.HasChange("remote_ep_group") {
		deleteErr := dc_endpoint_groups.Delete(vifClient, d.Get("remote_ep_group_id").(string)).Err
		if deleteErr != nil {
			log.Printf("[ERROR] Error deleting Dc Endpoint Group of Virtual Interface: %s", deleteErr)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_UPDATE"},
		Target:     []string{"ACTIVE", "ERROR"},
		Refresh:    refreshVifStatus(vifClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return diag.Errorf(
			"error waiting for Virtual Interface (%s) to become ACTIVE or ERROR: %s",
			d.Id(), stateErr)
	}

	return resourceVirtualInterfaceRead(ctx, d, meta)
}

func resourceVirtualInterfaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	vifClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Direct Connect client: %s", err)
	}

	err = virtual_interfaces.Delete(vifClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting Direct Connect: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_DELETE"},
		Target:     []string{"DELETED"},
		Refresh:    refreshVifStatus(vifClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return diag.Errorf(
			"error waiting for Virtual Interface (%s) to become ACTIVE or ERROR: %s",
			d.Id(), stateErr)
	}

	deleteErr := dc_endpoint_groups.Delete(vifClient, d.Get("remote_ep_group_id").(string)).Err
	if deleteErr != nil {
		log.Printf("[ERROR] Error deleting Dc Endpoint Group of Virtual Interface: %s", deleteErr)
	}

	return nil
}

func buildLinkInfos(d *schema.ResourceData) []virtual_interfaces.LinkInfoOpts {
	linkInfosData := d.Get("link_infos").(*schema.Set).List()
	linkInfos := make([]virtual_interfaces.LinkInfoOpts, len(linkInfosData))
	if len(linkInfosData) > 0 {
		for i, item := range linkInfosData {
			linkInfoData := item.(map[string]interface{})
			linkInfoOpts := virtual_interfaces.LinkInfoOpts{
				InterfaceGroupId:  linkInfoData["interface_group_id"].(string),
				HostingId:         linkInfoData["hosting_id"].(string),
				LocalGatewayV4Ip:  linkInfoData["local_gateway_v4_ip"].(string),
				LocalGatewayV6Ip:  linkInfoData["local_gateway_v6_ip"].(string),
				RemoteGatewayV4Ip: linkInfoData["remote_gateway_v4_ip"].(string),
				RemoteGatewayV6Ip: linkInfoData["remote_gateway_v6_ip"].(string),
				Vlan:              linkInfoData["vlan"].(int),
				BgpAsn:            linkInfoData["bgp_asn"].(int),
				BgpAsnDot:         linkInfoData["bgp_asn_dot"].(string),
			}
			linkInfos[i] = linkInfoOpts
		}
	}
	return linkInfos
}

func createRemoteEpGroup(name string, remoteEpGroups []interface{}, vifClient *golangsdk.ServiceClient) (*dc_endpoint_groups.DcEndpointGroup, error) {
	endpoints := make([]string, len(remoteEpGroups))
	for i, item := range remoteEpGroups {
		endpoints[i] = item.(string)
	}

	creatOpts := dc_endpoint_groups.CreateOpts{
		Name:      name,
		Type:      "cidr",
		Endpoints: endpoints,
	}
	return dc_endpoint_groups.Create(vifClient, creatOpts).Extract()
}

func getLinkInfos(vif *virtual_interfaces.VirtualInterface) []map[string]interface{} {
	linkInfos := make([]map[string]interface{}, len(vif.LinkInfos))
	if len(vif.LinkInfos) > 0 {
		for i, item := range vif.LinkInfos {
			linkInfos[i] = map[string]interface{}{
				"interface_group_id":   item.InterfaceGroupId,
				"hosting_id":           item.HostingId,
				"local_gateway_v4_ip":  item.LocalGatewayV4Ip,
				"local_gateway_v6_ip":  item.LocalGatewayV6Ip,
				"remote_gateway_v4_ip": item.RemoteGatewayV4Ip,
				"remote_gateway_v6_ip": item.RemoteGatewayV6Ip,
				"vlan":                 item.Vlan,
				"bgp_asn":              item.BgpAsn,
				"bgp_asn_dot":          item.BgpAsnDot,
			}
		}
	}
	return linkInfos
}

func refreshVifStatus(vifClient *golangsdk.ServiceClient, vifId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := virtual_interfaces.Get(vifClient, vifId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted Virtual Interface %s", vifId)
				return r, "DELETED", nil
			}
			return nil, "", err
		}

		return r, r.Status, nil
	}
}
