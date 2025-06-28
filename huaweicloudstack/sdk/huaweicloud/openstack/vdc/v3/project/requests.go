package project

import (
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

// CreateOpts 创建资源空间的北向接口支持的body参数定义
type CreateOpts struct {
	Project CreateProjectV31 `json:"project,omitempty" require:"true"`
}

type CreateProjectV31 struct {
	Name                  string   `json:"name,omitempty" require:"true"`
	DisplayName           string   `json:"display_name,omitempty" require:"false"`
	Description           string   `json:"description,omitempty" require:"false"`
	GroupId               string   `json:"group_id,omitempty" require:"false"`
	TenantId              string   `json:"tenant_id,omitempty" require:"true"`
	VrmProjectId          int64    `json:"vrm_project_id,omitempty" require:"false"`
	Regions               []Region `json:"regions,omitempty" require:"false"`
	ContractNumber        string   `json:"contract_number,omitempty" require:"false"`
	AttachmentId          string   `json:"attachment_id,omitempty" require:"false"`
	OwnerId               string   `json:"owner_id,omitempty" require:"false"`
	IsShared              string   `json:"is_shared,omitempty" require:"false"`
	IsBindExternalNetwork bool     `json:"is_bind_external_network" require:"false"`
	IsSupportHwsService   bool     `json:"is_support_hws_service" require:"false"`
	RealIamName           string   `json:"real_iam_name,omitempty" require:"false"`
}

type Region struct {
	RegionId    string       `json:"region_id,omitempty" require:"false"`
	Action      string       `json:"action,omitempty" require:"false"`
	Name        string       `json:"name,omitempty" require:"false"`
	Type        string       `json:"type,omitempty" require:"false"`
	Status      string       `json:"status,omitempty" require:"false"`
	CloudInfras []CloudInfra `json:"cloud_infras,omitempty" require:"false"`
}

type CloudInfra struct {
	CloudInfraId string `json:"cloud_infra_id,omitempty" require:"false"`

	Name           string          `json:"name,omitempty" require:"false"`
	Type           string          `json:"type,omitempty" require:"false"`
	Status         string          `json:"status,omitempty" require:"false"`
	AvailableZones []AvailableZone `json:"available_zones" require:"false"`
	Clusters       []Cluster       `json:"clusters,omitempty" require:"false"`
	Quotas         []Quota         `json:"quotas,omitempty" require:"false"`
}

type AvailableZone struct {
	AzId         string  `json:"az_id,omitempty" require:"false"`
	Name         string  `json:"name,omitempty" require:"false"`
	Type         string  `json:"type,omitempty" require:"false"`
	Status       string  `json:"status,omitempty" require:"false"`
	CloudInfraId string  `json:"cloud_infra_id,omitempty" require:"false"`
	ExtendParam  string  `json:"extend_param,omitempty" require:"false"`
	Quotas       []Quota `json:"quotas,omitempty" require:"false"`
}

type Cluster struct {
	ClusterId   string `json:"cluster_id,omitempty" require:"false"`
	ClusterName string `json:"cluster_name,omitempty" require:"false"`
}

type Quota struct {
	ServiceId string `json:"service_id,omitempty" require:"false"`
	Action    string `json:"action,omitempty" require:"false"`
	Resources []Dict `json:"resources,omitempty" require:"false"`
}

type Dict struct {
	Resource   string `json:"resource,omitempty" require:"false"`
	LocalLimit int64  `json:"local_limit,omitempty" require:"false"`
	OtherLimit int64  `json:"other_limit,omitempty" require:"false"`
	LocalUsed  int64  `json:"local_used,omitempty" require:"false"`
	OtherUsed  int64  `json:"other_used,omitempty" require:"false"`
	Min        int64  `json:"min,omitempty" require:"false"`
	Max        int64  `json:"max,omitempty" require:"false"`
}

func (opts CreateOpts) GetRequestBodyForCreateVdcProject() (map[string]interface{}, error) {
	requestBody, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return requestBody, nil
}

type CreateOptsBuilder interface {
	GetRequestBodyForCreateVdcProject() (map[string]interface{}, error)
}

// Create 创建角色北向接口调用方法
func Create(httpClient *golangsdk.ServiceClient, opts CreateOptsBuilder) (createResult CreateResult) {
	url := CreateVdcProjectURL(httpClient) // 获取创建的接口地址
	requestBody, err := opts.GetRequestBodyForCreateVdcProject()
	if err != nil {
		createResult.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{
		OkCodes:     []int{201},
		MoreHeaders: golangsdk.MoreHeaders,
	}
	_, createResult.Err = httpClient.Post(url, &requestBody, &createResult.Body, reqOpt)
	return
}

// DeleteOpts 删除资源空间的北向接口支持的query参数定义
type DeleteOpts struct {
	DeleteServiceTypes string `q:"delete_service_types,omitempty" require:"false"`
}

func (opts DeleteOpts) GetRequestQueryForDeleteVdcProject() (string, error) {
	queryParams, err := golangsdk.BuildQueryString(opts)
	return queryParams.String(), err
}

type DeleteOptsBuilder interface {
	GetRequestQueryForDeleteVdcProject() (string, error)
}

func Delete(client *golangsdk.ServiceClient, projectId string, deleteOptsBuilder DeleteOptsBuilder) (deleteResult DeleteResult) {
	url, err := DeleteVdcProjectURL(client, projectId)
	if err != nil {
		deleteResult.Err = err
		return
	}
	if deleteOptsBuilder != nil {
		query, err := deleteOptsBuilder.GetRequestQueryForDeleteVdcProject() // 构建查询参数为string类型
		if err != nil {
			deleteResult.Err = err
		}
		url += query // 将构建参数拼接到接口地址末尾
	}
	_, deleteResult.Err = client.Delete(url, &golangsdk.RequestOpts{
		OkCodes:     []int{204},
		MoreHeaders: golangsdk.MoreHeaders,
	})
	return
}

// UpdateOpts 更新角色的北向接口支持的body参数定义
type UpdateOpts struct {
	Project UpdateProject `json:"project,omitempty" require:"true"`
}

type UpdateProject struct {
	Name                string `json:"name,omitempty" require:"false"`
	DisplayName         string `json:"display_name,omitempty" require:"false"`
	Description         string `json:"description,omitempty" require:"false"`
	ContractNumber      string `json:"contract_number,omitempty" require:"false"`
	AttachmentId        string `json:"attachment_id,omitempty" require:"false"`
	OwnerId             string `json:"owner_id,omitempty" require:"false"`
	IsShared            string `json:"is_shared,omitempty" require:"false"`
	IsSupportHwsService bool   `json:"is_support_hws_service,omitempty" require:"false"`
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

func Update(client *golangsdk.ServiceClient, projectId string, opts UpdateOptsBuilder) (updateResult UpdateResult) {
	b, err := opts.GetRequestBodyForUpdateVdcRole()
	if err != nil {
		updateResult.Err = err
		return
	}

	url, err := UpdateVdcProjectURL(client, projectId)
	if err != nil {
		updateResult.Err = err
		return
	}
	_, updateResult.Err = client.Put(url, &b, &updateResult.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{204},
		MoreHeaders: golangsdk.MoreHeaders,
	})
	return
}

// GetOpts 查询资源空间详情的北向接口支持的query参数定义
type GetOpts struct {
	ServiceId string `q:"service_id,omitempty" require:"false"`
}

func (opts GetOpts) GetRequestQueryForGetVdcProject() (string, error) {
	queryParams, err := golangsdk.BuildQueryString(opts)
	return queryParams.String(), err
}

type GetOptsBuilder interface {
	GetRequestQueryForGetVdcProject() (string, error)
}

func Get(client *golangsdk.ServiceClient, projectId string, opts GetOptsBuilder) (getResult GetResult) {
	url, err := GetVdcProjectURL(client, projectId)
	if err != nil {
		getResult.Err = err
		return
	}
	if opts != nil {
		query, err := opts.GetRequestQueryForGetVdcProject() // 构建查询参数为string类型
		if err != nil {
			getResult.Err = err
		}
		url += query // 将构建参数拼接到接口地址末尾
	}
	_, getResult.Err = client.Get(url, &getResult.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: golangsdk.MoreHeaders,
	})
	return
}
