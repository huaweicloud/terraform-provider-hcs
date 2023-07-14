/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package premium_domains

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("host")
}

func resourceURL(c *golangsdk.ServiceClient, hostId string) string {
	return c.ServiceURL("host", hostId)
}

func protectStatusURL(c *golangsdk.ServiceClient, hostId string) string {
	return c.ServiceURL("host", hostId, "protect-status")
}
