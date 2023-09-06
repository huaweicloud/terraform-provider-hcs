package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/autoscaling/v1/policies"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func TestAccASPolicy_basic(t *testing.T) {
	var asPolicy policies.Policy
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_as_policy.acc_as_policy"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASPolicy_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASPolicyExists(resourceName, &asPolicy),
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
					resource.TestCheckResourceAttr(resourceName, "cool_down_time", "900"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_type", "SCHEDULED"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.operation", "ADD"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.instance_number", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "scaling_group_id", "hcs_as_group.acc_as_group", "id"),
				),
			},
			{
				Config: testASPolicy_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
					resource.TestCheckResourceAttr(resourceName, "cool_down_time", "900"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_type", "SCHEDULED"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.operation", "REMOVE"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.instance_number", "1"),
				),
			},
			{
				Config: testASPolicy_recurrence(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
					resource.TestCheckResourceAttr(resourceName, "cool_down_time", "900"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_type", "RECURRENCE"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.operation", "ADD"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.instance_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "scheduled_policy.0.launch_time", "07:00"),
					resource.TestCheckResourceAttr(resourceName, "scheduled_policy.0.recurrence_type", "Daily"),
					resource.TestCheckResourceAttrSet(resourceName, "scheduled_policy.0.start_time"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccASPolicy_Alarm(t *testing.T) {
	var asPolicy policies.Policy
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_as_policy.acc_as_policy"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASPolicy_alarm(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASPolicyExists(resourceName, &asPolicy),
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
					resource.TestCheckResourceAttr(resourceName, "cool_down_time", "600"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_type", "ALARM"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.operation", "ADD"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.instance_number", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "scaling_group_id", "hcs_as_group.acc_as_group", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "alarm_id", "hcs_ces_alarmrule.alarm_rule", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckASPolicyDestroy(s *terraform.State) error {
	cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	asClient, err := cfg.AutoscalingV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating autoscaling client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_as_policy" {
			continue
		}

		_, err := policies.Get(asClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("AS policy still exists")
		}
	}

	return nil
}

func testAccCheckASPolicyExists(n string, policy *policies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		asClient, err := cfg.AutoscalingV1Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating autoscaling client: %s", err)
		}

		found, err := policies.Get(asClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		policy = &found
		return nil
	}
}

func testASPolicy_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_ecs_compute_keypair" "acc_key" {
  name       = "%[2]s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "hcs_as_configuration" "acc_as_config"{
  scaling_configuration_name = "%[2]s"
  instance_config {
    image    = data.hcs_ims_images.test.images[0].id
    flavor   = data.hcs_ecs_compute_flavors.test.ids[0]
    key_name = hcs_ecs_compute_keypair.acc_key.id
    disk {
      size        = 40
      volume_type = "business_type_01"
      disk_type   = "SYS"
    }
  }
}

resource "hcs_as_group" "acc_as_group"{
  scaling_group_name       = "%[2]s"
  scaling_configuration_id = hcs_as_configuration.acc_as_config.id
  vpc_id                   = hcs_vpc.test.id
  networks {
    id = hcs_vpc_subnet.test.id
  }
  security_groups {
    id = hcs_networking_secgroup.test.id
  }
}
`, common.TestBaseComputeResources(rName), rName)
}

func testASPolicy_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_as_policy" "acc_as_policy"{
  scaling_policy_name = "%[2]s"
  scaling_policy_type = "SCHEDULED"
  scaling_group_id    = hcs_as_group.acc_as_group.id

  scaling_policy_action {
    operation       = "ADD"
    instance_number = 1
  }
  scheduled_policy {
    launch_time = "2099-12-22T12:00Z"
  }
}
`, testASPolicy_base(rName), rName)
}

func testASPolicy_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_as_policy" "acc_as_policy"{
  scaling_policy_name = "%[2]s"
  scaling_policy_type = "SCHEDULED"
  scaling_group_id    = hcs_as_group.acc_as_group.id
  cool_down_time      = 900

  scaling_policy_action {
    operation       = "REMOVE"
    instance_number = 1
  }
  scheduled_policy {
    launch_time = "2099-12-22T12:00Z"
  }
}
`, testASPolicy_base(rName), rName)
}

func testASPolicy_recurrence(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_as_policy" "acc_as_policy"{
  scaling_policy_name = "%[2]s"
  scaling_policy_type = "RECURRENCE"
  scaling_group_id    = hcs_as_group.acc_as_group.id

  scaling_policy_action {
    operation       = "ADD"
    instance_number = 1
  }
  scheduled_policy {
    launch_time     = "07:00"
    recurrence_type = "Daily"
    start_time      = "2099-11-22T12:00Z"
    end_time        = "2099-12-22T12:00Z"
  }
}
`, testASPolicy_base(rName), rName)
}

func testASPolicy_alarm(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ces_alarmrule" "alarm_rule" {
  alarm_name = "%[2]s"

  metric {
    namespace   = "SYS.AS"
    metric_name = "cpu_util"
    dimensions {
      name  = "AutoScalingGroup"
      value = hcs_as_group.acc_as_group.id
    }
  }
  condition {
    period              = 300
    filter              = "average"
    comparison_operator = ">="
    value               = 60
    unit                = "%%"
    count               = 1
    suppress_duration   = 300
  }
  alarm_actions {
    type              = "autoscaling"
    notification_list = []
  }
}

resource "hcs_as_policy" "acc_as_policy"{
  scaling_policy_name = "%[2]s"
  scaling_policy_type = "ALARM"
  scaling_group_id    = hcs_as_group.acc_as_group.id
  alarm_id            = hcs_ces_alarmrule.alarm_rule.id
  cool_down_time      = 600

  scaling_policy_action {
    operation       = "ADD"
    instance_number = 1
  }
}
`, testASPolicy_base(rName), rName)
}
