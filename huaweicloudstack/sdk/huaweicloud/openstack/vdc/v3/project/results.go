package project

import golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

type QueryProjectDetailV31Resp struct {
	Project QueryProjectDetailV31 `json:"project,omitempty"`
}

type QueryProjectDetailV31 struct {
	TenantId            string              `json:"tenant_id,omitempty" require:"false"`
	TenantName          string              `json:"tenant_name,omitempty" require:"false"`
	TenantType          string              `json:"tenant_type,omitempty" require:"false"`
	CreateUserId        string              `json:"create_user_id,omitempty" require:"false"`
	CreateUserName      string              `json:"create_user_name,omitempty" require:"false"`
	Description         string              `json:"description,omitempty" require:"false"`
	Enable              bool                `json:"enable,omitempty" require:"false"`
	DomainId            string              `json:"domain_id,omitempty" require:"false"`
	ContractNumber      string              `json:"contract_number,omitempty" require:"false"`
	IsShared            bool                `json:"is_shared,omitempty" require:"false"`
	Name                string              `json:"name,omitempty" require:"false"`
	IamProjectName      string              `json:"iam_project_name,omitempty" require:"false"`
	DisplayName         string              `json:"display_name,omitempty" require:"false"`
	Id                  string              `json:"id,omitempty" require:"false"`
	OwnerId             string              `json:"owner_id,omitempty" require:"false"`
	OwnerName           string              `json:"owner_name,omitempty" require:"false"`
	RegionName          string              `json:"region_name,omitempty" require:"false"`
	Regions             []ProjectRegionInfo `json:"regions,omitempty" require:"false"`
	QuotaUnitId         string              `json:"quota_unit_id,omitempty" require:"false"`
	AttachmentId        string              `json:"attachment_id,omitempty" require:"false"`
	AttachmentName      string              `json:"attachment_name,omitempty" require:"false"`
	AttachmentSize      int64               `json:"attachment_size,omitempty" require:"false"`
	IsSupportHwsService bool                `json:"is_support_hws_service,omitempty" require:"false"`
}

type ProjectRegionInfo struct {
	RegionId     string                  `json:"region_id,omitempty" require:"false"`
	RegionName   RegionName              `json:"region_name,omitempty" require:"false"`
	RegionType   string                  `json:"region_type,omitempty" require:"false"`
	RegionStatus string                  `json:"region_status,omitempty" require:"false"`
	CloudInfras  []ProjectCloudInfraInfo `json:"cloud_infras,omitempty" require:"false"`
}

type RegionName struct {
	ZhCn string `json:"zh_cn"`
	EnUs string `json:"en_us"`
}

type ProjectCloudInfraInfo struct {
	CloudInfraId string `json:"cloud_infra_id,omitempty" require:"false"`

	CloudInfraName   string `json:"cloud_infra_name,omitempty" require:"false"`
	CloudInfraType   string `json:"cloud_infra_type,omitempty" require:"false"`
	CloudInfraStatus string `json:"cloud_infra_status,omitempty" require:"false"`

	Azs    []ProjectAzInfo `json:"azs" require:"false"`
	Quotas []QuotaInfo     `json:"quotas,omitempty" require:"false"`
}

type ProjectAzInfo struct {
	AzId     string `json:"az_id,omitempty" require:"false"`
	AzName   string `json:"az_name,omitempty" require:"false"`
	AzStatus string `json:"az_status,omitempty" require:"false"`
}

type QuotaInfo struct {
	ServiceId string     `json:"service_id,omitempty" require:"false"`
	Action    string     `json:"action,omitempty" require:"false"`
	Resources []DictInfo `json:"resources,omitempty" require:"false"`
}

type DictInfo struct {
	Resource   string `json:"resource,omitempty" require:"false"`
	LocalLimit int64  `json:"local_limit,omitempty" require:"false"`
	OtherLimit int64  `json:"other_limit,omitempty" require:"false"`
	LocalUsed  int64  `json:"local_used,omitempty" require:"false"`
	OtherUsed  int64  `json:"other_used,omitempty" require:"false"`
}

type CreateProjectResponseV31 struct {
	Project ProjectId `json:"project,omitempty" require:"false"`
}

type ProjectId struct {
	Id string `json:"id,omitempty" require:"false"`
}

type CreateResult struct {
	golangsdk.Result
}

func (createResult CreateResult) Extract() (*ProjectId, error) {
	var result struct {
		Project *ProjectId `json:"project"`
	}
	err := createResult.Result.ExtractInto(&result)
	return result.Project, err
}

type UpdateResult struct {
	golangsdk.ErrResult
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type GetResult struct {
	golangsdk.ErrResult
}

func (getResult GetResult) Extract() (*QueryProjectDetailV31, error) {
	var result struct {
		Project *QueryProjectDetailV31 `json:"project"`
	}
	err := getResult.Result.ExtractInto(&result)
	return result.Project, err
}
