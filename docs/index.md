# Huawei Cloud Stack Provider

The Huawei Cloud Stack provider is used to interact with the many resources supported by Huawei Cloud Stack. The provider needs to be
configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

Terraform 0.13 and later:

```hcl
terraform {
  required_providers {
    hcs = {
      source  = "huaweicloud/hcs"
      version = "~> 2.3.0"
    }
  }
}

# Configure the Huawei Cloud Stack Provider
provider "hcs" {
  region       = "my-region-name"
  project_name = "my-project-name"
  cloud        = "mycloud.com"
  access_key   = "my-access-key"
  secret_key   = "my-secret-key"
}

# Create a VPC
resource "hcs_vpc" "example" {
  name = "my_vpc"
  cidr = "192.168.0.0/16"
}
```

## Authentication

The Huawei Cloud Stack provider offers a flexible means of providing credentials for authentication. The following methods are
supported, in this order, and explained below:

* Static credentials
* Environment variables

### Static credentials

!> **WARNING:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage
should this file ever be committed to a public version control system.

Static credentials can be provided by adding an `access_key` and `secret_key`
in-line in the provider block:

Usage:

```hcl
provider "hcs" {
  region       = "my-region-name"
  project_name = "my-project-name"
  cloud        = "mycloud.com"
  access_key   = "my-access-key"
  secret_key   = "my-secret-key"
}
```

### Environment variables

You can provide your credentials via the `HCS_ACCESS_KEY` and
`HCS_SECRET_KEY` environment variables, representing your Huawei Cloud Stack Access Key and Secret Key, respectively.

```hcl
provider "hcs" {}
```

Usage:

```sh
$ export HCS_ACCESS_KEY="anaccesskey"
$ export HCS_SECRET_KEY="asecretkey"
$ export HCS_REGION_NAME="cn-north-4"
$ export HCS_PROJECT_NAME="my-project-name"
$ export HCS_CLOUD="mycloud.com"
$ terraform plan
```

### Assume role

If provided with an agency, Terraform will attempt to assume this role using the supplied credentials.

-> **Note** It only can use one `assume_role` block, this means it can only be delegated to one user.
Usage:

```hcl
provider "hcs" {
  auth_url     = "https://iam-apigateway-proxy.my-cloud-name/v3"
  region       = "my-region-name"
  project_name = "my-project-name"
  cloud        = "my-cloud-name"
  insecure     = true

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

* `region` - (Optional) This is the Huawei Cloud Stack region. It must be provided when using `static credentials`
  authentication, but it can also be sourced from the `HCS_REGION_NAME` environment variables.

* `access_key` - (Optional) The access key of the Huawei Cloud Stack to use. If omitted, the `HCS_ACCESS_KEY` environment
  variable is used.

* `secret_key` - (Optional) The secret key of the Huawei Cloud Stack to use. If omitted, the `HCS_SECRET_KEY` environment
  variable is used.

* `project_name` - (Optional) The Name of the project to login with. If omitted, the `HCS_PROJECT_NAME` environment

* `cloud` - (Optional) The endpoint of the cloud provider. If omitted, the
  `HCS_CLOUD` environment variable is used. 

* `auth_url` - (Optional, Required before 1.14.0) The Identity authentication URL. If omitted, the
  `HCS_AUTH_URL` environment variable is used. Defaults to `https://iam-apigateway-proxy.{{region}}.{{cloud}}:443`.

* `insecure` - (Optional) Trust self-signed SSL certificates. If omitted, the
  `HCS_INSECURE` environment variable is used.

* `max_retries` - (Optional) This is the maximum number of times an API call is retried, in the case where requests are
  being throttled or experiencing transient failures. The delay between the subsequent API calls increases
  exponentially. The default value is `5`. If omitted, the `HCS_MAX_RETRIES` environment variable is used.

* `enterprise_project_id` - (Optional) Default Enterprise Project ID for supported resources. If omitted, the
  `HCS_ENTERPRISE_PROJECT_ID` environment variable is used.

* `domain_name` - (Optional) The tenant name of user account.
  If omitted, the `HCS_DOMAIN_NAME` environment variable is used.

* `assume_role` - (Optional) Configuration block for an assumed role.

  The `assume_role` block supports:
  * `agency_name` - (Required) The name of the agency for assume role. 
  If omitted, the `HCS_ASSUME_ROLE_AGENCY_NAME` environment variable is used.
  * `domain_name` - (Required) The name of the agency domain for assume role. 
  If omitted, the `HCS_ASSUME_ROLE_DOMAIN_NAME` environment variable is used.

* `endpoints` - (Optional) Configuration block in key/value pairs for customizing service endpoints.
  The [endpoints](#block--endpoints) block to support custom endpoints is documented below.
  An example provider configuration:

```hcl
provider "hcs" {
  ...
  endpoints = {
    ecs = "https://ecs-customizing-endpoint.com"
  }
}
```

<a name="block--endpoints"></a>
The `endpoints` block supports:

* `aom` - (Optional) Use this to override the default endpoint URL. It's used to customize **AOM** endpoints.

* `autoscaling` - (Optional) Use this to override the default endpoint URL. It's used to customize **AS** endpoints.

* `bms` - (Optional) Use this to override the default endpoint URL. It's used to customize **BMS** endpoints.

* `cce` - (Optional) Use this to override the default endpoint URL. It's used to customize **CCE** endpoints.

* `cfw` - (Optional) Use this to override the default endpoint URL. It's used to customize **CFW** endpoints.

* `codeartsreq` - (Optional) Use this to override the default endpoint URL. It's used to customize 
  **codearts req** endpoints.

* `codeartsrepo` - (Optional) Use this to override the default endpoint URL. It's used to customize
  **codearts repo** endpoints.

* `codeartspipeline` - (Optional) Use this to override the default endpoint URL. It's used to customize
  **codearts pipeline** endpoints.

* `csms` - (Optional) Use this to override the default endpoint URL. It's used to customize **CSMS** endpoints.

* `dcs` - (Optional) Use this to override the default endpoint URL. It's used to customize **DCS** endpoints.

* `dms` - (Optional) Use this to override the default endpoint URL. It's used to customize **DMS** endpoints.

* `dns` - (Optional) Use this to override the default endpoint URL. It's used to customize **DNS** endpoints.

* `dws` - (Optional) Use this to override the default endpoint URL. It's used to customize **DWS** endpoints.

* `ecs` - (Optional) Use this to override the default endpoint URL. It's used to customize **ECS** endpoints.

* `elb` - (Optional) Use this to override the default endpoint URL. It's used to customize **ELB** endpoints.

* `eps` - (Optional) Use this to override the default endpoint URL. It's used to customize **EPS** endpoints.

* `evs` - (Optional) Use this to override the default endpoint URL. It's used to customize **EVS** endpoints.

* `hss` - (Optional) Use this to override the default endpoint URL. It's used to customize **HSS** endpoints.

* `iam` - (Optional) Use this to override the default endpoint URL. It's used to customize **IAM** endpoints.

* `ims` - (Optional) Use this to override the default endpoint URL. It's used to customize **IMS** endpoints.

* `kms` - (Optional) Use this to override the default endpoint URL. It's used to customize **DEW** endpoints.

* `lts` - (Optional) Use this to override the default endpoint URL. It's used to customize **LTS** endpoints.

* `mrs` - (Optional) Use this to override the default endpoint URL. It's used to customize **MRS** endpoints.

* `nat` - (Optional) Use this to override the default endpoint URL. It's used to customize **NAT** endpoints.

* `opengauss` - (Optional) Use this to override the default endpoint URL. It's used to customize **OpenGauss**
  endpoints. It is **Required** when **opengaussv31** is not empty.

* `opengaussv31` - (Optional) Use this to override the default endpoint URL. It's used to customize **OpenGauss**
  endpoints. It is **Required** when **opengauss** is not empty.

* `roma` - (Optional) Use this to override the default endpoint URL. It's used to customize **ROMA Connect** endpoints.

* `rds` - (Optional) Use this to override the default endpoint URL. It's used to customize **RDS** endpoints.

* `secmaster` - (Optional) Use this to override the default endpoint URL. It's used to customize **SecMaster** endpoints.

* `servicestage` - (Optional) Use this to override the default endpoint URL. It's used to customize
  **servicestage** endpoints.

* `sfs` - (Optional) Use this to override the default endpoint URL. It's used to customize **SFS** endpoints.

* `sfs-turbo` - (Optional) Use this to override the default endpoint URL. It's used to customize **SFSTurbo** endpoints.

* `smn` - (Optional) Use this to override the default endpoint URL. It's used to customize **SMN** endpoints.

* `swr` - (Optional) Use this to override the default endpoint URL. It's used to customize **SWR** endpoints.

* `ucs` - (Optional) Use this to override the default endpoint URL. It's used to customize **UCS** endpoints.

* `vpc` - (Optional) Use this to override the default endpoint URL. It's used to customize **VPC** endpoints.

* `vpcep` - (Optional) Use this to override the default endpoint URL. It's used to customize **VPCEP** endpoints.

* `waf` - (Optional) Use this to override the default endpoint URL. It's used to customize **WAF** endpoints.

## Testing and Development

In order to run the Acceptance Tests for development, the following environment variables must also be set:

* `HCS_CLOUD` - The endpoint of the Huawei Cloud Stack to use.

* `HCS_PROJECT_NAME` - The project name of the Huawei Cloud Stack to use.

* `HCS_REGION_NAME` - The region in which to create the resources.

* `HCS_ACCESS_KEY` - The access key of the Huawei Cloud Stack to use.

* `HCS_SECRET_KEY` - The secret key of the Huawei Cloud Stack to use.

You should be able to use any Huawei Cloud Stack environment to develop on as long as the above environment variables are set.
