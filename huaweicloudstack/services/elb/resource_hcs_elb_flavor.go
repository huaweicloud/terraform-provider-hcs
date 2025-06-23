package elb

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/elb/v3/flavors"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
)

var flavorType = map[string]map[string]struct{}{
	"l4": {
		"cps": {}, "connection": {}, "bandwidth": {},
	},
	"l7": {
		"cps": {}, "connection": {}, "bandwidth": {}, "qps": {},
	},
}

func ResourceFlavorV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFlavorV3Create,
		ReadContext:   resourceFlavorV3Read,
		UpdateContext: resourceFlavorV3Update,
		DeleteContext: resourceFlavorV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v == "" {
						errs = append(errs, fmt.Errorf("%q must not be empty or null", key))
					}
					return
				},
			},
			"shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"l4", "l7",
				}, false),
			},
			"flavor_sold_out": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"info": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flavor_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{"bandwidth", "connection", "cps", "qps"},
								false,
							),
						},
						"value": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func buildFlavorInfo(infoRaw interface{}, qosName string, resourceType string) ([]flavors.Info, error) {
	// Support both TypeSet and TypeList
	var infos []interface{}
	switch v := infoRaw.(type) {
	case []interface{}:
		infos = v
	case *schema.Set:
		infos = v.List()
	default:
		return nil, fmt.Errorf("flavor name=%s info must be a list or set, got %T", qosName, infoRaw)
	}
	allowedTypes, ok := flavorType[resourceType]
	if !ok {
		return nil, fmt.Errorf("flavor name=%s invalid resource type: %s", qosName, resourceType)
	}
	var result []flavors.Info
	for _, v := range infos {
		mp, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("flavor name=%s info [%v] must be a map", qosName, v)
		}
		t, _ := mp["flavor_type"].(string)
		if _, exists := allowedTypes[t]; !exists {
			return nil, fmt.Errorf("flavor name=%s info flavor_type=%s not allowed for type=%s", qosName, t, resourceType)
		}
		valRaw, exists := mp["value"]
		if !exists {
			return nil, fmt.Errorf("flavor name=%s info missing value for flavor_type=%s", qosName, t)
		}
		var ivalue int
		switch val := valRaw.(type) {
		case int:
			ivalue = val
		case int64:
			ivalue = int(val)
		case float64: // Compatible with schema decode
			ivalue = int(val)
		default:
			return nil, fmt.Errorf("flavor name=%s value for flavor_type=%s must be int, got %T", qosName, t, valRaw)
		}
		if ivalue < 0 {
			return nil, fmt.Errorf("flavor name=%s value for flavor_type=%s must >= 0", qosName, t)
		}
		result = append(result, flavors.Info{
			FlavorType: t,
			Value:      ivalue,
		})
	}
	if len(result) == 0 {
		return []flavors.Info{}, nil
	}
	return result, nil
}

func resourceFlavorV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	resourceType := d.Get("type").(string)
	info, err := buildFlavorInfo(d.Get("info"), d.Get("name").(string), resourceType)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := flavors.CreateOpts{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
	}

	createOpts.Info = &info

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	flavor, err := flavors.Create(elbClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating flavor: %s", err)
	}

	d.SetId(flavor.ID)

	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForElbV3Flavor(ctx, elbClient, d.Id(), "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceFlavorV3Read(ctx, d, meta)
}

func resourceFlavorV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	flavor, err := flavors.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "flavor")
	}
	log.Printf("[DEBUG] Retrieved flavor %s: %#v", d.Id(), flavor)
	// Batch set base fields and collect errors
	mErr := multierror.Append(nil,
		d.Set("name", flavor.Name),
		d.Set("shared", flavor.Shared),
		d.Set("type", flavor.Type),
		d.Set("flavor_sold_out", flavor.SoldOut),
		d.Set("status", flavor.Status),
		d.Set("region", cfg.GetRegion(d)),
	)
	// Build info slice based on flavor type
	var info []map[string]interface{}
	switch flavor.Type {
	case "l4":
		if flavor.Info.Bandwidth != nil {
			info = append(info, map[string]interface{}{"flavor_type": "bandwidth", "value": *flavor.Info.Bandwidth})
		}
		if flavor.Info.Cps != nil {
			info = append(info, map[string]interface{}{"flavor_type": "cps", "value": *flavor.Info.Cps})
		}
		if flavor.Info.Connection != nil {
			info = append(info, map[string]interface{}{"flavor_type": "connection", "value": *flavor.Info.Connection})
		}
	case "l7":
		if flavor.Info.Bandwidth != nil {
			info = append(info, map[string]interface{}{"flavor_type": "bandwidth", "value": *flavor.Info.Bandwidth})
		}
		if flavor.Info.Cps != nil {
			info = append(info, map[string]interface{}{"flavor_type": "cps", "value": *flavor.Info.Cps})
		}
		if flavor.Info.Connection != nil {
			info = append(info, map[string]interface{}{"flavor_type": "connection", "value": *flavor.Info.Connection})
		}
		if flavor.Info.Qps != nil {
			info = append(info, map[string]interface{}{"flavor_type": "qps", "value": *flavor.Info.Qps})
		}
	default:
		info = nil
	}
	// Ensure the output order is stable
	sort.Slice(info, func(i, j int) bool {
		return info[i]["flavor_type"].(string) < info[j]["flavor_type"].(string)
	})
	mErr = multierror.Append(mErr, d.Set("info", info))
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Dedicated ELB flavor fields: %s", err)
	}
	return nil
}

func resourceFlavorV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	var updateOpts flavors.UpdateOpts
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("info") {
		resourceType := d.Get("type").(string)
		info, err := buildFlavorInfo(d.Get("info"), d.Get("name").(string), resourceType)
		if err != nil {
			return diag.FromErr(err)
		}
		updateOpts.Info = &info
	}

	log.Printf("[DEBUG] Updating flavor %s with options: %#v", d.Id(), updateOpts)
	_, err = flavors.Update(elbClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("unable to update flavor %s: %s", d.Id(), err)
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForElbV3Flavor(ctx, elbClient, d.Id(), "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceFlavorV3Read(ctx, d, meta)
}

func resourceFlavorV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	log.Printf("[DEBUG] Attempting to delete flavor %s", d.Id())
	err = flavors.Delete(elbClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("unable to delete flavor %s: %s", d.Id(), err)
	}

	// Wait for flavor to delete
	timeout := d.Timeout(schema.TimeoutDelete)
	err = waitForElbV3Flavor(ctx, elbClient, d.Id(), "DELETED", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForElbV3Flavor(ctx context.Context, elbClient *golangsdk.ServiceClient, id string, target string, pending []string,
	timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for flavor %s to become %s.", id, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceElbV3FlavorRefreshFunc(elbClient, id),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("error: flavor %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("error waiting for flavor %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceElbV3FlavorRefreshFunc(elbClient *golangsdk.ServiceClient, flavorID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		flavor, err := flavors.Get(elbClient, flavorID).Extract()
		if err != nil {
			return nil, "", err
		}

		// The flavor resource has no Status attribute, so a successful Get is the best we can do
		return flavor, "ACTIVE", nil
	}
}
