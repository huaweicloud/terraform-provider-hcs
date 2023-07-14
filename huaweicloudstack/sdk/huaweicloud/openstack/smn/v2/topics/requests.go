package topics

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/pagination"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOpsBuilder is used for creating topic parameters.
// any struct providing the parameters should implement this interface
type CreateOpsBuilder interface {
	ToTopicCreateMap() (map[string]interface{}, error)
}

// CreateOps is a struct that contains all the parameters.
type CreateOps struct {
	//Name of the topic to be created
	Name string `json:"name" required:"true"`

	//Topic display name
	DisplayName string `json:"display_name,omitempty"`

	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
}

func (ops CreateOps) ToTopicCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(ops, "")
}

// CreateOpsBuilder is used for updating topic parameters.
// any struct providing the parameters should implement this interface
type UpdateOpsBuilder interface {
	ToTopicUpdateMap() (map[string]interface{}, error)
}

// UpdateOps is a struct that contains all the parameters.
type UpdateOps struct {
	//Topic display name
	DisplayName string `json:"display_name,omitempty"`
}

func (ops UpdateOps) ToTopicUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(ops, "")
}

// Create a topic with given parameters.
func Create(client *golangsdk.ServiceClient, ops CreateOpsBuilder) (r CreateResult) {
	b, err := ops.ToTopicCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{201, 200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return
}

// Update a topic with given parameters.
func Update(client *golangsdk.ServiceClient, ops UpdateOpsBuilder, id string) (r UpdateResult) {
	b, err := ops.ToTopicUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(updateURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return
}

// delete a topic via id
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

// get a topic with detailed information by id
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

// list all the topics
func List(client *golangsdk.ServiceClient) (r ListResult) {
	pages, err := pagination.NewPager(client, listURL(client),
		func(r pagination.PageResult) pagination.Page {
			p := TopicPage{pagination.OffsetPageBase{PageResult: r}}
			return p
		}).AllPages()

	if err != nil {
		r.Err = err
		return
	}

	r.Body = pages.GetBody()
	return
}

type UpdatePoliciesOpts struct {
	// the value can be empty
	Value string `json:"value"`
}

// Update policies of the topic.
func UpdatePolicies(client *golangsdk.ServiceClient, ops UpdatePoliciesOpts, id, policyName string) (r UpdateResult) {
	b, err := golangsdk.BuildRequestBody(ops, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(updatePoliciesURL(client, id, policyName), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return
}

// get policies of the topic
func GetPolicies(client *golangsdk.ServiceClient, id, policyName string) (r GetPoliciesResult) {
	_, r.Err = client.Get(getPoliciesURL(client, id, policyName), &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}
