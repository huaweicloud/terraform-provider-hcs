package smn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/smn/v2/subscriptions"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getResourceSMNSubscription(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	smnClient, err := conf.SmnV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SMN client: %s", err)
	}

	foundList, err := subscriptions.List(smnClient).Extract()
	if err != nil {
		return nil, err
	}

	var subscription *subscriptions.SubscriptionGet
	urn := state.Primary.ID
	for i := range foundList {
		if foundList[i].SubscriptionUrn == urn {
			subscription = &foundList[i]
		}
	}

	if subscription == nil {
		return nil, fmt.Errorf("the subscription does not exist")
	}

	return subscription, nil
}

func TestAccSMNV2Subscription_basic(t *testing.T) {
	var subscription1 subscriptions.SubscriptionGet
	var subscription2 subscriptions.SubscriptionGet
	resourceName1 := "hcs_smn_subscription.subscription_1"
	resourceName2 := "hcs_smn_subscription.subscription_2"
	rName := acceptance.RandomAccResourceNameWithDash()

	rc1 := acceptance.InitResourceCheck(
		resourceName1,
		&subscription1,
		getResourceSMNSubscription,
	)

	rc2 := acceptance.InitResourceCheck(
		resourceName2,
		&subscription2,
		getResourceSMNSubscription,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSMNSubscriptionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSMNV2SubscriptionConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName1, "endpoint", "mailtest@gmail.com"),
					resource.TestCheckResourceAttr(resourceName2, "endpoint", "13600000000"),
				),
			},
			{
				ResourceName:      resourceName1,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSMNSubscriptionV2Destroy(s *terraform.State) error {
	cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	smnClient, err := cfg.SmnV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating SMN client: %s", err)
	}

	foundList, err := subscriptions.List(smnClient).Extract()
	if err != nil {
		return err
	}

	var subscription *subscriptions.SubscriptionGet
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_smn_subscription" {
			continue
		}

		urn := rs.Primary.ID
		for i := range foundList {
			if foundList[i].SubscriptionUrn == urn {
				subscription = &foundList[i]
			}
		}
		if subscription != nil {
			return fmt.Errorf("subscription still exists")
		}
	}

	return nil
}

func testAccSMNV2SubscriptionConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcs_smn_topic" "topic_1" {
  name         = "%s"
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
`, rName)
}
