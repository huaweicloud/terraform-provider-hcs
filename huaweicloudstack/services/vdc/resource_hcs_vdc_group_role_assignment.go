package vdc

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/group_role_assignment"
)

// @API EPS GET v3.2/groups/{group_id}/roles
// @API EPS PUT v3.2/groups/{group_id}/roles
// @API EPS PUT v1.0/enterprise-projects/{enterprise_project_id}/groups/{group_id}/role/{role_id}
// @API EPS DELETE v1.0/enterprise-projects/{enterprise_project_id}/groups/{group_id}/role/{role_id}
func ResourceVdcGroupRoleAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVdcGroupRoleAssignmentCreate,
		ReadContext:   resourceVdcGroupRoleAssignmentRead,
		DeleteContext: resourceVdcGroupRoleAssignmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVdcGroupRoleAssignmentImportState,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_assignment": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"project_id": {
							Type:     schema.TypeString, // all indicates that all resource spaces take effect, including those to be created in the future.
							Optional: true,
							ForceNew: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceGetRole(d *schema.ResourceData) ([]group_role_assignment.GroupRoleAssignmentRole, error) {
	var err error = nil
	roles := d.Get("role_assignment").(*schema.Set).List()
	var results []group_role_assignment.GroupRoleAssignmentRole
	for _, role := range roles {
		roleOpts := group_role_assignment.GroupRoleAssignmentRole{}
		c := role.(map[string]interface{})
		if val, ok := c["role_id"].(string); ok && val != "" {
			roleOpts.Id = val
		}

		roleOpts.Inherit = false

		targetCount := 0
		var domainId string
		if val, ok := c["domain_id"].(string); ok && val != "" {
			targetCount++
			roleOpts.TargetType = "domain"
			roleOpts.Targets = []map[string]string{{"id": val}}
			domainId = val
		}

		if val, ok := c["project_id"].(string); ok && val != "" {
			targetCount++
			roleOpts.TargetType = "project"
			roleOpts.Targets = []map[string]string{{"id": val}}
			if val == "all" {
				roleOpts.Inherit = true
				if len(domainId) < 1 {
					return nil, errors.New("when project_id is set to all, domain_id cannot be empty")
				}
				targetCount = 1
				roleOpts.TargetType = "domain"
				roleOpts.Targets = []map[string]string{{"id": domainId}}
			}
		}

		if val, ok := c["enterprise_project_id"].(string); ok && val != "" {
			targetCount++
			roleOpts.TargetType = "enterprise_project"
			roleOpts.Targets = []map[string]string{{"id": val}}
		}

		if targetCount != 1 {
			return nil, errors.New("when project is not set to all, only one of domain_id, project_id, and enterprise_project_id can be set")
		}

		results = append(results, roleOpts)
	}
	return results, err
}

func resourceVdcGroupRoleAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	vdcGroupClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating vdc group role assignment network v3 client: %s", err)
	}

	groupId := d.Get("group_id").(string)
	roles, roleErr := resourceGetRole(d)
	if roleErr != nil {
		return diag.Errorf("error to parse role assignment for vdc group when create: %s", roleErr)
	}

	for _, role := range roles {

		// 如果是企业项目
		if role.TargetType == "enterprise_project" {
			var targetId string
			target := role.Targets[0]
			if val, ok := target["id"]; ok {
				targetId = val
			}
			_, addErr := group_role_assignment.AddGroupRoleAssignmentForEps(vdcGroupClient, groupId, targetId, role.Id).ExtractJobStatus()

			if addErr != nil {
				return diag.Errorf("error to add role assignment for vdc group when create: %s", addErr)
			}

		} else { // 非企业项目
			roleOpts := group_role_assignment.GroupRoleAssignmentRole{
				Id:         role.Id,
				Inherit:    role.Inherit,
				TargetType: role.TargetType,
				Targets:    role.Targets,
			}
			opts := group_role_assignment.GroupRoleAssignmentReqParam{
				Action: "add",
				Roles:  []group_role_assignment.GroupRoleAssignmentRole{roleOpts},
			}
			_, addErr := group_role_assignment.AddOrDeleteGroupRoleAssignment(vdcGroupClient, groupId, opts).ExtractJobStatus()

			if addErr != nil {
				return diag.Errorf("error to add role assignment for vdc group when create: %s", addErr)
			}
		}
	}
	d.SetId(groupId)

	return resourceVdcGroupRoleAssignmentRead(ctx, d, meta)
}

func resourceVdcGroupRoleAssignmentList(_ context.Context, d *schema.ResourceData, meta interface{}) ([]group_role_assignment.Role, error) {
	hcsConfig := config.GetHcsConfig(meta)
	vdcClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))
	if err != nil {
		return nil, err
	}

	return group_role_assignment.GetVdcGroupRoleAssignmentAllRoles(vdcClient, d.Id())
}

func resourceVdcGroupRoleAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	allRoles, err := resourceVdcGroupRoleAssignmentList(ctx, d, meta)
	if err != nil {
		return diag.Errorf("error to retrieve vdc group role assignment when read: %s", err)
	}

	roles, roleErr := resourceGetRole(d)
	if roleErr != nil {
		return diag.Errorf("error to parse role assignment for vdc group when read: %s", roleErr)
	}

	result := make([]map[string]interface{}, 0)
	for _, role := range roles {
		hasData := false
		targetId, ok := role.Targets[0]["id"]
		tmpDomainId := ""
		tmpProjectId := ""
		tmpEnterpriseProjectId := ""
		for _, aRole := range allRoles {
			tmpDomainId = ""
			tmpProjectId = ""
			tmpEnterpriseProjectId = ""
			tmpTargetType := ""
			if aRole.AssignType == "GroupResource" {
				tmpTargetType = "enterprise_project"
				tmpEnterpriseProjectId = aRole.TargetID
			} else if aRole.AssignType == "GroupProject" {
				tmpTargetType = "project"
				tmpProjectId = aRole.TargetID
			} else {
				tmpTargetType = "domain"
				tmpDomainId = aRole.TargetID
				if aRole.Inherit {
					tmpProjectId = "all"
				}
			}

			if ok && aRole.ID == role.Id && aRole.Inherit == role.Inherit && aRole.TargetID == targetId && tmpTargetType == role.TargetType {
				hasData = true
			}

			if hasData {
				break
			}
		}

		if hasData {
			result = append(result, map[string]interface{}{
				"role_id":               role.Id,
				"domain_id":             tmpDomainId,
				"project_id":            tmpProjectId,
				"enterprise_project_id": tmpEnterpriseProjectId,
			})
		}
	}

	mErr := multierror.Append(nil,
		d.Set("group_id", d.Get("group_id").(string)),
		d.Set("role_assignment", result),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting vdc group role assignment fields: %s", mErr)
	}
	return nil
}

func resourceVdcGroupRoleAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	vdcGroupClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating vdc group role assignment network v3 client: %s", err)
	}

	groupId := d.Get("group_id").(string)
	roles, roleErr := resourceGetRole(d)
	if roleErr != nil {
		return diag.Errorf("error to parse role assignment for vdc group when delete: %s", roleErr)
	}

	for _, role := range roles {

		// 如果是企业项目
		if role.TargetType == "enterprise_project" {
			var targetId string
			target := role.Targets[0]
			if val, ok := target["id"]; ok {
				targetId = val
			}
			_, err1 := group_role_assignment.DeleteGroupRoleAssignmentForEps(vdcGroupClient, groupId, targetId, role.Id).ExtractJobStatus()

			if err1 != nil {
				return diag.Errorf("error to remove role assignment for vdc group when delete: %s", err1)
			}
		} else { // 非企业项目
			roleOpts := group_role_assignment.GroupRoleAssignmentRole{
				Id:         role.Id,
				Inherit:    role.Inherit,
				TargetType: role.TargetType,
				Targets:    role.Targets,
			}
			opts := group_role_assignment.GroupRoleAssignmentReqParam{
				Action: "delete",
				Roles:  []group_role_assignment.GroupRoleAssignmentRole{roleOpts},
			}
			_, addErr := group_role_assignment.AddOrDeleteGroupRoleAssignment(vdcGroupClient, groupId, opts).ExtractJobStatus()

			if addErr != nil {
				return diag.Errorf("error to remove role assignment for vdc group when delete: %s", addErr)
			}
		}
	}

	d.SetId("")
	return nil
}

func resourceVdcGroupRoleAssignmentImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	allRoles, err := resourceVdcGroupRoleAssignmentList(ctx, d, meta)
	if err != nil {
		return nil, fmt.Errorf("error to retrieve vdc group role assignment when import: %s", err)
	}

	result := make([]map[string]interface{}, 0)
	for _, role := range allRoles {
		tmpDomainId := ""
		tmpProjectId := ""
		tmpEnterpriseProjectId := ""
		if role.AssignType == "GroupResource" {
			tmpEnterpriseProjectId = role.TargetID
		} else if role.AssignType == "GroupProject" {
			tmpProjectId = role.TargetID
		} else {
			tmpDomainId = role.TargetID
			if role.Inherit {
				tmpProjectId = "all"
			}
		}

		result = append(result, map[string]interface{}{
			"role_id":               role.ID,
			"domain_id":             tmpDomainId,
			"project_id":            tmpProjectId,
			"enterprise_project_id": tmpEnterpriseProjectId,
		})
	}

	mErr := multierror.Append(nil,
		d.Set("group_id", d.Id()),
		d.Set("role_assignment", result),
	)

	if mErr.ErrorOrNil() != nil {
		return nil, fmt.Errorf("error setting vdc group role assignment fields when import: %s", mErr)
	}

	return []*schema.ResourceData{d}, nil
}
