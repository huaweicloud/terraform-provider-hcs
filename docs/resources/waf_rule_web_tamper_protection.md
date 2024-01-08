---
subcategory: "Web Application Firewall (WAF)"
---

# hcs_waf_rule_web_tamper_protection

Manages a WAF web tamper protection rule resource within HuaweiCloudStack.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The web tamper protection rule resource can be used in Cloud Mode, Dedicated Mode and ELB Mode.

## Example Usage

```hcl
variable "enterprise_project_id" {}
variable "policy_id" {}

resource "hcs_waf_rule_web_tamper_protection" "rule_1" {
  policy_id             = var.policy_id
  enterprise_project_id = var.enterprise_project_id
  domain                = "www.your-domain.com"
  path                  = "/payment"
  description           = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the WAF web tamper protection rules resource. If
  omitted, the provider-level region will be used. Changing this creates a new rule.

* `policy_id` - (Required, String, ForceNew) Specifies the WAF policy ID. Changing this creates a new rule.

* `domain` - (Required, String, ForceNew) Specifies the domain name. Changing this creates a new rule.

* `path` - (Required, String, ForceNew) Specifies the URL protected by the web tamper protection rule, excluding a
  domain name. Changing this creates a new rule.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of WAF tamper protection
  rule. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The rule ID in UUID format.

## Import

There are two ways to import WAF rule web tamper protection state.

* Using `policy_id` and `rule_id`, separated by a slash, e.g.

```bash
$ terraform import hcs_waf_rule_web_tamper_protection.test <policy_id>/<rule_id>
```

* Using `policy_id`, `rule_id` and `enterprise_project_id`, separated by slashes, e.g.

```bash
$ terraform import hcs_waf_rule_web_tamper_protection.test <policy_id>/<rule_id>/<enterprise_project_id>
```
