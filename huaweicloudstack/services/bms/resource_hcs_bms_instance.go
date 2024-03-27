package bms

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/bms/v1/baremetalservers"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v2/ports"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/logp"
)

func ResourceBmsInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBmsInstanceCreate,
		ReadContext:   resourceBmsInstanceRead,
		UpdateContext: resourceBmsInstanceUpdate,
		DeleteContext: resourceBmsInstanceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"nics": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// just stash the hash for state & diff comparisons
				StateFunc: utils.HashAndHexEncode,
			},
			"admin_pass": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
				ExactlyOneOf: []string{
					"admin_pass", "key_pair",
				},
			},
			"key_pair": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"admin_pass", "key_pair",
				},
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"eip_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"iptype", "bandwidth_size", "sharetype",
				},
			},
			"iptype": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"eip_id"},
				RequiredWith: []string{
					"sharetype", "bandwidth_size",
				},
			},
			"sharetype": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PER", "WHOLE",
				}, true),
				ConflictsWith: []string{"eip_id"},
				RequiredWith: []string{
					"iptype", "bandwidth_size",
				},
			},
			"bandwidth_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"eip_id"},
				RequiredWith: []string{
					"iptype", "sharetype",
				},
			},
			"data_disks": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 59,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"charging_mode": common.SchemaChargingMode([]string{}),
			"period_unit":   common.SchemaPeriodUnit([]string{}),
			"period":        common.SchemaPeriod([]string{}),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),

			"tags": common.TagsForceNewSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"agency_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"delete_eip_on_termination": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"delete_disks_on_termination": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceBmsInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	bmsClient, err := cfg.BmsV1Client(cfg.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack bms client: %s", err)
	}

	createOpts := &baremetalservers.CreateOpts{
		Name:      d.Get("name").(string),
		ImageRef:  d.Get("image_id").(string),
		FlavorRef: d.Get("flavor_id").(string),
		MetaData: baremetalservers.MetaData{
			OpSvcUserId: d.Get("user_id").(string),
			AgencyName:  d.Get("agency_name").(string),
		},
		UserData:         []byte(d.Get("user_data").(string)),
		AdminPass:        d.Get("admin_pass").(string),
		KeyName:          d.Get("key_pair").(string),
		VpcId:            d.Get("vpc_id").(string),
		SecurityGroups:   resourceBmsInstanceSecGroupsV1(d),
		AvailabilityZone: d.Get("availability_zone").(string),
		Nics:             resourceBmsInstanceNicsV1(d),
		DataVolumes:      resourceBmsInstanceDataVolumesV1(d),
		ExtendParam: baremetalservers.ServerExtendParam{
			ChargingMode:        normalizeChargingModeToNumber(d.Get("charging_mode").(string)),
			EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		},
	}

	var eipOpts baremetalservers.PublicIp
	var hasEIP bool
	if eipID, ok := d.GetOk("eip_id"); ok {
		hasEIP = true
		eipOpts.Id = eipID.(string)
	} else if eipType, ok := d.GetOk("iptype"); ok {
		hasEIP = true
		eipOpts.Eip = &baremetalservers.Eip{
			IpType: eipType.(string),
			BandWidth: baremetalservers.BandWidth{
				ShareType: d.Get("sharetype").(string),
				Size:      d.Get("bandwidth_size").(int),
			},
		}
	}
	if hasEIP {
		createOpts.PublicIp = &eipOpts
	}

	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		tagList := utils.ExpandResourceTagsString(tagRaw)
		createOpts.Tags = tagList
	}

	n, err := baremetalservers.Create(bmsClient, createOpts).ExtractJobResponse()
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack BMS server: %s", err)
	}
	err = baremetalservers.WaitForJobSuccess(bmsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID)
	if err != nil {
		return diag.FromErr(err)
	}

	entity, err := baremetalservers.GetJobEntity(bmsClient, n.JobID, "server_id")
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := entity.(string); ok {
		log.Printf("[INFO] BMS ID: %s", id)
		// Store the ID now
		d.SetId(id)
		return resourceBmsInstanceRead(ctx, d, meta)
	}
	return diag.Errorf("unexpected conversion error in resourceBmsInstanceCreate.")
}

func normalizeChargingModeToNumber(mode string) string {
	var chargingModeNum string
	switch mode {
	case "prePaid":
		chargingModeNum = "1"
	default:
		chargingModeNum = "0"
	}

	return chargingModeNum
}

func resourceBmsInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	bmsClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack bms client: %s", err)
	}

	opts := buildListOpts(d, cfg)
	server, err := baremetalservers.Get(bmsClient, d.Id(), opts).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "server")
	}
	if server.Status == "DELETED" {
		d.SetId("")
		return nil
	}

	logp.Printf("[DEBUG] Retrieved Server %s: %+v", d.Id(), server)

	d.Set("region", region)
	d.Set("name", server.Name)
	d.Set("image_id", server.Image.ID)
	d.Set("flavor_id", server.Flavor.ID)
	d.Set("host_id", server.HostID)

	// Set fixed and floating ip
	if eip := bmsPublicIP(server); eip != "" {
		d.Set("public_ip", eip)
	}
	nics := flattenBmsInstanceNicsV1(d, meta, server.Addresses)
	d.Set("nics", nics)

	d.Set("key_pair", server.KeyName)
	// Set security groups
	secGrpIds := []string{}
	for _, sg := range server.SecurityGroups {
		secGrpIds = append(secGrpIds, sg.ID)
	}
	d.Set("security_groups", secGrpIds)
	d.Set("status", server.Status)
	d.Set("user_id", server.Metadata.OpSvcUserId)
	d.Set("image_name", server.Metadata.ImageName)
	d.Set("vpc_id", server.Metadata.VpcID)
	d.Set("availability_zone", server.AvailabilityZone)
	d.Set("description", server.Description)
	d.Set("user_data", server.UserData)
	d.Set("enterprise_project_id", server.EnterpriseProjectID)
	// Set disk ids
	diskIds := []string{}
	for _, disk := range server.VolumeAttached {
		diskIds = append(diskIds, disk.ID)
	}
	d.Set("disk_ids", diskIds)
	return nil
}

func buildListOpts(d *schema.ResourceData, conf *config.HcsConfig) *baremetalservers.ListOpts {
	result := &baremetalservers.ListOpts{
		EnterpriseProjectID: conf.DataGetEnterpriseProjectID(d),
		Name:                d.Get("name").(string),
		Tags:                "__type_baremetal",
	}

	return result
}

func resourceBmsInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceBmsInstanceRead(ctx, d, meta)
}

func resourceBmsInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	bmsClient, err := cfg.BmsV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack compute client: %s", err)
	}

	deleteOpts := baremetalservers.DeleteOpts{
		Servers: []baremetalservers.Server{
			{Id: d.Id()},
		},
		DeletePublicIp: d.Get("delete_eip_on_termination").(bool),
		DeleteVolume:   d.Get("delete_disks_on_termination").(bool),
	}
	n, err := baremetalservers.Delete(bmsClient, deleteOpts).ExtractJobResponse()
	if err != nil {
		return diag.Errorf("error deleting server: %s", err)
	}

	if err := baremetalservers.WaitForJobSuccess(bmsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return diag.FromErr(err)
	}

	client, err := cfg.ComputeV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack bms client: %s", err)
	}

	opts := buildListOpts(d, cfg)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Deleting", "ACTIVE", "SHUTOFF"},
		Target:       []string{"DELETED"},
		Refresh:      waitForBmsInstanceDelete(client, d.Id(), opts),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf("Error deleting HuaweiCloudStack BMS instance: %s", err)
	}

	d.SetId("")
	return nil
}
func resourceBmsInstanceNicsV1(d *schema.ResourceData) []baremetalservers.Nic {
	var nicRequests []baremetalservers.Nic

	nics := d.Get("nics").([]interface{})
	for i := range nics {
		nic := nics[i].(map[string]interface{})
		nicRequest := baremetalservers.Nic{
			SubnetId:  nic["subnet_id"].(string),
			IpAddress: nic["ip_address"].(string),
		}

		nicRequests = append(nicRequests, nicRequest)
	}
	return nicRequests
}

func resourceBmsInstanceDataVolumesV1(d *schema.ResourceData) []baremetalservers.DataVolume {
	var volRequests []baremetalservers.DataVolume

	vols := d.Get("data_disks").([]interface{})
	for i := range vols {
		vol := vols[i].(map[string]interface{})
		volRequest := baremetalservers.DataVolume{
			VolumeType: vol["type"].(string),
			Size:       vol["size"].(int),
		}
		volRequests = append(volRequests, volRequest)
	}
	return volRequests
}

func resourceBmsInstanceSecGroupsV1(d *schema.ResourceData) []baremetalservers.SecurityGroup {
	rawSecGroups := d.Get("security_groups").(*schema.Set).List()
	secgroups := make([]baremetalservers.SecurityGroup, len(rawSecGroups))
	for i, raw := range rawSecGroups {
		secgroups[i] = baremetalservers.SecurityGroup{
			ID: raw.(string),
		}
	}
	return secgroups
}

func flattenBmsInstanceNicsV1(d *schema.ResourceData, meta interface{},
	addresses map[string][]baremetalservers.Address) []map[string]interface{} {

	config := config.GetHcsConfig(meta)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		logp.Printf("Error creating HuaweiCloudStack networking client: %s", err)
	}

	var network string
	nics := []map[string]interface{}{}
	// Loop through all networks and addresses.
	for _, addrs := range addresses {
		for _, addr := range addrs {
			// Skip if not fixed ip
			if addr.Type != "fixed" {
				continue
			}

			p, err := ports.Get(networkingClient, addr.PortID).Extract()
			if err != nil {
				network = ""
				logp.Printf("[DEBUG] flattenInstanceNicsV1: failed to fetch port %s", addr.PortID)
			} else {
				network = p.NetworkID
			}

			v := map[string]interface{}{
				"subnet_id":   network,
				"ip_address":  addr.Addr,
				"mac_address": addr.MacAddr,
				"port_id":     addr.PortID,
			}
			nics = append(nics, v)
		}
	}

	logp.Printf("[DEBUG] flattenInstanceNicsV1: %#v", nics)
	return nics
}

func bmsPublicIP(server *baremetalservers.CloudServer) string {
	var publicIP string

	for _, addresses := range server.Addresses {
		for _, addr := range addresses {
			if addr.Type == "floating" {
				publicIP = addr.Addr
				break
			}
		}
	}

	return publicIP
}

func waitForBmsInstanceDelete(client *golangsdk.ServiceClient, ServerId string, opts baremetalservers.ListOptsBuilder) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logp.Printf("[DEBUG] Attempting to delete HuaweiCloudStack BMS instance %s", ServerId)

		r, err := baremetalservers.Get(client, ServerId, opts).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloudStack BMS instance %s", ServerId)
				return r, "DELETED", nil
			}
			return r, "Deleting", err
		}

		return r, r.Status, nil
	}
}
