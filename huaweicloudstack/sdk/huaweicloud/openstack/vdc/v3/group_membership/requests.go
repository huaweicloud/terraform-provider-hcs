package group_membership

import (
	"fmt"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	url2 "net/url"
	"strings"
)

type GroupMembershipReqParam struct {
	// user group id
	GroupID string

	// user id
	UserID string
}

type GroupMembershipListReqParam struct {
	// user group id
	GroupID string

	// Start Index
	Start int `q:"start"`

	// Data per page
	Limit int `q:"limit"`
}

func (opts GroupMembershipListReqParam) hasQueryParameter() bool {
	return opts.Start != 0 || opts.Limit != 0
}

func (opts GroupMembershipListReqParam) ToListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

func GetGroupMembershipUsers(c *golangsdk.ServiceClient, opts GroupMembershipListReqParam) (r GroupMembershipUserListResult) {
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
	url := GroupMemberShipURL(c, opts.GroupID) + queryStr
	_, err := c.Get(url, &r.Body, nil)
	r.Err = err

	return
}

func GetGroupMemberShipAllUser(c *golangsdk.ServiceClient, groupId string) ([]User, error) {
	var allUsers []User
	var err error
	start := 0
	maxLoopCount := 20
	// setting the mast loop count is 20 times
	for i := 0; i < maxLoopCount; i++ {
		opts := GroupMembershipListReqParam{
			GroupID: groupId,
			Start:   start,
			Limit:   100,
		}

		users, total, err1 := GetGroupMembershipUsers(c, opts).Extract()
		if err1 != nil {
			err = err1
			break
		}

		allUsers = append(allUsers, users...)

		// the next page has data
		if start+opts.Limit < total {
			if i == maxLoopCount-1 {
				err = fmt.Errorf("the next page still contains data. Check whether the number of users in the user group is too large: %d", i)
				break
			}
			start = start + opts.Limit
		} else {
			break
		}
	}
	return allUsers, err
}

func AddGroupMembership(c *golangsdk.ServiceClient, opts GroupMembershipReqParam) (r CreateGroupMembershipResult) {

	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{204}}
	url := GroupMemberShipURL(c, url2.PathEscape(opts.GroupID)) + "/" + url2.PathEscape(opts.UserID)
	_, err := c.Put(url, map[string]string{}, &r.Body, reqOpt)
	r.Err = err

	return
}

func DeleteGroupMembership(c *golangsdk.ServiceClient, opts GroupMembershipReqParam) (r DeleteGroupMembershipResult) {

	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{204}}
	url := GroupMemberShipURL(c, url2.PathEscape(opts.GroupID)) + "/" + url2.PathEscape(opts.UserID)
	_, err := c.Delete(url, reqOpt)
	r.Err = err

	return
}
