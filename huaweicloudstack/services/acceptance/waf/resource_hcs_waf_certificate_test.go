package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/waf/v1/certificates"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	hwacceptance "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.WafV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF client: %s", err)
	}
	return certificates.GetWithEpsID(client, state.Primary.ID, state.Primary.Attributes["enterprise_project_id"]).Extract()
}

func TestAccWafCertificateV1_basic(t *testing.T) {
	var certificate certificates.Certificate
	resourceName := "hcs_waf_certificate.certificate_1"
	name := acceptance.RandomAccResourceName()

	rc := hwacceptance.InitResourceCheck(
		resourceName,
		&certificate,
		getResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafCertificateV1_conf(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				Config: testAccWafCertificateV1_conf_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", name)),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate", "private_key"},
			},
		},
	})
}

func TestAccWafCertificateV1_withEpsID(t *testing.T) {
	var certificate certificates.Certificate
	resourceName := "hcs_waf_certificate.certificate_1"
	name := acceptance.RandomAccResourceName()
	updateName := name + "_update"

	rc := hwacceptance.InitResourceCheck(
		resourceName,
		&certificate,
		getResourceFunc,
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
				Config: testAccWafCertificateV1_conf_withEpsID(name, acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				Config: testAccWafCertificateV1_conf_withEpsID(updateName, acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate", "private_key"},
				ImportStateIdFunc:       testWAFResourceImportState(resourceName),
			},
		},
	})
}

func testAccWafCertificateV1_conf(name string) string {
	return fmt.Sprintf(`

provider "hcs" {
  endpoints = {
    "iam" = "iam-apigateway-proxy.outfullstack.com"
    "vpc" = "vpc.cn-fullstack-1.outfullstack.com"
    "ecs" = "ecs.cn-fullstack-1.outfullstack.com"
    "ims" = "ims.cn-fullstack-1.outfullstack.com"

    "mrs" = "mrs.cn-fullstack-1.outfullstack.com"
    "waf" = "waf-api.cn-fullstack-1.outfullstack.com"
  }
}

resource "hcs_waf_certificate" "certificate_1" {
  name = "%s"

  certificate = <<EOT
-----BEGIN CERTIFICATE-----
MIIFszCCA5ugAwIBAgIUKrTehAfpjNDrCg2J25S6qmZ6oMAwDQYJKoZIhvcNAQEL
BQAwaTELMAkGA1UEBhMCdGYxCzAJBgNVBAgMAnRmMQswCQYDVQQHDAJ0ZjELMAkG
A1UECgwCdGYxCzAJBgNVBAsMAnRmMQswCQYDVQQDDAJ0ZjEZMBcGCSqGSIb3DQEJ
ARYKdGZAMTYzLmNvbTAeFw0yMzEyMTkwODQ1NTlaFw0zMzEyMTYwODQ1NTlaMGkx
CzAJBgNVBAYTAnRmMQswCQYDVQQIDAJ0ZjELMAkGA1UEBwwCdGYxCzAJBgNVBAoM
AnRmMQswCQYDVQQLDAJ0ZjELMAkGA1UEAwwCdGYxGTAXBgkqhkiG9w0BCQEWCnRm
QDE2My5jb20wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDtXOqiaxN5
EbNCKBz1KhCFaSZ7yQ7K3/JPSbHPp2CAa4A87Bqs0On3fHGwEzjFXTOYZPB2+rwA
v9pm+6dIWGbp5szbFScoQAuRLXMPfUFW+O6TZmJtT3rrlLE1NKmAfTvPFkMUkjHT
289vIrCgTIBWTVSUmBVT7xGev743VIBUURp2071GXdAuaUom849DOk82BxOfYowl
QiYLpwrWnYH83DcvxjGaeSBX/xAY6hqF9jNsYw+nuDR3owuNg0N7l216G5HIoMvW
i8UihJ8Bv3zaFL5wbWDvU0qu3sAzbRrq46KuyihAJwQcaB1nWAkGaCRU3K/RuCi6
sKP7nXGBldnXNVJiYRcKhiZ+ZWEQQcir2zhSUHSP0v3yMZdJvMVK1WSRu3eUyJs6
ViW7bprUoyoQ8IHR2AWv8iApdtLpcHTZuaiX3EH9R1bFlYVV7Ak9ABIq0Wg7IwkE
J4y9Gdd/vwM/jVOjwC++OColmEffb7aG/99C+q5oa7O+LgglCxD7JMJ8+REWza8Z
ZarG6SkT8xdSXxqjGXOrYQcHqeh+UA5p7T3RnzIYhGCHjfXBCM8AH9E+qGhwy1BR
fd+cvwoj6MvYS1PwfClncxGlPS9sOQsCXPzJ4GJHa88JHYELcBhwIk1QAZAOrUD2
igaM9ZwEVrCdxvuYylFi1xYDy3rbsOwU+QIDAQABo1MwUTAdBgNVHQ4EFgQUFMuJ
cLmDRaokf4YXOMj9hfX3DOwwHwYDVR0jBBgwFoAUFMuJcLmDRaokf4YXOMj9hfX3
DOwwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAgEAEk88R6ERC6Er
IMuXMCHOl9j3KTp6TiaCnWUWsHpWF9qsULkDbcW5ikoaB96qie2ELIRxXlL3joC+
q5n0NshaRZ0c5SsEv5HcuG6djwTXX17z5jHhbCSMVzt6i24D9Hsv7f0JZl6AZ6OO
IPK7T+SoRZQ3WiNNWUBhGUt0zV0YHEPt25CWyAwLtBeg5Z3nmh46cKtcNXA6wpuj
BfLunyew4sKGhXMFdtftpkLOX9/RRaHSVRianm9JFQ5un5XQ5+87G5ePbeQ4xelo
Q8mVrbIKnpnkOONq/JHwDgjp/7XbWQpvPzbUO9dB11tatWbxaQR3xgRB5X6Hm54q
IES5d61rSjQKZIbbIlF51HqKb0TX6tIcItlwCrGcKDfkpmgYy3PI/AnrIrTwrlpX
9hUHx/LZCGYIG1jmNVWtBeHhPfMdxTOxwhpGZesgQtOeoIB8JxLvEZZQZcsfFF5k
Bl99QVqEgidR1jFCwVMlopPBBTjaCoqZcaS2PH0zBcrCuLLkZEWdEXeRmISpwQwJ
Z/Pvg/mBu5xg42C42wbCqUx33jiiDSMZlwRK4kbOXw7laYt4Xz1ly9Qcrmudzc/V
c5c/G15hEsI1avK/OC4FOQZcv+ZEGgmIUM528HPHGkoRO4rbCFSYTFdnxsPJE8TJ
WiXOtq+AqUUEGfrMipeoqUTpP5S5V70=
-----END CERTIFICATE-----
EOT

  private_key = <<EOT
-----BEGIN PRIVATE KEY-----
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
oqUbExjgBr9xSdGJbv4SBYskYS9e4i0=
-----END PRIVATE KEY-----
EOT
}
`, name)
}

func testAccWafCertificateV1_conf_update(name string) string {
	return fmt.Sprintf(`

  endpoints = {
    "iam" = "iam-apigateway-proxy.outfullstack.com"
    "vpc" = "vpc.cn-fullstack-1.outfullstack.com"
    "ecs" = "ecs.cn-fullstack-1.outfullstack.com"
    "ims" = "ims.cn-fullstack-1.outfullstack.com"

    "mrs" = "mrs.cn-fullstack-1.outfullstack.com"
    "waf" = "waf-api.cn-fullstack-1.outfullstack.com"
  }

resource "hcs_waf_certificate" "certificate_1" {
  name = "%s_update"

  certificate = <<EOT
-----BEGIN CERTIFICATE-----
MIIFszCCA5ugAwIBAgIUKrTehAfpjNDrCg2J25S6qmZ6oMAwDQYJKoZIhvcNAQEL
BQAwaTELMAkGA1UEBhMCdGYxCzAJBgNVBAgMAnRmMQswCQYDVQQHDAJ0ZjELMAkG
A1UECgwCdGYxCzAJBgNVBAsMAnRmMQswCQYDVQQDDAJ0ZjEZMBcGCSqGSIb3DQEJ
ARYKdGZAMTYzLmNvbTAeFw0yMzEyMTkwODQ1NTlaFw0zMzEyMTYwODQ1NTlaMGkx
CzAJBgNVBAYTAnRmMQswCQYDVQQIDAJ0ZjELMAkGA1UEBwwCdGYxCzAJBgNVBAoM
AnRmMQswCQYDVQQLDAJ0ZjELMAkGA1UEAwwCdGYxGTAXBgkqhkiG9w0BCQEWCnRm
QDE2My5jb20wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDtXOqiaxN5
EbNCKBz1KhCFaSZ7yQ7K3/JPSbHPp2CAa4A87Bqs0On3fHGwEzjFXTOYZPB2+rwA
v9pm+6dIWGbp5szbFScoQAuRLXMPfUFW+O6TZmJtT3rrlLE1NKmAfTvPFkMUkjHT
289vIrCgTIBWTVSUmBVT7xGev743VIBUURp2071GXdAuaUom849DOk82BxOfYowl
QiYLpwrWnYH83DcvxjGaeSBX/xAY6hqF9jNsYw+nuDR3owuNg0N7l216G5HIoMvW
i8UihJ8Bv3zaFL5wbWDvU0qu3sAzbRrq46KuyihAJwQcaB1nWAkGaCRU3K/RuCi6
sKP7nXGBldnXNVJiYRcKhiZ+ZWEQQcir2zhSUHSP0v3yMZdJvMVK1WSRu3eUyJs6
ViW7bprUoyoQ8IHR2AWv8iApdtLpcHTZuaiX3EH9R1bFlYVV7Ak9ABIq0Wg7IwkE
J4y9Gdd/vwM/jVOjwC++OColmEffb7aG/99C+q5oa7O+LgglCxD7JMJ8+REWza8Z
ZarG6SkT8xdSXxqjGXOrYQcHqeh+UA5p7T3RnzIYhGCHjfXBCM8AH9E+qGhwy1BR
fd+cvwoj6MvYS1PwfClncxGlPS9sOQsCXPzJ4GJHa88JHYELcBhwIk1QAZAOrUD2
igaM9ZwEVrCdxvuYylFi1xYDy3rbsOwU+QIDAQABo1MwUTAdBgNVHQ4EFgQUFMuJ
cLmDRaokf4YXOMj9hfX3DOwwHwYDVR0jBBgwFoAUFMuJcLmDRaokf4YXOMj9hfX3
DOwwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAgEAEk88R6ERC6Er
IMuXMCHOl9j3KTp6TiaCnWUWsHpWF9qsULkDbcW5ikoaB96qie2ELIRxXlL3joC+
q5n0NshaRZ0c5SsEv5HcuG6djwTXX17z5jHhbCSMVzt6i24D9Hsv7f0JZl6AZ6OO
IPK7T+SoRZQ3WiNNWUBhGUt0zV0YHEPt25CWyAwLtBeg5Z3nmh46cKtcNXA6wpuj
BfLunyew4sKGhXMFdtftpkLOX9/RRaHSVRianm9JFQ5un5XQ5+87G5ePbeQ4xelo
Q8mVrbIKnpnkOONq/JHwDgjp/7XbWQpvPzbUO9dB11tatWbxaQR3xgRB5X6Hm54q
IES5d61rSjQKZIbbIlF51HqKb0TX6tIcItlwCrGcKDfkpmgYy3PI/AnrIrTwrlpX
9hUHx/LZCGYIG1jmNVWtBeHhPfMdxTOxwhpGZesgQtOeoIB8JxLvEZZQZcsfFF5k
Bl99QVqEgidR1jFCwVMlopPBBTjaCoqZcaS2PH0zBcrCuLLkZEWdEXeRmISpwQwJ
Z/Pvg/mBu5xg42C42wbCqUx33jiiDSMZlwRK4kbOXw7laYt4Xz1ly9Qcrmudzc/V
c5c/G15hEsI1avK/OC4FOQZcv+ZEGgmIUM528HPHGkoRO4rbCFSYTFdnxsPJE8TJ
WiXOtq+AqUUEGfrMipeoqUTpP5S5V70=
-----END CERTIFICATE-----
EOT

  private_key = <<EOT
-----BEGIN PRIVATE KEY-----
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
oqUbExjgBr9xSdGJbv4SBYskYS9e4i0=
-----END PRIVATE KEY-----
EOT
}
`, name)
}

func testAccWafCertificateV1_conf_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%s

resource "hcs_waf_certificate" "certificate_1" {
  name                  = "%s"
  enterprise_project_id = "%s"

  certificate = <<EOT
-----BEGIN CERTIFICATE-----
MIIFszCCA5ugAwIBAgIUKrTehAfpjNDrCg2J25S6qmZ6oMAwDQYJKoZIhvcNAQEL
BQAwaTELMAkGA1UEBhMCdGYxCzAJBgNVBAgMAnRmMQswCQYDVQQHDAJ0ZjELMAkG
A1UECgwCdGYxCzAJBgNVBAsMAnRmMQswCQYDVQQDDAJ0ZjEZMBcGCSqGSIb3DQEJ
ARYKdGZAMTYzLmNvbTAeFw0yMzEyMTkwODQ1NTlaFw0zMzEyMTYwODQ1NTlaMGkx
CzAJBgNVBAYTAnRmMQswCQYDVQQIDAJ0ZjELMAkGA1UEBwwCdGYxCzAJBgNVBAoM
AnRmMQswCQYDVQQLDAJ0ZjELMAkGA1UEAwwCdGYxGTAXBgkqhkiG9w0BCQEWCnRm
QDE2My5jb20wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDtXOqiaxN5
EbNCKBz1KhCFaSZ7yQ7K3/JPSbHPp2CAa4A87Bqs0On3fHGwEzjFXTOYZPB2+rwA
v9pm+6dIWGbp5szbFScoQAuRLXMPfUFW+O6TZmJtT3rrlLE1NKmAfTvPFkMUkjHT
289vIrCgTIBWTVSUmBVT7xGev743VIBUURp2071GXdAuaUom849DOk82BxOfYowl
QiYLpwrWnYH83DcvxjGaeSBX/xAY6hqF9jNsYw+nuDR3owuNg0N7l216G5HIoMvW
i8UihJ8Bv3zaFL5wbWDvU0qu3sAzbRrq46KuyihAJwQcaB1nWAkGaCRU3K/RuCi6
sKP7nXGBldnXNVJiYRcKhiZ+ZWEQQcir2zhSUHSP0v3yMZdJvMVK1WSRu3eUyJs6
ViW7bprUoyoQ8IHR2AWv8iApdtLpcHTZuaiX3EH9R1bFlYVV7Ak9ABIq0Wg7IwkE
J4y9Gdd/vwM/jVOjwC++OColmEffb7aG/99C+q5oa7O+LgglCxD7JMJ8+REWza8Z
ZarG6SkT8xdSXxqjGXOrYQcHqeh+UA5p7T3RnzIYhGCHjfXBCM8AH9E+qGhwy1BR
fd+cvwoj6MvYS1PwfClncxGlPS9sOQsCXPzJ4GJHa88JHYELcBhwIk1QAZAOrUD2
igaM9ZwEVrCdxvuYylFi1xYDy3rbsOwU+QIDAQABo1MwUTAdBgNVHQ4EFgQUFMuJ
cLmDRaokf4YXOMj9hfX3DOwwHwYDVR0jBBgwFoAUFMuJcLmDRaokf4YXOMj9hfX3
DOwwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAgEAEk88R6ERC6Er
IMuXMCHOl9j3KTp6TiaCnWUWsHpWF9qsULkDbcW5ikoaB96qie2ELIRxXlL3joC+
q5n0NshaRZ0c5SsEv5HcuG6djwTXX17z5jHhbCSMVzt6i24D9Hsv7f0JZl6AZ6OO
IPK7T+SoRZQ3WiNNWUBhGUt0zV0YHEPt25CWyAwLtBeg5Z3nmh46cKtcNXA6wpuj
BfLunyew4sKGhXMFdtftpkLOX9/RRaHSVRianm9JFQ5un5XQ5+87G5ePbeQ4xelo
Q8mVrbIKnpnkOONq/JHwDgjp/7XbWQpvPzbUO9dB11tatWbxaQR3xgRB5X6Hm54q
IES5d61rSjQKZIbbIlF51HqKb0TX6tIcItlwCrGcKDfkpmgYy3PI/AnrIrTwrlpX
9hUHx/LZCGYIG1jmNVWtBeHhPfMdxTOxwhpGZesgQtOeoIB8JxLvEZZQZcsfFF5k
Bl99QVqEgidR1jFCwVMlopPBBTjaCoqZcaS2PH0zBcrCuLLkZEWdEXeRmISpwQwJ
Z/Pvg/mBu5xg42C42wbCqUx33jiiDSMZlwRK4kbOXw7laYt4Xz1ly9Qcrmudzc/V
c5c/G15hEsI1avK/OC4FOQZcv+ZEGgmIUM528HPHGkoRO4rbCFSYTFdnxsPJE8TJ
WiXOtq+AqUUEGfrMipeoqUTpP5S5V70=
-----END CERTIFICATE-----
EOT

  private_key = <<EOT
-----BEGIN PRIVATE KEY-----
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
oqUbExjgBr9xSdGJbv4SBYskYS9e4i0=
-----END PRIVATE KEY-----
EOT

  depends_on = [
    hcs_waf_dedicated_instance.instance_1
  ]
}
`, testAccWafDedicatedInstance_epsId(name, epsID), name, epsID)
}
