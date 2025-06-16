package role

import (
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

// ListOpts 北向接口支持的查询参数定义
type ListOpts struct {
	DomainId    string `q:"domain_id"`
	IsSystem    string `q:"is_system"`
	FineGrained bool   `q:"fine_grained"`
	Start       int    `q:"start"`
	Limit       int    `q:"limit"`
}

func (opts ListOpts) GetVdcRoleQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

type ListOptsBuilder interface {
	GetVdcRoleQuery() (string, error)
}

func List(httpClient *golangsdk.ServiceClient, listOptsBuilder ListOptsBuilder) (listResult ListResult) {
	url := GetVdcRoleURL(httpClient) // 获取列表的接口地址
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
