package group

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

type CreatResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type GroupModel struct {
	ID          string `json:"id"`
	VdcId       string `json:"vdc_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r commonResult) ToExtract() (GroupModel, error) {
	var a struct {
		Group GroupModel `json:"Group"`
	}
	err := r.Result.ExtractInto(&a)
	return a.Group, err
}

type UpdateResult struct {
	golangsdk.Result
}

func (r UpdateResult) ToExtract() ([]int, error) {
	var code []int
	err := r.Result.ExtractInto(&code)
	return code, err
}

type DeleteResult struct {
	golangsdk.ErrResult
}

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
