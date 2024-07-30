# Create a WAF dedicated mode domain

To run, configure your HuaweiCloudStack provider as described in the
[document](https://registry.terraform.io/providers/huaweicloud/hcs/latest/docs).

This example assumes that you have created a random password. If you want to use key-pair and do not have one, please
visit the
[document](https://registry.terraform.io/providers/huaweicloud/hcs/latest/docs/resources/ecs_keypair)
to create a key-pair.

## Usage

```shell
terraform init
terraform plan
terraform apply
terraform destroy
```

The creation of the WAF dedicated instance takes about 5 minutes. After the creation is successful, the WAF policy and
dedicated mode domain start to be created.

## Requirements

| Name             | Version   |
|------------------|-----------|
| terraform        | >= 0.12.0 |
| huaweicloudstack | >= 2.4.0  |
