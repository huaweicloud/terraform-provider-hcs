package servicestage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/servicestage"
)

func getV3ComponentFunc(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("servicestage", acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ServiceStage client: %s", err)
	}
	return servicestage.QueryV3Component(client, state.Primary.Attributes["application_id"], state.Primary.ID)
}

func TestAccV3Component_basic(t *testing.T) {
	var (
		component interface{}

		resourceName = "hcs_servicestage_component.test"
		rc           = acceptance.InitResourceCheck(resourceName, &component, getV3ComponentFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Make sure the networks of the CCE cluster, CSE engine and ELB loadbalancer are same.
			acceptance.TestAccPreCheckCceClusterId(t) // Make sure at least one of node exist.
			acceptance.TestAccPreCheckCSEMicroserviceEngineID(t)
			// Two different JAR packages need to be provided.
			acceptance.TestAccPreCheckServiceStageJarPkgStorageURLs(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3Component_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "application_id", "hcs_servicestage_application.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "environment_id", "hcs_servicestage_environment.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.1"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "replica", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "source"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "runtime_stack.0.deploy_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "runtime_stack.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "runtime_stack.0.type"),
					resource.TestCheckResourceAttr(resourceName, "refer_resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccV3Component_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "application_id", "hcs_servicestage_application.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "environment_id", "hcs_servicestage_environment.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.2"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Asia/Shanghai"),
					resource.TestCheckResourceAttrSet(resourceName, "source"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "runtime_stack.0.deploy_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "runtime_stack.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "runtime_stack.0.type"),
					resource.TestCheckResourceAttr(resourceName, "refer_resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "baar"),
					resource.TestCheckResourceAttr(resourceName, "tags.new_key", "value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV3ComponentImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"tags",
					"source_origin",
					"build_origin",
					"deploy_strategy.0.rolling_release_origin",
					"command_origin",
					"tomcat_opts_origin",
					"update_strategy_origin",
				},
			},
		},
	})
}

func testAccV3ComponentImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var applicationId, resourceId string
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of compnent is not found in the tfstate", resourceName)
		}
		applicationId = rs.Primary.Attributes["application_id"]
		resourceId = rs.Primary.ID
		if applicationId == "" || resourceId == "" {
			return "", fmt.Errorf("the component ID is not exist or application ID is missing")
		}
		return fmt.Sprintf("%s/%s", applicationId, resourceId), nil
	}
}

func testAccV3Component_base(name string) string {
	return fmt.Sprintf(`
data "hcs_availability_zones" "test" {}

data "hcs_cce_clusters" "test" {
  cluster_id = "%[1]s"
}

data "hcs_ecs_compute_instance" "test" {
  instance_id = "%[2]s"
}

resource "hcs_servicestage_environment" "test" {
  name                  = "%[3]s"
  vpc_id                = try(data.hcs_cce_clusters.test.clusters[0].vpc_id, "")
  enterprise_project_id = "0"
}

resource "hcs_servicestage_application" "test" {
  name                  = "%[4]s"
  enterprise_project_id = "0"
}
`, acceptance.HCS_CCE_CLUSTER_ID, acceptance.HCS_ECS_INSTANCE_ID, acceptance.HCS_CSE_MICROSERVICE_ENGINE_ID, name)
}

func testAccV3Component_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_servicestage_component" "test" {
  application_id = hcs_servicestage_application.test.id
  environment_id = hcs_servicestage_environment.test.id
  name           = "%[2]s"
  version        = "1.0.1"
  description    = "Created by terraform script"
  replica        = 1

  source = jsonencode({
    "kind": "package",
    "storage": "obs",
    "url": try(element(split(",", "%[3]s"), 0), "")
  })

  runtime_stack {
    name        = "OpenJDK17"
    version     = "1.4.6"
    type        = "Java"
    deploy_mode = "virtualmachine"
  }

  refer_resources {
    id   = "%[4]s"
    type = "ecs"
  }

  tags = {
    foo = "bar"
    key = "value"
  }

  lifecycle {
    ignore_changes = [refer_resources, source]
  }
}
`, testAccV3Component_base(name),
		name,
		acceptance.HCS_SERVICESTAGE_JAR_PKG_STORAGE_URLS,
		acceptance.HCS_ECS_INSTANCE_ID)
}

func testAccV3Component_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_servicestage_component" "test" {
  application_id = hcs_servicestage_application.test.id
  environment_id = hcs_servicestage_environment.test.id
  name           = "%[2]s"
  version        = "1.0.1"
  description    = "Updated by terraform script"
  replica        = 2

  source = jsonencode({
    "kind": "package",
    "storage": "obs",
    "url": try(element(split(",", "%[3]s"), 0), "")
  })

  runtime_stack {
    name        = "OpenJDK17"
    version     = "1.4.6"
    type        = "Java"
    deploy_mode = "virtualmachine"
  }

  refer_resources {
    id   = "%[4]s"
    type = "ecs"
  }

  tags = {
    foo     = "baar"
    new_key = "value"
  }

  lifecycle {
    ignore_changes = [refer_resources, source]
  }
}
`, testAccV3Component_base(name),
		name,
		acceptance.HCS_SERVICESTAGE_JAR_PKG_STORAGE_URLS,
		acceptance.HCS_ECS_INSTANCE_ID)
}

func TestAccV3Component_yaml(t *testing.T) {
	var (
		component interface{}

		resourceName = "hcs_servicestage_component.test"
		rc           = acceptance.InitResourceCheck(resourceName, &component, getV3ComponentFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Make sure at least one of node exist.
			acceptance.TestAccPreCheckCceClusterId(t)
			// At least one of JAR package must be provided.
			acceptance.TestAccPreCheckServiceStageJarPkgStorageURLs(t, 1)
		},
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3Component_yaml_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "application_id", "hcs_servicestage_application.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "environment_id", "hcs_servicestage_environment.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.1"),
					resource.TestCheckResourceAttr(resourceName, "config_mode", "yaml"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "replica", "2"),
					resource.TestCheckResourceAttr(resourceName, "build", ""),
					resource.TestCheckResourceAttrSet(resourceName, "source"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.0.deploy_mode", "container"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.0.name", "Docker"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.0.type", "Docker"),
					resource.TestCheckResourceAttr(resourceName, "runtime_stack.0.version", "1.0"),
					resource.TestCheckResourceAttr(resourceName, "refer_resources.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "update_strategy"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV3ComponentImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"workload_content",
					"tags",
					"source_origin",
					"build_origin",
					"deploy_strategy.0.rolling_release_origin",
					"command_origin",
					"tomcat_opts_origin",
					"update_strategy_origin",
				},
			},
		},
	})
}

func testAccV3Component_yaml_base(name string) string {
	return fmt.Sprintf(`
data "hcs_availability_zones" "test" {}

data "hcs_cce_clusters" "test" {
  cluster_id = "%[1]s"
}

resource "hcs_servicestage_application" "test" {
  name                  = "%[2]s"
  enterprise_project_id = "0"
}

resource "hcs_servicestage_environment" "test" {
  name                  = "%[2]s"
  vpc_id                = try(data.hcs_cce_clusters.test.clusters[0].vpc_id, "")
  enterprise_project_id = "0"
}
`, acceptance.HCS_CCE_CLUSTER_ID, name)
}

func testAccV3Component_yaml_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_servicestage_component" "test" {
  name           = "%[2]s"
  description    = "Created by terraform script"
  version        = "1.0.1"
  environment_id = hcs_servicestage_environment.test.id
  application_id = hcs_servicestage_application.test.id

  runtime_stack {
    deploy_mode = "container"
    name        = "Docker"
    type        = "Docker"
    version     = "1.0"
  }

  source = jsonencode({
    kind    = "image"
    storage = "swr"
    url     = try(element(split(",", "%[3]s"), 0), "")
  })

  refer_resources {
    type       = "cce"
    id         = try(data.hcs_cce_clusters.test.clusters[0].id, "")
    parameters = jsonencode({
      type      = "VirtualMachine"
      namespace = "default"
    })
  }

  config_mode      = "yaml"
  workload_content = jsonencode({
    apiVersion = "apps/v1"
    kind       = "Deployment"
    metadata   = {
      name      = "%[2]s"
      namespace = "default"
    }
    spec = {
      selector = {}
      template = {
        metadata = {}
        spec = {
          imagePullSecrets = [
            {
              name = "default-secret",
            }
          ]
          terminationGracePeriodSeconds = 30
          volumes                       = []
          restartPolicy                 = "Always"
          dnsPolicy                     = "ClusterFirst"
          containers                    = [
            {
              image           = try(element(split(",", "%[3]s"), 0), "")
              name            = "%[2]s"
              imagePullPolicy = "Always"
              resources       = {
                requests = {
                  cpu    = "0"
                  memory = "0"
                }
                limits = {
                  cpu    = "0"
                  memory = "0"
                }
              }
              ports = [
                {
                  containerPort = 8080,
                  protocol      = "TCP"
                }
              ]
            }
          ]
        }
      }
    }
    strategy = {
      type = "RollingUpdate"
      rollingUpdate = {
        maxSurge       = 0
        maxUnavailable = 1
      }
    }
    replicas = 2
  })
}
`, testAccV3Component_yaml_base(name), name, acceptance.HCS_SERVICESTAGE_JAR_PKG_STORAGE_URLS)
}
