---
subcategory: "ServiceStage"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_servicestagev3_component"
description: |-
  Manages a component resource within HuaweiCloudStack.
---

# hcs_servicestagev3_component

Manages a component resource within HuaweiCloudStack.

## Example Usage

```hcl
variable "component_name" {}
variable "application_id" {}
variable "environment_id" {}
variable "obs_package_url" {}
variable "associated_ecs_id" {}

resource "hcs_servicestagev3_component" "test" {
  name           = var.component_name
  application_id = var.application_id
  environment_id = var.environment_id
  description    = "Created by terraform script"
  version        = "1.2.3"
  replica        = 1

  runtime_stack {
    name        = "OpenJDK17"
    version     = "1.4.6"
    type        = "Java"
    deploy_mode = "virtualmachine"
  }

  source = jsonencode({
    "kind" : "package",
    "storage" : "obs",
    "url" : var.obs_package_url,
  })

  tomcat_opts = jsonencode({
    "server_xml" = ""
  })

  refer_resources {
    id   = var.associated_ecs_id
    type = "ecs"
  }

  auto_lts_config {
    enable       = true
    lts_log_path = ["/root/*", "/tmp/*"]
  }
  
  tags = {
    foo = "bar"
    key = "value"
  }

  lifecycle {
    ignore_changes = [refer_resources, source]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the component is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `application_id` - (Required, String, ForceNew) Specifies the application ID to which the component belongs.  
  Changing this will create a new resource.

* `environment_id` - (Required, String, ForceNew) Specifies the environment ID where the component is deployed.  
  Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the component.  
  The valid length is limited from `2` to `64`, only letters, digits, hyphens (-) and underscores (_) are allowed.
  The name must start with a letter and end with a letter or a digit.  
  Changing this will create a new resource.

* `runtime_stack` - (Required, List, ForceNew) Specifies the configuration of the runtime stack.  
  The [runtime_stack](#servicestage_v3_component_runtime_stack) structure is documented below.  
  Changing this will create a new resource.

* `source` - (Required, String) Specifies the source configuration of the component, in JSON format.  
  The valid keys are as follows:
  + **kind**. The valid values are **package**, **image**.
  + **version**. 
  + **url**. 
  + **storage**. The valid values are **devcloud**, **obs**, **http**, **swr**.
  + **http_username**. This parameter can only be used in HCS **8.5.0** and **later** version.
  + **http_password**. This parameter can only be used in HCS **8.5.0** and **later** version.
  + **http_command**. This parameter can only be used in HCS **8.5.0** and **later** version.
  + **http_headers**. This parameter can only be used in HCS **8.5.0** and **later** version.

* `version` - (Required, String) Specifies the version of the component.  
  The format is **{number}.{number}.{number}** or **{number}.{number}.{number}.{number}**, e.g. **1.0.1**.

* `refer_resources` - (Required, List) Specifies the configuration of the reference resources.  
  The [refer_resources](#servicestage_v3_component_refer_resources) structure is documented below.

* `workload_content` - (Optional, String) Specifies the workload content of the component.

* `config_mode` - (Optional, String) Specifies the config mode of the component.
  + **ui**. This value can only be used in HCS **8.5.0** and **later** version.
  + **yaml**. This value can only be used in HCS **8.5.0** and **later** version.
  + **template**. This value can only be used in HCS **8.5.1** and **later** version.

* `description` - (Optional, String) Specifies the description of the component.  
  The value can contain a maximum of `128` characters.

  -> The value of the `description` cannot be set to empty value by updating.

* `build` - (Optional, String) Specifies the build configuration of the component, in JSON format.  
  The valid keys are as follows:
  + **parameters**
    - **artifact_namespace**
    - **cluster_id**
    - **node_label_selector**
    - **environment_id**

* `replica` - (Optional, Int, ForceNew) Specifies the replica number of the component.  
  Changing this will create a new resource.

* `limit_cpu` - (Optional, Float) Specifies the maximum number of the CPU limit.  
  The unit is **Core**.

* `limit_memory` - (Optional, Float) Specifies the maximum number of the memory limit.  
  The unit is **GiB**.

* `request_cpu` - (Optional, Float) Specifies the number of the CPU request resources.  
  The unit is **Core**.

* `request_memory` - (Optional, Float) Specifies the number of the memory request resources.  
  The unit is **GiB**.

* `envs` - (Optional, List) Specifies the configuration of the environment variables.  
  The [envs](#servicestage_v3_component_envs) structure is documented below.

* `storages` - (Optional, List) Specifies the storage configuration.  
  The [storages](#servicestage_v3_component_storages) structure is documented below.

* `deploy_strategy` - (Optional, List) Specifies the configuration of the deploy strategy.  
  The [deploy_strategy](#servicestage_v3_component_deploy_strategy) structure is documented below.

* `update_strategy` - (Optional, String) Specifies the configuration of the update strategy, in JSON format.

  -> **NOTE:** This parameter can only be used in HCS **8.5.0** and **later** version.

* `command` - (Optional, String) Specifies the start commands of the component, in JSON format.  
  The valid keys are as follows:
  + **command**
  + **args**

* `post_start` - (Optional, List) Specifies the post start configuration.  
  The [post_start](#servicestage_v3_component_lifecycle) structure is documented below.

* `pre_stop` - (Optional, List) Specifies the pre stop configuration.  
  The [pre_stop](#servicestage_v3_component_lifecycle) structure is documented below.

* `mesher` - (Optional, List) Specifies the configuration of the access mesher.  
  The [mesher](#servicestage_v3_component_mesher) structure is documented below.

* `timezone` - (Optional, String) Specifies the time zone in which the component runs, e.g. **Asia/Shanghai**.

* `jvm_opts` - (Optional, String) Specifies the JVM parameters of the component. e.g. **-Xms256m -Xmx1024m**.  
  If there are multiple parameters, separate them by spaces.  
  If this parameter is left blank, the default value is used.

* `tomcat_opts` - (Optional, String) Specifies the configuration of the tomcat server, in JSON format.  
  The valid keys are as follows:
  + **server_xml**

* `is_system_logging` - (Optional, Bool) Whether to disable the log recording function of the framework.

  -> **NOTE:** This parameter can only be used in HCS **8.6.0** and **later** version.

* `auto_lts_config` - (Optional, List) Specifies the configuration of lts.  
  The [auto_lts_config](#servicestage_v3_component_auto_lts_config) structure is documented below.

  -> **NOTE:** This parameter can only be used in HCS **8.5.1** and **later** version.

* `logs` - (Optional, List) Specifies the configuration of the logs collection.  
  The [logs](#servicestage_v3_component_logs) structure is documented below.

* `custom_metric` - (Optional, List) Specifies the configuration of the monitor metric.  
  The [custom_metric](#servicestage_v3_component_custom_metric) structure is documented below.

* `affinity` - (Optional, List) Specifies the affinity configuration of the component.  
  The [affinity](#servicestage_v3_component_affinity) structure is documented below.

* `anti_affinity` - (Optional, List) Specifies the anti-affinity configuration of the component.  
  The [anti_affinity](#servicestage_v3_component_affinity) structure is documented below.

* `liveness_probe` - (Optional, List) Specifies the liveness probe configuration of the component.  
  The [liveness_probe](#servicestage_v3_component_probe) structure is documented below.

* `readiness_probe` - (Optional, List) Specifies the readiness probe configuration of the component.  
  The [readiness_probe](#servicestage_v3_component_probe) structure is documented below.

* `external_accesses` - (Optional, List) Specifies the configuration of the external accesses.  
  The [external_accesses](#servicestage_v3_component_external_accesses) structure is documented below.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the component.

<a name="servicestage_v3_component_runtime_stack"></a>
The `runtime_stack` block supports:

* `name` - (Required, String, ForceNew) Specifies the stack name.  
  Changing this will create a new resource.

* `type` - (Required, String, ForceNew) Specifies the stack type.  
  Changing this will create a new resource.

* `deploy_mode` - (Required, String, ForceNew) Specifies the deploy mode of the stack.  
  Changing this will create a new resource.

* `version` - (Optional, String, ForceNew) Specifies the stack version.  
  Changing this will create a new resource.

<a name="servicestage_v3_component_refer_resources"></a>
The `refer_resources` block supports:

* `id` - (Required, String) Specifies the resource ID.

* `type` - (Required, String) Specifies the resource type.

* `parameters` - (Optional, String) Specifies the resource parameters, in JSON format.  
  The valid keys are as follows:
  + **namespace**.
  + **name**.
  + **capacity**. This parameter can only be used in HCS **8.5.1** and **later** version.
  + **type**. This parameter can only be used in HCS **8.5.1** and **later** version.
  + **class**. This parameter can only be used in HCS **8.5.1** and **later** version.
  + **obs_volume_type**. This parameter can only be used in HCS **8.5.1** and **later** version.
  + **access_mode**. This parameter can only be used in HCS **8.5.1** and **later** version.
  + **availableZone**. This parameter can only be used in HCS **8.5.1** and **later** version.
  + **volume_type**. This parameter can only be used in HCS **8.5.1** and **later** version.

<a name="servicestage_v3_component_envs"></a>
The `envs` block supports:

* `name` - (Required, String) Specifies the name of the environment variable.

* `value` - (Optional, String) Specifies the value of the environment variable.

<a name="servicestage_v3_component_storages"></a>
The `storages` block supports:

* `type` - (Required, String) Specifies the type of the data storage.
  + **HostPath**: Host path for local disk mounting.
  + **EmptyDir**: Temporary directory for local disk mounting.
  + **ConfigMap**: Configuration item for local disk mounting.
  + **Secret**: Secrets for local disk mounting.
  + **PersistentVolumeClaim**: Cloud storage mounting.

* `name` - (Required, String) Specifies the name of the disk where the data is stored.  
  Only lowercase letters, digits, and hyphens (-) are allowed and must start and end with a lowercase letter or digit.

* `parameters` - (Required, String) Specifies the information corresponding to the specific types of data storage,
  in JSON format.  
  The valid keys are as follows:
  + **type**. The valid values are **HostPath**, **EmptyDir**, **ConfigMap**, **Secret** and **PersistentVolumeClaim**.
  + **name**.
  + **parameters**.
    - **path**
    - **name**
    - **default_mode**
    - **medium**
  + **mounts**.
    - **path**
    - **sub_path**
    - **read_only**

* `mounts` - (Required, List) Specifies the configuration of the disk mounts.  
  The [mounts](#servicestage_v3_component_storage_mounts) structure is documented below.

<a name="servicestage_v3_component_storage_mounts"></a>
The `mounts` block supports:

* `path` - (Required, String) Specifies the mount path.

* `sub_path` - (Required, String) Specifies the sub mount path.

* `read_only` - (Required, Bool) Specifies whether the disk mount is read-only.

<a name="servicestage_v3_component_deploy_strategy"></a>
The `deploy_strategy` block supports:

* `type` - (Required, String) Specifies the deploy type.
  + **OneBatchRelease**: Single-batch upgrade.
  + **RollingRelease**: Rolling deployment and upgrade.
  + **GrayRelease**: Dark launch upgrade.

* `rolling_release` - (Optional, String) Specifies the rolling release parameters, in JSON format.  
  Required if the `type` is **RollingRelease**.  
  The valid keys are as follows:
  + **batches**

* `gray_release` - (Optional, String) Specifies the gray release parameters, in JSON format.  
  Required if the `type` is **GrayRelease**.  
  The valid keys are as follows:
  + **type**. The valid values are **weight** and **content**
  + **first_batch_weight**
  + **first_batch_replica**
  + **remaining_batch**
  + **deployment_mode**. The valid values are **1**, **3**, **4**.
  + **replica_surge_mode**. The valid values are **mirror**, **extra**, **no_usage**.
  + **rule_match_mode**
  + **rules**
    - **type**
    - **key**
    - **value**
    - **condition**

<a name="servicestage_v3_component_lifecycle"></a>
The `post_start` and `pre_stop` blocks support:

* `type` - (Required, String) Specifies the processing method.
  + **http**
  + **command**

* `scheme` - (Optional, String) Specifies the HTTP request type.
  + **HTTP**
  + **HTTPS**

  This parameter is only available when the `type` is set to `http`.

* `host` - (Optional, String) Specifies the host (IP) of the lifecycle configuration.  
  If this parameter is left blank, the pod IP address is used.  
  This parameter is only available when the `type` is set to `http`.

* `port` - (Optional, String) Specifies the port number of the lifecycle configuration.
  This parameter is only available when the `type` is set to `http`.

* `path` - (Optional, String) Specifies the request path of the lifecycle configuration.
  This parameter is only available when the `type` is set to `http`.

* `command` - (Optional, List) Specifies the command list of the lifecycle configuration.
  This parameter is only available when the `type` is set to `command`.

<a name="servicestage_v3_component_mesher"></a>
The `mesher` block supports:

* `port` - (Required, Int) Specifies the process listening port.

<a name="servicestage_v3_component_auto_lts_config"></a>
The `logs` block supports:

* `enable` - (Required, String) Whether enable the LTS auto config.

* `lts_log_path` - (Optional, List) Specifies the log path of the LTS, e.g. **/tmp/***.

<a name="servicestage_v3_component_logs"></a>
The `logs` block supports:

* `log_path` - (Required, String) Specifies the log path of the container, e.g. **/tmp**.

* `rotate` - (Required, String) Specifies the interval for dumping logs.
  + **Hourly**
  + **Daily**
  + **Weekly**

* `host_path` - (Required, String) Specifies the mounted host path, e.g. **/tmp**.

* `host_extend_path` - (Required, String) Specifies the extension path of the host.
  + **None**: the extended path is not used.
  + **PodUID**: extend the host path based on the pod ID.
  + **PodName**: extend the host path based on the pod name.
  + **PodUID/ContainerName**: extend the host path based on the pod ID and container name.
  + **PodName/ContainerName**: extend the host path based on the pod name and container name.

<a name="servicestage_v3_component_custom_metric"></a>
The `custom_metric` block supports:

* `path` - (Required, String) Specifies the collection path, such as **./metrics**.

* `port` - (Required, Int) Specifies the collection port, such as **9090**.

* `dimensions` - (Required, String) Specifies the monitoring dimension, such as **cpu_usage**, **mem_usage** or
  **cpu_usage,mem_usage** (separated by a comma).

<a name="servicestage_v3_component_affinity"></a>
The `affinity` and `anti_affinity` blocks support:

* `condition` - (Required, String) Specifies the condition type of the (anti) affinity rule.

* `kind` - (Required, String) Specifies the kind of the (anti) affinity rule.

* `match_expressions` - (Required, List) Specifies the list of the match rules for (anti) affinity.  
  The [match_expressions](#servicestage_v3_component_affinity_match_expressions) structure is documented below.

* `weight` - (Optional, Int) Specifies the weight of the (anti) affinity rule.  
  The valid value is range from `1` to `100`.

<a name="servicestage_v3_component_affinity_match_expressions"></a>
The `match_expressions` block supports:

* `key` - (Required, String) Specifies the key of the match rule.

* `operation` - (Required, String) Specifies the operation of the match rule.

* `value` - (Required, String) Specifies the value of the match rule.

<a name="servicestage_v3_component_probe"></a>
The `liveness_probe` and `readiness_probe` blocks support:

* `type` - (Required, String) Specifies the type of the probe.
  + **http**
  + **tcp**
  + **command**

* `delay` - (Required, Int) Specifies the delay time of the probe.

* `timeout` - (Required, Int) Specifies the timeout of the probe.

* `scheme` - (Optional, String) Specifies the scheme type of the probe.
  + **HTTP**
  + **HTTPS**

  This parameter is only available when the `type` is set to `http`.

* `host` - (Optional, String) Specifies the host of the probe.  
  Defaults to pod ID, also custom IP address can be specified.  
  This parameter is only available when the `type` is set to `http`.

* `port` - (Optional, Int) Specifies the port of the probe.  
  This parameter is only available when the `type` is set to `tcp` or `http`.

* `path` - (Optional, String) Specifies the path of the probe.  
  This parameter is only available when the `type` is set to `http`.

* `command` - (Optional, List) Specifies the command list of the probe.  
  This parameter is only available when the `type` is set to `command`.

<a name="servicestage_v3_component_external_accesses"></a>
The `external_accesses` block supports:

* `protocol` - (Required, String) Specifies the protocol of the external access.

* `address` - (Optional, String) Specifies the address of the external access.

* `forward_port` - (Optional, Int) Specifies the forward port of the external access.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, in UUID format.

* `status` - The status of the component.
  + **RUNNING**
  + **PENDING**

* `created_at` - The creation time of the component, in RFC3339 format.

* `updated_at` - The latest update time of the component, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.

* `update` - Default is 20 minutes.

* `delete` - Default is 5 minutes.

## Import

Components can be imported using `application_id` and `id` separated by a slash e.g.

```bash
$ terraform import hcs_servicestagev3_component.test <application_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to attributes missing from the API
response, security or some other reason.
The missing attribute is `workload_content`, `tags`, `refer_resources`, `source`.
It is generally recommended running `terraform plan` after importing resource.
You can decide if changes should be applied to resource, or the definition should be updated to align with the resource.
Also you can ignore changes as below.

```hcl
resource "hcs_servicestagev3_component" "test" {
  ...

  lifecycle {
    ignore_changes = [
      refer_resources, source,
    ]
  }
}
```
