package sqlfilter

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

type JobResponse struct {
	JobID string `json:"job_id"`
}

type JobResult struct {
	golangsdk.Result
}

func (r JobResult) ExtractJobResponse() (*JobResponse, error) {
	var job JobResponse
	err := r.ExtractInto(&job)
	return &job, err
}

type SqlFilter struct {
	SwitchStatus string `json:"switch_status"`
}

type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() (*SqlFilter, error) {
	var sqlFilter SqlFilter
	err := r.ExtractInto(&sqlFilter)
	return &sqlFilter, err
}
