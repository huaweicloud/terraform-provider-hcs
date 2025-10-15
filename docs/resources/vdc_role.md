---
subcategory: "Virtual Data Center (VDC)"
---

# hcs_vdc_role

Manages a VDC role resource within Huawei Cloud Stack.

-> **NOTE:** Supported from ManageOne version 8.6.0 onwards.

## Example Usage

```hcl
variable "role_name" {}

resource "hcs_vdc_role" "test" {
  name   = var.role_name
  type   = "AX"
  policy = <<EOF
    {
      "Depends": [],
      "Statement": [
        {
          "Action": [
            "moscapp:home:console",
            "moscapp:system:console",
            "moscapp:resource-tenant:console",
            "moscapp:ai-dashboard:console",
            "moscapp:organization-vdc:console",
            "moscapp:resource-cloud-resources:console",
            "moscapp:ai-dashboard:list",
            "moscapp:organization-vdc:list",
            "moscapp:resource-cloud-resources:list"
          ],
          "Effect": "Deny"
        }
      ],
      "Version": "1.1"
    }
    EOF
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) Specifies the name of the VDC custom role.

* `description` - (Optional, string) Specifies the description of the VDC custom role.

* `type` - (Optional, string, ForceNew) Specifies the display mode of the VDC custom role. Valid options are as follows:

    * AX: Global services.
    * XA: Regional services.

* `policy` - (Required, String) Specifies the content of the VDC custom role in JSON format.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The VDC custom role ID.

## import

VDC roles can be imported using the `id`, e.g.

```bash
$ terraform import hcs_vdc_role.role1 fa163eebcccbe1c10baa324fc930c75a
```
