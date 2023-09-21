resource "hcs_smn_topic" "topic_1" {
  name         = "topic_1"
  display_name = "topc_name1"
}

resource "hcs_smn_subscription" "subscription_1" {
  topic_urn    = hcs_smn_topic.topic_1.id
  protocol      = "email"
  endpoint      = "test@huawei.com"
  remark       = "remark_test"
}