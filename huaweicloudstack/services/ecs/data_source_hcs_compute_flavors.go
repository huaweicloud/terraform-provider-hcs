package ecs

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/ecs/v1/flavors"
)

func DataSourceEcsFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEcsFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cpu_core_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"memory_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     FlavorsRefSchema(),
			},
		},
	}
}

func FlavorsRefSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vcpus": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_boot_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceEcsFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := config.GetHcsConfig(meta)
	region := conf.GetRegion(d)
	ecsClient, err := conf.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	listOpts := &flavors.ListOpts{
		AvailabilityZone: d.Get("availability_zone").(string),
	}

	pages, err := flavors.List(ecsClient, listOpts).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}

	allFlavors, err := flavors.ExtractFlavors(pages)
	if err != nil {
		return diag.Errorf("unable to retrieve flavors: %s ", err)
	}

	cpu := d.Get("cpu_core_count").(int)
	mem := int64(d.Get("memory_size").(int)) * 1024

	var ids []string
	var resultFlavors []interface{}
	for _, flavor := range allFlavors {
		vCpu, _ := strconv.Atoi(flavor.Vcpus)
		if cpu > 0 && vCpu != cpu {
			continue
		}

		if mem > 0 && flavor.Ram != mem {
			continue
		}

		ids = append(ids, flavor.ID)
		resultFlavors = append(resultFlavors, flattenFlavor(&flavor))
	}

	if len(ids) < 1 {
		return diag.Errorf("your query returned no results, please change your search criteria and try again.")
	}

	d.SetId(hashcode.Strings(ids))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("ids", ids),
		d.Set("flavors", resultFlavors),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenFlavor(flavor *flavors.Flavor) map[string]interface{} {
	res := map[string]interface{}{
		"id":            flavor.ID,
		"name":          flavor.Name,
		"ram":           flavor.Ram,
		"vcpus":         flavor.Vcpus,
		"ext_boot_type": flavor.OsExtraSpecs.ExtBootType,
	}
	return res
}
