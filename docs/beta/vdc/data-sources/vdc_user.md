---
subcategory: "Enterprise Projects (EPS)"
---

# hcs_vdc_user

Use this data source to get the list of the vdc user.

## Example Usage

```hcl
data "hcs_vdc_user" "userList" {
    vdc_id = "xxx-xxx-xxx"
    name = "xxx"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Username. Fuzzy search by the name field is supported. The value can contain 1 to 128 characters.

* `vdc_id` - (Required, String) VDC ID, which is used to query the user list in the specified VDC. The value can contain 1 to 36 characters, including only lowercase letters, digits, and hyphens (-).


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - User ID.
  
* `domain_id` - Tenant ID.

* `name` - Username.

* `display_name` - Alias of a user.

* `enabled` - Whether a user is enabled. The value can be true (default value) or false. If the value is false, the user is disabled and cannot be used for login.

* `description` - User description.

* `vdc_id` - VDC ID.

* `ldap_id` - LDAP server configuration ID.

* `create_at` - Creation time.

* `auth_type` - User type. The value can be LOCAL_AUTH (local authentication), SAML_AUTH (SAML authentication), LDAP_AUTH (LDAP authentication), or MACHINE_USER (machine-machine user).

* `access_mode` - Access mode. The value can be default (console access or programming access), console (console access), or programmatic (programming access).

* `top_vdc_id` - ID of the first-level VDC.
