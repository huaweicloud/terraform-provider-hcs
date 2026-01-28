---
subcategory: "Object Storage Service (OBS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_obs_bucket_acl"
description: |-
  Manages an OBS bucket acl resource within HuaweiCloudStack. 
---

# hcs_obs_bucket_acl

Manages an OBS bucket acl resource within HuaweiCloudStack.

-> **NOTE:** When creating or updating the OBS bucket acl, the original bucket acl will be overwritten. When deleting
the OBS bucket acl, the full permissions of the bucket owner will be set, and the other permissions will be removed.

## Example Usage

```hcl
variable "bucket" {}
variable "account1" {}
variable "account2" {}

resource "hcs_obs_bucket_acl" "test" {
  bucket = var.bucket

  owner_permission {
    access_to_bucket = ["READ", "WRITE"]
    access_to_acl    = ["READ_ACP", "WRITE_ACP"]
  }

  account_permission {
    access_to_bucket = ["READ", "WRITE"]
    access_to_acl    = ["READ_ACP", "WRITE_ACP"]
    account_id       = var.account1
  }

  account_permission {
    access_to_bucket = ["READ"]
    access_to_acl    = ["READ_ACP", "WRITE_ACP"]
    account_id       = var.account2
  }

  public_permission {
    access_to_bucket = ["READ", "WRITE"]
  }

  log_delivery_user_permission {
    access_to_bucket = ["READ", "WRITE"]
    access_to_acl    = ["READ_ACP", "WRITE_ACP"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

  Changing this parameter will create a new resource.

* `bucket` - (Required, String, ForceNew) Specifies the name of the bucket to which to set the acl.

  Changing this parameter will create a new resource.

* `owner_permission` - (Optional, List) Specifies the bucket owner permission. If omitted, the current obs bucket acl
  owner permission will not be changed.

  The [permissions](#obs_bucket_acl_permissions) structure is documented below.

* `public_permission` - (Optional, List) Specifies the public permission.

  The [permissions](#obs_bucket_acl_permissions) structure is documented below.

* `log_delivery_user_permission` - (Optional, List) Specifies the log delivery user permission.

  The [permissions](#obs_bucket_acl_permissions) structure is documented below.

* `account_permission` - (Optional, List) Specifies the account permissions.

  The [account_permission](#obs_bucket_acl_account_permission) structure is documented below.

<a name="obs_bucket_acl_permissions"></a>
The `permissions` block supports:

* `access_to_bucket` - (Optional, List) Specifies the access to bucket.  
  The valid values are as follows:
  + **READ**
  + **WRITE**

* `access_to_acl` - (Optional, List) Specifies the access to acl.  
  The valid values are as follows:
  + **READ_ACP**. Allows the grantee to read the object's ACL property.
  + **WRITE_ACP**. Allows the grantee to update the object's ACL property.

<a name="obs_bucket_acl_account_permission"></a>
The `account_permission` block supports:

* `account_id` - (Required, String) Specifies the account id to authorize. The account id cannot be the bucket owner,
  and must be unique.

* `access_to_bucket` - (Optional, List) Specifies the access to bucket.  
  The valid values are as follows:
  + **READ**
  + **WRITE**

* `access_to_acl` - (Optional, List) Specifies the access to acl.  
  The valid values are as follows:
  + **READ_ACP**. Allows the grantee to read the object's ACL property.
  + **WRITE_ACP**. Allows the grantee to update the object's ACL property.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The name of the bucket.

## Import

The obs bucket acl can be imported using the `bucket`, e.g.

```bash
$ terraform import hcs_obs_bucket_acl.test <bucket>
```
