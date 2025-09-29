---
subcategory: "VDC"
---

# hcs_vdc_user

Manages a VDC user within Huawei Cloud Stack.

## Example Usage

```hcl
variable "vdc_id" {}

variable "user_password" {}

resource "hcs_vdc_user" "user01" {
  vdc_id   = var.vdc_id
  name     = "Username"
  password = var.user_password
}
```

## Argument Reference

The following arguments are supported:

* `vdc_id` - (Required, String) VDC ID. The value can contain 1 to 36 characters, including only lowercase letters,
  digits, and hyphens (-). Once set, the value cannot be modified.

* `name` - (Required, String) Username. The value can contain only special characters @._-\, letters (case-insensitive),
  and digits. It cannot start with `op_svc`, `paas_op`, or \ and cannot end with \. In addition, it cannot be `admin`,
  `power_user`, or `guest`. The value can contain 4 to 32 characters. If the third-party user is of the LDAP AD type,
  the username cannot exceed 20 characters. Once set, the value cannot be modified.

* `password` - (Optional, String) Password. The value can contain at least three types of the following characters:
  uppercase letters, lowercase letters, digits, and special characters (except < and >). The password must contain
  special characters and cannot contain the username or username spelled backwards. The value can contain 8 to 32
  characters. For `LOCAL_AUTH` or `MACHINE_USER` authentication type users, the password cannot be empty.

* `display_name` - (Optional, String) Alias of a user. The value can contain 0 to 128 characters. The following special
  characters are not allowed: ><

* `auth_type` - (Optional, String) User type. The value can be `LOCAL_AUTH` (default value), `SAML_AUTH`, `LDAP_AUTH`,
  or `MACHINE_USER`. The parameter value cannot be changed. Exception: If access_mode is set to `programmatic`,
  auth_type must be set to `MACHINE_USER`.
  > [!NOTE]
  >
  > * When auth_type is set to `MACHINE_USER`, access_mode must be set to `programmatic`.
  > * When auth_type is set to `LOCAL_AUTH` or `MACHINE_USER`, you need to specify a value for password.
  > * When auth_type is set to `SAML_AUTH` or `LDAP_AUTH`, do not specify password.

* `enabled` - (Optional, Boolean) Whether a user is enabled. The value can be `true` (default value) or `false`. If the
  value is `false`, the user is disabled.

* `description` - (Optional, String) Description. The value cannot contain the following characters: >< The value can
  contain 0 to 255 characters.

* `access_mode` - (Optional, String) Access mode. The value can be `default` (default value. It indicates console access
  or programming access), `console` (console access), or `programmatic` (programming access).
  > [!NOTE]
  >
  > * When access_mode is set to `programmatic`, auth_type must be set to `MACHINE_USER`.
  > * If access_mode is set to `programmatic`, it cannot be changed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - User ID.

## Import

Users can be imported using the id, e.g.

```
terraform import hcs_vdc_user.user02 ed35bb2dada543d5977069780e98b2c3
```
