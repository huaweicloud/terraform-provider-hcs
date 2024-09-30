package instances

import (
	"github.com/chnsz/golangsdk"
)

type commonResult struct {
	golangsdk.Result
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Service.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Service.
type GetResult struct {
	commonResult
}

type RomaInstance struct {
	ID                  string    `json:"instance_id"`
	Name                string    `json:"instance_name"`
	SiteId              string    `json:"site_id"`
	Description         string    `json:"description"`
	FlavorId            string    `json:"flavor_id"`
	FlavorType          string    `json:"flavor_type"`
	ProjectId           string    `json:"project_id"`
	AvailableZoneIds    []string  `json:"available_zone_ids"`
	VpcId               string    `json:"vpc_id"`
	SubnetId            string    `json:"subnet_id"`
	SecurityGroupId     string    `json:"security_group_id"`
	CpuArch             string    `json:"cpu_arch"`
	Status              string    `json:"status"`
	PublicIpId          string    `json:"publicip_id"`
	PublicIpAddress     string    `json:"publicip_address"`
	PublicIpEnable      bool      `json:"publicip_enable"`
	ConnectAddress      string    `json:"connect_address"`
	ChargeType          string    `json:"charge_type"`
	Bandwidths          int       `json:"bandwidths"`
	Ipv6Enable          bool      `json:"ipv6_enable"`
	MaintainBegin       string    `json:"maintain_begin"`
	MaintainEnd         string    `json:"maintain_end"`
	EnterpriseProjectId string    `json:"enterprise_project_id"`
	Resources           Resources `json:"resources"`
	Ipv6ConnectAddress  string    `json:"ipv6_connect_address"`
	RocketmqEnable      bool      `json:"rocketmq_enable"`
	ExternalElbEnable   bool      `json:"external_elb_enable"`
	ExternalElbId       string    `json:"external_elb_id"`
	ExternalElbAddress  string    `json:"external_elb_address"`
	ExternalEipBound    string    `json:"external_eip_bound"`
	ExternalEipId       string    `json:"external_eip_id"`
	ExternalEipAddress  string    `json:"external_eip_address"`
	CreateTime          string    `json:"create_time"`
	UpdateTime          string    `json:"update_time"`
}

type Resources struct {
	Compose Compose `json:"compose"`
	Abm     Abm     `json:"abm"`
	Bfs     Bfs     `json:"bfs"`
	Lb      Lb      `json:"lb"`
	Apic    Apic    `json:"apic"`
	Fdi     Fdi     `json:"fdi"`
	Link    Link    `json:"link"`
	Mqs     Mqs     `json:"mqs"`
}

type Mqs struct {
	ID                    string `json:"id"`
	Enable                bool   `json:"enable"`
	RetentionPolicy       string `json:"retention_policy"`
	SslEnable             bool   `json:"ssl_enable"`
	TraceEnable           bool   `json:"trace_enable"`
	VpcClientPlain        bool   `json:"vpc_client_plain"`
	PartitionNum          int    `json:"partition_num"`
	Specification         string `json:"specification"`
	PrivateConnectAddress string `json:"private_connect_address"`
	PublicConnectAddress  string `json:"public_connect_address"`
	PrivateRestfulAddress string `json:"private_restful_address"`
	PublicRestfulAddress  string `json:"public_restful_address"`
}

type Compose struct {
	ID     string `json:"id"`
	Enable bool   `json:"enable"`
}

type Abm struct {
	ID     string `json:"id"`
	Enable bool   `json:"enable"`
}

type Bfs struct {
	ID     string `json:"id"`
	Enable bool   `json:"enable"`
}

type Lb struct {
	ID     string `json:"id"`
	Enable bool   `json:"enable"`
}

type Apic struct {
	ID     string `json:"id"`
	Enable bool   `json:"enable"`
}

type Fdi struct {
	ID     string `json:"id"`
	Enable bool   `json:"enable"`
}

type Link struct {
	ID     string `json:"id"`
	Enable bool   `json:"enable"`
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}

func (r commonResult) Extract() (*RomaInstance, error) {
	var s RomaInstance
	err := r.ExtractInto(&s)
	return &s, err
}
