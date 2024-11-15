package clone

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"log"
)

type OsCloneOpts struct {
	CloneNum     int               `json:"clone_num,omitempty"`
	CloneType    string            `json:"clone_type,omitempty"`
	RetainPasswd *bool             `json:"retain_passwd,omitempty"`
	KeyPair      string            `json:"key_name,omitempty"`
	MetaData     map[string]string `json:"metadata,omitempty"`
	Name         string            `json:"name,omitempty"`
	UserData     string            `json:"user_data,omitempty"`
	PowerOn      *bool             `json:"power_on" required:"true"`
	Postfix      string            `json:"postfix,omitempty"`
	Nics         []NicsOpts        `json:"nics,omitempty"`
	VpcId        string            `json:"vpc_id,omitempty"`
	AdminPass    string            `json:"admin_password,omitempty"`
}

type NicsOpts struct {
	SubnetId       string               `json:"subnet_id" required:"true"`
	IpAddress      string               `json:"ip_address,omitempty"`
	IpAddressV6    string               `json:"ip_address_v6,omitempty"`
	Ipv6Enable     bool                 `json:"ipv6_enable,omitempty"`
	SecurityGroups []SecurityGroupsOpts `json:"security_groups,omitempty"`
}

type SecurityGroupsOpts struct {
	Id string `json:"id" required:"true"`
}

func CloneVm(c *golangsdk.ServiceClient, id string, opts OsCloneOpts) (r JobResult) {
	b, err := golangsdk.BuildRequestBody(opts, "os-clone")
	if err != nil {
		log.Printf("erros is: %q", err)
		return
	}
	log.Printf("[DEBUG] clone vm url:%q", cloneURL(c, id))
	_, r.Err = c.Post(cloneURL(c, id), b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

// Get retrieves a particular Server based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 203},
	})
	return
}

func getURL(sc *golangsdk.ServiceClient, serverID string) string {
	return sc.ServiceURL("cloudservers", serverID)
}
