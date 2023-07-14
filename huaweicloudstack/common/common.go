/* Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights resvered. */
/*
The common package defines some common functions, which are mainly used for the functions of the following services.

The difference between common package and utils:
1. Common functions under common are related to the project, and common functions are placed here.
2. Utils are some stored tool functions, which are not related to the project.
   Such as: date conversion, type conversion.
*/
package common

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/bss/v2/orders"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/bss/v2/resources"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/eps/v1/enterpriseprojects"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/eips"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
)

// ErrorResp is the response when API failed
type ErrorResp struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func ParseErrorMsg(body []byte) (ErrorResp, error) {
	resp := ErrorResp{}
	err := json.Unmarshal(body, &resp)
	return resp, err
}

// GetRegion returns the region that was specified ina the resource. If a
// region was not set, the provider-level region is checked. The provider-level
// region can either be set by the region argument or by HW_REGION_NAME.
func GetRegion(d *schema.ResourceData, config *config.Config) string {
	if v, ok := d.GetOk("region"); ok {
		return v.(string)
	}

	return config.Region
}

// GetEnterpriseProjectID returns the enterprise_project_id that was specified in the resource.
// If it was not set, the provider-level value is checked. The provider-level value can
// either be set by the `enterprise_project_id` argument or by HW_ENTERPRISE_PROJECT_ID.
func GetEnterpriseProjectID(d *schema.ResourceData, config *config.Config) string {
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		return v.(string)
	}

	return config.EnterpriseProjectID
}

func MigrateEnterpriseProject(client *golangsdk.ServiceClient, region, targetEPSId, resourceType, resourceID string) error {
	if targetEPSId == "" {
		targetEPSId = "0"
	} else {
		// check enterprise_project_id existed
		if result := enterpriseprojects.Get(client, targetEPSId); result.Err != nil {
			return fmt.Errorf("failed to query the target enterprise project %s: %s", targetEPSId, result.Err)
		}
	}

	migrateOpts := enterpriseprojects.MigrateResourceOpts{
		RegionId:     region,
		ProjectId:    client.ProjectID,
		ResourceType: resourceType,
		ResourceId:   resourceID,
	}
	migrateResult := enterpriseprojects.Migrate(client, migrateOpts, targetEPSId)
	if err := migrateResult.Err; err != nil {
		return fmt.Errorf("failed to migrate %s to enterprise project %s, err: %s", resourceID, targetEPSId, err)
	}

	return nil
}

// GetEipIDbyAddress returns the EIP ID of address when success.
func GetEipIDbyAddress(client *golangsdk.ServiceClient, address, epsID string) (string, error) {
	listOpts := &eips.ListOpts{
		PublicIp:            []string{address},
		EnterpriseProjectId: epsID,
	}
	pages, err := eips.List(client, listOpts).AllPages()
	if err != nil {
		return "", err
	}

	allEips, err := eips.ExtractPublicIPs(pages)
	if err != nil {
		return "", fmtp.Errorf("Unable to retrieve eips: %s ", err)
	}

	total := len(allEips)
	if total == 0 {
		return "", fmtp.Errorf("queried none results with %s", address)
	} else if total > 1 {
		return "", fmtp.Errorf("queried more results with %s", address)
	}

	return allEips[0].ID, nil
}

// CheckDeleted checks the error to see if it's a 404 (Not Found) and, if so,
// sets the resource ID to the empty string instead of throwing an error.
func CheckDeleted(d *schema.ResourceData, err error, msg string) error {
	if _, ok := err.(golangsdk.ErrDefault404); ok {
		d.SetId("")
		return nil
	}

	return fmtp.Errorf("%s: %s", msg, err)
}

// CheckDeletedDiag checks the error to see if it's a 404 (Not Found) and, if so,
// sets the resource ID to the empty string instead of throwing an error.
func CheckDeletedDiag(d *schema.ResourceData, err error, msg string) diag.Diagnostics {
	var statusCode int

	// check if the error is raised by **golangsdk**
	if _, ok := err.(golangsdk.ErrDefault404); ok {
		statusCode = http.StatusNotFound
	} else if responseErr, ok := err.(*sdkerr.ServiceResponseError); ok {
		// check if the error is raised by **huaweicloudstack-sdk-go-v3**
		statusCode = responseErr.StatusCode
	}

	if statusCode == http.StatusNotFound {
		resourceID := d.Id()
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("the resource %s is gone and will be removed in Terraform state.", resourceID),
			},
		}
	}

	return fmtp.DiagErrorf("%s: %s", msg, err)
}

// UnsubscribePrePaidResource impl the action of unsubscribe resource
func UnsubscribePrePaidResource(d *schema.ResourceData, config *config.Config, resourceIDs []string) error {
	bssV2Client, err := config.BssV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud bss V2 client: %s", err)
	}

	unsubscribeOpts := orders.UnsubscribeOpts{
		ResourceIds:     resourceIDs,
		UnsubscribeType: 1,
	}
	_, err = orders.Unsubscribe(bssV2Client, unsubscribeOpts).Extract()
	return err
}

func CheckForRetryableError(err error) *resource.RetryError {
	switch errCode := err.(type) {
	case golangsdk.ErrDefault500:
		return resource.RetryableError(err)
	case golangsdk.ErrUnexpectedResponseCode:
		switch errCode.Actual {
		case 409, 503:
			return resource.RetryableError(err)
		default:
			return resource.NonRetryableError(err)
		}
	default:
		return resource.NonRetryableError(err)
	}
}

func WaitOrderComplete(ctx context.Context, client *golangsdk.ServiceClient, orderId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"3", "6"}, // 3: Processing; 6: Pending payment.
		Target:       []string{"5"},      // 5: Completed.
		Refresh:      refreshOrderStatusFunc(client, orderId),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the order (%s) to complete payment: %#v", orderId, err)
	}
	return nil
}

func refreshOrderStatusFunc(client *golangsdk.ServiceClient, orderId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := orders.Get(client, orderId).Extract()
		if err != nil {
			return nil, "Error", err
		}
		return r, strconv.Itoa(r.OrderInfo.Status), nil
	}
}

// WaitOrderResourceComplete is the method to wait for the resource to be generated.
// Notes: Note that this method needs to be used in conjunction with method "WaitOrderComplete", because the ID of some
// resources may not be generated when the order is not completed.
func WaitOrderResourceComplete(ctx context.Context, client *golangsdk.ServiceClient, orderId string,
	timeout time.Duration) (string, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DONE"},
		Refresh:      refreshOrderResourceStatusFunc(client, orderId),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	res, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "", fmt.Errorf("error while waiting for the order (%s) to complete: %#v", orderId, err)
	}

	r := res.(resources.Resource)
	return r.ResourceId, nil
}

func refreshOrderResourceStatusFunc(client *golangsdk.ServiceClient, orderId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		listOpts := resources.ListOpts{
			OrderId:          orderId,
			OnlyMainResource: 1,
		}
		resp, err := resources.List(client, listOpts)
		if err != nil || resp == nil {
			return nil, "ERROR", fmt.Errorf("error waiting for the order (%s) to complete: %#v", orderId, err)
		}
		if resp.TotalCount < 1 {
			return nil, "PENDING", nil
		}
		return resp.Resources[0], "DONE", nil
	}
}

func CaseInsensitiveFunc() schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return strings.EqualFold(old, new)
	}
}

// GetAutoPay is a method to return whether order is auto pay according to the user input.
// auto_pay parameter inputs and returns:
//
//	false: false
//	true, empty: true
//
// Before using this function, make sure the parameter behavior is auto pay (the default value is "true").
func GetAutoPay(d *schema.ResourceData) string {
	if val, ok := d.GetOk("auto_pay"); ok && val.(string) == "false" {
		return "false"
	}
	return "true"
}

func UpdateAutoRenew(c *golangsdk.ServiceClient, enabled, resourceId string) error {
	if enabled == "true" {
		return resources.EnableAutoRenew(c, resourceId)
	}
	return resources.DisableAutoRenew(c, resourceId)
}

func HasFilledOpt(d *schema.ResourceData, param string) bool {
	_, b := d.GetOk(param)
	return b
}
