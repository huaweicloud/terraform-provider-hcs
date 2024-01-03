package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	domains "github.com/chnsz/golangsdk/openstack/waf_hw/v1/premium_domains"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	hwacceptance "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getWafDedicateDomainResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.WafDedicatedV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF dedicated client: %s", err)
	}

	epsID := state.Primary.Attributes["enterprise_project_id"]
	return domains.GetWithEpsID(client, state.Primary.ID, epsID)
}

func TestAccWafDedicateDomainV1_basic(t *testing.T) {
	var obj interface{}

	randName := acceptance.RandomAccResourceName()
	resourceName := "hcs_waf_dedicated_domain.domain_1"

	rc := hwacceptance.InitResourceCheck(
		resourceName,
		&obj,
		getWafDedicateDomainResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedDomainV1_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "domain", fmt.Sprintf("www.%s.com", randName)),
					resource.TestCheckResourceAttr(resourceName, "proxy", "false"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.1"),
					resource.TestCheckResourceAttr(resourceName, "cipher", "cipher_1"),
					resource.TestCheckResourceAttr(resourceName, "protect_status", "1"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.0.14"),
					resource.TestCheckResourceAttr(resourceName, "server.0.type", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "custom_page.0.http_return_code", "404"),
					resource.TestCheckResourceAttr(resourceName, "custom_page.0.block_page_type", "application/json"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key1", "$time_local"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key2", "$tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.status", "false"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.connection_timeout", "50"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.read_timeout", "200"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.write_timeout", "200"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.ip_tags.0", "ip_tag"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.ip_tags.1", "$remote_addr"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.session_tag", "session_tag"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.user_tag", "user_tag"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_page.0.page_content"),
					resource.TestCheckResourceAttrSet(resourceName, "server.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "certificate_id"),
					resource.TestCheckResourceAttrSet(resourceName, "certificate_name"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttrSet(resourceName, "protect_status"),
					resource.TestCheckResourceAttrSet(resourceName, "protocol"),
					resource.TestCheckResourceAttrSet(resourceName, "tls"),
					resource.TestCheckResourceAttrSet(resourceName, "cipher"),
					resource.TestCheckResourceAttrSet(resourceName, "alarm_page.template_name"),
					resource.TestCheckResourceAttrSet(resourceName, "compliance_certification.pci_3ds"),
					resource.TestCheckResourceAttrSet(resourceName, "compliance_certification.pci_dss"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_protection.0.error_threshold"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_protection.0.error_percentage"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_protection.0.initial_downtime"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_protection.0.multiplier_for_consecutive_breakdowns"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_protection.0.pending_url_request_threshold"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_protection.0.duration"),
				),
			},
			{
				Config: testAccWafDedicatedDomainV1_update1(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.2"),
					resource.TestCheckResourceAttr(resourceName, "cipher", "cipher_2"),
					resource.TestCheckResourceAttr(resourceName, "pci_3ds", "true"),
					resource.TestCheckResourceAttr(resourceName, "pci_dss", "true"),
					resource.TestCheckResourceAttr(resourceName, "protect_status", "0"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", "${http_host}/error.html"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8443"),
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.0.14"),
					resource.TestCheckResourceAttr(resourceName, "server.1.address", "119.8.0.15"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key2", "$request_length"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key3", "$remote_addr"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.error_threshold", "1000"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.error_percentage", "87.5"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.initial_downtime", "200"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.multiplier_for_consecutive_breakdowns", "5"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.pending_url_request_threshold", "7000"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.duration", "10000"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.status", "true"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.connection_timeout", "100"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.read_timeout", "1000"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.write_timeout", "1000"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.ip_tags.0", "ip_tag_update"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.ip_tags.1", "ip_tag_another"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.session_tag", "session_tag_update"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.user_tag", "user_tag_update"),
				),
			},
			{
				Config: testAccWafDedicatedDomainV1_update2(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.error_threshold", "2147483647"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.error_percentage", "99"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.initial_downtime", "2147483647"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.multiplier_for_consecutive_breakdowns", "2147483647"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.pending_url_request_threshold", "2147483647"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.duration", "2147483647"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.status", "false"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.connection_timeout", "180"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.read_timeout", "3600"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.write_timeout", "3600"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.ip_tags.0", "ip_tag_update"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.ip_tags.1", "ip_tag_another"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.session_tag", "session_tag_update"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.user_tag", "user_tag_update"),
				),
			},
			{
				Config: testAccWafDedicatedDomainV1_update3(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.error_threshold", "0"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.error_percentage", "0"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.initial_downtime", "0"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.multiplier_for_consecutive_breakdowns", "0"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.pending_url_request_threshold", "0"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.duration", "0"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.status", "true"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.connection_timeout", "0"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.read_timeout", "0"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.write_timeout", "0"),
				),
			},
			{
				Config: testAccWafDedicatedDomainV1_policy(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "domain", fmt.Sprintf("www.%s.com", randName)),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.2"),
					resource.TestCheckResourceAttr(resourceName, "protect_status", "0"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "server.0.type", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.0.14"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"keep_policy"},
			},
		},
	})
}

func TestAccWafDedicateDomainV1_withEpsID(t *testing.T) {
	var obj interface{}

	randName := acceptance.RandomAccResourceName()
	resourceName := "hcs_waf_dedicated_domain.domain_1"

	rc := hwacceptance.InitResourceCheck(
		resourceName,
		&obj,
		getWafDedicateDomainResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedDomainV1_basic_withEpsID(randName, acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "domain", fmt.Sprintf("www.%s.com", randName)),
					resource.TestCheckResourceAttr(resourceName, "proxy", "false"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.1"),
					resource.TestCheckResourceAttr(resourceName, "cipher", "cipher_1"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", "${http_host}/error.html"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.0.14"),
					resource.TestCheckResourceAttr(resourceName, "server.0.type", "ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "server.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "certificate_id"),
					resource.TestCheckResourceAttrSet(resourceName, "certificate_name"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttrSet(resourceName, "protect_status"),
					resource.TestCheckResourceAttrSet(resourceName, "protocol"),
					resource.TestCheckResourceAttrSet(resourceName, "tls"),
					resource.TestCheckResourceAttrSet(resourceName, "cipher"),
					resource.TestCheckResourceAttrSet(resourceName, "alarm_page.template_name"),
					resource.TestCheckResourceAttrSet(resourceName, "compliance_certification.pci_3ds"),
					resource.TestCheckResourceAttrSet(resourceName, "compliance_certification.pci_dss"),
				),
			},
			{
				Config: testAccWafDedicatedDomainV1_update_withEpsID(randName, acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.2"),
					resource.TestCheckResourceAttr(resourceName, "cipher", "cipher_2"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", ""),
					resource.TestCheckResourceAttr(resourceName, "pci_3ds", "true"),
					resource.TestCheckResourceAttr(resourceName, "pci_dss", "true"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8443"),
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.0.14"),
					resource.TestCheckResourceAttr(resourceName, "server.1.address", "119.8.0.15"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key2", "$request_length"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key3", "$remote_addr"),
				),
			},
			{
				Config: testAccWafDedicatedDomainV1_policy_withEpsID(randName, acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "domain", fmt.Sprintf("www.%s.com", randName)),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.2"),
					resource.TestCheckResourceAttr(resourceName, "protect_status", "0"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "server.0.type", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.0.14"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"keep_policy"},
				ImportStateIdFunc:       testWAFResourceImportState(resourceName),
			},
		},
	})
}

func testAccWafDedicatedDomainV1_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_waf_dedicated_domain" "domain_1" {
  domain         = "www.%s.com"
  certificate_id = hcs_waf_certificate.certificate_1.id
  keep_policy    = false
  proxy          = false
  tls            = "TLS v1.1"
  cipher         = "cipher_1"
  protect_status = 1

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
    type            = "ipv4"
    vpc_id          = hcs_vpc.test.id
  }

  custom_page {
    http_return_code = "404"
    block_page_type  = "application/json"
    page_content     = <<EOF
{
  "event_id": "$${waf_event_id}",
  "error_msg": "error message"
}
EOF
  }

  forward_header_map = {
    "key1" = "$time_local"
    "key2" = "$tenant_id"
  }

  timeout_settings {
    connection_timeout = 50
    read_timeout       = 200
    write_timeout      = 200
  }

  traffic_mark {
    ip_tags     = ["ip_tag", "$remote_addr"]
    session_tag = "session_tag"
    user_tag    = "user_tag"
  }

  depends_on = [
    hcs_waf_certificate.certificate_1
  ]
}
`, testAccWafCertificateV1_conf(name), name)
}

func testAccWafDedicatedDomainV1_update1(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_waf_dedicated_domain" "domain_1" {
  domain         = "www.%s.com"
  certificate_id = hcs_waf_certificate.certificate_1.id
  keep_policy    = false
  proxy          = true
  tls            = "TLS v1.2"
  cipher         = "cipher_2"
  pci_3ds        = true
  pci_dss        = true
  protect_status = 0
  redirect_url   = "$${http_host}/error.html"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8443
    type            = "ipv4"
    vpc_id          = hcs_vpc.test.id
  }

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.15"
    port            = 8443
    type            = "ipv4"
    vpc_id          = hcs_vpc.test.id
  }

  forward_header_map = {
    "key2" = "$request_length"
    "key3" = "$remote_addr"
  }

  connection_protection {
    error_threshold                       = 1000
    error_percentage                      = 87.5
    initial_downtime                      = 200
    multiplier_for_consecutive_breakdowns = 5
    pending_url_request_threshold         = 7000
    duration                              = 10000
    status                                = true
  }

  timeout_settings {
    connection_timeout = 100
    read_timeout       = 1000
    write_timeout      = 1000
  }

  traffic_mark {
    ip_tags     = ["ip_tag_update", "ip_tag_another"]
    session_tag = "session_tag_update"
    user_tag    = "user_tag_update"
  }

  depends_on = [
    hcs_waf_certificate.certificate_1
  ]
}
`, testAccWafCertificateV1_conf(name), name)
}

func testAccWafDedicatedDomainV1_update2(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_waf_dedicated_domain" "domain_1" {
  domain         = "www.%s.com"
  certificate_id = hcs_waf_certificate.certificate_1.id
  keep_policy    = false
  proxy          = true
  tls            = "TLS v1.2"
  cipher         = "cipher_2"
  pci_3ds        = true
  pci_dss        = true
  protect_status = 0
  redirect_url   = "$${http_host}/error.html"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8443
    type            = "ipv4"
    vpc_id          = hcs_vpc.test.id
  }

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.15"
    port            = 8443
    type            = "ipv4"
    vpc_id          = hcs_vpc.test.id
  }

  forward_header_map = {
    "key2" = "$request_length"
    "key3" = "$remote_addr"
  }

  connection_protection {
    error_threshold                       = 2147483647
    error_percentage                      = 99
    initial_downtime                      = 2147483647
    multiplier_for_consecutive_breakdowns = 2147483647
    pending_url_request_threshold         = 2147483647
    duration                              = 2147483647
    status                                = false
  }

  timeout_settings {
    connection_timeout = 180
    read_timeout       = 3600
    write_timeout      = 3600
  }

  depends_on = [
    hcs_waf_certificate.certificate_1
  ]
}
`, testAccWafCertificateV1_conf(name), name)
}

func testAccWafDedicatedDomainV1_update3(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_waf_dedicated_domain" "domain_1" {
  domain         = "www.%s.com"
  certificate_id = hcs_waf_certificate.certificate_1.id
  keep_policy    = false
  proxy          = true
  tls            = "TLS v1.2"
  cipher         = "cipher_2"
  pci_3ds        = true
  pci_dss        = true
  protect_status = 0
  redirect_url   = "$${http_host}/error.html"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8443
    type            = "ipv4"
    vpc_id          = hcs_vpc.test.id
  }

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.15"
    port            = 8443
    type            = "ipv4"
    vpc_id          = hcs_vpc.test.id
  }

  forward_header_map = {
    "key2" = "$request_length"
    "key3" = "$remote_addr"
  }

  connection_protection {
    error_threshold                       = 0
    error_percentage                      = 0
    initial_downtime                      = 0
    multiplier_for_consecutive_breakdowns = 0
    pending_url_request_threshold         = 0
    duration                              = 0
    status                                = true
  }

  timeout_settings {
    connection_timeout = 0
    read_timeout       = 0
    write_timeout      = 0
  }

  depends_on = [
    hcs_waf_certificate.certificate_1
  ]
}
`, testAccWafCertificateV1_conf(name), name)
}

func testAccWafDedicatedDomainV1_policy(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_waf_policy" "policy_1" {
  name = "%s"

  depends_on = [
    hcs_waf_dedicated_instance.instance_1
  ]
}

resource "hcs_waf_dedicated_domain" "domain_1" {
  domain         = "www.%s.com"
  certificate_id = hcs_waf_certificate.certificate_1.id
  policy_id      = hcs_waf_policy.policy_1.id
  keep_policy    = true
  proxy          = true
  tls            = "TLS v1.2"
  protect_status = 0

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
    type            = "ipv4"
    vpc_id          = hcs_vpc.test.id
  }

  depends_on = [
    hcs_waf_certificate.certificate_1
  ]
}
`, testAccWafCertificateV1_conf(name), name, name)
}

func testAccWafDedicatedDomainV1_basic_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%s

resource "hcs_waf_dedicated_domain" "domain_1" {
  domain                = "www.%s.com"
  certificate_id        = hcs_waf_certificate.certificate_1.id
  keep_policy           = false
  proxy                 = false
  tls                   = "TLS v1.1"
  cipher                = "cipher_1"
  redirect_url          = "$${http_host}/error.html"
  enterprise_project_id = "%s"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
    type            = "ipv4"
    vpc_id          = hcs_vpc.test.id
  }
}
`, testAccWafCertificateV1_conf_withEpsID(name, epsID), name, epsID)
}

func testAccWafDedicatedDomainV1_update_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%s

resource "hcs_waf_dedicated_domain" "domain_1" {
  domain                = "www.%s.com"
  certificate_id        = hcs_waf_certificate.certificate_1.id
  keep_policy           = false
  proxy                 = true
  tls                   = "TLS v1.2"
  cipher                = "cipher_2"
  pci_3ds               = true
  pci_dss               = true
  enterprise_project_id = "%s"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8443
    type            = "ipv4"
    vpc_id          = hcs_vpc.test.id
  }

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.15"
    port            = 8443
    type            = "ipv4"
    vpc_id          = hcs_vpc.test.id
  }

  forward_header_map = {
    "key2" = "$request_length"
    "key3" = "$remote_addr"
  }
}
`, testAccWafCertificateV1_conf_withEpsID(name, epsID), name, epsID)
}

func testAccWafDedicatedDomainV1_policy_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_waf_policy" "policy_1" {
  name                  = "%[2]s"
  enterprise_project_id = "%[3]s"

  depends_on = [
    hcs_waf_certificate.certificate_1
  ]
}

resource "hcs_waf_dedicated_domain" "domain_1" {
  domain                = "www.%[2]s.com"
  certificate_id        = hcs_waf_certificate.certificate_1.id
  policy_id             = hcs_waf_policy.policy_1.id
  keep_policy           = true
  proxy                 = true
  tls                   = "TLS v1.2"
  protect_status        = 0
  enterprise_project_id = "%[3]s"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
    type            = "ipv4"
    vpc_id          = hcs_vpc.test.id
  }
}
`, testAccWafCertificateV1_conf_withEpsID(name, epsID), name, epsID)
}
