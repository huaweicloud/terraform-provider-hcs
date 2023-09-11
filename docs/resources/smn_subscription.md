---
subcategory: "Simple Message Notification (SMN)"
---

# hcs_smn_subscription

Manages an SMN subscription resource within HuaweiCloudStack.

## Example Usage

```hcl
resource "hcs_smn_topic" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
}

resource "hcs_smn_subscription" "subscription_1" {
  topic_urn = hcs_smn_topic.topic_1.id
  endpoint  = "mailtest@gmail.com"
  protocol  = "email"
  remark    = "O&M"
}

resource "hcs_smn_subscription" "subscription_2" {
  topic_urn = hcs_smn_topic.topic_1.id
  endpoint  = "13600000000"
  protocol  = "sms"
  remark    = "O&M"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the SMN subscription resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `topic_urn` - (Required, String, ForceNew) Specifies the resource identifier of a topic, which is unique.
  Changing this parameter will create a new resource.

* `protocol` - (Required, String, ForceNew) Specifies the protocol of the message endpoint. Currently, email, sms are supported. Changing this parameter will create a new resource.

* `endpoint` - (Required, String, ForceNew) Message endpoint. Changing this parameter will create a new resource.
  + **For an email subscription**, the endpoint is an mail address.
  + **For an SMS message subscription**, the endpoint is a phone number,
    the format is \[+\]\[country code\]\[phone number\], e.g. +86185xxxx0000.

* `remark` - (Optional, String, ForceNew) Remark information. The remarks must be a UTF-8-coded character string
  containing 128 bytes. Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the subscription urn.

* `subscription_urn` - Resource identifier of a subscription, which is unique.

* `owner` - Project ID of the topic creator.

* `status` - Subscription status.
  + **0**: indicates that the subscription is not confirmed.
  + **1**: indicates that the subscription is confirmed.
  + **3**: indicates that the subscription is canceled.

## Import

SMN subscription can be imported using the `id` (subscription urn), e.g.

```
$ terraform import hcs_smn_subscription.subscription_1 urn:smn:cn-north-4:0970dd7a1300f5672ff2c003c60ae115:topic_1:a2aa5a1f66df494184f4e108398de1a6
```
