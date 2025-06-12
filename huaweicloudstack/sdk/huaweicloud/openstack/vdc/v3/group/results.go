package vdc

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

// Group is the structure that represents the details of the forward Group.
type Group struct {
	ID          string `json:"id"`
	DomainId    string `json:"domain_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	VdcId       string `json:"vdc_id"`
	CreateAt    int    `json:"create_at"`
}

type ListResult struct {
	commonResult
}

type commonResult struct {
	golangsdk.Result
}

func (r ListResult) Extract() ([]Group, int, error) {
	var s struct {
		Groups []Group `json:"groups"`
		Total  int     `json:"total"`
	}
	err := r.ExtractInto(&s)
	return s.Groups, s.Total, err
}
