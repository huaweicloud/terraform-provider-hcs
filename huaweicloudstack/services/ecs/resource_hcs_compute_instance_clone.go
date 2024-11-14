package ecs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/ecs/v1/clone"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/ims"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"log"
	"time"
)

func ResourceComputeInstanceClone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeInstanceClone,
		ReadContext:   resourceComputeInstanceCloneRead,
		DeleteContext: resourceComputeInstanceCloneDelete,
		UpdateContext: resourceComputeInstanceCloneUpdate,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"power_on": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"retain_passwd": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"key_pair": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"admin_pass": {
				Type:          schema.TypeString,
				Sensitive:     true,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"key_pair"},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"network": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"vpc_id"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"fixed_ip_v4": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"ipv6_enable": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"fixed_ip_v6": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 5,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_group_id": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceComputeInstanceClone(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	ecsClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute V1 client: %s", err)
	}
	serverId := d.Get("instance_id").(string)
	server, err := clone.Get(ecsClient, serverId).Extract()
	if err != nil {
		return diag.Errorf("query compute instance by id %s failed, error: %s", serverId, err)
	}
	imageId := server.Image.ID
	imageV2Client, err := cfg.ImageV2Client(region)
	img, err := ims.GetCloudImage(imageV2Client, imageId)
	if err != nil {
		return diag.Errorf("error query image by id %s: %s", imageId, err)
	}
	osCloneOpts := buildOsCloneReqBody(d)

	if img.CloudInit == "True" {
		if img.OsType == "Windows" {
			if d.Get("key_pair") == "" && d.Get("admin_pass") == "" {
				return diag.Errorf("Cloud init has been installed in the image of the original server. The administrator password must be provided for the clone server")
			}
			if d.Get("key_pair") != "" {
				osCloneOpts.KeyPair = d.Get("key_pair").(string)
			} else if d.Get("admin_pass") != "" {
				osCloneOpts.AdminPass = d.Get("admin_pass").(string)
			}
		} else {
			if d.Get("key_pair") != "" {
				osCloneOpts.KeyPair = d.Get("key_pair").(string)
			} else if d.Get("admin_pass") != "" {
				osCloneOpts.AdminPass = d.Get("admin_pass").(string)
			} else {
				var retainPasswd = true
				osCloneOpts.RetainPasswd = &retainPasswd
			}
		}
	} else {
		var retainPasswd = true
		osCloneOpts.RetainPasswd = &retainPasswd
	}
	log.Printf("build os clone req success.")
	n, err := clone.CloneVm(ecsClient, serverId, osCloneOpts).ExtractJobResponse()
	if err != nil {
		return diag.Errorf("error cloning server: %s", err)
	}
	if err := clone.WaitForJobSuccess(ecsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return diag.Errorf("waiting clone for instance (%s) failed: %s", serverId, err)
	}
	name, ok := d.GetOk("name")
	if ok {
		d.SetId(serverId + "-" + name.(string))
	} else {
		d.SetId(serverId + "-" + utils.RandomString(8))
	}
	log.Printf("clone server success.")
	return nil
}

func buildOsCloneReqBody(d *schema.ResourceData) clone.OsCloneOpts {
	powerOn := d.Get("power_on").(bool)
	retainPasswd := d.Get("retain_passwd").(bool)
	osCloneOpts := clone.OsCloneOpts{
		PowerOn:      &powerOn,
		RetainPasswd: &retainPasswd,
	}
	if cloneType, ok := d.GetOk("clone_type"); ok {
		osCloneOpts.CloneType = cloneType.(string)
	}
	if name, ok := d.GetOk("name"); ok {
		osCloneOpts.Name = name.(string)
	}
	if postFix, ok := d.GetOk("postfix"); ok {
		osCloneOpts.Postfix = postFix.(string)
	}
	if vpcId, ok := d.GetOk("vpc_id"); ok {
		osCloneOpts.VpcId = vpcId.(string)
		osCloneOpts.Nics = buildNicsReqBody(d)
	}
	return osCloneOpts
}

func buildNicsReqBody(d *schema.ResourceData) []clone.NicsOpts {
	var nicsRequests []clone.NicsOpts
	nics := d.Get("network").([]interface{})
	for i := range nics {
		nic := nics[i].(map[string]interface{})
		var secGroupsReq []clone.SecurityGroupsOpts
		secGroups := nic["security_group_ids"].([]interface{})
		for j := range secGroups {
			secGroupId := secGroups[j].(map[string]interface{})["security_group_id"].(string)
			secGroupReq := clone.SecurityGroupsOpts{
				Id: secGroupId,
			}
			secGroupsReq = append(secGroupsReq, secGroupReq)
		}
		nicReq := clone.NicsOpts{
			SecurityGroups: secGroupsReq,
		}
		subnetId, ok := nic["subnet_id"]
		if ok {
			nicReq.SubnetId = subnetId.(string)
		}
		ipAddress, ok := nic["fixed_ip_v4"]
		if ok {
			nicReq.IpAddress = ipAddress.(string)
		}
		ipAddressV6, ok := nic["fixed_ip_v6"]
		if ok {
			nicReq.IpAddressV6 = ipAddressV6.(string)
		}
		ipv6Enable, ok := nic["ipv6_enable"]
		if ok {
			nicReq.Ipv6Enable = ipv6Enable.(bool)
		}
		nicsRequests = append(nicsRequests, nicReq)
	}
	return nicsRequests
}

func resourceComputeInstanceCloneRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name, ok := d.GetOk("name")
	if ok {
		d.Set("name", name)
	} else {
		d.Set("name", utils.RandomString(8))
	}
	return nil
}

func resourceComputeInstanceCloneDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceComputeInstanceCloneUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
