---
subcategory: "Web Application Firewall (WAF)"
---

# hcs_waf_dedicated_domain

Manages a dedicated mode domain resource within HuaweiCloudStack.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The dedicated mode domain name resource can be used in Dedicated Mode and ELB Mode.

## Example Usage

```hcl
variable "certificated_id" {}
variable "vpc_id" {}
variable "enterprise_project_id" {}

resource "hcs_waf_dedicated_domain" "domain_1" {
  domain                = "www.example.com"
  certificate_id        = var.certificated_id
  enterprise_project_id = var.enterprise_project_id
  protect_status        = 1

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "192.168.1.100"
    port            = 8080
    type            = "ipv4"
    vpc_id          = var.vpc_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the dedicated mode domain resource. If omitted,
  the provider-level region will be used. Changing this setting will push a new domain.

* `domain` - (Required, String, ForceNew) Specifies the protected domain name or IP address (port allowed). For example,
  `www.example.com` or `*.example.com` or `www.example.com:89`. Changing this creates a new domain.

* `server` - (Required, List, ForceNew) The server configuration list of the domain. A maximum of 80 can be configured.
  The [server](#waf_server) object structure is documented below. Changing this creates a new domain.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of WAF dedicated domain.
  Changing this parameter will create a new resource.

* `certificate_id` - (Optional, String) Specifies the certificate ID. This parameter is mandatory when `client_protocol`
  is set to HTTPS.

* `policy_id` - (Optional, String) Specifies the policy ID associated with the domain. If not specified, a new policy
  will be created automatically.

* `proxy` - (Optional, Bool) Specifies whether a proxy is configured. Default value is `false`.

  -> **NOTE:** WAF forwards only HTTP/S traffic. So WAF cannot serve your non-HTTP/S traffic, such as UDP, SMTP, FTP,
  and basically all other non-HTTP/S traffic. If a proxy such as public network ELB (or Nginx) has been used, set
  proxy `true` to ensure that the WAF security policy takes effect for the real source IP address.

* `keep_policy` - (Optional, Bool) Specifies whether to retain the policy when deleting a domain name.
  Defaults to `true`.

* `protect_status` - (Optional, Int) The protection status of domain, `0`: suspended, `1`: enabled.
  Default value is `0`.

* `tls` - (Optional, String) Specifies the minimum required TLS version. The options include `TLS v1.0`, `TLS v1.1`,
  `TLS v1.2`.

* `cipher` - (Optional, String) Specifies the cipher suite of domain. The options include `cipher_1`, `cipher_2`,
  `cipher_default`.

* `pci_3ds` - (Optional, Bool) Specifies the status of the PCI 3DS compliance certification check. The options
  include `true` and `false`. This parameter must be used together with tls and cipher.

  -> **NOTE:** Tls must be set to TLS v1.2, and cipher must be set to cipher_2. The PCI 3DS compliance certification
  check cannot be disabled after being enabled.

* `pci_dss` - (Optional, Bool) Specifies the status of the PCI DSS compliance certification check. The options
  include `true` and `false`. This parameter must be used together with tls and cipher.

  -> **NOTE:** Tls must be set to TLS v1.2, and cipher must be set to cipher_2.

* `connection_protection` - (Optional, List) Specifies the connection protection configuration to let WAF protect your
  origin servers from being crashed when WAF detects a large number of 502/504 error codes or pending requests.
  Only supports one protection configuration.
  The [connection_protection](#DedicatedDomain_connection_protection) structure is documented below.

<a name="waf_server"></a>
The `server` block supports:

* `client_protocol` - (Required, String, ForceNew) Protocol type of the client. The options include `HTTP` and `HTTPS`.
  Changing this creates a new service.

* `server_protocol` - (Required, String, ForceNew) Protocol used by WAF to forward client requests to the server. The
  options include `HTTP` and `HTTPS`. Changing this creates a new service.

* `vpc_id` - (Required, String, ForceNew) The id of the vpc used by the server. Changing this creates a service.

* `type` - (Required, String, ForceNew) Server network type, IPv4 or IPv6. Valid values are: `ipv4` and `ipv6`. Changing
  this creates a new service.

* `address` - (Required, String, ForceNew) IP address or domain name of the web server that the client accesses. For
  example, `192.168.1.1` or `www.example.com`. Changing this creates a new service.

* `port` - (Required, Int, ForceNew) Port number used by the web server. The value ranges from 0 to 65535. Changing this
  creates a new service.

<a name="DedicatedDomain_connection_protection"></a>
The `connection_protection` block supports:

* `error_threshold` - (Optional, Int) Specifies the 502/504 error threshold for every 30 seconds. Valid value ranges
  from `0` to `2,147,483,647`.

* `error_percentage` - (Optional, Float) Specifies the 502/504 error percentage. A breakdown protection is triggered
  when the 502/504 error threshold and percentage threshold have been reached. Valid value ranges from `0` to `99`.

* `initial_downtime` - (Optional, Int) Specifies the breakdown duration (s) when the breakdown is triggered for the
  first time. Valid value ranges from `0` to `2,147,483,647`.

* `multiplier_for_consecutive_breakdowns` - (Optional, Int) Specifies the maximum multiplier for consecutive breakdowns
  that occur within an hour. Valid value ranges from `0` to `2,147,483,647`.
  For example: Assume that you set the initial downtime to `180s` and the maximum multiplier to `3`. If the breakdown
  protection is triggered for the second time, the website downtime is 360s (180s X 2).
  If the breakdown protection is triggered for the third or fourth time, the website downtime is 540s (180s x 3).
  The breakdowns are calculated every one hour.

* `pending_url_request_threshold` - (Optional, Int) Specifies the pending URL request threshold. Connection protection
  is triggered when the number of read URL requests reaches the threshold you configure. Valid value ranges from `0` to
  `2,147,483,647`.

* `duration` - (Optional, Int) Specifies the protection duration (s) for connection protection. During this period, WAF
  stops forwarding website requests. Valid value ranges from `0` to `2,147,483,647`.

* `status` - (Optional, Bool) Specifies whether to enable connection protection. Defaults to **false**.

## Attribute Reference

The following attributes are exported:

* `id` - ID of the domain.

* `certificate_name` - The name of the certificate used by the domain name.

* `access_status` - Whether a domain name is connected to WAF. Valid values are:
    + `0` - The domain name is not connected to WAF,
    + `1` - The domain name is connected to WAF.

* `protocol` - The protocol type of the client. The options are `HTTP` and `HTTPS`.

* `compliance_certification` - The compliance certifications of the domain, values are:
    + `pci_dss` - The status of domain PCI DSS, `true`: enabled, `false`: disabled.
    + `pci_3ds` - The status of domain PCI 3DS, `true`: enabled, `false`: disabled.

* `alarm_page` - The alarm page of domain. Valid values are:
    + `template_name` - The template of alarm page, values are: `default`, `custom` and `redirection`.
    + `redirect_url` - The redirection URL when `template_name` is set to `redirection`.

* `traffic_identifier` - The traffic identifier of domain. Valid values are:
    + `ip_tag` - The IP tag of traffic identifier.
    + `session_tag` - The session tag of traffic identifier.
    + `user_tag` - The user tag of traffic identifier.

## Import

There are two ways to import WAF dedicated domain state.

* Using the `id`, e.g.

```bash
$ terraform import hcs_waf_dedicated_domain.test <id>
```

* Using `id` and `enterprise_project_id`, separated by a slash, e.g.

```bash
$ terraform import hcs_waf_dedicated_domain.test <id>/<enterprise_project_id>
```
