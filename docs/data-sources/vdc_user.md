---
subcategory: "Enterprise Projects (EPS)"
---

# hcs_vdc_user

Use this data source to get the list of the vdc user.

## Example Usage

```hcl
data "hcs_vdc_user" "userList" {
    vdc_id = "xxx"
    name = "xxx"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (必填, String)用户名，按照名称字段模糊搜索。长度：1到128个字符。

* `vdc_id` - (必填, String) 用于查询指定VDC下的用户列表，VDC id，只能包含小写字母、数字、中划线，长度在1-36之间。


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - 用户id。
  
* `domain_id` - 租户ID。

* `name` - 用户名，长度在5-32之间，首位为字母或下划线_，不能为纯数字，且不能包含除下划线_以外的特殊字符。

* `display_name` - 用户别名，由除">"和"<"之外的字符组成，0-128个字符。

* `enabled` - 用户状态，true或false，默认为true。为false时，表示用户处于停用状态，无法登录。

* `description` - 说明，长度小于等于255。

* `vdc_id` - VDC ID。

* `ldap_id` - LDAP服务器配置id。

* `create_at` - 创建时间。

* `auth_type` - 用户类型; 分别有"LOCAL_AUTH"、"SAML_AUTH"、"LDAP_AUTH"、"MACHINE_USER", 表示本地认证、SAML认证、LDAP认证、机机用户。

* `access_mode` - 访问方式; 分别有 "default"、"programmatic"、"console", 表示控制台访问、编程访问、控制台访问。

* `top_vdc_id` - 一级VDC id，只能包含小写字母、数字、中划线，长度在1-36之间。
