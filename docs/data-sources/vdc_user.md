---
subcategory: "Virtual Data Center (VDC)"
---

# hcs_vdc_user

Use this data source to obtain VDC user information.

## Example Usage

```hcl
variable "vdc_id" {}

data "hcs_vdc_user" "user" {
  vdc_id = var.vdc_id
  name   = "Username"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Username.

* `vdc_id` - (Required, String) VDC ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - User ID.

* `domain_id` - Tenant ID.

* `display_name` - Alias of a user.

* `enabled` - Whether a user is enabled. The value can be `true` (default value) or `false`. If the value is `false`,
  the user is disabled and cannot be used for login.

* `description` - User description.

* `ldap_id` - LDAP server configuration ID.

* `create_at` - Time when the user was created.

* `auth_type` - User type. The value can be `LOCAL_AUTH` (local authentication), `SAML_AUTH` (SAML authentication),
  `LDAP_AUTH` (LDAP authentication), or `MACHINE_USER` (machine-machine user).

* `access_mode` - Access mode. The value can be `default` (console access or programming access), `console` (console
  access), or `programmatic` (programming access).

* `top_vdc_id` - ID of the first-level VDC.
