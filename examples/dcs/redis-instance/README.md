# Create DCS Redis instances

In this example, we will create different types of Redis instances.
And show how to configure the backup policy, IP whitelist and tags.

The types of Redis instances are Single-node, Master/Standby, Redis Cluster, Proxy Cluster, and read-write separation.
Examples are as follows:

* Create a Single-node Redis instance.
* Create Master/Standby Redis instances.
* Create Proxy Cluster Redis instances.
* Create Cluster Redis instances, and configure Backup Policy, Whitelists and Tags.

To run, configure your HuaweiCloudStack provider as described in the
[document](https://registry.terraform.io/providers/huaweicloud/hcs/latest/docs).

This example assumes that you have created a random password. If you want to use key-pair and do not have one, please
visit the
[document](https://registry.terraform.io/providers/huaweicloud/hcs/latest/docs/resources/ecs_keypair)
to create a key-pair.

You can refer to the following configuration to write scripts according to your needs also.

## Configuration

### Configure Backup Policy

Configure a backup policy, I need to back up from 02:00 to 04:00 every day and keep the backup data for 3 days.

```hcl
resource "hcs_dcs_instance" "instance" {
  ...
  backup_policy {
    backup_type = "auto"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [4]
    save_days   = 1
  }
  ...
}
```

### Configure IP Whitelists

Configure IP Whitelists, only the IP in the whitelist can access Redis instances.

```hcl
resource "hcs_dcs_instance" "instance" {
  ...
  whitelists {
    group_name = "group_1"
    ip_address = ["192.168.1.0/24", "192.168.19.0/24"]
  }
  whitelists {
    group_name = "group_2"
    ip_address = ["10.11.3.0/24"]
  }
...
}
```

### Rename the original Redis Commands

Rename the original command to make redis more secure.

```hcl
resource "hcs_dcs_instance" "instance" {
  ...
  rename_commands = {
    command  = "command001"
    keys     = "keys001"
    flushall = "flushall001"
    flushdb  = "flushdb001"
    hgetall  = "hgetall001"
  }
  ...
}
```

### Configure Tags for Redis instance

```hcl
resource "hcs_dcs_instance" "instance" {
  ...
  tags = {
    "level": "A",
    "yourKey": ""  # the value of tag can be empty.
  }
  ...
}
```

## Usage

```shell
terraform init
terraform plan
terraform apply
terraform destroy
```

## Requirements

| Name             | Version   |
|------------------|-----------|
| terraform        | >= 0.12.0 |
| huaweicloudstack | >= 2.4.2  |
