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
	DomainId string             `json:"domain_id,omitempty" require:"true"`
	Role     RequestBodyVdcRole `json:"role,omitempty" require:"true"`
}

// UpdateOpts 更新角色的北向接口支持的body参数定义
type UpdateOpts struct {
	DomainId string             `json:"domain_id,omitempty" require:"true"`
	Role     RequestBodyVdcRole `json:"role,omitempty" require:"true"`
}

type UpdateOptsBuilder interface {
	GetRequestBodyForUpdateVdcRole() (map[string]interface{}, error)
}

func (opts UpdateOpts) GetRequestBodyForUpdateVdcRole() (map[string]interface{}, error) {
	requestBody, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return requestBody, nil
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
		MoreHeaders: golangsdk.MoreHeaders,
	}
	_, createResult.Err = httpClient.Post(url, &requestBody, &createResult.Body, reqOpt)
	return
}

func Delete(client *golangsdk.ServiceClient, roleId string) (deleteResult DeleteResult) {
	url, err := DeleteVdcRoleURL(client, roleId)
	if err != nil {
		deleteResult.Err = err
		return
	}
	_, deleteResult.Err = client.Delete(url, &golangsdk.RequestOpts{
		OkCodes:     []int{204},
		MoreHeaders: golangsdk.MoreHeaders,
	})
	return
}

func Update(client *golangsdk.ServiceClient, roleId string, opts UpdateOptsBuilder) (updateResult UpdateResult) {
	b, err := opts.GetRequestBodyForUpdateVdcRole()
	if err != nil {
		updateResult.Err = err
		return
	}

	url, err := UpdateVdcRoleURL(client, roleId)
	if err != nil {
		updateResult.Err = err
		return
	}
	_, updateResult.Err = client.Put(url, &b, &updateResult.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: golangsdk.MoreHeaders,
	})
	return
}

func Get(client *golangsdk.ServiceClient, roleId string) (getResult GetResult) {

	url, err := GetVdcRoleURL(client, roleId)
	if err != nil {
		getResult.Err = err
		return
	}
	_, getResult.Err = client.Get(url, &getResult.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: golangsdk.MoreHeaders,
	})
	return
}
