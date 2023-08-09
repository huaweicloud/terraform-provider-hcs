# huaweicloudstack Provider

The huaweicloudstack provider is used to interact with the many resources supported by huaweicloudstack. The provider needs to be
configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

Terraform 0.13 and later:

```hcl
terraform {
  required_providers {
    huaweicloudstack = {
      source  = "huaweicloudstack/huaweicloudstack"
      version = "~> 2.3.0"
    }
  }
}

# Configure the huaweicloudstack Provider
provider "huaweicloudstack" {
  region     = "cn-north-4"
  cloud      = "mycloud.com"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}

# Create a VPC
resource "hcs_vpc" "example" {
  name = "my_vpc"
  cidr = "192.168.0.0/16"
}
```

Terraform 0.12 and earlier:

```hcl
# Configure the huaweicloudstack Provider
provider "huaweicloudstack" {
  version    = "~> 2.3.0"
  region     = "cn-north-4"
  cloud      = "mycloud.com"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}

# Create a VPC
resource "hcs_vpc" "example" {
  name = "my_vpc"
  cidr = "192.168.0.0/16"
}
```

## Authentication

The Huawei Cloud provider offers a flexible means of providing credentials for authentication. The following methods are
supported, in this order, and explained below:

* Static credentials
* Environment variables
* Shared configuration file
* ECS Instance Metadata Service

The Huawei Cloud Provider supports assuming role with IAM agency, either in the provider configuration
block parameter assume_role or shared configuration file.

### Static credentials

!> **WARNING:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage
should this file ever be committed to a public version control system.

Static credentials can be provided by adding an `access_key` and `secret_key`
in-line in the provider block:

Usage:

```hcl
provider "huaweicloudstack" {
  region     = "cn-north-4"
  cloud      = "mycloud.com"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}
```

### Environment variables

You can provide your credentials via the `HCS_ACCESS_KEY` and
`HCS_SECRET_KEY` environment variables, representing your Huawei Cloud Access Key and Secret Key, respectively.

```hcl
provider "huaweicloudstack" {}
```

Usage:

```sh
$ export HCS_ACCESS_KEY="anaccesskey"
$ export HCS_SECRET_KEY="asecretkey"
$ export HCS_REGION_NAME="cn-north-4"
$ export HCS_CLOUD="mycloud.com"
$ terraform plan
```

### Shared Configuration File

You can use a
[huaweicloudstack CLI configuration file](https://support.huaweicloudstack.com/intl/en-us/usermanual-hcli/hcli_03_002.html)
to specify your credentials. You need to specify a location in the Terraform configuration by providing the
`shared_configuration_file` argument or using the `HCS_SHARED_CONFIGURATION_FILE` environment variable.
This method also supports a `profile` configuration and matching `HCS_PROFILE` environment variable:

!> **NOTE:** The CLI configuration file can not be used directly by terraform, you need to copy it to another
path and replace the AccessKey and SecretKey with yours as they are encrypted which terraform can not read.

Usage:

```terraform
provider "huaweicloudstack" {
  shared_config_file = "/home/tf_user/.hcloud/config.json"
  profile            = "customprofile"
}
```

### ECS Instance Metadata Service

If you're running Terraform from an ECS instance with Agency configured, Terraform will just ask
[the metadata API](https://support.huaweicloudstack.com/intl/en-us/usermanual-ecs/ecs_03_0166.html)
for credentials.

This is a preferred approach over any other when running in ECS as you can avoid
hard coding credentials. Instead these are leased on-the-fly by Terraform
which reduces the chance of leakage.

Usage:

```hcl
provider "huaweicloudstack" {
  region = "cn-north-4"
}
```

### Assume role

If provided with an IAM agency, Terraform will attempt to assume this role using the supplied credentials.

Usage:

```hcl
provider "huaweicloudstack" {
  region     = "cn-north-4"
  cloud      = "mycloud.com"
  access_key = "my-access-key"
  secret_key = "my-secret-key"

  assume_role {
    agency_name = "agency"
    domain_name = "agency_domain"
  }
}
```

## Configuration Reference

The following arguments are supported:

* `region` - (Optional) This is the Huawei Cloud region. It must be provided when using `static credentials`
  authentication, but it can also be sourced from the `HCS_REGION_NAME` environment variables.

* `access_key` - (Optional) The access key of the huaweicloudstack to use. If omitted, the `HCS_ACCESS_KEY` environment
  variable is used.

* `secret_key` - (Optional) The secret key of the huaweicloudstack to use. If omitted, the `HCS_SECRET_KEY` environment
  variable is used.

* `shared_config_file` - (Optional) The path to the shared config file. If omitted, the `HCS_SHARED_CONFIG_FILE` environment
  variable is used.

* `profile` - (Optional) The profile name as set in the shared config file. If omitted, the `HCS_PROFILE` environment
  variable is used. Defaults to the `current` profile in the shared config file.

[//]: # (* `assume_role` - &#40;Optional&#41; Configuration block for an assumed role. See below. Only one assume_role)

[//]: # (  block may be in the configuration.)

[//]: # ()
[//]: # (* `project_name` - &#40;Optional&#41; The Name of the project to login with. If omitted, the `HCS_PROJECT_NAME` environment)

[//]: # (  variable or `region` is used.)

[//]: # ()
[//]: # (* `domain_name` - &#40;Optional&#41; The [Account name]&#40;https://support.huaweicloudstack.com/en-us/usermanual-iam/iam_01_0552.html&#41;)

[//]: # (  of IAM to scope to. If omitted, the `HCS_DOMAIN_NAME` environment variable is used.)

[//]: # ()
[//]: # (* `security_token` - &#40;Optional&#41; The security token to authenticate with a)

[//]: # (  [temporary security credential]&#40;https://support.huaweicloudstack.com/intl/en-us/iam_faq/iam_01_0620.html&#41;. If omitted,)

[//]: # (  the `HCS_SECURITY_TOKEN` environment variable is used.)

* `cloud` - (Optional) The endpoint of the cloud provider. If omitted, the
  `HCS_CLOUD` environment variable is used. 

* `auth_url` - (Optional, Required before 1.14.0) The Identity authentication URL. If omitted, the
  `HCS_AUTH_URL` environment variable is used. Defaults to `https://iam.{{region}}.{{cloud}}:443/v3`.

* `insecure` - (Optional) Trust self-signed SSL certificates. If omitted, the
  `HCS_INSECURE` environment variable is used.

* `max_retries` - (Optional) This is the maximum number of times an API call is retried, in the case where requests are
  being throttled or experiencing transient failures. The delay between the subsequent API calls increases
  exponentially. The default value is `5`. If omitted, the `HCS_MAX_RETRIES` environment variable is used.

* `enterprise_project_id` - (Optional) Default Enterprise Project ID for supported resources. Please see the
  documentation
  at [EPS](https://registry.terraform.io/providers/huaweicloudstack/huaweicloudstack/latest/docs/data-sources/enterprise_project).
  If omitted, the `HCS_ENTERPRISE_PROJECT_ID` environment variable is used.

[//]: # (* `regional` - &#40;Optional&#41; Whether the service endpoints are regional. The default value is `false`.)

* `endpoints` - (Optional) Configuration block in key/value pairs for customizing service endpoints. The following
  endpoints support to be customized: autoscaling, ecs, ims, vpc, nat, evs, obs, sfs, cce, rds, dds, iam. An example
  provider configuration:

```hcl
provider "huaweicloudstack" {
  ...
  endpoints = {
    ecs = "https://ecs-customizing-endpoint.com"
  }
}
```

The `assume_role` block supports:

* `agency_name` - (Required) The name of the agency for assume role.
  If omitted, the `HCS_ASSUME_ROLE_AGENCY_NAME` environment variable is used.

* `domain_name` - (Required) The name of the agency domain for assume role.
  If omitted, the `HCS_ASSUME_ROLE_DOMAIN_NAME` environment variable is used.

## Testing and Development

In order to run the Acceptance Tests for development, the following environment variables must also be set:

* `HCS_REGION_NAME` - The region in which to create the resources.

* `HCS_ACCESS_KEY` - The access key of the huaweicloudstack to use.

* `HCS_SECRET_KEY` - The secret key of the huaweicloudstack to use.

You should be able to use any huaweicloudstack environment to develop on as long as the above environment variables are set.
