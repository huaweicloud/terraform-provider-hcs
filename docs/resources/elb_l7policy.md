---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
---

# hcs_elb_l7policy

Manages an ELB L7 Policy resource within HCS.

## Example Usage

### Redirect To Pool

```hcl
variable listener_id {}
variable pool_id {}

resource "hcs_elb_l7policy" "policy_1" {
  name             = "policy_1"
  description      = "test description"
  listener_id      = var.listener_id
  action           = "REDIRECT_TO_POOL"
  redirect_pool_id = var.pool_id
}
```

### Redirect To URL

```hcl
variable listener_id {}
variable pool_id {}

resource "hcs_elb_l7policy" "policy_1" {
  name          = "policy_1"
  description   = "test description"
  listener_id   = var.listener_id
  action        = "REDIRECT_TO_URL"
  redirect_url_config {
    status_code = "302"
    protocol    = "$${protocol}"
    host        = "$${host}"
    port        = "20999"
    path        = "$${path}"
    query       = "$${query}&addition=info"
  }
}
```

Note: the first symbol in $$ is an escape character.

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the L7 Policy resource. If omitted, the
  provider-level region will be used. Changing this creates a new L7 Policy.

* `name` - (Optional, String) Human-readable name for the L7 Policy. Does not have to be unique.

* `description` - (Optional, String) Human-readable description for the L7 Policy.

* `listener_id` - (Required, String, ForceNew) The Listener on which the L7 Policy will be associated with. Changing
  this creates a new L7 Policy.

* `action` - (Required, String) Action for the forwarding policy. Must be one of REDIRECT_TO_POOL, 
  REDIRECT_TO_LISTENER, REDIRECT_TO_URL.

* `redirect_pool_id` - (Optional, String) ID of the backend server group that traffic is forwarded to. 
  Must be specified when the value of action is REDIRECT_TO_POOL.

* `redirect_listener_id` - (Optional, String) ID of the listener to which requests are forwarded. Must be 
  specified when the value of action is REDIRECT_TO_LISTENER.

* `redirect_url_config` - (Optional, RedirectUrlConfig Object) URL that traffic is forwarded to. Must be specified
  when the value of action is REDIRECT_TO_URL.


### RedirectUrlConfig Object

The following arguments are supported:

* `status_code` - (Required during creation | Optional during update, String) Return code after redirection. The value can be 301, 302, 303, 307, or 308.
  The value consists of 1 to 16 characters. 

* `protocol` - (Optional, String) Redirection protocol. The value can be HTTP, HTTPS, or `${protocol}`. The default 
  value is `${protocol}`, which indicates that the original value is used. The value consists of 1 to 36 characters.

* `host` - (Optional, String) Name of the host to which requests are directed. The value consists of 1 to 128 
  characters and contains letters, digits, hyphens (-), and periods (.). It must start with a letter or digit.
  The default value is `${host}`, which indicates that the original value is used.

* `port` - (Optional, String) Port that packets are redirected to. The value ranges from 1 to 65535. The default 
  value is `${port}`, which indicates that the original value is inherited. Minimum length: 1 character, Maximum
  length: 16 characters.

* `path` - (Optional, String) Redirection path. The value consists of 1 to 128 characters and contains letters,
  digits, and the following special characters: `_~';@^-%#&$.*+?,=!:|/()[]{}`. The value must start with a
  slash (/). The default value is `${path}`, which indicates that the original value is used.

* `query` - (Optional, String) Query string for redirection. The value can contain only letters, digits, and
  special characters `` !$&'()*+,-./:;=?@^_` ``. The default value is `${query}`, which indicates that
  the original value is inherited. Minimum length: 0 characters, Maximum length: 128 characters.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the L7 policy.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `update` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

ELB policy can be imported using the policy ID, e.g.

```
$ terraform import hcs_elb_l7policy.policy_1 5c20fdad-7288-11eb-b817-0255ac10158b
```
