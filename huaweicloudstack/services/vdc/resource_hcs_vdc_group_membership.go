package vdc

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/group_membership"
)

func ResourceVdcGroupMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVdcGroupMembershipCreate,
		ReadContext:   resourceVdcGroupMembershipRead,
		UpdateContext: resourceVdcGroupMembershipUpdate,
		DeleteContext: resourceVdcGroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVdcGroupMembershipInstanceImportState,
		},

		Schema: map[string]*schema.Schema{
			"group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"users": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceVdcGroupMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	vdcGroupClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating vdc user group membership network v3 client when create: %s", err)
	}

	groupId := d.Get("group").(string)
	userIds := d.Get("users").(*schema.Set).List()

	for _, id := range userIds {
		opts := group_membership.GroupMembershipReqParam{
			GroupID: groupId,
			UserID:  id.(string),
		}
		_, addErr := group_membership.AddGroupMembership(vdcGroupClient, opts).ExtractJobStatus()

		if addErr != nil {
			return diag.Errorf("error to add user to user group when created: %s", addErr)
		}
	}

	d.SetId(groupId)

	return resourceVdcGroupMembershipRead(ctx, d, meta)
}

func resourceVdcGroupMembershipUserList(_ context.Context, d *schema.ResourceData, meta interface{}) (map[string]bool, error) {
	hcsConfig := config.GetHcsConfig(meta)
	vdcGroupClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating vdc user group membership network v3 client: %s", err)
	}

	allUsers, err1 := group_membership.GetGroupMemberShipAllUser(vdcGroupClient, d.Id())
	if err1 != nil {
		return nil, fmt.Errorf("error to retrieve vdc user group users: %s", err)
	}

	userMap := make(map[string]bool)
	for _, user := range allUsers {
		userMap[user.ID] = true
	}

	return userMap, nil
}

func resourceVdcGroupMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	users, err := resourceVdcGroupMembershipUserList(ctx, d, meta)
	if err != nil {
		return diag.Errorf("error to retrieve vdc user group users in read: %s", err)
	}

	userIds := d.Get("users").(*schema.Set).List()
	var realUserIds []string
	for _, id := range userIds {
		if _, ok := users[id.(string)]; ok {
			realUserIds = append(realUserIds, id.(string))
		}
	}

	mErr := multierror.Append(nil,
		d.Set("group", d.Get("group").(string)),
		d.Set("users", realUserIds),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting vdc group users fields in read: %s", mErr)
	}
	return nil
}

func resourceVdcGroupMembershipUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	vdcGroupClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating vdc user group membership network v3 client when updated: %s", err)
	}

	groupId := d.Get("group").(string)
	if d.HasChange("users") {
		var (
			oldUserIds map[string]bool
			newUserIds map[string]bool
			allUsers   map[string]bool
			allUserErr error
		)
		allUsers, allUserErr = resourceVdcGroupMembershipUserList(ctx, d, meta)
		if allUserErr != nil {
			return diag.Errorf("error to retrieve vdc user group users when updated: %s", allUserErr)
		}

		oldVal, newVal := d.GetChange("users")
		if oldVal != nil {
			oldIds := oldVal.(*schema.Set).List()
			oldUserIds = make(map[string]bool)
			for _, raw := range oldIds {
				oldUserIds[raw.(string)] = true
			}
		}
		if newVal != nil {
			newIds := newVal.(*schema.Set).List()
			newUserIds = make(map[string]bool)
			for _, raw := range newIds {
				newUserIds[raw.(string)] = true
			}
		}
		// before delete
		for id, _ := range oldUserIds {
			// The new collection doesn't have
			if _, newOk := newUserIds[id]; !newOk {
				if _, ok := allUsers[id]; ok {
					opts := group_membership.GroupMembershipReqParam{
						GroupID: groupId,
						UserID:  id,
					}
					_, delErr := group_membership.DeleteGroupMembership(vdcGroupClient, opts).ExtractJobStatus()
					if delErr != nil {
						return diag.Errorf("error to delete user for user group when updated: %s", delErr)
					}
				}
			}
		}

		allUsers, allUserErr = resourceVdcGroupMembershipUserList(ctx, d, meta)
		if allUserErr != nil {
			return diag.Errorf("error to retrieve vdc user group users when updated: %s", allUserErr)
		}

		// after add
		for id, _ := range newUserIds {
			if _, ok := allUsers[id]; !ok {
				opts := group_membership.GroupMembershipReqParam{
					GroupID: groupId,
					UserID:  id,
				}
				_, addErr := group_membership.AddGroupMembership(vdcGroupClient, opts).ExtractJobStatus()
				if addErr != nil {
					return diag.Errorf("error to add user for user group when updated: %s", addErr)
				}
			}
		}
	}

	return resourceVdcGroupMembershipRead(ctx, d, meta)
}

func resourceVdcGroupMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	vdcGroupClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating vdc user group membership network v3 client when deleted: %s", err)
	}

	allUsers, err := resourceVdcGroupMembershipUserList(ctx, d, meta)
	if err != nil {
		return diag.Errorf("error to retrieve vdc user group users when deleted: %s", err)
	}

	groupId := d.Get("group").(string)
	userIds := d.Get("users").(*schema.Set).List()

	for _, id := range userIds {
		// has data, then delete
		if _, ok := allUsers[id.(string)]; ok {
			opts := group_membership.GroupMembershipReqParam{
				GroupID: groupId,
				UserID:  id.(string),
			}
			_, delErr := group_membership.DeleteGroupMembership(vdcGroupClient, opts).ExtractJobStatus()
			if delErr != nil {
				return diag.Errorf("error to delete user for user group when deleted: %s", delErr)
			}
		}
	}

	d.SetId("")

	return nil
}

func resourceVdcGroupMembershipInstanceImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	allUsers, err := resourceVdcGroupMembershipUserList(ctx, d, meta)
	if err != nil {
		return nil, fmt.Errorf("error to retrieve vdc user group users when imported: %s", err)
	}

	var userIds []string
	for userId, _ := range allUsers {
		userIds = append(userIds, userId)
	}

	mErr := multierror.Append(nil,
		d.Set("group", d.Id()),
		d.Set("users", userIds),
	)
	if mErr.ErrorOrNil() != nil {
		return nil, fmt.Errorf("error setting vdc group users fields when imported: %s", mErr)
	}

	return []*schema.ResourceData{d}, nil
}
