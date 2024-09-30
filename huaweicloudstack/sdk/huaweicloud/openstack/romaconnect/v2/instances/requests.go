package instances

import (
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToRomaInstanceCreateMap() (map[string]interface{}, error)
}

type CreateOpts struct {
	Name                  string   `json:"name" required:"true"`
	Description           string   `json:"description" required:"true"`
	ProductId             string   `json:"product_id" required:"true"`
	AvailableZones        []string `json:"available_zones" required:"true"`
	EnterpriseProjectId   string   `json:"ep_id" required:"false"`
	VpcId                 string   `json:"vpc_id" required:"true"`
	SubnetId              string   `json:"subnet_id" required:"true"`
	SecurityGroupId       string   `json:"security_group_id" required:"true"`
	Ipv6Enable            *bool    `json:"ipv6_enable" required:"true"`
	EnableAll             *bool    `json:"enable_all" required:"true"`
	EipId                 string   `json:"eip_id" required:"false"`
	EntranceBandwidthSize int      `json:"entrance_bandwidth_size" required:"false"`
	Mqs                   MqsOpts  `json:"mqs" required:"true"`
	MaintainBegin         string   `json:"maintain_begin" required:"false"`
	MaintainEnd           string   `json:"maintain_end" required:"false"`
	CpuArchitecture       string   `json:"cpu_architecture" required:"true"`
}

type MqsOpts struct {
	EngineVersion   string `json:"engine_version" required:"false"`
	RocketMqEnable  bool   `json:"rocketmq_enable" required:"false"`
	EnablePublicIp  bool   `json:"enable_publicip" required:"false"`
	EnableAcl       bool   `json:"enable_acl" required:"false"`
	SslEnable       bool   `json:"ssl_enable" required:"false"`
	RetentionPolicy string `json:"retention_policy" required:"false"`
	TraceEnable     bool   `json:"trace_enable" required:"false"`
	VpcClientPlain  bool   `json:"vpc_client_plain" required:"false"`
	ConnectorEnable bool   `json:"connector_enable" required:"false"`
}

// ToRomaInstanceCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToRomaInstanceCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new ROMA Connect instances based on the values in CreateOpts.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRomaInstanceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get retrieves the ROMA Connect instances with the provided ID.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)
	return
}

// Delete will delete the existing ROMA Connect instances with the provided ID.
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}
