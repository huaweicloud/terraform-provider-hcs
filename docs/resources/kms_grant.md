---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_kms_grant"
description: ""
---

# hcs_kms_grant

-> **NOTE:** This resource can only be used in HCS **8.5.0** and **later** version.

Users can create authorizations for other IAM users or accounts,
granting them permission to use their own master key (CMK),
and a maximum of 100 authorizations can be created under one master key.

## Example Usage

```HCL
variable "key_id" {}
variable "user_id" {}
variable "retiring_principal" {}

resource "hcs_kms_grant" "test" {
  key_id             = var.key_id
  type               = "user"
  grantee_principal  = var.user_id
  operations         = ["create-datakey", "encrypt-datakey"]
  retiring_principal = var.retiring_principal
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `key_id` - (Required, String, ForceNew) Specifies the ID of the KMS key.

  Changing this parameter will create a new resource.

* `grantee_principal` - (Required, String, ForceNew) Specifies the ID of the authorized user or account.

  Changing this parameter will create a new resource.

* `operations` - (Required, List, ForceNew) Specifies the list of granted operations. A value containing only
  **create-grant** is invalid.
  The valid values are as follows:
  + **create-datakey**
  + **create-datakey-without-plaintext**
  + **encrypt-datakey**
  + **decrypt-datakey**
  + **describe-key**
  + **create-grant**
  + **retire-grant**
  + **encrypt-data**
  + **decrypt-data**
  
  Changing this parameter will create a new resource.

* `name` - (Optional, String, ForceNew) Specifies the grant name.  
  It must be `1` to `255` characters long, start with a letter, and contain only letters (case-sensitive),
  digits, hyphens (-), underscores (_), and slash(/).

  Changing this parameter will create a new resource.

* `type` - (Optional, String, ForceNew) Specifies the authorization type. Defaults to **user**.
  The valid values are as follows:
  + **user**
  + **domain**

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `creator` - The ID of the creator.  

## Import

The kms grant can be imported using `key_id`, `grant_id`, separated by slashes, e.g.

```bash
$ terraform import hcs_kms_grant.test <key_id>/<grant_id>
```
