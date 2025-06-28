package user

import golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

type commonResult struct {
	golangsdk.Result
}

type CreatResult struct {
	commonResult
}

type UpdateResult struct {
	golangsdk.Result
}

func (r UpdateResult) ToExtract() ([]int, error) {
	var code []int
	err := r.Result.ExtractInto(&code)
	return code, err
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type GetResult struct {
	commonResult
}

func (r commonResult) ToExtract() (VdcUserModel, error) {
	var a struct {
		User VdcUserModel `json:"User"`
	}
	err := r.Result.ExtractInto(&a)
	return a.User, err
}

var AuthType map[string]string = map[string]string{
	"0": "LOCAL_AUTH",
	"1": "SAML_AUTH",
	"2": "LDAP_AUTH",
	"4": "MACHINE_USER",
}

var AccessMode map[string]string = map[string]string{
	"0": "default",
	"1": "console",
	"2": "programmatic",
}

type VdcUserModel struct {
	ID          string `json:"id"`
	VdcId       string `json:"vdc_id"`
	TopVdcId    string `json:"top_vdc_id"`
	DomainId    string `json:"domain_id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	AuthType    string `json:"auth_type"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	AccessMode  string `json:"access_mode"`
	LdapId      string `json:"ldap_id"`
	CreateAt    int64  `json:"create_at"`
}

type UserList struct {
	Users []VdcUserModel `json:"users"`
	Total int            `json:"total"`
}

type ListResult struct {
	golangsdk.Result
}

func (r ListResult) Extract() (UserList, error) {
	var a struct {
		Users []VdcUserModel `json:"users"`
		Total int            `json:"total"`
	}
	err := r.Result.ExtractInto(&a)
	return a, err
}
