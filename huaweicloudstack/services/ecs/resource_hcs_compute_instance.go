package ecs

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/helper/hashcode"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/compute/v2/extensions/secgroups"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/ecs/v1/block_devices"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/ecs/v1/cloudservers"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/ecs/v1/flavors"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/ecs/v1/powers"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/evs/v2/cloudvolumes"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/ims/v2/cloudimages"
	groups "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/security/securitygroups"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/subnets"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v2/ports"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/evs"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

var (
	powerActionMap = map[string]string{
		"ON":     "os-start",
		"OFF":    "os-stop",
		"REBOOT": "reboot",
	}
)

func ResourceComputeInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeInstanceCreate,
		ReadContext:   resourceComputeInstanceRead,
		UpdateContext: resourceComputeInstanceUpdate,
		DeleteContext: resourceComputeInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceComputeInstanceImportState,
		},

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
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("HW_IMAGE_ID", nil),
			},
			"image_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("HW_IMAGE_NAME", nil),
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("HW_FLAVOR_ID", nil),
				Description: "schema: Required",
			},
			"ext_boot_type": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Deprecated",
				ValidateFunc: validation.StringInSlice([]string{
					"LocalDisk", "Volume",
				}, false),
			},
			"flavor_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("HW_FLAVOR_NAME", nil),
				Description: "schema: Computed",
			},
			"admin_pass": {
				Type:          schema.TypeString,
				Sensitive:     true,
				Optional:      true,
				ConflictsWith: []string{"key_pair"},
			},
			"key_pair": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"security_groups": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Description:   "schema: Computed",
				ConflictsWith: []string{"security_group_ids"},
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"network": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 12,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "schema: Required",
						},
						"port": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "schema: Computed",
						},
						"ipv6_enable": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"fixed_ip_v4": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"source_dest_check": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"fixed_ip_v6": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "schema: Computed",
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_network": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"system_disk_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"system_disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"encrypt_cipher": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"AES256-XTS", "SM4-XTS",
				}, false),
			},
			"data_disks": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 23,
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
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"encrypt_cipher": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"AES256-XTS", "SM4-XTS",
							}, false),
						},
					},
				},
			},
			"scheduler_hints": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"fault_domain": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "schema: Internal",
						},
						"tenancy": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"deh_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
				Set: resourceComputeSchedulerHintsHash,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// just stash the hash for state & diff comparisons
				StateFunc: utils.HashAndHexEncode,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"delete_disks_on_termination": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"delete_eip_on_termination": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"eip_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"eip_type", "bandwidth"},
			},
			"eip_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"eip_id"},
				RequiredWith:  []string{"bandwidth"},
			},
			"bandwidth": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				MaxItems:      1,
				ConflictsWith: []string{"eip_id"},
				RequiredWith:  []string{"eip_type"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"share_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"PER", "WHOLE",
							}, true),
						},
						"id": {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"bandwidth.0.size", "bandwidth.0.charge_mode"},
						},
						"size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							RequiredWith: []string{"bandwidth.0.charge_mode"},
						},
						"charge_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							RequiredWith: []string{"bandwidth.0.size"},
						},
					},
				},
			},
			"user_id": { // required if in prePaid charging mode with key_pair.
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"power_action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				// If you want to support more actions, please update powerActionMap simultaneously.
				ValidateFunc: validation.StringInSlice([]string{
					"ON", "OFF", "REBOOT", "FORCE-OFF", "FORCE-REBOOT",
				}, false),
			},
			"volume_attached": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"boot_index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pci_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encrypt_cipher": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"system_disk_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_ip_v4": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_ip_v6": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	ecsClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute V1 client: %s", err)
	}
	imsClient, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating image client: %s", err)
	}
	nicClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating networking client: %s", err)
	}

	// Determines the Image ID using the following rules:
	// If an image_id was specified, use it.
	// If an image_name was specified, look up the image ID, report if error.
	imageId, err := getImageIDFromConfig(imsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	flavor, err := getFlavor(ecsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	ecsV11Client, err := cfg.ComputeV11Client(region)
	if err != nil {
		return diag.Errorf("error creating compute V1.1 client: %s", err)
	}
	vpcClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating networking V1 client: %s", err)
	}

	vpcId, err := getVpcID(vpcClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	secGroups, err := resourceInstanceSecGroupIdsV1(vpcClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := &cloudservers.CreateOpts{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		ImageRef:         imageId,
		FlavorRef:        flavor["id"].(string),
		KeyName:          d.Get("key_pair").(string),
		VpcId:            vpcId,
		SecurityGroups:   secGroups,
		AvailabilityZone: d.Get("availability_zone").(string),
		RootVolume:       resourceInstanceRootVolume(d, flavor["ext_boot_type"].(string)),
		DataVolumes:      resourceInstanceDataVolumes(d),
		Nics:             buildInstanceNicsRequest(d),
		PublicIp:         buildInstancePublicIPRequest(d),
		UserData:         []byte(d.Get("user_data").(string)),
	}

	if tags, ok := d.GetOk("tags"); ok {
		if !checkTags(tags.(map[string]interface{})) {
			return diag.Errorf("tags check failed")
		}
		tagList := utils.ExpandResourceTagsString(tags.(map[string]interface{}))
		for _, tag := range tagList {
			createOpts.Tags = append(createOpts.Tags, tag.(string))
		}
	}

	var extendParam cloudservers.ServerExtendParam
	epsID := cfg.GetEnterpriseProjectID(d)
	if epsID != "" {
		extendParam.EnterpriseProjectId = epsID
	}

	extBootType := flavor["ext_boot_type"].(string)
	if extBootType == "LocalDisk" {
		extendParam.Image_Boot = true
	}
	if extendParam != (cloudservers.ServerExtendParam{}) {
		createOpts.ExtendParam = &extendParam
	}
	schedulerHintsRaw := d.Get("scheduler_hints").(*schema.Set).List()
	if len(schedulerHintsRaw) > 0 {
		log.Printf("[DEBUG] schedulerhints: %+v", schedulerHintsRaw)
		schedulerHints := resourceInstanceSchedulerHintsV1(schedulerHintsRaw[0].(map[string]interface{}))
		createOpts.SchedulerHints = &schedulerHints
	}

	// Create an instance in the shutdown state.
	if action, ok := d.GetOk("power_action"); ok {
		action := action.(string)
		var PowerOn bool
		if action == "ON" {
			PowerOn = true
		} else if action == "OFF" {
			PowerOn = false
		} else {
			log.Printf("[ERROR] the power action (%s) is invalid after instance created, the value of power_action must be ON or OFF.", action)
			return diag.Errorf("the power action (%s) is invalid after instance created, the value of power_action must be ON or OFF.", action)
		}
		createOpts.PowerOn = &PowerOn
	}

	log.Printf("[DEBUG] ECS create options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.AdminPass = d.Get("admin_pass").(string)

	n, err := cloudservers.Create(ecsV11Client, createOpts).ExtractJobResponse()
	if err != nil {
		return diag.Errorf("error creating server: %s", err)
	}
	if err := cloudservers.WaitForJobSuccess(ecsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return diag.FromErr(err)
	}
	serverId, err := cloudservers.GetJobEntity(ecsClient, n.JobID, "server_id")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(serverId.(string))

	// get the original value of source_dest_check in script
	originalNetworks := d.Get("network").([]interface{})
	sourceDestChecks := make([]bool, len(originalNetworks))
	var flag bool

	for i, v := range originalNetworks {
		nic := v.(map[string]interface{})
		sourceDestChecks[i] = nic["source_dest_check"].(bool)
		if !flag && !sourceDestChecks[i] {
			flag = true
		}
	}

	if flag {
		// Get the instance network and address information
		server, err := cloudservers.Get(ecsClient, d.Id()).Extract()
		if err != nil {
			return diag.Errorf("error retrieving compute instance: %s", d.Id())
		}
		networks, err := flattenInstanceNetworks(d, meta, server)
		if err != nil {
			return diag.FromErr(err)
		}

		for i, nic := range networks {
			nicPort := nic["port"].(string)
			if nicPort == "" {
				continue
			}

			if !sourceDestChecks[i] {
				if err := disableSourceDestCheck(nicClient, nicPort); err != nil {
					return diag.Errorf("error disabling source dest check on port(%s) of instance(%s): %s", nicPort, d.Id(), err)
				}
			}
		}
	}

	return resourceComputeInstanceRead(ctx, d, meta)
}

func resourceComputeInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	ecsClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute V1 client: %s", err)
	}
	blockStorageClient, err := cfg.BlockStorageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating evs client: %s", err)
	}
	imsClient, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating image client: %s", err)
	}

	server, err := cloudservers.Get(ecsClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving compute instance")
	} else if server.Status == "DELETED" {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] retrieved compute instance %s: %+v", d.Id(), server)
	// Set some attributes
	d.Set("region", region)
	d.Set("enterprise_project_id", server.EnterpriseProjectID)
	d.Set("availability_zone", server.AvailabilityZone)
	d.Set("name", server.Name)
	d.Set("description", server.Description)
	d.Set("status", server.Status)
	d.Set("created_at", server.Created.Format(time.RFC3339))
	d.Set("updated_at", server.Updated.Format(time.RFC3339))

	flavorInfo := server.Flavor
	d.Set("flavor_id", flavorInfo.ID)
	d.Set("flavor_name", flavorInfo.Name)

	if server.Status == "ACTIVE" {
		d.Set("power_action", "ON")
	} else if server.Status == "SHUTOFF" {
		// The server instance is in the shutdown state. The local setting can be OFF or FORCE-OFF.
		if d.Get("power_action") != "OFF" && d.Get("power_action") != "FORCE-OFF" {
			d.Set("power_action", "OFF")
		}
	}

	// Set the instance's image information appropriately
	if err := setImageInformation(d, imsClient, server.Image.ID); err != nil {
		return diag.FromErr(err)
	}

	if server.KeyName != "" {
		d.Set("key_pair", server.KeyName)
	}
	if eip := computePublicIP(server); eip != "" {
		d.Set("public_ip", eip)
	}

	// Get the instance network and address information
	networks, err := flattenInstanceNetworks(d, meta, server)
	if err != nil {
		return diag.FromErr(err)
	}
	// Determine the best IPv4 and IPv6 addresses to access the instance with
	hostv4, hostv6 := getInstanceAccessAddresses(networks)

	// update hostv4/6 by AccessIPv4/v6
	// currently, AccessIPv4/v6 are Reserved in HuaweiCloudStack
	if server.AccessIPv4 != "" {
		hostv4 = server.AccessIPv4
	}
	if server.AccessIPv6 != "" {
		hostv6 = server.AccessIPv6
	}

	d.Set("network", networks)
	d.Set("access_ip_v4", hostv4)
	d.Set("access_ip_v6", hostv6)

	// Determine the best IP address to use for SSH connectivity.
	// Prefer IPv4 over IPv6.
	var preferredSSHAddress string
	if hostv4 != "" {
		preferredSSHAddress = hostv4
	} else if hostv6 != "" {
		preferredSSHAddress = hostv6
	}

	if preferredSSHAddress != "" {
		// Initialize the connection info
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": preferredSSHAddress,
		})
	}

	secGrpNames := []string{}
	for _, sg := range server.SecurityGroups {
		secGrpNames = append(secGrpNames, sg.Name)
	}
	d.Set("security_groups", secGrpNames)

	secGrpIDs := make([]string, len(server.SecurityGroups))
	for i, sg := range server.SecurityGroups {
		secGrpIDs[i] = sg.ID
	}
	d.Set("security_group_ids", secGrpIDs)

	// Set volume attached
	if len(server.VolumeAttached) > 0 {
		bds := make([]map[string]interface{}, len(server.VolumeAttached))
		for i, b := range server.VolumeAttached {
			// retrieve volume `size` and `type`
			volumeInfo, err := cloudvolumes.Get(blockStorageClient, b.ID).Extract()
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[DEBUG] retrieved volume %s: %#v", b.ID, volumeInfo)

			// retrieve volume `pci_address`
			va, err := block_devices.Get(ecsClient, d.Id(), b.ID).Extract()
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[DEBUG] retrieved block device %s: %#v", b.ID, va)

			bds[i] = map[string]interface{}{
				"volume_id":   b.ID,
				"size":        volumeInfo.Size,
				"type":        volumeInfo.VolumeType,
				"boot_index":  va.BootIndex,
				"pci_address": va.PciAddress,
				"kms_key_id":  volumeInfo.Metadata.SystemCmkID,
			}

			if va.BootIndex == 0 {
				d.Set("system_disk_id", b.ID)
				d.Set("system_disk_size", volumeInfo.Size)
			}
		}
		d.Set("volume_attached", bds)
	}

	// set scheduler_hints
	osHints := server.OsSchedulerHints
	if len(osHints.Group) > 0 {
		schedulerHints := make([]map[string]interface{}, len(osHints.Group))
		for i, v := range osHints.Group {
			schedulerHints[i] = map[string]interface{}{
				"group": v,
			}
		}
		d.Set("scheduler_hints", schedulerHints)
	}
	d.Set("tags", flattenTagsToMap(server.Tags))
	return nil
}

func resourceComputeInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	computeClient, err := cfg.ComputeV2Client(region)
	if err != nil {
		return diag.Errorf("error creating compute V2 client: %s", err)
	}
	ecsClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute V1 client: %s", err)
	}
	ecsV11Client, err := cfg.ComputeV11Client(region)
	if err != nil {
		return diag.Errorf("error creating compute V1.1 client: %s", err)
	}

	if d.HasChanges("name", "description") {
		var updateOpts cloudservers.UpdateOpts
		updateOpts.Name = d.Get("name").(string)
		description := d.Get("description").(string)
		updateOpts.Description = &description

		err := cloudservers.Update(ecsClient, d.Id(), updateOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error updating server: %s", err)
		}
	}

	if d.HasChanges("security_group_ids", "security_groups") {
		var oldSGRaw interface{}
		var newSGRaw interface{}
		if d.HasChange("security_group_ids") {
			oldSGRaw, newSGRaw = d.GetChange("security_group_ids")
		} else {
			oldSGRaw, newSGRaw = d.GetChange("security_groups")
		}
		oldSGSet := oldSGRaw.(*schema.Set)
		newSGSet := newSGRaw.(*schema.Set)
		secgroupsToAdd := newSGSet.Difference(oldSGSet)
		secgroupsToRemove := oldSGSet.Difference(newSGSet)
		log.Printf("[DEBUG] security groups to add: %v", secgroupsToAdd)
		log.Printf("[DEBUG] security groups to remove: %v", secgroupsToRemove)

		for _, g := range secgroupsToRemove.List() {
			err := secgroups.RemoveServer(computeClient, d.Id(), g.(string)).ExtractErr()
			if err != nil && err.Error() != "EOF" {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					continue
				}
				return diag.Errorf("error removing security group (%s) from server (%s): %s", g, d.Id(), err)
			}
			log.Printf("[DEBUG] removed security group (%s) from instance (%s)", g, d.Id())
		}

		for _, g := range secgroupsToAdd.List() {
			err := secgroups.AddServer(computeClient, d.Id(), g.(string)).ExtractErr()
			if err != nil && err.Error() != "EOF" {
				return diag.Errorf("error adding security group (%s) to server (%s): %s", g, d.Id(), err)
			}
			log.Printf("[DEBUG] added security group (%s) to instance (%s)", g, d.Id())
		}
	}

	if d.HasChange("admin_pass") {
		if newPwd, ok := d.Get("admin_pass").(string); ok {
			err := cloudservers.ChangeAdminPassword(ecsClient, d.Id(), newPwd).ExtractErr()
			if err != nil {
				return diag.Errorf("error changing admin password of server (%s): %s", d.Id(), err)
			}
		}
	}

	if d.HasChanges("flavor_id", "flavor_name") {
		newFlavorId, err := getFlavorID(d)
		if err != nil {
			return diag.FromErr(err)
		}

		extendParam := &cloudservers.ResizeExtendParam{
			AutoPay: common.GetAutoPay(d),
		}
		resizeOpts := &cloudservers.ResizeOpts{
			FlavorRef:   newFlavorId,
			Mode:        "withStopServer",
			ExtendParam: extendParam,
		}
		log.Printf("[DEBUG] resize configuration: %#v", resizeOpts)
		job, err := cloudservers.Resize(ecsV11Client, resizeOpts, d.Id()).ExtractJobResponse()
		if err != nil {
			return diag.Errorf("error resizing server: %s", err)
		}

		if err := cloudservers.WaitForJobSuccess(ecsClient, int(d.Timeout(schema.TimeoutUpdate)/time.Second), job.JobID); err != nil {
			return diag.Errorf("error waiting for instance (%s) to be resized: %s", d.Id(), err)
		}
	}

	if d.HasChange("network") {
		var err error
		nicClient, err := cfg.NetworkingV2Client(region)
		if err != nil {
			return diag.Errorf("error creating networking client: %s", err)
		}

		if err := updateSourceDestCheck(d, nicClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {

		tagErr := UpdateResourceTags(computeClient, d, "servers", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of instance:%s, err:%s", d.Id(), tagErr)
		}
	}

	if d.HasChange("enterprise_project_id") {
		epsClient, err := cfg.EnterpriseProjectClient(region)
		if err != nil {
			return diag.Errorf("error creating EPS client: %s", err)
		}

		if err := migrateEnterpriseProject(ctx, d, ecsClient, epsClient, region); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("system_disk_size") {
		extendOpts := cloudvolumes.ExtendOpts{
			SizeOpts: cloudvolumes.ExtendSizeOpts{
				NewSize: d.Get("system_disk_size").(int),
			},
		}

		evsV2Client, err := cfg.BlockStorageV2Client(region)
		if err != nil {
			return diag.Errorf("error creating evs V2 client: %s", err)
		}

		systemDiskID := d.Get("system_disk_id").(string)

		cloudvolumes.ExtendSize(evsV2Client, systemDiskID, extendOpts)

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"extending"},
			Target:     []string{"available", "in-use"},
			Refresh:    evs.VolumeV2StateRefreshFunc(evsV2Client, systemDiskID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"error waiting for hcs_compute_instance system disk %s to become ready: %s", systemDiskID, err)
		}
	}

	// update the key_pair before power action
	if d.HasChange("key_pair") {
		kmsClient, err := cfg.KmsV3Client(region)
		if err != nil {
			return diag.Errorf("error creating KMS v3 client: %s", err)
		}

		o, n := d.GetChange("key_pair")
		keyPairOpts := &common.KeypairAuthOpts{
			InstanceID:       d.Id(),
			InUsedKeyPair:    o.(string),
			NewKeyPair:       n.(string),
			InUsedPrivateKey: d.Get("private_key").(string),
			Password:         d.Get("admin_pass").(string),
			Timeout:          d.Timeout(schema.TimeoutUpdate),
		}
		if err := common.UpdateEcsInstanceKeyPair(ctx, ecsClient, kmsClient, keyPairOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	// The instance power status update needs to be done at the end
	if d.HasChange("power_action") {
		action := d.Get("power_action").(string)
		if err = doPowerAction(ecsClient, d, action); err != nil {
			return diag.Errorf("Doing power action (%s) for instance (%s) failed: %s", action, d.Id(), err)
		}
	}

	return resourceComputeInstanceRead(ctx, d, meta)
}

func resourceComputeInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	ecsClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}
	serverRequests := []cloudservers.Server{
		{Id: d.Id()},
	}

	deleteOpts := cloudservers.DeleteOpts{
		Servers:        serverRequests,
		DeleteVolume:   d.Get("delete_disks_on_termination").(bool),
		DeletePublicIP: d.Get("delete_eip_on_termination").(bool),
	}

	n, err := cloudservers.Delete(ecsClient, deleteOpts).ExtractJobResponse()
	if err != nil {
		return diag.Errorf("error deleting server: %s", err)
	}

	if err := cloudservers.WaitForJobSuccess(ecsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return diag.FromErr(err)
	}

	// Instance may still exist after Order/Job succeed.
	pending := []string{"ACTIVE", "SHUTOFF"}
	target := []string{"DELETED", "SOFT_DELETED"}
	deleteTimeout := d.Timeout(schema.TimeoutDelete)
	if err := waitForServerTargetState(ctx, ecsClient, d.Id(), pending, target, deleteTimeout); err != nil {
		return diag.Errorf("State waiting timeout: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceComputeInstanceImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	ecsClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating compute client: %s", err)
	}

	server, err := cloudservers.Get(ecsClient, d.Id()).Extract()
	if err != nil {
		return nil, common.CheckDeleted(d, err, "compute instance")
	}

	allInstanceNics, err := getInstanceAddresses(d, meta, server)
	if err != nil {
		return nil, fmt.Errorf("error fetching networks of compute instance %s: %s", d.Id(), err)
	}

	networks := []map[string]interface{}{}
	for _, nic := range allInstanceNics {
		v := map[string]interface{}{
			"uuid":              nic.NetworkID,
			"port":              nic.PortID,
			"fixed_ip_v4":       nic.FixedIPv4,
			"fixed_ip_v6":       nic.FixedIPv6,
			"ipv6_enable":       nic.FixedIPv6 != "",
			"source_dest_check": nic.SourceDestCheck,
			"mac":               nic.MAC,
		}
		networks = append(networks, v)
	}

	log.Printf("[DEBUG] flatten Instance Networks: %#v", networks)
	d.Set("network", networks)
	d.Set("tags", flattenTagsToMap(server.Tags))
	return []*schema.ResourceData{d}, nil
}

// ServerV1StateRefreshFunc returns a resource.StateRefreshFunc that is used to watch an HuaweiCloudStack instance.
func ServerV1StateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := cloudservers.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return s, "DELETED", nil
			}
			return nil, "", err
		}

		// get fault message when status is ERROR
		if s.Status == "ERROR" {
			fault := fmt.Errorf("error code: %d, message: %s", s.Fault.Code, s.Fault.Message)
			return s, "ERROR", fault
		}
		return s, s.Status, nil
	}
}

func resourceInstanceSecGroupIdsV1(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]cloudservers.SecurityGroup, error) {
	if v, ok := d.GetOk("security_group_ids"); ok {
		rawSecGroups := v.(*schema.Set).List()
		secGroups := make([]cloudservers.SecurityGroup, len(rawSecGroups))
		for i, raw := range rawSecGroups {
			secGroups[i] = cloudservers.SecurityGroup{
				ID: raw.(string),
			}
		}
		return secGroups, nil
	}

	rawSecGroups := d.Get("security_groups").(*schema.Set).List()
	secGroups := make([]cloudservers.SecurityGroup, 0, len(rawSecGroups))

	opt := groups.ListOpts{
		EnterpriseProjectId: "all_granted_eps",
	}
	pages, err := groups.List(client, opt).AllPages()
	if err != nil {
		return nil, err
	}
	resp, err := groups.ExtractSecurityGroups(pages)
	if err != nil {
		return nil, err
	}

	for _, raw := range rawSecGroups {
		secName := raw.(string)
		for _, secGroup := range resp {
			if secName == secGroup.Name {
				secGroups = append(secGroups, cloudservers.SecurityGroup{
					ID: secGroup.ID,
				})
				break
			}
		}
	}
	if len(secGroups) != len(rawSecGroups) {
		return nil, fmt.Errorf("the list contains invalid security groups (num: %d), please check your entry",
			len(rawSecGroups)-len(secGroups))
	}

	return secGroups, nil
}

func getOpSvcUserID(d *schema.ResourceData, conf *config.HcsConfig) string {
	if v, ok := d.GetOk("user_id"); ok {
		return v.(string)
	}
	return conf.UserID
}

func validateComputeInstanceConfig(d *schema.ResourceData, conf *config.HcsConfig) error {
	_, hasSSH := d.GetOk("key_pair")
	if hasSSH {
		if getOpSvcUserID(d, conf) == "" {
			return fmt.Errorf("user_id must be specified when the ECS is logged in using an SSH key")
		}
	}

	return nil
}

func buildInstanceNicsRequest(d *schema.ResourceData) []cloudservers.Nic {
	var nicRequests []cloudservers.Nic

	networks := d.Get("network").([]interface{})
	for _, v := range networks {
		network := v.(map[string]interface{})
		nicRequest := cloudservers.Nic{
			SubnetId:   network["uuid"].(string),
			IpAddress:  network["fixed_ip_v4"].(string),
			Ipv6Enable: network["ipv6_enable"].(bool),
		}

		nicRequests = append(nicRequests, nicRequest)
	}
	return nicRequests
}

func buildInstancePublicIPRequest(d *schema.ResourceData) *cloudservers.PublicIp {
	if v, ok := d.GetOk("eip_id"); ok {
		return &cloudservers.PublicIp{
			Id: v.(string),
		}
	}

	bandWidthRaw := d.Get("bandwidth").([]interface{})
	if len(bandWidthRaw) != 1 {
		return nil
	}

	bandWidth := bandWidthRaw[0].(map[string]interface{})
	bwOpts := cloudservers.BandWidth{
		ShareType:  bandWidth["share_type"].(string),
		Id:         bandWidth["id"].(string),
		ChargeMode: bandWidth["charge_mode"].(string),
		Size:       bandWidth["size"].(int),
	}

	return &cloudservers.PublicIp{
		Eip: &cloudservers.Eip{
			IpType:    d.Get("eip_type").(string),
			BandWidth: &bwOpts,
		},
		DeleteOnTermination: d.Get("delete_eip_on_termination").(bool),
	}
}

func resourceInstanceSchedulerHintsV1(schedulerHintsRaw map[string]interface{}) cloudservers.SchedulerHints {
	schedulerHints := cloudservers.SchedulerHints{
		Group:           schedulerHintsRaw["group"].(string),
		FaultDomain:     schedulerHintsRaw["fault_domain"].(string),
		Tenancy:         schedulerHintsRaw["tenancy"].(string),
		DedicatedHostID: schedulerHintsRaw["deh_id"].(string),
	}

	return schedulerHints
}

func getImage(client *golangsdk.ServiceClient, id, name string) (*cloudimages.Image, error) {
	listOpts := &cloudimages.ListOpts{
		ID:    id,
		Name:  name,
		Limit: 1,
	}
	allPages, err := cloudimages.List(client, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("unable to query images: %s", err)
	}

	allImages, err := cloudimages.ExtractImages(allPages)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve images: %s", err)
	}

	if len(allImages) < 1 {
		return nil, fmt.Errorf("unable to find images %s: Maybe not existed", id)
	}

	img := allImages[0]
	if id != "" && img.ID != id {
		return nil, fmt.Errorf("unexpected images ID")
	}
	if name != "" && img.Name != name {
		return nil, fmt.Errorf("unexpected images Name")
	}
	log.Printf("[DEBUG] retrieved image %s: %#v", id, img)
	return &img, nil
}

func getImageIDFromConfig(imsClient *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {
	if imageID := d.Get("image_id").(string); imageID != "" {
		return imageID, nil
	}

	if imageName := d.Get("image_name").(string); imageName != "" {
		img, err := getImage(imsClient, "", imageName)
		if err != nil {
			return "", err
		}
		return img.ID, nil
	}

	return "", fmt.Errorf("neither a boot device, image ID, or image name were able to be determined")
}

func setImageInformation(d *schema.ResourceData, imsClient *golangsdk.ServiceClient, imageID string) error {
	if imageID != "" {
		d.Set("image_id", imageID)
		image, err := getImage(imsClient, imageID, "")
		if err != nil {
			// If the image name can't be found, set the value to "Image not found".
			// The most likely scenario is that the image no longer exists in the Image Service
			// but the instance still has a record from when it existed.
			d.Set("image_name", "Image not found")
			return nil
		}
		d.Set("image_name", image.Name)
	}

	return nil
}

// computePublicIP get the first floating address
func computePublicIP(server *cloudservers.CloudServer) string {
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

func getFlavorID(d *schema.ResourceData) (string, error) {
	var flavorID string

	// both flavor_id and flavor_name are the same value
	if v1, ok := d.GetOk("flavor_id"); ok {
		flavorID = v1.(string)
	} else if v2, ok := d.GetOk("flavor_name"); ok {
		flavorID = v2.(string)
	}

	if flavorID == "" {
		return "", fmt.Errorf("missing required argument: the `flavor_id` must be specified")
	}
	return flavorID, nil
}

func getFlavor(client *golangsdk.ServiceClient, d *schema.ResourceData) (map[string]interface{}, error) {
	var flavorID string

	// both flavor_id and flavor_name are the same value
	if v1, ok := d.GetOk("flavor_id"); ok {
		flavorID = v1.(string)
	} else if v2, ok := d.GetOk("flavor_name"); ok {
		flavorID = v2.(string)
	}
	var resultFlavor map[string]interface{}
	if flavorID == "" {
		return resultFlavor, fmt.Errorf("missing required argument: the `flavor_id` must be specified")
	}
	listOpts := &flavors.ListOpts{
		AvailabilityZone: d.Get("availability_zone").(string),
	}

	pages, err := flavors.List(client, listOpts).AllPages()
	if err != nil {
		return resultFlavor, fmt.Errorf("flavor result is empty")
	}
	allFlavors, err := flavors.ExtractFlavors(pages)
	for _, flavor := range allFlavors {
		if flavorID != flavor.ID {
			continue
		}
		resultFlavor = flattenFlavor(&flavor)
	}
	return resultFlavor, nil
}

func getVpcID(client *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {
	var networkID string

	networks := d.Get("network").([]interface{})
	if len(networks) > 0 {
		// all networks belongs to one VPC
		network := networks[0].(map[string]interface{})
		networkID = network["uuid"].(string)
	}

	if networkID == "" {
		return "", fmt.Errorf("network ID should not be empty")
	}

	subnet, err := subnets.Get(client, networkID).Extract()
	if err != nil {
		return "", fmt.Errorf("error retrieving subnets: %s", err)
	}

	return subnet.VPC_ID, nil
}

func resourceComputeSchedulerHintsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["group"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["group"].(string)))
	}

	if m["tenancy"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["tenancy"].(string)))
	}

	if m["deh_id"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["deh_id"].(string)))
	}

	return hashcode.String(buf.String())
}

func waitForServerTargetState(ctx context.Context, client *golangsdk.ServiceClient, instanceID string, pending, target []string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      pending,
		Target:       target,
		Refresh:      ServerV1StateRefreshFunc(client, instanceID),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for instance (%s) to become target state (%v): %s", instanceID, target, err)
	}
	return nil
}

// doPowerAction is a method for instance power doing shutdown, startup and reboot actions.
func doPowerAction(client *golangsdk.ServiceClient, d *schema.ResourceData, action string) error {
	var jobResp *cloudservers.JobResponse
	powerOpts := powers.PowerOpts{
		Servers: []powers.ServerInfo{
			{ID: d.Id()},
		},
	}
	// In the reboot structure, Type is a required option.
	// Since the type of power off and reboot is 'SOFT' by default, setting this value has solved the power structural
	// compatibility problem between optional and required.
	if action != "ON" {
		powerOpts.Type = "SOFT"
	}
	if strings.HasPrefix(action, "FORCE-") {
		powerOpts.Type = "HARD"
		action = strings.TrimPrefix(action, "FORCE-")
	}
	op, ok := powerActionMap[action]
	if !ok {
		return fmt.Errorf("the powerMap does not contain option (%s)", action)
	}
	jobResp, err := powers.PowerAction(client, powerOpts, op).ExtractJobResponse()
	if err != nil {
		return fmt.Errorf("doing power action (%s) for instance (%s) failed: %s", action, d.Id(), err)
	}

	// The time of the power on/off and reboot is usually between 15 and 35 seconds.
	timeout := 3 * time.Minute
	if err := cloudservers.WaitForJobSuccess(client, int(timeout/time.Second), jobResp.JobID); err != nil {
		return fmt.Errorf("waiting power action (%s) for instance (%s) failed: %s", action, d.Id(), err)
	}
	return nil
}

func disableSourceDestCheck(networkClient *golangsdk.ServiceClient, portID string) error {
	// Update the allowed-address-pairs of the port to 1.1.1.1/0
	// to disable the source/destination check
	portpairs := []ports.AddressPair{
		{
			IPAddress: "1.1.1.1/0",
		},
	}
	portUpdateOpts := ports.UpdateOpts{
		AllowedAddressPairs: &portpairs,
	}

	_, err := ports.Update(networkClient, portID, portUpdateOpts).Extract()
	return err
}

func enableSourceDestCheck(networkClient *golangsdk.ServiceClient, portID string) error {
	// cancle all allowed-address-pairs to enable the source/destination check
	portpairs := make([]ports.AddressPair, 0)
	portUpdateOpts := ports.UpdateOpts{
		AllowedAddressPairs: &portpairs,
	}

	_, err := ports.Update(networkClient, portID, portUpdateOpts).Extract()
	return err
}

func updateSourceDestCheck(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var err error

	networks := d.Get("network").([]interface{})
	for i, v := range networks {
		nic := v.(map[string]interface{})
		nicPort := nic["port"].(string)
		if nicPort == "" {
			continue
		}

		if d.HasChange(fmt.Sprintf("network.%d.source_dest_check", i)) {
			sourceDestCheck := nic["source_dest_check"].(bool)
			if !sourceDestCheck {
				err = disableSourceDestCheck(client, nicPort)
			} else {
				err = enableSourceDestCheck(client, nicPort)
			}

			if err != nil {
				return fmt.Errorf("error updating source_dest_check on port(%s) of instance(%s) failed: %s", nicPort, d.Id(), err)
			}
		}
	}

	return nil
}

func shouldUnsubscribeEIP(d *schema.ResourceData) bool {
	deleteEIP := d.Get("delete_eip_on_termination").(bool)
	eipAddr := d.Get("public_ip").(string)
	eipType := d.Get("eip_type").(string)
	_, sharebw := d.GetOk("bandwidth.0.id")

	return deleteEIP && eipAddr != "" && eipType != "" && !sharebw
}

func resourceInstanceRootVolume(d *schema.ResourceData, bootType string) cloudservers.RootVolume {
	if bootType == "LocalDisk" {
		log.Printf("[INFO] extBootType is: %s, no need config root valume param.", bootType)
		return cloudservers.RootVolume{}
	}
	diskType := d.Get("system_disk_type").(string)
	if diskType == "" {
		diskType = "business_type_01"
	}
	volRequest := cloudservers.RootVolume{
		VolumeType: diskType,
		Size:       d.Get("system_disk_size").(int),
	}
	if d.Get("kms_key_id") != "" {
		encryptioninfo := cloudservers.VolumeEncryptInfo{
			CmkId:  d.Get("kms_key_id").(string),
			Cipher: d.Get("encrypt_cipher").(string),
		}
		volRequest.EncryptionInfo = &encryptioninfo
	}
	return volRequest
}

func resourceInstanceDataVolumes(d *schema.ResourceData) []cloudservers.DataVolume {
	var volRequests []cloudservers.DataVolume

	vols := d.Get("data_disks").([]interface{})
	for i := range vols {
		vol := vols[i].(map[string]interface{})
		volRequest := cloudservers.DataVolume{
			VolumeType: vol["type"].(string),
			Size:       vol["size"].(int),
		}
		if vol["snapshot_id"] != "" {
			extendparam := cloudservers.VolumeExtendParam{
				SnapshotId: vol["snapshot_id"].(string),
			}
			volRequest.Extendparam = &extendparam
		}

		if vol["kms_key_id"] != "" {
			encryptioninfo := cloudservers.VolumeEncryptInfo{
				CmkId:  vol["kms_key_id"].(string),
				Cipher: vol["encrypt_cipher"].(string),
			}
			volRequest.EncryptionInfo = &encryptioninfo
		}

		volRequests = append(volRequests, volRequest)
	}
	return volRequests
}

func migrateEnterpriseProject(ctx context.Context, d *schema.ResourceData,
	ecsClient, epsClient *golangsdk.ServiceClient, region string) error {
	resourceID := d.Id()
	targetEPSId := d.Get("enterprise_project_id").(string)

	if err := common.MigrateEnterpriseProject(epsClient, region, targetEPSId, "ecs", resourceID); err != nil {
		return err
	}

	// wait for the Enterprise Project ID changed
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Success"},
		Refresh:      waitForEnterpriseProjectIdChanged(ecsClient, resourceID, targetEPSId),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for migrating Enterprise Project ID: %s", err)
	}

	return nil
}

func waitForEnterpriseProjectIdChanged(client *golangsdk.ServiceClient, instanceID, epsID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := cloudservers.Get(client, instanceID).Extract()
		if err != nil {
			return nil, "ERROR", err
		}

		// get fault message when status is ERROR
		if s.Status == "ERROR" {
			fault := fmt.Errorf("error code: %d, message: %s", s.Fault.Code, s.Fault.Message)
			return s, "ERROR", fault
		}

		if s.EnterpriseProjectID == epsID {
			return s, "Success", nil
		}
		return s, "Pending", nil
	}
}

func checkTags(tagMap map[string]interface{}) bool {
	if len(tagMap) > 10 {
		return false
	}
	for k, v := range tagMap {
		if len(k) > 36 || len(v.(string)) > 43 {
			return false
		}
		keyReg, _ := regexp.Compile("^[a-zA-Z0-9\u4e00-\u9fa5_-]+$")
		if !keyReg.MatchString(k) {
			return false
		}
		valueReg, _ := regexp.Compile("^[a-zA-Z0-9\u4e00-\u9fa5._-]+$")
		if !valueReg.MatchString(v.(string)) {
			return false
		}
	}
	return true
}

func UpdateResourceTags(conn *golangsdk.ServiceClient, d *schema.ResourceData, resourceType, id string) error {
	oRaw, nRaw := d.GetChange("tags")
	oMap := oRaw.(map[string]interface{})
	nMap := nRaw.(map[string]interface{})
	if !checkTags(nMap) {
		return fmt.Errorf("tags check failed")
	}
	var oTags []string
	for k, v := range oMap {
		oTags = append(oTags, fmt.Sprintf("%s.%s", k, v))
	}
	var nTags []string
	for k, v := range nMap {
		nTags = append(nTags, fmt.Sprintf("%s.%s", k, v))
	}
	// remove old tags
	if len(oTags) > 0 {
		oTags = utils.RemoveDuplicateElem(oTags)
		err := utils.DeleteResourceTagsWithKeys(conn, oTags, resourceType, id)
		if err != nil {
			return err
		}
	}

	// set new tags
	if len(nTags) > 0 {
		nTags = utils.RemoveDuplicateElem(nTags)
		err := utils.CreateResourceTagsWithKeys(conn, nTags, resourceType, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func flattenTagsToMap(tags []string) map[string]string {
	result := make(map[string]string)
	for _, tagStr := range tags {
		tag := strings.SplitN(tagStr, ".", 2)
		if len(tag) == 2 {
			result[tag[0]] = tag[1]
		}
	}
	return result
}
