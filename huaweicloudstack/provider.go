package huaweicloudstack

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cce"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cfw"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/as"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/bms"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/dns"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/ecs"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/eip"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/elb"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/eps"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/evs"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/ims"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/nat"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/smn"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/vpc"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/vpcep"
)

const (
	defaultCloud string = "myhuaweicloud.com"
)

// Provider returns a schema.Provider for HuaweiCloudStack.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["region"],
				InputDefault: "cn-north-1",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCS_REGION_NAME",
					"OS_REGION_NAME",
				}, nil),
			},

			"access_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["access_key"],
				RequiredWith: []string{"secret_key"},
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCS_ACCESS_KEY",
					"OS_ACCESS_KEY",
				}, nil),
			},

			"secret_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["secret_key"],
				RequiredWith: []string{"access_key"},
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCS_SECRET_KEY",
					"OS_SECRET_KEY",
				}, nil),
			},

			"security_token": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["security_token"],
				RequiredWith: []string{"access_key"},
				DefaultFunc:  schema.EnvDefaultFunc("HCS_SECURITY_TOKEN", nil),
			},

			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["domain_id"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCS_DOMAIN_ID",
					"OS_DOMAIN_ID",
					"OS_USER_DOMAIN_ID",
					"OS_PROJECT_DOMAIN_ID",
				}, ""),
			},

			"domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["domain_name"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCS_DOMAIN_NAME",
					"OS_DOMAIN_NAME",
					"OS_USER_DOMAIN_NAME",
					"OS_PROJECT_DOMAIN_NAME",
				}, ""),
			},

			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["user_name"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCS_USER_NAME",
					"OS_USERNAME",
				}, ""),
			},

			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["user_id"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCS_USER_ID",
					"OS_USER_ID",
				}, ""),
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: descriptions["password"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCS_USER_PASSWORD",
					"OS_PASSWORD",
				}, ""),
			},

			"assume_role": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agency_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: descriptions["assume_role_agency_name"],
							DefaultFunc: schema.EnvDefaultFunc("HCS_ASSUME_ROLE_AGENCY_NAME", nil),
						},
						"domain_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: descriptions["assume_role_domain_name"],
							DefaultFunc: schema.EnvDefaultFunc("HCS_ASSUME_ROLE_DOMAIN_NAME", nil),
						},
					},
				},
			},

			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["project_id"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCS_PROJECT_ID",
					"OS_PROJECT_ID",
				}, nil),
			},

			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["project_name"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCS_PROJECT_NAME",
					"OS_PROJECT_NAME",
				}, nil),
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["tenant_id"],
				DefaultFunc: schema.EnvDefaultFunc("OS_TENANT_ID", ""),
			},

			"tenant_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["tenant_name"],
				DefaultFunc: schema.EnvDefaultFunc("OS_TENANT_NAME", ""),
			},

			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["token"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCS_AUTH_TOKEN",
					"OS_AUTH_TOKEN",
				}, ""),
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: descriptions["insecure"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCS_INSECURE",
					"OS_INSECURE",
				}, false),
			},

			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CACERT", ""),
				Description: descriptions["cacert_file"],
			},

			"cert": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CERT", ""),
				Description: descriptions["cert"],
			},

			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_KEY", ""),
				Description: descriptions["key"],
			},

			"agency_name": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("OS_AGENCY_NAME", nil),
				Description:  descriptions["agency_name"],
				RequiredWith: []string{"agency_domain_name"},
			},

			"agency_domain_name": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("OS_AGENCY_DOMAIN_NAME", nil),
				Description:  descriptions["agency_domain_name"],
				RequiredWith: []string{"agency_name"},
			},

			"delegated_project": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DELEGATED_PROJECT", ""),
				Description: descriptions["delegated_project"],
			},

			"auth_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["auth_url"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCS_AUTH_URL",
					"OS_AUTH_URL",
				}, nil),
			},

			"cloud": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["cloud"],
				DefaultFunc: schema.EnvDefaultFunc("HCS_CLOUD", ""),
			},

			"endpoints": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: descriptions["endpoints"],
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"regional": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: descriptions["regional"],
			},

			"shared_config_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["shared_config_file"],
				DefaultFunc: schema.EnvDefaultFunc("HCS_SHARED_CONFIG_FILE", ""),
			},

			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["profile"],
				DefaultFunc: schema.EnvDefaultFunc("HCS_PROFILE", ""),
			},

			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["enterprise_project_id"],
				DefaultFunc: schema.EnvDefaultFunc("HCS_ENTERPRISE_PROJECT_ID", ""),
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: descriptions["max_retries"],
				DefaultFunc: schema.EnvDefaultFunc("HCS_MAX_RETRIES", 5),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"hcs_enterprise_project":    eps.DataSourceEnterpriseProject(),
			"hcs_vpcep_public_services": vpcep.DataSourceVPCEPPublicServices(),

			"hcs_cfw_firewalls": cfw.DataSourceFirewalls(),

			"hcs_vpc_bandwidth": eip.DataSourceBandWidth(),
			"hcs_vpc_eip":       eip.DataSourceVpcEip(),
			"hcs_vpc_eips":      eip.DataSourceVpcEips(),

			"hcs_evs_volumes": evs.DataSourceEvsVolumesV2(),

			"hcs_nat_gateway": nat.DataSourcePublicGateway(),
			"hcs_smn_topics":  smn.DataSourceTopics(),

			"hcs_ims_images": ims.DataSourceImagesImages(),

			"hcs_vpcs":                   vpc.DataSourceVpcs(),
			"hcs_vpc_subnets":            vpc.DataSourceVpcSubnets(),
			"hcs_vpc_peering_connection": vpc.DataSourceVpcPeeringConnectionV2(),

			"hcs_networking_secgroup":  vpc.DataSourceNetworkingSecGroup(),
			"hcs_networking_secgroups": vpc.DataSourceNetworkingSecGroups(),

			"hcs_availability_zones": DataSourceAvailabilityZones(),
			"hcs_as_configurations":  as.DataSourceASConfigurations(),
			"hcs_as_groups":          as.DataSourceASGroups(),

			"hcs_bms_flavors": bms.DataSourceBmsFlavors(),

			"hcs_cce_cluster":        cce.DataSourceCCEClusterV3(),
			"hcs_cce_clusters":       cce.DataSourceCCEClusters(),
			"hcs_cce_addon_template": cce.DataSourceAddonTemplate(),
			"hcs_cce_node_pool":      cce.DataSourceCCENodePoolV3(),
			"hcs_cce_node":           cce.DataSourceNode(),
			"hcs_cce_nodes":          cce.DataSourceNodes(),

			"hcs_ecs_compute_flavors":      ecs.DataSourceEcsFlavors(),
			"hcs_ecs_compute_instance":     ecs.DataSourceComputeInstance(),
			"hcs_ecs_compute_instances":    ecs.DataSourceComputeInstances(),
			"hcs_ecs_compute_servergroups": ecs.DataSourceComputeServerGroups(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"hcs_cce_addon":       cce.ResourceAddon(),
			"hcs_cce_cluster":     cce.ResourceCluster(),
			"hcs_cce_namespace":   cce.ResourceCCENamespaceV1(),
			"hcs_cce_node":        cce.ResourceNode(),
			"hcs_cce_node_attach": cce.ResourceNodeAttach(),
			"hcs_cce_node_pool":   cce.ResourceNodePool(),
			"hcs_cce_pvc":         cce.ResourceCcePersistentVolumeClaimsV1(),

			"hcs_cfw_address_group":        cfw.ResourceAddressGroup(),
			"hcs_cfw_address_group_member": cfw.ResourceAddressGroupMember(),
			"hcs_cfw_black_white_list":     cfw.ResourceBlackWhiteList(),
			"hcs_cfw_eip_protection":       cfw.ResourceEipProtection(),
			"hcs_cfw_protection_rule":      cfw.ResourceProtectionRule(),
			"hcs_cfw_service_group_member": cfw.ResourceServiceGroupMember(),
			"hcs_cfw_service_group":        cfw.ResourceServiceGroup(),

			"hcs_enterprise_project": eps.ResourceEnterpriseProject(),
			"hcs_dns_recordset":      dns.ResourceDNSRecordset(),
			"hcs_dns_zone":           dns.ResourceDNSZone(),
			"hcs_vpcep_approval":     vpcep.ResourceVPCEndpointApproval(),
			"hcs_vpcep_endpoint":     vpcep.ResourceVPCEndpoint(),
			"hcs_vpcep_service":      vpcep.ResourceVPCEndpointService(),

			"hcs_vpc_bandwidth":           eip.ResourceVpcBandWidthV2(),
			"hcs_vpc_eip":                 eip.ResourceVpcEIPV1(),
			"hcs_vpc_eip_associate":       eip.ResourceEIPAssociate(),
			"hcs_vpc_bandwidth_associate": eip.ResourceBandWidthAssociate(),

			"hcs_vpc_bandwidth_v2": eip.ResourceVpcBandWidthV2(),
			"hcs_vpc_eip_v1":       eip.ResourceVpcEIPV1(),

			"hcs_evs_volume":   evs.ResourceEvsVolume(),
			"hcs_evs_snapshot": evs.ResourceEvsSnapshotV2(),

			"hcs_elb_certificate":     elb.ResourceCertificateV3(),
			"hcs_elb_l7policy":        elb.ResourceL7PolicyV3(),
			"hcs_elb_l7rule":          elb.ResourceL7RuleV3(),
			"hcs_elb_listener":        elb.ResourceListenerV3(),
			"hcs_elb_loadbalancer":    elb.ResourceLoadBalancerV3(),
			"hcs_elb_member":          elb.ResourceMemberV3(),
			"hcs_elb_monitor":         elb.ResourceMonitorV3(),
			"hcs_elb_pool":            elb.ResourcePoolV3(),
			"hcs_elb_security_policy": elb.ResourceSecurityPolicy(),

			"hcs_ecs_compute_volume_attach":    ecs.ResourceComputeVolumeAttach(),
			"hcs_ecs_compute_server_group":     ecs.ResourceComputeServerGroup(),
			"hcs_ecs_compute_interface_attach": ecs.ResourceComputeInterfaceAttach(),
			"hcs_ecs_compute_instance":         ecs.ResourceComputeInstance(),

			// Legacy
			"hcs_networking_eip_associate": eip.ResourceEIPAssociate(),

			"hcs_nat_gateway":   nat.ResourcePublicGateway(),
			"hcs_nat_snat_rule": nat.ResourcePublicSnatRule(),
			"hcs_nat_dnat_rule": nat.ResourcePublicDnatRule(),

			"hcs_smn_topic":            smn.ResourceTopic(),
			"hcs_smn_subscription":     smn.ResourceSubscription(),
			"hcs_smn_message_template": smn.ResourceSmnMessageTemplate(),
			"hcs_smn_topic_v2":         smn.ResourceTopic(),
			"hcs_smn_subscription_v2":  smn.ResourceSubscription(),

			"hcs_as_bandwidth_policy": as.ResourceASBandWidthPolicy(),
			"hcs_as_configuration":    as.ResourceASConfiguration(),
			"hcs_as_group":            as.ResourceASGroup(),
			"hcs_as_instance_attach":  as.ResourceASInstanceAttach(),
			"hcs_as_lifecycle_hook":   as.ResourceASLifecycleHook(),
			"hcs_as_notification":     as.ResourceAsNotification(),
			"hcs_as_policy":           as.ResourceASPolicy(),

			"hcs_ims_image":                ims.ResourceImsImage(),
			"hcs_ims_image_share":          ims.ResourceImsImageShare(),
			"hcs_ims_image_share_accepter": ims.ResourceImsImageShareAccepter(),

			"hcs_vpc":        vpc.ResourceVirtualPrivateCloudV1(),
			"hcs_vpc_subnet": vpc.ResourceVpcSubnetV1(),

			"hcs_vpc_route_v2":                    vpc.ResourceVPCRouteV2(),
			"hcs_vpc_v1":                          vpc.ResourceVirtualPrivateCloudV1(),
			"hcs_vpc_subnet_v1":                   vpc.ResourceVpcSubnetV1(),
			"hcs_networking_vip":                  vpc.ResourceNetworkingVip(),
			"hcs_networking_vip_associate":        vpc.ResourceNetworkingVIPAssociateV2(),
			"hcs_vpc_peering_connection":          vpc.ResourceVpcPeeringConnectionV2(),
			"hcs_vpc_peering_connection_accepter": vpc.ResourceVpcPeeringConnectionAccepterV2(),

			"hcs_network_acl":              ResourceNetworkACL(),
			"hcs_network_acl_rule":         ResourceNetworkACLRule(),
			"hcs_networking_secgroup":      vpc.ResourceNetworkingSecGroup(),
			"hcs_networking_secgroup_rule": vpc.ResourceNetworkingSecGroupRule(),

			"hcs_bms_instance": bms.ResourceBmsInstance(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11 cc
			terraformVersion = "0.11+compatible"
		}

		return configureProvider(ctx, d, terraformVersion)
	}

	return provider
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"auth_url": "The Identity authentication URL.",

		"region": "The HuaweiCloudStack region to connect to.",

		"user_name": "Username to login with.",

		"user_id": "User ID to login with.",

		"project_id": "The ID of the project to login with.",

		"project_name": "The name of the project to login with.",

		"tenant_id": "The ID of the Tenant (Identity v2) to login with.",

		"tenant_name": "The name of the Tenant (Identity v2) to login with.",

		"password": "Password to login with.",

		"token": "Authentication token to use as an alternative to username/password.",

		"domain_id": "The ID of the Domain to scope to.",

		"domain_name": "The name of the Domain to scope to.",

		"access_key":     "The access key of the HuaweiCloudStack to use.",
		"secret_key":     "The secret key of the HuaweiCloudStack to use.",
		"security_token": "The security token to authenticate with a temporary security credential.",

		"insecure": "Trust self-signed certificates.",

		"cacert_file": "A Custom CA certificate.",

		"cert": "A client certificate to authenticate with.",

		"key": "A client private key to authenticate with.",

		"agency_name": "The name of agency",

		"agency_domain_name": "The name of domain who created the agency (Identity v3).",

		"delegated_project": "The name of delegated project (Identity v3).",

		"assume_role_agency_name": "The name of agency for assume role.",

		"assume_role_domain_name": "The name of domain for assume role.",

		"cloud": "The endpoint of cloud provider, defaults to myhuaweicloud.com",

		"endpoints": "The custom endpoints used to override the default endpoint URL.",

		"regional": "Whether the service endpoints are regional",

		"shared_config_file": "The path to the shared config file. If not set, the default is ~/.hcloud/config.json.",

		"profile": "The profile name as set in the shared config file.",

		"max_retries": "How many times HTTP connection should be retried until giving up.",

		"enterprise_project_id": "enterprise project id",
	}
}

func configureProvider(_ context.Context, d *schema.ResourceData, terraformVersion string) (interface{},
	diag.Diagnostics) {
	var tenantName, tenantID, delegatedProject, identityEndpoint string
	region := d.Get("region").(string)
	isRegional := d.Get("regional").(bool)
	// different from hws, there is no default "cloud" in hcs, throw if not provided.
	cloud := d.Get("cloud").(string)

	// project_name is prior to tenant_name
	// if neither of them was set, use region as the default project
	if v, ok := d.GetOk("project_name"); ok && v.(string) != "" {
		tenantName = v.(string)
	} else if v, ok := d.GetOk("tenant_name"); ok && v.(string) != "" {
		tenantName = v.(string)
	} else {
		tenantName = region
	}

	// project_id is prior to tenant_id
	if v, ok := d.GetOk("project_id"); ok && v.(string) != "" {
		tenantID = v.(string)
	} else {
		tenantID = d.Get("tenant_id").(string)
	}

	// Use region as delegated_project if it's not set
	if v, ok := d.GetOk("delegated_project"); ok && v.(string) != "" {
		delegatedProject = v.(string)
	} else {
		delegatedProject = region
	}

	// use auth_url as identityEndpoint if specified
	if v, ok := d.GetOk("auth_url"); ok {
		identityEndpoint = v.(string)
	} else {
		// use cloud as basis for identityEndpoint
		identityEndpoint = fmt.Sprintf("https://iam-apigateway-proxy.%s:443/v3", cloud)
	}

	hcsConfig := config.HcsConfig{
		Config: config.Config{
			AccessKey:           d.Get("access_key").(string),
			SecretKey:           d.Get("secret_key").(string),
			CACertFile:          d.Get("cacert_file").(string),
			ClientCertFile:      d.Get("cert").(string),
			ClientKeyFile:       d.Get("key").(string),
			DomainID:            d.Get("domain_id").(string),
			DomainName:          d.Get("domain_name").(string),
			IdentityEndpoint:    identityEndpoint,
			Insecure:            d.Get("insecure").(bool),
			Password:            d.Get("password").(string),
			Token:               d.Get("token").(string),
			SecurityToken:       d.Get("security_token").(string),
			Region:              region,
			TenantID:            tenantID,
			TenantName:          tenantName,
			Username:            d.Get("user_name").(string),
			UserID:              d.Get("user_id").(string),
			AgencyName:          d.Get("agency_name").(string),
			AgencyDomainName:    d.Get("agency_domain_name").(string),
			DelegatedProject:    delegatedProject,
			Cloud:               cloud,
			RegionClient:        isRegional,
			MaxRetries:          d.Get("max_retries").(int),
			EnterpriseProjectID: d.Get("enterprise_project_id").(string),
			SharedConfigFile:    d.Get("shared_config_file").(string),
			Profile:             d.Get("profile").(string),
			TerraformVersion:    terraformVersion,
			RegionProjectIDMap:  make(map[string]string),
			RPLock:              new(sync.Mutex),
			SecurityKeyLock:     new(sync.Mutex),
		},
	}

	// Save hcsConfig to config.Config for extend
	hcsConfig.Metadata = &hcsConfig

	// get assume role
	assumeRoleList := d.Get("assume_role").([]interface{})
	if len(assumeRoleList) == 1 {
		assumeRole := assumeRoleList[0].(map[string]interface{})
		hcsConfig.AssumeRoleAgency = assumeRole["agency_name"].(string)
		hcsConfig.AssumeRoleDomain = assumeRole["domain_name"].(string)
	}

	// get custom endpoints
	endpoints, err := flattenProviderEndpoints(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	hcsConfig.Endpoints = endpoints

	if err := hcsConfig.LoadAndValidate(); err != nil {
		return nil, diag.FromErr(err)
	}

	return &hcsConfig.Config, nil
}

func flattenProviderEndpoints(d *schema.ResourceData) (map[string]string, error) {
	endpoints := d.Get("endpoints").(map[string]interface{})
	epMap := make(map[string]string)

	for key, val := range endpoints {
		endpoint := strings.TrimSpace(val.(string))
		// check empty string
		if endpoint == "" {
			return nil, fmt.Errorf("the value of customer endpoint %s must be specified", key)
		}

		// add prefix "https://" and suffix "/"
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", endpoint)
		}
		if !strings.HasSuffix(endpoint, "/") {
			endpoint = fmt.Sprintf("%s/", endpoint)
		}
		epMap[key] = endpoint
	}

	// unify the endpoint which has multiple versions
	for key := range endpoints {
		ep, ok := epMap[key]
		if !ok {
			continue
		}

		multiKeys := config.GetServiceDerivedCatalogKeys(key)
		for _, k := range multiKeys {
			epMap[k] = ep
		}
	}

	log.Printf("[DEBUG] customer endpoints: %+v", epMap)
	return epMap, nil
}

func getCloudDomain(cloud, region string) string {
	// first, use the specified value
	if cloud != "" {
		return cloud
	}

	return defaultCloud
}

func isGlobalIamEndpoint(domain, region string, isRegional bool) bool {
	if !isRegional && domain == defaultCloud && !strings.HasPrefix(region, "eu-west-1") {
		return true
	}

	return false
}
