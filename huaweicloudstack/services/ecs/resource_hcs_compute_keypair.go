package ecs

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/compute/v2/extensions/keypairs"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/logp"
)

func ResourceComputeKeypairV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeKeypairV2Create,
		ReadContext:   resourceComputeKeypairV2Read,
		DeleteContext: resourceComputeKeypairV2Delete,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"key_file"},
			},
			"key_file": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeKeypairV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	computeClient, err := cfg.ComputeV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	pk, isExist := d.GetOk("public_key")
	createOpts := keypairs.CreateOpts{
		Name:      d.Get("name").(string),
		PublicKey: pk.(string),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	kp, err := keypairs.Create(computeClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating HuaweiCloud keypair: %s", err)
	}

	d.SetId(kp.Name)

	if !isExist {
		fp := getKeyFilePath(d)
		if err = utils.WriteToPemFile(fp, kp.PrivateKey); err != nil {
			return diag.Errorf("Unable to generate private key: %s", err)
		}
		d.Set("key_file", fp)
	}

	return resourceComputeKeypairV2Read(ctx, d, meta)
}

func getKeyFilePath(d *schema.ResourceData) string {
	if path, ok := d.GetOk("key_file"); ok {
		return path.(string)
	}
	keypairName := d.Get("name").(string)
	return fmt.Sprintf("%s.pem", keypairName)
}

func resourceComputeKeypairV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	computeClient, err := cfg.ComputeV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	kp, err := keypairs.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "keypair")
	}

	d.Set("name", kp.Name)
	d.Set("public_key", kp.PublicKey)
	d.Set("region", cfg.GetRegion(d))

	return nil
}

func resourceComputeKeypairV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	computeClient, err := cfg.ComputeV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	err = keypairs.Delete(computeClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("Error deleting HuaweiCloud keypair: %s", err)
	}
	d.SetId("")
	return nil
}
