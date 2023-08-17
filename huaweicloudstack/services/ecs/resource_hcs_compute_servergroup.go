package ecs

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/ecs/v1/cloudservers"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/ecs/v1/servergroups"
)

func ResourceComputeServerGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeServerGroupCreate,
		ReadContext:   resourceComputeServerGroupRead,
		UpdateContext: resourceComputeServerGroupUpdate,
		DeleteContext: resourceComputeServerGroupDelete,
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
			"policies": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "schema: Required",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"members": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func buildServerGroupPolicies(d *schema.ResourceData) []string {
	rawPolicies := d.Get("policies").([]interface{})
	policies := make([]string, len(rawPolicies))
	for i, raw := range rawPolicies {
		policies[i] = raw.(string)
	}
	return policies
}

func resourceComputeServerGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	ecsClient, err := cfg.ComputeV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	createOpts := servergroups.CreateOpts{
		Name:     d.Get("name").(string),
		Policies: buildServerGroupPolicies(d),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	newSG, err := servergroups.Create(ecsClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating ECS server group: %s", err)
	}

	d.SetId(newSG.ID)

	membersToAdd := d.Get("members").(*schema.Set)
	for _, v := range membersToAdd.List() {
		instanceId := v.(string)
		// The ECS instances do not support other operations when binding server groups.
		config.MutexKV.Lock(instanceId)

		var addMemberOpts servergroups.MemberOpts
		addMemberOpts.InstanceID = instanceId
		err := servergroups.UpdateMember(ecsClient, addMemberOpts, "add_member", d.Id()).ExtractErr()
		// Release the ECS instance after the binding operation is complete whether it success or not.
		config.MutexKV.Unlock(instanceId)
		if err != nil {
			return diag.Errorf("error binding instance %s to ECS server group: %s", instanceId, err)
		}
	}

	return resourceComputeServerGroupRead(ctx, d, meta)
}

func resourceComputeServerGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	ecsClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	sg, err := servergroups.Get(ecsClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "server group")
	}

	log.Printf("[DEBUG] Retrieved server group %s: %+v", d.Id(), sg)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", sg.Name),
		d.Set("members", sg.Members),
		d.Set("policies", sg.Policies),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting server group fields: %s", err)
	}
	return nil
}

func resourceComputeServerGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	ecsClient, err := cfg.ComputeV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	if d.HasChange("members") {
		oldMembers, newMembers := d.GetChange("members")
		oldMemberSet, newMemberSet := oldMembers.(*schema.Set), newMembers.(*schema.Set)
		membersToAdd := newMemberSet.Difference(oldMemberSet)
		membersToRemove := oldMemberSet.Difference(newMemberSet)

		for _, v := range membersToAdd.List() {
			var addMemberOpts servergroups.MemberOpts
			instanceId := v.(string)
			// The ECS instances do not support other operations when binding server groups.
			config.MutexKV.Lock(instanceId)
			addMemberOpts.InstanceID = instanceId
			err = servergroups.UpdateMember(ecsClient, addMemberOpts, "add_member", d.Id()).ExtractErr()
			// Release the ECS instance ID after the binding operation is complete whether it success or not.
			config.MutexKV.Unlock(instanceId)
			if err != nil {
				return diag.Errorf("error binding instance %s to server group: %s", instanceId, err)
			}
		}

		for _, v := range membersToRemove.List() {
			instanceId := v.(string)
			server, err := cloudservers.Get(ecsClient, instanceId).Extract()
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					log.Printf("[WARN] the compute %s is not exist, ignore to remove it from the group", instanceId)
					continue
				}
				log.Printf("[WARN] failed to retrieve compute %s: %s, try to remove it from the group", instanceId, err)
			} else if server.Status == "DELETED" {
				log.Printf("[WARN] the compute %s was removed, ignore to remove it from the group", instanceId)
				continue
			}

			var removeMemberOpts servergroups.MemberOpts
			// Any operations are not supported when an ECS instance is unbound from a server group.
			config.MutexKV.Lock(instanceId)
			removeMemberOpts.InstanceID = instanceId
			err = servergroups.UpdateMember(ecsClient, removeMemberOpts, "remove_member", d.Id()).ExtractErr()
			// Release the ECS instance ID after the unbinding operation is complete whether it success or not.
			config.MutexKV.Unlock(instanceId)
			if err != nil {
				return diag.Errorf("error unbinding instance %s from ECS server group: %s", instanceId, err)
			}
		}
	}

	return resourceComputeServerGroupRead(ctx, d, meta)
}

func LockAll(ids []interface{}) {
	for _, instanceId := range ids {
		config.MutexKV.Lock(instanceId.(string))
	}
}

func UnlockAll(ids []interface{}) {
	for _, instanceId := range ids {
		config.MutexKV.Unlock(instanceId.(string))
	}
}

func resourceComputeServerGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	ecsClient, err := cfg.ComputeV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	members := d.Get("members").(*schema.Set).List()
	// Make sure that no other operations on the ECS instance are performed during the unbinding process.
	LockAll(members)

	log.Printf("[DEBUG] Deleting server group %s", d.Id())
	err = servergroups.Delete(ecsClient, d.Id()).ExtractErr()
	UnlockAll(members)
	if err != nil {
		return diag.Errorf("error deleting server group: %s", err)
	}

	return nil
}
