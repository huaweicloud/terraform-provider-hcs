package baremetalservers

import (
	"encoding/base64"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

type CreateOpts struct {
	ImageRef string `json:"imageRef" required:"true"`

	FlavorRef string `json:"flavorRef" required:"true"`

	Name string `json:"name" required:"true"`

	MetaData MetaData `json:"metadata" required:"true"`

	UserData []byte `json:"-"`

	AdminPass string `json:"adminPass,omitempty"`

	KeyName string `json:"key_name,omitempty"`

	SecurityGroups []SecurityGroup `json:"security_groups,omitempty"`

	Nics []Nic `json:"nics" required:"true"`

	AvailabilityZone string `json:"availability_zone" required:"true"`

	VpcId string `json:"vpcid" required:"true"`

	PublicIp *PublicIp `json:"publicip,omitempty"`

	Count int `json:"count,omitempty"`

	DataVolumes []DataVolume `json:"data_volumes,omitempty"`

	ExtendParam ServerExtendParam `json:"extendparam,omitempty"`

	Tags []string `json:"tags,omitempty"`
}

type MetaData struct {
	OpSvcUserId string `json:"op_svc_userid" required:"true"`
	BYOL        string `json:"BYOL,omitempty"`
	AdminPass   string `json:"admin_pass,omitempty"`
}

type SecurityGroup struct {
	ID string `json:"id" required:"true"`
}

type Nic struct {
	SubnetId  string `json:"subnet_id" required:"true"`
	IpAddress string `json:"ip_address,omitempty"`
}

type PublicIp struct {
	Id  string `json:"id,omitempty"`
	Eip *Eip   `json:"eip,omitempty"`
}

type DataVolume struct {
	VolumeType  string            `json:"volumetype" required:"true"`
	Size        int               `json:"size" required:"true"`
	Shareable   bool              `json:"shareable,omitempty"`
	extendparam map[string]string `json:"extendparam,omitempty"`
	metadata    map[string]string `json:"metadata,omitempty"`
}

type ServerExtendParam struct {
	ChargingMode        string `json:"chargingMode,omitempty"`
	RegionID            string `json:"regionID,omitempty"`
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
}

type Eip struct {
	IpType    string    `json:"iptype" required:"true"`
	BandWidth BandWidth `json:"bandwidth" required:"true"`
}

type BandWidth struct {
	Name      string `json:"name,omitempty"`
	ShareType string `json:"sharetype" required:"true"`
	bwId      string `json:"bwId,omitempty"`
	Size      int    `json:"size" required:"true"`
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToServerCreateMap() (map[string]interface{}, error)
}

// ToServerCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToServerCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if opts.UserData != nil {
		var userData string
		if _, err := base64.StdEncoding.DecodeString(string(opts.UserData)); err != nil {
			userData = base64.StdEncoding.EncodeToString(opts.UserData)
		} else {
			userData = string(opts.UserData)
		}
		b["user_data"] = &userData
	}

	return map[string]interface{}{"server": b}, nil
}

// CreatePrePaid requests a server to be provisioned to the user in the current tenant.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r JobResult) {
	reqBody, err := opts.ToServerCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), reqBody, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

// Get retrieves a particular Server based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

func Delete(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r JobResult) {
	reqBody, err := opts.ToServerCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(deleteURL(client), reqBody, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}
