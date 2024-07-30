# Create a Kafka cluster

To run, configure your HuaweiCloudStack provider as described in the
[document](https://registry.terraform.io/providers/huaweicloud/hcs/latest/docs).

## The Kafka instance configuration

| Attributes    | Value                    |
|---------------|--------------------------|
| Instance Type | cluster                  |
| Version       | 2.7                      |
| Flavor        | redis.ha.xu1.tiny.r2.128 |
| Broker Number | 3                        |
| storage_space | 600 GB                   |
| IO Type       | ultra-high I/O           |

### FAQ

- **How to obtain the TPS, and number of partitions of the Kafka instances?**

  The **TPS**, and the **number of partitions** are included in the flavor and do not need to be
  specified when creating a resource.

  The expected `flavor_id` can be obtained in the following way.

  ```hcl
  data "hcs_dms_kafka_flavors" "test" {
    type      = "cluster"
    flavor_id = "redis.ha.xu1.tiny.r2.128"
  }
  ```

## Usage

```shell
terraform init
terraform plan
terraform apply
terraform destroy
```

It takes about 20 to 50 minutes to create a Kafka cluster depending on the flavor and broker number.

## Requirements

| Name             | Version   |
|------------------|-----------|
| terraform        | >= 0.12.0 |
| huaweicloudstack | >= 2.4.2  |
