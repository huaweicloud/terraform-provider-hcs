# CCE Addon Templates

Addon support configuration input depending on addon type and version. This page contains description of addon
arguments. You can get up-to-date reference of addon arguments for your cluster using data source
[`hcs_cce_addon_template`](https://registry.terraform.io/providers/huaweicloud/hcs/latest/docs/data-sources/cce_addon_template)
.

Following addon templates exist in the addon template list:

- [`autoscaler`](#autoscaler)
- [`coredns`](#coredns)
- [`everest`](#everest)
- [`metrics-server`](#metrics-server)
- [`gpu-beta`](#gpu-beta)

All addons accept `basic` and some can accept `custom`, `flavor` input values.
It is recommended to use `basic_json`, `custom_json` and `flavor_json` for more flexible input.

## Example Usage

### Use basic_json, custom_json and flavor_json

```hcl
variable "cluster_id" {}
variable "tenant_id" {}

data "hcs_cce_addon_template" "autoscaler" {
  cluster_id = var.cluster_id
  name       = "autoscaler"
  version    = "1.25.21"
}

resource "hcs_cce_addon" "autoscaler" {
  cluster_id    = var.cluster_id
  template_name = "autoscaler"
  version       = "1.25.21"

  values {
    basic_json  = jsonencode(jsondecode(data.hcs_cce_addon_template.autoscaler.spec).basic)
    custom_json = jsonencode(merge(
      jsondecode(data.hcs_cce_addon_template.autoscaler.spec).parameters.custom,
      {
        cluster_id = var.cluster_id
        tenant_id  = var.tenant_id
      }
    ))
    flavor_json = jsonencode(jsondecode(data.hcs_cce_addon_template.autoscaler.spec).parameters.flavor2)
  }
}

```