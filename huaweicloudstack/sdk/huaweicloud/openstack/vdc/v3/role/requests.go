package role

import (
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

// ListOpts 查询角色列表的北向接口支持的query参数定义
type ListOpts struct {
	DomainId    string `q:"domain_id"`
	IsSystem    string `q:"is_system"`
	FineGrained bool   `q:"fine_grained"`
	Start       int    `q:"start"`
	Limit       int    `q:"limit"`
}

func (opts ListOpts) GetVdcRoleQuery() (string, error) {
	queryParams, err := golangsdk.BuildQueryString(opts)
	return queryParams.String(), err
}

type ListOptsBuilder interface {
	GetVdcRoleQuery() (string, error)
}

func List(httpClient *golangsdk.ServiceClient, listOptsBuilder ListOptsBuilder) (listResult ListResult) {
	url := ListVdcRoleURL(httpClient) // 获取列表的接口地址
	if listOptsBuilder != nil {
		query, err := listOptsBuilder.GetVdcRoleQuery() // 构建查询参数为string类型
		if err != nil {
			listResult.Err = err
		}
		url += query // 将构建参数拼接到接口地址末尾
	}

	// 使用get方法发送请求
	_, listResult.Err = httpClient.Get(url, &listResult.Body, nil)
	return
}

// CreateOpts 创建角色的北向接口支持的body参数定义
type CreateOpts struct {
	DomainId string             `json:"domain_id,omitempty" require:"true"` // 创建接口要求必须传domain_id，若用户在resource中指定了domain_id，则以填的为准，否则从provider中获取。
	Role     RequestBodyVdcRole `json:"role,omitempty" require:"true"`
}

type RequestBodyVdcRole struct {
	DisplayName string     `json:"display_name,omitempty" require:"true"` // name对应接口中的display_name字段
	Type        string     `json:"type,omitempty" require:"true"`
	Description string     `json:"description" require:"false"`
	Policy      PolicyBase `json:"policy,omitempty" require:"true"`
	Tag         string     `json:"tag"`
}

func (opts CreateOpts) GetRequestBodyForCreateVdcRole() (map[string]interface{}, error) {
	requestBody, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return requestBody, nil
}

type CreateOptsBuilder interface {
	GetRequestBodyForCreateVdcRole() (map[string]interface{}, error)
}

// Create 创建角色北向接口调用方法
func Create(httpClient *golangsdk.ServiceClient, opts CreateOptsBuilder) (createResult CreateResult) {
	url := CreateVdcRoleURL(httpClient) // 获取创建的接口地址
	requestBody, err := opts.GetRequestBodyForCreateVdcRole()
	if err != nil {
		createResult.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}
	_, createResult.Err = httpClient.Post(url, &requestBody, &createResult.Body, reqOpt)
	return
}

func Delete(client *golangsdk.ServiceClient, roleID string) (r DeleteResult) {
	_, r.Err = client.Delete(DeleteVdcRoleURL(client, roleID), &golangsdk.RequestOpts{
		OkCodes:     []int{204},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	})
	return
}

func Update(client *golangsdk.ServiceClient, roleID string, opts CreateOptsBuilder) (updateResult UpdateResult) {
	b, err := opts.GetRequestBodyForCreateVdcRole()
	if err != nil {
		updateResult.Err = err
		return
	}
	_, updateResult.Err = client.Put(UpdateVdcRoleURL(client, roleID), &b, &updateResult.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	})
	return
}

func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(GetVdcRoleURL(client, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	})
	return
}
