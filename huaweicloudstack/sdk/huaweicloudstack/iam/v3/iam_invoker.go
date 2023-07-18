package v3

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/invoker"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/response"
)

type GetDomainByUserTokenInvoker struct {
	*invoker.BaseInvoker
}

func (i *GetDomainByUserTokenInvoker) Invoke() (*response.DefaultHttpResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

type GetProjectDetailInvoker struct {
	*invoker.BaseInvoker
}

func (i *GetProjectDetailInvoker) Invoke() (*response.DefaultHttpResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

type GetProjectListInvoker struct {
	*invoker.BaseInvoker
}

func (i *GetProjectListInvoker) Invoke() (*response.DefaultHttpResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

type GetUserDetailInvoker struct {
	*invoker.BaseInvoker
}

func (i *GetUserDetailInvoker) Invoke() (*response.DefaultHttpResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

type GetUserListInvoker struct {
	*invoker.BaseInvoker
}

func (i *GetUserListInvoker) Invoke() (*response.DefaultHttpResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

type GetTokenInvoker struct {
	*invoker.BaseInvoker
}

func (i *GetTokenInvoker) Invoke() (*response.DefaultHttpResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}
