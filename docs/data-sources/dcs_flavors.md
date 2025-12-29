---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_dcs_flavors"
description: |-
  Use this data source to get a list of available DCS flavors.
---

# hcs_dcs_flavors

Use this data source to get a list of available DCS flavors.

## Example Usage

```hcl
data "hcs_dcs_flavors" "flavors" {
  capacity = "4"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the DCS flavors.
  If omitted, the provider-level region will be used.

* `capacity` - (Optional, Float) The total memory of the cache, in GB.
  + **Redis4.0, Redis5.0 and Redis6.0**: Stand-alone and active/standby type instance values:
    `0.125`, `0.25`, `0.5`, `1`, `2`, `4`, `8`, `16`, `32` and `64`.
    Cluster instance specifications support `4`,`8`,`16`,`24`, `32`, `48`, `64`, `96`, `128`, `192`,
    `256`, `384`, `512`, `768` and `1024`.
  + **Memcached**: Stand-alone and active/standby type instance values: `2`, `4`, `8`, `16`, `32` and `64`.

* `engine` - (Optional, String) The engine of the cache instance. Defaults to **Redis**.  
  The valid values are as follows:  
  + **Redis**
  + **Memcached**

* `engine_version` - (Optional, String) The version of a cache engine. This parameter is **Required** when the engine
  is **Redis**. The valid values are as follows:  
  + **3.0**
  + **4.0**
  + **5.0**
  + **6.0**

* `cache_mode` - (Optional, String) The mode of a cache engine. The valid values are as follows:
  + `single` - Single-node.
  + `ha` - Master/Standby.
  + `cluster` - Redis Cluster.
  + `proxy` - Proxy Cluster. Redis6.0 not support this mode.
  + `ha_rw_split` - Read/Write splitting. Redis6.0 not support this mode.
  
* `name` - (Optional, String) The flavor name of the cache instance.

* `cpu_architecture` - (Optional, String) The CPU architecture of cache instance.
  The valid values are as follows:  
  + **x86_64**
  + **aarch64**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `flavors` - The list of DCS flavors. The [flavors](#rds_flavors_attr) object structure is documented below.

<a name="rds_flavors_attr"></a>
The `flavors` block supports:

* `name` - The flavor name of the cache instance.

* `cache_mode` - The mode of a cache instance.

* `engine` - The engine of the cache instance.
  + **Redis**
  + **Memcached**

* `engine_versions` - The versions of the specification.

* `cpu_architecture` - The CPU architecture of cache instance.
  + **x86_64**
  + **aarch64**

* `capacity` - The total memory of the cache, in GB.

* `available_zones` - The list of available zones where the cache specification can be used.

* `ip_count` - The number of IP addresses corresponding to the specifications.
