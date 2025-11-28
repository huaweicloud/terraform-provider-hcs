---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_drs_availability_zones"
description: |-
  Use this data source to query availability zones where DRS jobs can be created within HuaweiCloudStack.
---

# hcs_drs_availability_zones

-> **NOTE:** This resource can only be used in HCS **8.5.0** and **later** version.

Use this data source to query availability zones where DRS jobs can be created within HuaweiCloudStack.

## Example Usage

```hcl
data "hcs_drs_availability_zones" "test" {
  engine_type = "mysql"
  type        = "migration"
  direction   = "up"
  node_type   = "high"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `engine_type` - (Required, String) Specifies the DRS job engine type. The options are as follows:
  + **mysql**: Used for data migration and data synchronization from MySQL to MySQL.
  + **mongodb**: Used for data migration from MongoDB to DDS.
  + **redis-to-gaussredis**: Used for data migration from Redis to GeminiDB Redis.
  + **rediscluster-to-gaussredis**: Used for data migration from Redis Cluster to GeminiDB Redis.
  + **cloudDataGuard-mysql**: Used for disaster recovery from MySQL to MySQL.
  + **gaussdbv5**: Used for data synchronization from GaussDB Distributed to GaussDB Distributed.
  + **mysql-to-kafka**: Used from data synchronization from MySQL to Kafka.
  + **taurus-to-kafka**: Used from data synchronization from TaurusDB to Kafka.
  + **gaussdbv5ha-to-kafka**: Used for data synchronization from GaussDB Centralized to Kafka.
  + **postgresql**: Used from data synchronization from PostgreSQL to PostgreSQL.
  + **oracle-to-gaussdbv5**: Used for data synchronization from Oracle to GaussDB Distributed.

* `type` - (Required, String) Specifies the job type. The options are as follows:
  + **migration**: Online Migration.
  + **sync**: Data Synchronization.
  + **cloudDataGuard**: Disaster Recovery.

* `direction` - (Required, String) Specifies the direction of data flow. The options are as follows:
  + **up**: To the cloud. The destination database must be a database in the current cloud.
  + **down**: Out of the cloud. The source database must be a database in the current cloud.
  + **non-dbs**: Self-built database.

* `node_type` - (Required, String) Specifies the node type of the job instance. The options are as follows:
  + **micro**: extremely small specification.
  + **small**: small specification.
  + **medium**: medium specification.
  + **high**: large specification.

* `multi_write` - (Optional, Bool) Specifies whether it is dual-AZ disaster recovery.

## Attribute Reference

In addition to all arguments above, the following attributes are supported:

* `id` - The data source ID.

* `names` - The names of availability zone.
