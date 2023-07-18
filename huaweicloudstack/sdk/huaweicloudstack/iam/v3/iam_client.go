package v3

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/invoker"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/response"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/v3/model"
)

type IamClient struct {
	HcClient *core.HttpClient
}

func NewIamClient(hcClient *core.HttpClient) *IamClient {
	return &IamClient{HcClient: hcClient}
}

func IamClientBuilder() *core.HttpClientBuilder {
	builder := core.NewHcsHttpClientBuilder()
	return builder
}

// GetDomainByUserToken GetDomainByUserToken
func (c *IamClient) GetDomainByUserToken(request *model.GetDomainByUserTokenRequest) (*response.DefaultHttpResponse, error) {
	requestDef := GenReqDefForGetDomainByUserToken()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

// GetDomainByUserTokenInvoker  GetDomainByUserTokenInvoker
func (c *IamClient) GetDomainByUserTokenInvoker(request *model.GetDomainByUserTokenRequest) *GetDomainByUserTokenInvoker {
	requestDef := GenReqDefForGetDomainByUserToken()
	return &GetDomainByUserTokenInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// GetProjectDetail GetProjectDetail
func (c *IamClient) GetProjectDetail(request *model.GetProjectDetailRequest) (*response.DefaultHttpResponse, error) {
	requestDef := GenReqDefForGetProjectDetail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

// GetProjectDetailInvoker GetProjectDetailInvoker
func (c *IamClient) GetProjectDetailInvoker(request *model.GetProjectDetailRequest) *GetProjectDetailInvoker {
	requestDef := GenReqDefForGetProjectDetail()
	return &GetProjectDetailInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// GetProjectList GetProjectList
func (c *IamClient) GetProjectList(request *model.GetProjectListRequest) (*response.DefaultHttpResponse, error) {
	requestDef := GenReqDefForGetProjectList()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

// GetProjectListInvoker GetProjectListInvoker
func (c *IamClient) GetProjectListInvoker(request *model.GetProjectListRequest) *GetProjectListInvoker {
	requestDef := GenReqDefForGetProjectList()
	return &GetProjectListInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// GetUserDetail GetUserDetail
func (c *IamClient) GetUserDetail(request *model.GetUserDetailRequest) (*response.DefaultHttpResponse, error) {
	requestDef := GenReqDefForGetUserDetail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

// GetUserDetailInvoker GetUserDetailInvoker
func (c *IamClient) GetUserDetailInvoker(request *model.GetUserDetailRequest) *GetUserDetailInvoker {
	requestDef := GenReqDefForGetUserDetail()
	return &GetUserDetailInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// GetUserList GetUserList
func (c *IamClient) GetUserList(request *model.GetUserListRequest) (*response.DefaultHttpResponse, error) {
	requestDef := GenReqDefForGetUserList()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

// GetUserListInvoker GetUserListInvoker
func (c *IamClient) GetUserListInvoker(request *model.GetUserListRequest) *GetUserListInvoker {
	requestDef := GenReqDefForGetUserList()
	return &GetUserListInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// GetToken GetToken
func (c *IamClient) GetToken(request *model.GetTokenRequest) (*response.DefaultHttpResponse, error) {
	requestDef := GenReqDefForGetToken()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

// GetTokenInvoker GetTokenInvoker
func (c *IamClient) GetTokenInvoker(request *model.GetTokenRequest) *GetTokenInvoker {
	requestDef := GenReqDefForGetToken()
	return &GetTokenInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
