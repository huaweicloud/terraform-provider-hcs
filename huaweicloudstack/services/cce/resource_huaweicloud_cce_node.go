package cce

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/cce/v3/clusters"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/cce/v3/nodes"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/common/tags"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func ResourceNode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNodeCreate,
		ReadContext:   resourceNodeRead,
		UpdateContext: resourceNodeUpdate,
		DeleteContext: resourceNodeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNodeImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"os": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"key_pair": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"password", "key_pair"},
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ForceNew:     true,
				ExactlyOneOf: []string{"password", "key_pair"},
			},
			"private_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"root_volume": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"volumetype": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"hw_passthrough": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: "schema: Internal",
						},
						"extend_param": {
							Type:       schema.TypeString,
							Optional:   true,
							ForceNew:   true,
							Deprecated: "use extend_params instead",
						},
						"extend_params": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					}},
			},
			"data_volumes": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"volumetype": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"hw_passthrough": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: "schema: Internal",
						},
						"extend_param": {
							Type:       schema.TypeString,
							Optional:   true,
							ForceNew:   true,
							Deprecated: "use extend_params instead",
						},
						"extend_params": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					}},
			},
			"storage": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"selectors": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										Default:  "evs",
									},
									"match_label_size": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"match_label_volume_type": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"match_label_metadata_encrypted": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"match_label_metadata_cmkid": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"match_label_count": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"groups": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"cce_managed": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
									"selector_names": {
										Type:     schema.TypeList,
										Required: true,
										ForceNew: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"virtual_spaces": {
										Type:     schema.TypeList,
										Required: true,
										ForceNew: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
													ForceNew: true,
												},
												"size": {
													Type:     schema.TypeString,
													Required: true,
													ForceNew: true,
												},
												"lvm_lv_type": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"lvm_path": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"runtime_lv_type": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
											},
										},
									},
								},
							},
						},
					}},
			},
			"taints": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"effect": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"eip_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"iptype", "bandwidth_charge_mode", "bandwidth_size", "sharetype",
				},
			},
			"iptype": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				RequiredWith: []string{
					"iptype", "bandwidth_size", "sharetype",
				},
			},
			"bandwidth_charge_mode": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"eip_id"},
			},
			"sharetype": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				RequiredWith: []string{
					"iptype", "bandwidth_size", "sharetype",
				},
			},
			"bandwidth_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				RequiredWith: []string{
					"iptype", "bandwidth_size", "sharetype",
				},
			},
			"runtime": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"docker", "containerd",
				}, false),
			},
			"ecs_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ecs_performance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "schema: Internal",
			},
			"product_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "schema: Internal",
			},
			"max_pods": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"public_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "schema: Internal",
			},
			"preinstall": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				StateFunc: utils.DecodeHashAndHexEncode,
			},
			"postinstall": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				StateFunc: utils.DecodeHashAndHexEncode,
			},
			"labels": { //(k8s_tags)
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			//(node/ecs_tags)
			"tags": common.TagsForceNewSchema(),
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "schema: Internal",
			},

			"extend_param": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"fixed_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"keep_ecs": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "schema: Internal",
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildResourceNodeAnnotations(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("annotations").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func buildResourceNodeK8sTags(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("labels").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func buildResourceNodeTags(d *schema.ResourceData) []tags.ResourceTag {
	tagRaw := d.Get("tags").(map[string]interface{})
	return utils.ExpandResourceTags(tagRaw)
}

func buildResourceNodeRootVolume(d *schema.ResourceData) nodes.VolumeSpec {
	var root nodes.VolumeSpec
	volumeRaw := d.Get("root_volume").([]interface{})
	if len(volumeRaw) == 1 {
		rawMap := volumeRaw[0].(map[string]interface{})
		root.Size = rawMap["size"].(int)
		root.VolumeType = rawMap["volumetype"].(string)
		root.HwPassthrough = rawMap["hw_passthrough"].(bool)
		root.ExtendParam = rawMap["extend_params"].(map[string]interface{})

		if rawMap["kms_key_id"].(string) != "" {
			metadata := nodes.VolumeMetadata{
				SystemEncrypted: "1",
				SystemCmkid:     rawMap["kms_key_id"].(string),
			}
			root.Metadata = &metadata
		}
	}

	return root
}

func buildResourceNodeDataVolume(d *schema.ResourceData) []nodes.VolumeSpec {
	volumeRaw := d.Get("data_volumes").([]interface{})
	volumes := make([]nodes.VolumeSpec, len(volumeRaw))
	for i, raw := range volumeRaw {
		rawMap := raw.(map[string]interface{})
		volumes[i] = nodes.VolumeSpec{
			Size:          rawMap["size"].(int),
			VolumeType:    rawMap["volumetype"].(string),
			HwPassthrough: rawMap["hw_passthrough"].(bool),
			ExtendParam:   rawMap["extend_params"].(map[string]interface{}),
		}
		if rawMap["kms_key_id"].(string) != "" {
			metadata := nodes.VolumeMetadata{
				SystemEncrypted: "1",
				SystemCmkid:     rawMap["kms_key_id"].(string),
			}
			volumes[i].Metadata = &metadata
		}
	}
	return volumes
}

func buildResourceNodeTaint(d *schema.ResourceData) []nodes.TaintSpec {
	taintRaw := d.Get("taints").([]interface{})
	taints := make([]nodes.TaintSpec, len(taintRaw))
	for i, raw := range taintRaw {
		rawMap := raw.(map[string]interface{})
		taints[i] = nodes.TaintSpec{
			Key:    rawMap["key"].(string),
			Value:  rawMap["value"].(string),
			Effect: rawMap["effect"].(string),
		}
	}
	return taints
}

func buildResourceNodeEipIDs(d *schema.ResourceData) []string {
	if v, ok := d.GetOk("eip_id"); ok {
		return []string{v.(string)}
	}
	return nil
}

func buildResourceNodeExtendParam(d *schema.ResourceData) map[string]interface{} {
	extendParam := make(map[string]interface{})
	if v, ok := d.GetOk("extend_param"); ok {
		for key, val := range v.(map[string]interface{}) {
			extendParam[key] = val.(string)
		}
		if v, ok := extendParam["periodNum"]; ok {
			periodNum, err := strconv.Atoi(v.(string))
			if err != nil {
				log.Printf("[WARNING] PeriodNum %s invalid, Type conversion error: %s", v.(string), err)
			}
			extendParam["periodNum"] = periodNum
		}
	}

	if v, ok := d.GetOk("ecs_performance_type"); ok {
		extendParam["ecs:performancetype"] = v.(string)
	}
	if v, ok := d.GetOk("max_pods"); ok {
		extendParam["maxPods"] = v.(int)
	}
	if v, ok := d.GetOk("product_id"); ok {
		extendParam["productID"] = v.(string)
	}
	if v, ok := d.GetOk("public_key"); ok {
		extendParam["publicKey"] = v.(string)
	}
	if v, ok := d.GetOk("preinstall"); ok {
		extendParam["alpha.cce/preInstall"] = utils.TryBase64EncodeString(v.(string))
	}
	if v, ok := d.GetOk("postinstall"); ok {
		extendParam["alpha.cce/postInstall"] = utils.TryBase64EncodeString(v.(string))
	}

	return extendParam
}

func buildResourceNodeStorage(d *schema.ResourceData) *nodes.StorageSpec {
	v, ok := d.GetOk("storage")
	if !ok {
		return nil
	}

	var storageSpec nodes.StorageSpec
	storageSpecRaw := v.([]interface{})
	storageSpecRawMap := storageSpecRaw[0].(map[string]interface{})
	storageSelectorSpecRaw := storageSpecRawMap["selectors"].([]interface{})
	storageGroupSpecRaw := storageSpecRawMap["groups"].([]interface{})

	var selectors []nodes.StorageSelectorsSpec
	for _, s := range storageSelectorSpecRaw {
		sMap := s.(map[string]interface{})
		selector := nodes.StorageSelectorsSpec{
			Name:        sMap["name"].(string),
			StorageType: sMap["type"].(string),
			MatchLabels: nodes.MatchLabelsSpec{
				Size:              sMap["match_label_size"].(string),
				VolumeType:        sMap["match_label_volume_type"].(string),
				MetadataEncrypted: sMap["match_label_metadata_encrypted"].(string),
				MetadataCmkid:     sMap["match_label_metadata_cmkid"].(string),
				Count:             sMap["match_label_count"].(string),
			},
		}
		selectors = append(selectors, selector)
	}
	storageSpec.StorageSelectors = selectors

	var groups []nodes.StorageGroupsSpec
	for _, g := range storageGroupSpecRaw {
		gMap := g.(map[string]interface{})
		group := nodes.StorageGroupsSpec{
			Name:          gMap["name"].(string),
			CceManaged:    gMap["cce_managed"].(bool),
			SelectorNames: utils.ExpandToStringList(gMap["selector_names"].([]interface{})),
		}

		virtualSpacesRaw := gMap["virtual_spaces"].([]interface{})
		virtualSpaces := make([]nodes.VirtualSpacesSpec, 0, len(virtualSpacesRaw))
		for _, v := range virtualSpacesRaw {
			virtualSpaceMap := v.(map[string]interface{})
			virtualSpace := nodes.VirtualSpacesSpec{
				Name: virtualSpaceMap["name"].(string),
				Size: virtualSpaceMap["size"].(string),
			}

			if virtualSpaceMap["lvm_lv_type"].(string) != "" {
				lvmConfig := nodes.LVMConfigSpec{
					LvType: virtualSpaceMap["lvm_lv_type"].(string),
					Path:   virtualSpaceMap["lvm_path"].(string),
				}
				virtualSpace.LVMConfig = &lvmConfig
			}

			if virtualSpaceMap["runtime_lv_type"].(string) != "" {
				runtimeConfig := nodes.RuntimeConfigSpec{
					LvType: virtualSpaceMap["runtime_lv_type"].(string),
				}
				virtualSpace.RuntimeConfig = &runtimeConfig
			}

			virtualSpaces = append(virtualSpaces, virtualSpace)
		}
		group.VirtualSpaces = virtualSpaces

		groups = append(groups, group)
	}

	storageSpec.StorageGroups = groups
	return &storageSpec
}

func buildResourceNodePublicIP(d *schema.ResourceData) nodes.PublicIPSpec {
	// eipCount must be specified when bandwidth_size parameters was set
	eipCount := 0
	if _, ok := d.GetOk("bandwidth_size"); ok {
		eipCount = 1
	}

	res := nodes.PublicIPSpec{
		Ids:   buildResourceNodeEipIDs(d),
		Count: eipCount,
		Eip: nodes.EipSpec{
			IpType: d.Get("iptype").(string),
			Bandwidth: nodes.BandwidthOpts{
				ChargeMode: d.Get("bandwidth_charge_mode").(string),
				Size:       d.Get("bandwidth_size").(int),
				ShareType:  d.Get("sharetype").(string),
			},
		},
	}

	return res
}

func buildResourceNodeNicSpec(d *schema.ResourceData) nodes.NodeNicSpec {
	res := nodes.NodeNicSpec{
		PrimaryNic: nodes.PrimaryNic{
			SubnetId: d.Get("subnet_id").(string),
		},
	}

	if v, ok := d.GetOk("fixed_ip"); ok {
		res.PrimaryNic.FixedIps = []string{v.(string)}
	}

	return res
}

func buildResourceNodeLoginSpec(d *schema.ResourceData) (nodes.LoginSpec, error) {
	var loginSpec nodes.LoginSpec
	if v, ok := d.GetOk("key_pair"); ok {
		loginSpec = nodes.LoginSpec{
			SshKey: v.(string),
		}
	} else {
		password, err := utils.TryPasswordEncrypt(d.Get("password").(string))
		if err != nil {
			return loginSpec, err
		}
		loginSpec = nodes.LoginSpec{
			UserPassword: nodes.UserPassword{
				Username: "root",
				Password: password,
			},
		}
	}

	return loginSpec, nil
}

func resourceNodeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	nodeClient, err := cfg.CceV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE Node client: %s", err)
	}

	clusterId := d.Get("cluster_id").(string)
	createOpts := nodes.CreateOpts{
		Kind:       "Node",
		ApiVersion: "v3",
		Metadata: nodes.CreateMetaData{
			Name:        d.Get("name").(string),
			Annotations: buildResourceNodeAnnotations(d),
		},
		Spec: nodes.Spec{
			Flavor:      d.Get("flavor_id").(string),
			Az:          d.Get("availability_zone").(string),
			Os:          d.Get("os").(string),
			RootVolume:  buildResourceNodeRootVolume(d),
			DataVolumes: buildResourceNodeDataVolume(d),
			Storage:     buildResourceNodeStorage(d),
			PublicIP:    buildResourceNodePublicIP(d),
			Count:       1,
			NodeNicSpec: buildResourceNodeNicSpec(d),
			EcsGroupID:  d.Get("ecs_group_id").(string),
			ExtendParam: buildResourceNodeExtendParam(d),
			Taints:      buildResourceNodeTaint(d),
			K8sTags:     buildResourceNodeK8sTags(d),
			UserTags:    buildResourceNodeTags(d),
		},
	}

	if v, ok := d.GetOk("runtime"); ok {
		createOpts.Spec.RunTime = &nodes.RunTimeSpec{
			Name: v.(string),
		}
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	// Add loginSpec here so it wouldn't go in the above log entry
	loginSpec, err := buildResourceNodeLoginSpec(d)
	if err != nil {
		diag.FromErr(err)
	}
	createOpts.Spec.Login = loginSpec

	s, err := nodes.Create(nodeClient, clusterId, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating Node: %s", err)
	}

	// The completion of the creation of the underlying resource (ECS) corresponding to the CCE node does not mean that
	// the creation of the CCE node is completed.
	nodeID, err := getResourceIDFromJob(ctx, nodeClient, s.Status.JobID, "CreateNode", "CreateNodeVM",
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(nodeID)

	log.Printf("[DEBUG] Waiting for CCE Node (%s) to become available", s.Metadata.Name)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Build", "Installing"},
		Target:       []string{"Active"},
		Refresh:      waitForNodeActive(nodeClient, clusterId, nodeID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error creating CCE Node: %s", err)
	}

	return resourceNodeRead(ctx, d, meta)
}

func resourceNodeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	nodeClient, err := cfg.CceV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CCE Node client: %s", err)
	}
	clusterId := d.Get("cluster_id").(string)
	s, err := nodes.Get(nodeClient, clusterId, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CCE Node")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", s.Metadata.Name),
		d.Set("flavor_id", s.Spec.Flavor),
		d.Set("availability_zone", s.Spec.Az),
		d.Set("os", s.Spec.Os),
		d.Set("key_pair", s.Spec.Login.SshKey),
		d.Set("subnet_id", s.Spec.NodeNicSpec.PrimaryNic.SubnetId),
		d.Set("ecs_group_id", s.Spec.EcsGroupID),
		d.Set("server_id", s.Status.ServerID),
		d.Set("private_ip", s.Status.PrivateIP),
		d.Set("public_ip", s.Status.PublicIP),
		d.Set("status", s.Status.Phase),
		d.Set("root_volume", flattenResourceNodeRootVolume(s.Spec.RootVolume)),
		d.Set("data_volumes", flattenResourceNodeDataVolume(s.Spec.DataVolumes)),
	)

	if s.Spec.RunTime != nil {
		mErr = multierror.Append(mErr, d.Set("runtime", s.Spec.RunTime.Name))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting CCE Node fields: %s", err)
	}
	return nil
}

func flattenResourceNodeRootVolume(rootVolume nodes.VolumeSpec) []map[string]interface{} {
	res := []map[string]interface{}{
		{
			"size":           rootVolume.Size,
			"volumetype":     rootVolume.VolumeType,
			"hw_passthrough": rootVolume.HwPassthrough,
			"extend_params":  rootVolume.ExtendParam,
			"extend_param":   "",
		},
	}
	if rootVolume.Metadata != nil {
		res[0]["kms_key_id"] = rootVolume.Metadata.SystemCmkid
	}

	return res
}

func flattenResourceNodeDataVolume(dataVolumes []nodes.VolumeSpec) []map[string]interface{} {
	if len(dataVolumes) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(dataVolumes))
	for i, v := range dataVolumes {
		res[i] = map[string]interface{}{
			"size":           v.Size,
			"volumetype":     v.VolumeType,
			"hw_passthrough": v.HwPassthrough,
			"extend_params":  v.ExtendParam,
			"extend_param":   "",
		}

		if v.Metadata != nil {
			res[i]["kms_key_id"] = v.Metadata.SystemCmkid
		}
	}

	return res
}

func resourceNodeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	nodeClient, err := cfg.CceV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	var updateOpts nodes.UpdateOpts
	updateOpts.Metadata.Name = d.Get("name").(string)

	clusterId := d.Get("cluster_id").(string)
	_, err = nodes.Update(nodeClient, clusterId, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating cce node: %s", err)
	}

	return resourceNodeRead(ctx, d, meta)
}

func resourceNodeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	nodeClient, err := cfg.CceV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	clusterId := d.Get("cluster_id").(string)
	// remove node without deleting ecs
	if d.Get("keep_ecs").(bool) {
		loginSpec, err := buildResourceNodeLoginSpec(d)
		if err != nil {
			diag.FromErr(err)
		}

		removeOpts := nodes.RemoveOpts{
			Spec: nodes.RemoveNodeSpec{
				Login: loginSpec,
				Nodes: []nodes.NodeItem{
					{
						Uid: d.Id(),
					},
				},
			},
		}

		err = nodes.Remove(nodeClient, clusterId, removeOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error removing CCE node: %s", err)
		}
	} else {
		err = nodes.Delete(nodeClient, clusterId, d.Id()).ExtractErr()
		if err != nil {
			return diag.Errorf("error deleting CCE node: %s", err)
		}
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Deleting"},
		Target:       []string{"Deleted"},
		Refresh:      waitForNodeDelete(nodeClient, clusterId, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting CCE Node: %s", err)
	}

	d.SetId("")
	return nil
}

func getResourceIDFromJob(ctx context.Context, client *golangsdk.ServiceClient, jobID, jobType, subJobType string,
	timeout time.Duration) (string, error) {

	stateJob := &resource.StateChangeConf{
		Pending:      []string{"Initializing", "Running"},
		Target:       []string{"Success"},
		Refresh:      waitForJobStatus(client, jobID),
		Timeout:      timeout,
		Delay:        120 * time.Second,
		PollInterval: 20 * time.Second,
	}

	v, err := stateJob.WaitForStateContext(ctx)
	if err != nil {
		if job, ok := v.(*nodes.Job); ok {
			return "", fmt.Errorf("error waiting for job (%s) to become success: %s, reason: %s",
				jobID, err, job.Status.Reason)
		}

		return "", fmt.Errorf("error waiting for job (%s) to become success: %s", jobID, err)
	}

	job := v.(*nodes.Job)
	if len(job.Spec.SubJobs) == 0 {
		return "", fmt.Errorf("error fetching sub jobs from %s", jobID)
	}

	var subJobID string
	var refreshJob bool
	for _, s := range job.Spec.SubJobs {
		// postPaid: should get details of sub job ID
		if s.Spec.Type == jobType {
			subJobID = s.Metadata.ID
			refreshJob = true
			break
		}
	}

	if refreshJob {
		job, err = nodes.GetJobDetails(client, subJobID).ExtractJob()
		if err != nil {
			return "", fmt.Errorf("error fetching sub Job %s: %s", subJobID, err)
		}
	}

	var nodeId string
	for _, s := range job.Spec.SubJobs {
		if s.Spec.Type == subJobType {
			nodeId = s.Spec.ResourceID
			break
		}
	}
	if nodeId == "" {
		return "", fmt.Errorf("error fetching %s Job resource id", subJobType)
	}
	return nodeId, nil
}

func waitForNodeActive(cceClient *golangsdk.ServiceClient, clusterId, nodeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := nodes.Get(cceClient, clusterId, nodeId).Extract()
		if err != nil {
			return nil, "", err
		}

		return n, n.Status.Phase, nil
	}
}

func waitForNodeDelete(cceClient *golangsdk.ServiceClient, clusterId, nodeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete CCE Node %s", nodeId)

		r, err := nodes.Get(cceClient, clusterId, nodeId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted CCE Node %s", nodeId)
				return r, "Deleted", nil
			}
			return r, "Deleting", err
		}

		return r, r.Status.Phase, nil
	}
}

func waitForClusterAvailable(cceClient *golangsdk.ServiceClient, clusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[INFO] Waiting for CCE Cluster %s to be available", clusterId)
		n, err := clusters.Get(cceClient, clusterId).Extract()

		if err != nil {
			return nil, "", err
		}

		return n, n.Status.Phase, nil
	}
}

func waitForJobStatus(cceClient *golangsdk.ServiceClient, jobID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		job, err := nodes.GetJobDetails(cceClient, jobID).ExtractJob()
		if err != nil {
			return nil, "", err
		}

		return job, job.Status.Phase, nil
	}
}

func resourceNodeImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format specified for CCE Node. Format must be <cluster id>/<node id>")
		return nil, err
	}

	clusterID := parts[0]
	nodeID := parts[1]

	d.SetId(nodeID)
	err := d.Set("cluster_id", clusterID)
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}
