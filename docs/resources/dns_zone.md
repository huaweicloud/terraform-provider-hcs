---
subcategory: "Domain Name Service (DNS)"
---

# hcs_dns_zone

Manages a DNS zone in the HuaweiCloudStack DNS Service.

## Example Usage

### Create a public DNS zone

```hcl
resource "hcs_dns_zone" "my_public_zone" {
  name        = "example.com."
  email       = "jdoe@example.com"
  description = "An example zone"
  ttl         = 3000
  zone_type   = "public"
}
```

### Create a private DNS zone

```hcl
resource "hcs_dns_zone" "my_private_zone" {
  name        = "1.example.com."
  email       = "jdoe@example.com"
  description = "An example zone"
  ttl         = 3000
  zone_type   = "private"

  router {
    router_id = "2c1fe4bd-ebad-44ca-ae9d-e94e63847b75"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the DNS zone. If omitted, the `region` argument
  of the provider will be used. Changing this creates a new DNS zone.

* `name` - (Required, String, ForceNew) The name of the zone. Note the `.` at the end of the name. Changing this creates
  a new DNS zone.

* `email` - (Optional, String) The email address of the administrator managing the zone.

* `zone_type` - (Optional, String, ForceNew) The type of zone. Can either be `public` or `private`. Changing this
  creates a new DNS zone.

* `router` - (Optional, List) Router configuration block which is required if zone_type is private. The router
  structure is documented below.

* `ttl` - (Optional, Int) The time to live (TTL) of the zone.

* `description` - (Optional, String) A description of the zone.

The `router` block supports:

* `router_id` - (Required, String) ID of the associated VPC.

* `router_region` - (Optional, String) The region of the VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `masters` - An array of master DNS servers.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `update` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

This resource can be imported by specifying the zone ID:

```
$ terraform import hcs_dns_zone.zone_1 <zone_id>
```
