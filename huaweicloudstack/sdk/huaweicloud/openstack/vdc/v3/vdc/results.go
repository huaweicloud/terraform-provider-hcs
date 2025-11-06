package vdc

import golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

type VdcModel struct {
	ID          string `json:"id"`
	VdcId       string `json:"vdc_id"`
	TopVdcId    string `json:"top_vdc_id"`
	DomainId    string `json:"domain_id"`
	DomainName  string `json:"domain_name"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

type VdcList struct {
	Vdcs  []VdcModel `json:"vdcs"`
	Total int        `json:"total"`
}

type ListResult struct {
	golangsdk.Result
}

func (r ListResult) Extract() (VdcList, error) {
	var a struct {
		Vdcs  []VdcModel `json:"vdcs"`
		Total int        `json:"total"`
	}
	err := r.Result.ExtractInto(&a)
	return a, err
}
