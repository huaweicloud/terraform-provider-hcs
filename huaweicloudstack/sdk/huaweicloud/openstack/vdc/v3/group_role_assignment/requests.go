package group_role_assignment

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"strings"
)

type GroupRoleAssignmentReqParam struct {
	Action string                    `json:"action"`
	Roles  []GroupRoleAssignmentRole `json:"roles"`
}

type GroupRoleAssignmentRole struct {
	Id         string              `json:"id"`
	Inherit    bool                `json:"inherit"`
	TargetType string              `json:"target_type"`
	Targets    []map[string]string `json:"targets"`
}

type GroupRoleAssignmentListReqParam struct {
	// vdc group id
	GroupID string

	// start index
	Start int `q:"start"`

	// pre page num
	Limit int `q:"limit"`
}

func (opts GroupRoleAssignmentListReqParam) hasQueryParameter() bool {

	return opts.Start != 0 || opts.Limit != 0
}

func (opts GroupRoleAssignmentListReqParam) ToListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

func GetGroupRoleAssignmentList(c *golangsdk.ServiceClient, opts GroupRoleAssignmentListReqParam) (r GroupRoleAssignmentListResult) {
	queryStr := ""
	if opts.hasQueryParameter() {
		query, err := opts.ToListQuery()
		if err != nil {
			r.Err = err
			return
		}
		queryStr = query
	}
	if len(queryStr) > 0 && !strings.HasPrefix(queryStr, "?") {
		queryStr = "?" + queryStr
	}
	url := GroupRoleAssignmentURL(c, opts.GroupID) + queryStr
	reqOpt := &golangsdk.RequestOpts{MoreHeaders: golangsdk.MoreHeaders}
	_, r.Err = c.Get(url, &r.Body, reqOpt)

	return
}

func GetVdcGroupRoleAssignmentAllRoles(c *golangsdk.ServiceClient, groupId string) ([]Role, error) {

	var err error
	var allRoles []Role
	start := 0
	for {
		opts := GroupRoleAssignmentListReqParam{
			GroupID: groupId,
			Start:   start,
			Limit:   1000,
		}

		roles, total, err1 := GetGroupRoleAssignmentList(c, opts).Extract()
		if err1 != nil {
			err = err1
			break
		}

		allRoles = append(allRoles, roles...)

		// 是否有下一页数据
		if start+opts.Limit < total {
			start = start + opts.Limit
		} else {
			break
		}
	}

	return allRoles, err
}

func AddOrDeleteGroupRoleAssignment(c *golangsdk.ServiceClient, groupId string, opts GroupRoleAssignmentReqParam) (r CreateGroupRoleAssignmentResult) {

	reqBody, err := golangsdk.BuildRequestBody(opts, "group")
	if err != nil {
		r.Err = err
		return
	}

	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{204}, MoreHeaders: golangsdk.MoreHeaders}
	_, r.Err = c.Put(GroupRoleAssignmentURL(c, groupId), reqBody, &r.Body, reqOpt)

	return
}

func AddGroupRoleAssignmentForEps(c *golangsdk.ServiceClient, groupId string, epId string, roleId string) (r CreateGroupRoleAssignmentResult) {

	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{204}, MoreHeaders: golangsdk.MoreHeaders}
	_, r.Err = c.Put(GroupRoleAssignmentEpsURL(c, epId, groupId, roleId), map[string]interface{}{}, &r.Body, reqOpt)

	return
}

func DeleteGroupRoleAssignmentForEps(c *golangsdk.ServiceClient, groupId string, epId string, roleId string) (r DeleteGroupRoleAssignmentResult) {

	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{204}, MoreHeaders: golangsdk.MoreHeaders}
	_, r.Err = c.Delete(GroupRoleAssignmentEpsURL(c, epId, groupId, roleId), reqOpt)

	return
}
