---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_dcs_account"
description: |-
  Manages a DCS account resource within HuaweiCloudStack.
---

# hcs_dcs_account

Manages a DCS account resource within HuaweiCloudStack.

## Example Usage

```hcl
variable "instance_id" {}

resource "hcs_dcs_account" "test" {
  instance_id      = var.instance_id
  account_name     = "user"
  account_role     = "read"
  account_password = "Terraform@123"
  description      = "add account"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the DCS instance.

  Changing this creates a new resource.

* `account_name` - (Required, String, ForceNew) Specifies the name of the account.

  Changing this creates a new resource.

* `account_password` - (Required, String) Specifies the password of the account.

* `account_role` - (Required, String) Specifies the role of the account.  
  The valid values are as follows:
  + **read**: The account has read-only privilege.
  + **write**: The account has read and write privilege.

* `description` - (Optional, String) Specifies the description of the account.

* `account_read_policy` - (Optional, String) Specifies the read requests for the specified account are forwarded to
  the primary node or secondary node. Default to **null**.
  + **master**: Read requests are forwarded to the primary node.
  + **replica**: Read requests are forwarded to the secondary node.
  + **master-replica**: Read requests are forwarded to the primary and secondary nodes.

  -> **Note** This parameter only supported by proxy clusters and read-write splitting instances.

  -> **Note** This parameter only can be used in HCS **8.6.0** and **later** version.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the ID of the account.

* `account_type` - Indicates the type of the account. The value can be **normal** or **default**.

* `status` - Indicates the status of the account.
  The valid values are as follows:  
  + **CREATING**
  + **AVAILABLE**
  + **CREATEFAILED**
  + **DELETED**
  + **DELETEFAILED**
  + **DELETING**
  + **UPDATING** 
  + **ERROR**

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The DCS account can be imported using the DCS instance ID and the DCS account ID separated by a slash, e.g.

```bash
$ terraform import hcs_dcs_account.test <instance_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `account_password`.
It is generally recommended running `terraform plan` after importing the account.
You can then decide if changes should be applied to the account, or the resource definition should be updated to
align with the account. Also you can ignore changes as below.

```hcl
resource "hcs_dcs_account" "test" {
    ...

  lifecycle {
    ignore_changes = [
      account_password,
    ]
  }
}
```
