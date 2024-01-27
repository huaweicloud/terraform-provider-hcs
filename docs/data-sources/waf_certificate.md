---
subcategory: "Web Application Firewall (WAF)"
---

# hcs_waf_certificate

Get the certificate in the WAF, including the one pushed from SCM.

## Example Usage

```hcl
variable enterprise_project_id {}

data "hcs_waf_certificate" "certificate_1" {
  name                  = "certificate name"
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to obtain the WAF. If omitted, the provider-level region will be
  used.

* `name` - (Required, String) The name of certificate. The value is case sensitive and supports fuzzy matching.

  -> **NOTE:** The certificate name is not unique. Only returns the last created one when matched multiple certificates.

* `expire_status` - (Optional, Int) The expire status of certificate. Defaults is `0`. The value can be:
  + `0`: not expire
  + `1`: has expired
  + `2`: wil expired soon

* `enterprise_project_id` - (Optional, String) The enterprise project ID of WAF certificate.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The certificate ID in UUID format.

* `expiration` - Indicates the time when the certificate expires.
