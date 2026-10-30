package rds

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

// updateInstanceFieldParams defines the parameters used to update a single field of an RDS instance.
//
// Simplified from the upstream version: checkJobExpression and checkOrderExpression are removed,
// because the updates handled here (failover strategy, ops window) are synchronous and their
// responses contain no jobId. If an asynchronous update is needed later, add these fields back
// and port checkRDSInstanceJobFinish from the upstream rds package.
type updateInstanceFieldParams struct {
	httpUrl             string
	httpMethod          string
	pathParams          map[string]string
	updateBodyParams    interface{}
	isRetry             bool
	timeout             string
	isWaitInstanceReady bool
}

// updateRdsInstanceField sends a PUT request to update a single field of an RDS instance.
//
// Depends on functions that already exist in the local rds package:
//   - handleMultiOperationsError (common.go)
//   - rdsInstanceStateRefreshFunc (resource_hcs_rds_pg_database_privilege.go)
//   - common.RetryContextWithWaitForState, utils.FlattenResponse
//
// isWaitInstanceReady: set to true when called from UpdateContext. The upstream UpdateContext may
// return while the instance is still in UPDATING state, so we must wait for ACTIVE before issuing
// the next PUT. CreateContext does not need it, as the upstream create already waits for ACTIVE.
func updateRdsInstanceField(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	params updateInstanceFieldParams) (interface{}, error) {
	updatePath := client.Endpoint + params.httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	for pathParam, value := range params.pathParams {
		updatePath = strings.ReplaceAll(updatePath, fmt.Sprintf("{%s}", pathParam), value)
	}

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         params.updateBodyParams,
	}

	var res interface{}
	var err error
	if params.isRetry {
		retryFunc := func() (interface{}, bool, error) {
			r, err := client.Request(params.httpMethod, updatePath, &updateOpt)
			shouldRetry, err := handleMultiOperationsError(err)
			return r, shouldRetry, err
		}
		res, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rdsInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(params.timeout),
			DelayTimeout: 10 * time.Second,
			PollInterval: 10 * time.Second,
		})
	} else {
		res, err = client.Request(params.httpMethod, updatePath, &updateOpt)
	}
	if err != nil {
		return nil, err
	}

	updateRespBody, err := utils.FlattenResponse(res.(*http.Response))
	if err != nil && err.Error() != "EOF" {
		return nil, err
	}

	if params.isWaitInstanceReady {
		stateConf := &resource.StateChangeConf{
			Target:       []string{"ACTIVE"},
			Refresh:      rdsInstanceStateRefreshFunc(client, d.Id()),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        1 * time.Second,
			PollInterval: 10 * time.Second,
		}
		if _, err = stateConf.WaitForStateContext(ctx); err != nil {
			return nil, fmt.Errorf("error waiting for instance (%s) to ready: %s", d.Id(), err)
		}
	}

	return updateRespBody, nil
}
