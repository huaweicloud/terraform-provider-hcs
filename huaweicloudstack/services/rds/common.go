package rds

import (
	"encoding/json"
	"fmt"

	"github.com/jmespath/go-jmespath"

	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

var (
	retryErrCodes = map[string]struct{}{
		"DBS.201202":   {},
		"DBS.200011":   {},
		"DBS.200018":   {},
		"DBS.200019":   {},
		"DBS.200047":   {},
		"DBS.200076":   {},
		"DBS.200611":   {},
		"DBS.200080":   {},
		"DBS.200463":   {}, // create replica instance
		"DBS.201015":   {},
		"DBS.201206":   {},
		"DBS.212033":   {}, // http response code is 403
		"DBS.280011":   {},
		"DBS.280343":   {},
		"DBS.280816":   {},
		"DBS.01010337": {},
		"DBS.01280030": {}, // instance status is illegal
	}
)

// The RDS instance is limited to only one operation at a time.
// In addition to locking and waiting between multiple operations, a retry method is required to ensure that the
// request can be executed correctly.
func handleMultiOperationsError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("error_code||errCode", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		if _, ok = retryErrCodes[errorCode.(string)]; ok {
			return true, err
		}
	}
	if errCode, ok := err.(golangsdk.ErrDefault409); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("errCode||error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		if _, ok = retryErrCodes[errorCode.(string)]; ok {
			// The operation failed to execute and needs to be executed again, because other operations are
			// currently in progress.
			return true, err
		}
	}
	if errCode, ok := err.(golangsdk.ErrDefault403); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("errCode||error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		if _, ok = retryErrCodes[errorCode.(string)]; ok {
			// The operation failed to execute and needs to be executed again, because other operations are
			// currently in progress.
			return true, err
		}
	}
	if errCode, ok := err.(golangsdk.ErrDefault500); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		// if the error code is RDS.0005, it indicates that the SSL is changed, and the db is rebooted
		if errorCode.(string) == "RDS.0005" {
			return true, err
		}
	}
	// Operation execution failed due to some resource or server issues, no need to try again.
	return false, err
}
