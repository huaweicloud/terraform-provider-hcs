package baremetalservers

import (
	"encoding/base64"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

type CreateOpts struct {
	ImageRef string `json:"imageRef" required:"true"`

	FlavorRef string `json:"flavorRef" required:"true"`

	Name string `json:"name" required:"true"`

	MetaData MetaData `json:"metadata"`

	UserData []byte `json:"-"`

	AdminPass string `json:"admin_password,omitempty"`

	KeyName string `json:"key_name,omitempty"`

	SecurityGroups []SecurityGroup `json:"security_groups"`

	Nics []Nic `json:"nics" required:"true"`

	AvailabilityZone string `json:"availability_zone" required:"true"`

	VpcId string `json:"vpcid" required:"true"`

	PublicIp *PublicIp `json:"publicip,omitempty"`

	Count int `json:"count,omitempty"`

	DataVolumes []DataVolume `json:"data_volumes,omitempty"`

	ExtendParam ServerExtendParam `json:"extendparam,omitempty"`

	Tags []interface{} `json:"tags,omitempty"`
}

type MetaData struct {
	OpSvcUserId string `json:"op_svc_userid,omitempty"`
	BYOL        string `json:"BYOL,omitempty"`
	AdminPass   string `json:"admin_pass,omitempty"`
	AgencyName  string `json:"agency_name,omitempty"`
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
	Extendparam map[string]string `json:"extendparam,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
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

type DeleteOpts struct {
	Servers        []Server `json:"servers" required:"true"`
	DeletePublicIp bool     `json:"delete_publicip"`
	DeleteVolume   bool     `json:"delete_volume"`
}

type Server struct {
	Id string `json:"id" required:"true"`
}

// ToServerDeleteMap assembles a request body based on the contents of a
// DeleteOpts.
func (opts DeleteOpts) ToServerDeleteMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToServerListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned. Marker and Limit are used
// for pagination.
type ListOpts struct {
	// Name of the server as a string; can be queried with regular expressions.
	// Realize that ?name=bob returns both bob and bobb. If you need to match bob
	// only, you can use a regular expression matching the syntax of the
	// underlying database server implemented for Compute.
	Name string `q:"name"`

	// Status is the value of the status of the server so that you can filter on
	// "ACTIVE" for example.
	Status string `q:"status"`

	// Specifies the BMS' id.
	ID string `q:"id"`

	// Specifies the BMS' tags.
	Tags string `q:"tags"`

	// Expected field to be returned.
	ExpectFields string `q:"expect_fields"`

	// Specifies the BMS that is bound to an enterprise project.
	EnterpriseProjectID string `q:"enterprise_project_id"`

	// Specifies the maximum number of ECSs on one page.
	// Each page contains 25 BMSs by default, and a maximum of 1000 BMSs are returned.
	Limit int `q:"limit"`

	// Specifies a page number. The default value is 1.
	// The value must be greater than or equal to 0. If the value is 0, the first page is displayed.
	Offset int `q:"offset"`
}

// ToServerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToServerListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
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
func Get(client *golangsdk.ServiceClient, id string, opts ListOptsBuilder) (r GetResult) {
	url := getURL(client, id)
	if opts != nil {
		query, err := opts.ToServerListQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 203},
	})
	return
}

func Delete(client *golangsdk.ServiceClient, opts DeleteOpts) (r JobResult) {
	reqBody, err := opts.ToServerDeleteMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(deleteURL(client), reqBody, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}
