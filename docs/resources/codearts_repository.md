---
subcategory: "CodeArts"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_codearts_repository"
description: |-
  Manages a CodeArts repository resource within HuaweiCloudStack.
---

# hcs_codearts_repository

-> **NOTE:** This resource can only be used in HCS **8.5.0** and **later** version.

Manages a CodeArts repository resource within HuaweiCloudStack.

## Example Usage

```hcl
variable "project_name" {}
variable "repository_name" {}

resource "hcs_codearts_project" "test" {
  name = var.project_name
  type = "scrum"
}

resource "hcs_codearts_repository" "test" {
  project_id = hcs_codearts_project.test.id

  name             = var.repository_name
  description      = "terraform test"
  gitignore_id     = "Go"
  enable_readme    = 0 
  visibility_level = 20
  license_id       = 2 
  import_members   = 0 
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) The repository name.

  Changing this parameter will create a new resource.

* `project_id` - (Required, String, ForceNew) The project ID for Codehub service.

  Changing this parameter will create a new resource.

* `visibility_level` - (Optional, Int, ForceNew) The visibility level.  
  The valid values are as follows:
  + **0**: Private.
  + **20**: Public read-only.

  Defaults to `0`. Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) The repository description.

  Changing this parameter will create a new resource.

* `import_url` - (Optional, String, ForceNew) The HTTPS address of the template repository encrypted using Base64.

  Changing this parameter will create a new resource.

* `gitignore_id` - (Optional, String, ForceNew) The program language type for generating `.gitignore` files.

  Changing this parameter will create a new resource.

* `license_id` - (Optional, Int, ForceNew) The license ID for public repository.
  **1**: Apache_License_v2.0.txt
  **2**: MIT_License.txt
  **3**: BSD_3Clause.txt
  **4**: Eclipse_Public_License_v1.0.txt
  **5**: BSD_2Clause.txt
  **6**: GNU_General_Public_License_v2.0.txt
  **7**: GNU_General_Public_License_v3.0.txt
  **8**: GNU_Affero_General_Public_License_v3.0.txt
  **9**: GNU_Lesser_General_Public_License.txt
  **10**: GNU_Lesser_General_Public_License_v3.0.txt
  **11**: Mozilla_Public_License_v2.0.txt
  **12**: The_Unlicense.txt

  Changing this parameter will create a new resource.

* `enable_readme` - (Optional, Int, ForceNew) Whether to generate the `README.md` file.  
  The valid values are as follows:
  + **0**: Disable.
  + **1**: Enable.

  Changing this parameter will create a new resource.

* `import_members` - (Optional, Int, ForceNew) Whether to import the project members.  
  The valid values are as follows:
  + **0**: Do not import members.
  + **1**: Import members.

  Defaults to `1`. Changing this parameter will create a new resource.

* `template_id` - (Optional, String, ForceNew) The ID of the copied template.

  Changing this parameter will create a new resource.

* `caller` - (Optional, String, ForceNew) The ID of the caller.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `https_url` - The HTTPS URL that used to the fork repository.

* `ssh_url` - The SSH URL that used to the fork repository.

* `web_url` - The web URL, accessing this URL will redirect to the repository detail page.

* `lfs_size` - The LFS capacity, in MB. If the capacity is greater than `1,024`M, the unit is GB.

* `status` - The repository status.  
  The valid values are as follows:
  + **0**: Normal.
  + **3**: Frozen.
  + **4**: Closed.

* `created_at` - The creation time.

* `updated_at` - The last update time.

* `repository_id` - The repository short ID.

* `repository_uuid` - The repository UUID.

* `repository_size` - The total size of the repository, in MB. If the capacity is greater than `1,024`M, the unit is GB.

* `creator_name` - The username of the creator. When the user is a tenant, the username is equal to the tenant name.

* `domain_name` - The tenant name of the creator.

* `group_name` - The repository group name(the segment between the domain name and the repository name in the 
  clone URL). Example: **git@repo.alpha.devcloud.intest.com:Demo00228/testword.git**. `group_name` is **Demo00228**.

* `iam_user_uuid` - The UUID of the IAM user. 

* `is_owner` - Whether the current user is the creator of the repository.
  + **1**: Yes.
  + **0**: No.

* `project_is_deleted` - Whether the project was deleted.

* `star` - Indicate whether the repository is stored.

* `user_role` - The user permissions in the repository.
  + **20**: Read-only member.
  + **30**: Regular member.
  + **40**: Administrator.

## Import

The repository can be imported using the `id`, e.g.

```bash
$ terraform import hcs_codearts_repository.test 0ce123456a00f2591fabc00385ff1234
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `name`, `description`, `gitignore_id`, `enable_readme`, `license_id` and
`import_members`. It is generally recommended running `terraform plan` after importing the repository.
You can then decide if changes should be applied to the repository, or the resource definition should be updated to
align with the repository. Also you can ignore changes as below.

```hcl
resource "hcs_codearts_repository" "test" {
  ...

  lifecycle {
    ignore_changes = [
      name, license_id,
    ]
  }
}
```
