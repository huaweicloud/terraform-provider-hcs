package aom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	aom "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/aom/v2/model"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func getAlarmRuleResourceFunc(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.HcAomV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}
	response, err := c.ShowAlarmRule(&aom.ShowAlarmRuleRequest{AlarmRuleId: state.Primary.ID})
	if err != nil {
		return nil, fmt.Errorf("error retrieving AOM alarm rule: %s", state.Primary.ID)
	}

	allRules := *response.Thresholds
	if len(allRules) != 1 {
		return nil, fmt.Errorf("error retrieving AOM alarm rule %s", state.Primary.ID)
	}
	rule := allRules[0]
	return rule, nil
}

func TestAccAOMAlarmRule_basic(t *testing.T) {
	var ar aom.QueryAlarmResult
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_aom_alarm_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ar,
		getAlarmRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAOMAlarmRule_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test rule"),
					resource.TestCheckResourceAttr(resourceName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_level", "2"),
					resource.TestCheckResourceAttr(resourceName, "dimensions.0.name", "hostID"),
					resource.TestCheckResourceAttrPair(resourceName, "dimensions.0.value", "hcs_ecs_compute_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "comparison_operator", ">"),
					resource.TestCheckResourceAttr(resourceName, "period", "300000"),
					resource.TestCheckResourceAttr(resourceName, "threshold", "2"),
					resource.TestCheckResourceAttr(resourceName, "evaluation_periods", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAOMAlarmRule_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test rule update"),
					resource.TestCheckResourceAttr(resourceName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_level", "3"),
					resource.TestCheckResourceAttr(resourceName, "comparison_operator", ">="),
					resource.TestCheckResourceAttr(resourceName, "period", "60000"),
					resource.TestCheckResourceAttr(resourceName, "threshold", "3"),
					resource.TestCheckResourceAttr(resourceName, "evaluation_periods", "2"),
				),
			},
		},
	})
}

func TestAccAOMAlarmRule_period(t *testing.T) {
	var ar aom.QueryAlarmResult
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_aom_alarm_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ar,
		getAlarmRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAOMAlarmRule_period(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test acc rule"),
					resource.TestCheckResourceAttr(resourceName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_level", "2"),
					resource.TestCheckResourceAttr(resourceName, "dimensions.0.name", "hostID"),
					resource.TestCheckResourceAttr(resourceName, "dimensions.0.value", "test1"),
					resource.TestCheckResourceAttr(resourceName, "comparison_operator", ">="),
					resource.TestCheckResourceAttr(resourceName, "period", "60000"),
					resource.TestCheckResourceAttr(resourceName, "statistic", "maximum"),
					resource.TestCheckResourceAttr(resourceName, "threshold", "3"),
					resource.TestCheckResourceAttr(resourceName, "evaluation_periods", "5"),
				),
			},
			{
				Config: testAOMAlarmRule_period_update1(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "period", "300000"),
					resource.TestCheckResourceAttr(resourceName, "statistic", "minimum"),
					resource.TestCheckResourceAttr(resourceName, "comparison_operator", ">"),
				),
			},
			{
				Config: testAOMAlarmRule_period_update2(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "period", "900000"),
					resource.TestCheckResourceAttr(resourceName, "statistic", "average"),
					resource.TestCheckResourceAttr(resourceName, "comparison_operator", "="),
				),
			},
			{
				Config: testAOMAlarmRule_period_update3(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "period", "3600000"),
					resource.TestCheckResourceAttr(resourceName, "statistic", "sum"),
					resource.TestCheckResourceAttr(resourceName, "comparison_operator", "<="),
				),
			},
		},
	})
}

func TestAccAOMAlarmRule_level(t *testing.T) {
	var ar aom.QueryAlarmResult
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_aom_alarm_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ar,
		getAlarmRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAOMAlarmRule_level(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test acc rule"),
					resource.TestCheckResourceAttr(resourceName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_level", "1"),
					resource.TestCheckResourceAttr(resourceName, "dimensions.0.name", "hostID"),
					resource.TestCheckResourceAttr(resourceName, "dimensions.0.value", "test1"),
					resource.TestCheckResourceAttr(resourceName, "comparison_operator", ">="),
					resource.TestCheckResourceAttr(resourceName, "period", "60000"),
					resource.TestCheckResourceAttr(resourceName, "statistic", "maximum"),
					resource.TestCheckResourceAttr(resourceName, "threshold", "3"),
					resource.TestCheckResourceAttr(resourceName, "evaluation_periods", "1"),
				),
			},
			{
				Config: testAOMAlarmRule_level_update1(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "alarm_level", "2"),
					resource.TestCheckResourceAttr(resourceName, "evaluation_periods", "2"),
				),
			},
			{
				Config: testAOMAlarmRule_level_update2(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "alarm_level", "3"),
					resource.TestCheckResourceAttr(resourceName, "evaluation_periods", "3"),
				),
			},
			{
				Config: testAOMAlarmRule_level_update3(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "alarm_level", "4"),
					resource.TestCheckResourceAttr(resourceName, "evaluation_periods", "4"),
				),
			},
		},
	})
}

func testAOMAlarmRule_base(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ecs_compute_instance" "test" {
  name               = "ecs-%s"
  image_id           = data.hcs_ims_images.test.id
  flavor_id          = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids = [hcs_networking_secgroup.test.id]
  availability_zone  = data.hcs_availability_zones.test.names[0]

  network {
    uuid = hcs_vpc_subnet.test.id
  }
}
`, common.TestBaseComputeResources(rName), rName)
}

func testAOMAlarmRule_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_aom_alarm_rule" "test" {
  name                 = "%s"
  alarm_level          = 2
  alarm_action_enabled = false
  description          = "test rule"

  namespace   = "PAAS.NODE"
  metric_name = "cupUsage"

  dimensions {
    name  = "hostID"
    value = hcs_ecs_compute_instance.test.id
  }

  comparison_operator = ">"
  period              = 300000
  statistic           = "average"
  threshold           = 2
  unit                = "Percent"
  evaluation_periods  = 3
}
`, testAOMAlarmRule_base(rName), rName)
}

func testAOMAlarmRule_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_aom_alarm_rule" "test" {
  name                 = "%s"
  alarm_level          = 3
  alarm_action_enabled = false
  description          = "test rule update"

  namespace   = "PAAS.NODE"
  metric_name = "cupUsage"

  dimensions {
    name  = "hostID"
    value = hcs_ecs_compute_instance.test.id
  }

  comparison_operator = ">="
  period              = 60000
  statistic           = "average"
  threshold           = 3
  unit                = "Percent"
  evaluation_periods  = 2
}
`, testAOMAlarmRule_base(rName), rName)
}

func testAOMAlarmRule_period(rName string) string {
	return fmt.Sprintf(`
resource "hcs_aom_alarm_rule" "alarm_rule" {
  name        = "%s"
  alarm_level = 2
  description = "test acc rule"

  namespace   = "PAAS.NODE"
  metric_name = "cupUsage"

  dimensions {
    name  = "hostID"
    value = "test1"
  }

  comparison_operator = ">="
  period              = 60000 
  statistic           = "maximum"
  threshold           = 3
  unit                = "32"
  evaluation_periods  = 5
}
`, rName)
}

func testAOMAlarmRule_period_update1(rName string) string {
	return fmt.Sprintf(`
resource "hcs_aom_alarm_rule" "alarm_rule" {
  name        = "%s"
  alarm_level = 2
  description = "test acc rule"

  namespace   = "PAAS.NODE"
  metric_name = "cupUsage"

  dimensions {
    name  = "hostID"
    value = "test1"
  }

  comparison_operator = ">"
  period              = 300000 
  statistic           = "minimum"
  threshold           = 3
  unit                = "32"
  evaluation_periods  = 5
}
`, rName)
}

func testAOMAlarmRule_period_update2(rName string) string {
	return fmt.Sprintf(`
resource "hcs_aom_alarm_rule" "alarm_rule" {
  name        = "%s"
  alarm_level = 2
  description = "test acc rule"

  namespace   = "PAAS.NODE"
  metric_name = "cupUsage"

  dimensions {
    name  = "hostID"
    value = "test1"
  }

  comparison_operator = "="
  period              = 900000 
  statistic           = "average"
  threshold           = 3
  unit                = "32"
  evaluation_periods  = 5
}
`, rName)
}

func testAOMAlarmRule_period_update3(rName string) string {
	return fmt.Sprintf(`
resource "hcs_aom_alarm_rule" "alarm_rule" {
  name        = "%s"
  alarm_level = 2
  description = "test acc rule"

  namespace   = "PAAS.NODE"
  metric_name = "cupUsage"

  dimensions {
    name  = "hostID"
    value = "test1"
  }

  comparison_operator = "<="
  period              = 3600000
  statistic           = "sum"
  threshold           = 3
  unit                = "32"
  evaluation_periods  = 5
}
`, rName)
}

func testAOMAlarmRule_level(rName string) string {
	return fmt.Sprintf(`
resource "hcs_aom_alarm_rule" "alarm_rule" {
  name        = "%s"
  alarm_level = 3
  description = "test acc rule"

  namespace   = "PAAS.NODE"
  metric_name = "cupUsage"

  dimensions {
    name  = "hostID"
    value = "test1"
  }

  comparison_operator = ">="
  period              = 60000 
  statistic           = "maximum"
  threshold           = 3
  unit                = "32"
  evaluation_periods  = 3
}
`, rName)
}

func testAOMAlarmRule_level_update1(rName string) string {
	return fmt.Sprintf(`
resource "hcs_aom_alarm_rule" "alarm_rule" {
  name        = "%s"
  alarm_level = 2
  description = "test acc rule"

  namespace   = "PAAS.NODE"
  metric_name = "cupUsage"

  dimensions {
    name  = "hostID"
    value = "test1"
  }

  comparison_operator = ">="
  period              = 60000 
  statistic           = "maximum"
  threshold           = 3
  unit                = "32"
  evaluation_periods  = 2
}
`, rName)
}

func testAOMAlarmRule_level_update2(rName string) string {
	return fmt.Sprintf(`
resource "hcs_aom_alarm_rule" "alarm_rule" {
  name        = "%s"
  alarm_level = 2
  description = "test acc rule"

  namespace   = "PAAS.NODE"
  metric_name = "cupUsage"

  dimensions {
    name  = "hostID"
    value = "test1"
  }

  comparison_operator = ">="
  period              = 60000 
  statistic           = "maximum"
  threshold           = 3
  unit                = "32"
  evaluation_periods  = 5
}
`, rName)
}

func testAOMAlarmRule_level_update3(rName string) string {
	return fmt.Sprintf(`
resource "hcs_aom_alarm_rule" "alarm_rule" {
  name        = "%s"
  alarm_level = 4
  description = "test acc rule"

  namespace   = "PAAS.NODE"
  metric_name = "cupUsage"

  dimensions {
    name  = "hostID"
    value = "test1"
  }

  comparison_operator = ">="
  period              = 60000 
  statistic           = "maximum"
  threshold           = 3
  unit                = "32"
  evaluation_periods  = 4
}
`, rName)
}
