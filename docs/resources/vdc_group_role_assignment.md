---
subcategory: "VDC"
---

# hcs_vdc_group_role_assignment

Manages a VDC user group role assignment within Huawei Cloud Stack.

> [!NOTE]
>
> Supported from ManageOne version 8.5.1 onwards.

## Example Usage

```hcl

variable "domain_id" {
  default = "de432e279e104f7ab9329c4d8f66b7d4"
}
variable "vdc_id" {
  default = "a18c2ce0-5379-4b34-8a12-eee47f5cfa89"
}
variable  "project_id" {
  default = "60d745a7498e40f0a3c346ec53f7b101"
}
variable  "enterprise_project_id" {
  default = "b7b19bbc-2826-44d2-aa6a-f8f4e023e10f"
}

resource "hcs_vdc_role" "role1" {
  domain_id = var.domain_id
  name = "roleName1"
  type = "AX"
  policy = {
    "Depends" : [],
    "Statement" : [
      {
        "Action" : [
          "ecs:cloudServers:start",
          "ecs:cloudServers:list"
        ],
        "Effect" : "Allow"
      }
    ],
    "Version" : "1.1"
  }
}

resource "hcs_vdc_group" "group01" {
  vdc_id       = var.vdc_id
  name         = "Usergroup1"
}

// Adding Tenant Authorization to a User Group
resource "hcs_vdc_group_role_assignment" "vdc_group_role_assignment1" {
  group_id = hcs_vdc_group.group01.id
  role_assignment {
    domain_id = var.domain_id
    role_id = hcs_vdc_role.role1.id
  }
}

// Adding Resource Space Permissions to a User Group
resource "hcs_vdc_group_role_assignment" "vdc_group_role_assignment2" {
  group_id = hcs_vdc_group.group01.id
  role_assignment {
    project_id = var.project_id
    role_id = hcs_vdc_role.role1.id
  }
}

// Grant the all resource space to the user group.
resource "hcs_vdc_group_role_assignment" "vdc_group_role_assignment3" {
  group_id = hcs_vdc_group.group01.id
  role_assignment {
    domain_id = var.domain_id
    project_id = "all"
    role_id = hcs_vdc_role.role1.id
  }
}

// Adding Enterprise Project Authorization to a User Group
resource "hcs_vdc_group_role_assignment" "vdc_group_role_assignment4" {
  group_id = hcs_vdc_group.group01.id
  role_assignment {
    enterprise_project_id = var.enterprise_project_id
    role_id = hcs_vdc_role.role1.id
  }
}

```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String, ForceNew) User group ID

* `role_assignment` - (Required, Set, ForceNew) Permission.
  * `role_id` - (Required, String, ForceNew) Role ID.
  * `domain_id` - (Optional, String, ForceNew) Tenant ID.
  * `project_id` - (Optional, String, ForceNew) Resource space ID.
  * `enterprise_project_id` - (Optional, String, ForceNew) Enterprise project ID.

>[!NOTE]
>
> domain_id project_id and enterprise_project_id, only one of the values can be set and cannot be empty. When project_id is set to all, domain_id is mandatory.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - User group ID.

## Import

VDC user group role assignment can be imported using the `id`, e.g.

```
$ terraform import hcs_vdc_group_role_assignment.vdc_group_role_assignment 3b002f5e4aae407082630a00d2ac0f40
```