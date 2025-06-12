---
subcategory: "Enterprise Projects (EPS)"
---

# hcs_vdc_user

Manages an vdc user within HuaweiCloudStack.

## Example Usage

```hcl
resource "hcs_vdc_user" "user01" {
  vdc_id = "xxx"
  name = "userName"
  password = "xxx"
  display_name = "displayName"
  enabled = true
  description = "descriptionInfo"
  auth_type = "LOCAL_AUTH"
  access_mode = "default"
}
```

## Argument Reference

The following arguments are supported:

* `vdc_id` - (Required, String) VDC ID. The value can contain 1 to 36 characters, including only lowercase letters, digits, and hyphens (-).

* `name` - (Required, String) Username. The value can contain only special characters @._-, letters (case-insensitive), and digits. It cannot start with op_svc, paas_op, or \ and cannot end with \. In addition, it cannot be admin, power_user, or guest. The value can contain 4 to 32 characters. If the third-party user is of the LDAP AD type, the username cannot exceed 20 characters.

* `password` - (Optional, String) Password. The value can contain at least three types of the following characters: uppercase letters, lowercase letters, digits, and special characters (except < and >). The password must contain special characters and cannot contain the username or username spelled backwards. The value can contain 8 to 32 characters.

* `display_name` - (Optional, String) Alias of a user. The value can contain 0 to 128 characters. The following special characters are not allowed: ><

* `auth_type` - (Optional, String) User type. The value can be LOCAL_AUTH (default value), SAML_AUTH, LDAP_AUTH, or MACHINE_USER.

* `enabled` - (Optional, Boolean) Whether a user is enabled. The value can be true (default value) or false. If the value is false, the user is disabled.

* `description` - (Optional, String) Description. The value cannot contain the following characters: >< The value can contain 0 to 255 characters.

* `access_mode` - (Optional, String) Access mode. The value can be default (default value. It indicates console access or programming access), console (console access), or programmatic (programming access).

* `ldap_id` - (Optional, String) LDAP server configuration ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - User ID.

## Import

Instances can be imported by their `id`. For example,

```
terraform import hcs_vdc_user.user01 xxx
```

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minute.
* `update` - Default is 5 minute.
* `delete` - Default is 5 minute.