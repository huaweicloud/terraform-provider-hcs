package config

import (
	"fmt"
)

// ServiceCatalog defines a struct which was used to generate a service client for huaweicloudstack.
// the endpoint likes https://{Name}.{Region}.myhuaweicloud.com/{Version}/{project_id}/{ResourceBase}
// For more information, please refer to Config.NewServiceClient
type ServiceCatalog struct {
	Name             string
	Version          string
	Scope            string
	Admin            bool
	ResourceBase     string
	WithOutProjectID bool
	Product          string
}

// multiCatalogKeys is a map of primary and derived catalog keys for services with multiple clients.
// If we add another version of a service client, don't forget to update it.
var multiCatalogKeys = map[string][]string{
	"iam-apigateway-proxy": {"identity", "iam_no_version"},
	"bss":                  {"bssv2"},
	"ecs":                  {"ecsv21", "ecsv11"},
	"evs":                  {"evsv21"},
	"cce":                  {"ccev1", "cce_addon"},
	"cci":                  {"cciv1_bata"},
	"vpc":                  {"networkv2", "vpcv3", "fwv2"},
	"elb":                  {"elbv2", "elbv3"},
	"dns":                  {"dns_region"},
	"kms":                  {"kmsv1", "kmsv3"},
	"mrs":                  {"mrsv2"},
	"rds":                  {"rdsv1"},
	"waf":                  {"waf-dedicated"},
	"geminidb":             {"geminidbv31"},
	"dli":                  {"dliv2"},
	"dcs":                  {"dcsv1"},
	"dis":                  {"disv3"},
	"dms":                  {"dmsv2"},
	"dws":                  {"dwsv2"},
	"apig":                 {"apigv2"},
	"modelarts":            {"modelartsv2"},
	"servicestage":         {"servicestagev2"},
	"smn":                  {"smn-tag"},
	"ces":                  {"cesv2"},
	"ims":                  {"imsv1"},
	"roma":                 {"romav2"},
}

// GetServiceDerivedCatalogKeys returns the derived catalog keys of a service.
func GetServiceDerivedCatalogKeys(mainKey string) []string {
	return multiCatalogKeys[mainKey]
}

var allServiceCatalog = map[string]ServiceCatalog{
	// catalog for global service
	// identity is used for openstack keystone APIs
	"identity": {
		Name:             "iam-apigateway-proxy",
		Version:          "v3",
		Admin:            true,
		WithOutProjectID: true,
		Product:          "IAM",
	},
	"iam_no_version": {
		Name:             "iam-apigateway-proxy",
		Version:          "",
		Admin:            true,
		WithOutProjectID: true,
		Product:          "IAM",
	},
	// iam is used for huaweicloudstack IAM APIs
	"iam": {
		Name:             "iam-apigateway-proxy",
		Version:          "v3.0",
		Admin:            true,
		WithOutProjectID: true,
		Product:          "IAM",
	},
	"cdn": {
		Name:             "cdn",
		Version:          "v1.0",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "CDN",
	},
	"eps": {
		Name:             "sc",
		Version:          "v1.0",
		Scope:            "global",
		Admin:            true,
		WithOutProjectID: true,
		Product:          "EPS",
	},

	// ******* catalog for vdcUser *******
	"vdc": {
		Name:             "sc",
		Scope:            "global",
		Admin:            true,
		WithOutProjectID: true,
		Product:          "VDC",
	},
	"bss": {
		Name:             "bss",
		Version:          "v1.0",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "BSS",
	},
	"bssv2": {
		Name:             "bss",
		Version:          "v2",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "BSS",
	},

	// ******* catalog for Compute *******
	"ecs": {
		Name:    "ecs",
		Version: "v1",
		Product: "ECS",
	},
	"ecsv11": {
		Name:    "ecs",
		Version: "v1.1",
		Product: "ECS",
	},
	"ecsv21": {
		Name:    "ecs",
		Version: "v2.1",
		Product: "ECS",
	},
	"autoscaling": {
		Name:    "as",
		Version: "autoscaling-api/v1",
		Product: "AS",
	},
	"imsv1": {
		Name:             "ims",
		Version:          "v1",
		WithOutProjectID: true,
		Product:          "IMS",
	},
	"ims": {
		Name:             "ims",
		Version:          "v2",
		WithOutProjectID: true,
		Product:          "IMS",
	},
	"ccev1": {
		Name:             "cce",
		Version:          "api/v1",
		WithOutProjectID: true,
		Product:          "CCE",
	},
	"cce": {
		Name:    "cce",
		Version: "api/v3/projects",
		Product: "CCE",
	},
	"cce_addon": {
		Name:             "cce",
		Version:          "api/v3",
		WithOutProjectID: true,
		Product:          "CCE",
	},
	"aom": {
		Name:    "aom",
		Version: "svcstg/icmgr/v1",
		Product: "AOM",
	},
	"cci": {
		Name:             "cci",
		Version:          "api/v1",
		WithOutProjectID: true,
		Product:          "CCI",
	},
	"cciv1_bata": {
		Name:             "cci",
		Version:          "apis/networking.cci.io/v1beta1",
		WithOutProjectID: true,
		Product:          "CCI",
	},
	"fgs": {
		Name:    "functiongraph",
		Version: "v2",
		Product: "FunctionGraph",
	},
	"swr": {
		Name:             "swr-api",
		Version:          "v2",
		WithOutProjectID: true,
		Product:          "SWR",
	},
	"bms": {
		Name:    "bms",
		Version: "v1",
		Product: "BMS",
	},
	"aos": {
		Name:    "aos",
		Version: "v1",
		Product: "AOS",
	},

	// ******* catalog for storage ******
	"evs": {
		Name:    "evs",
		Version: "v2",
		Product: "EVS",
	},
	"evsv21": {
		Name:    "evs",
		Version: "v2.1",
		Product: "EVS",
	},
	"sfs": {
		Name:    "sfs",
		Version: "v2",
		Product: "SFS",
	},
	"sfs-turbo": {
		Name:    "sfs-turbo",
		Version: "v1",
		Product: "SFSTurbo",
	},
	"cbh": {
		Name:    "cbh",
		Version: "v1",
		Product: "CBH",
	},
	"cbr": {
		Name:    "cbr",
		Version: "v3",
		Product: "CBR",
	},
	"csbs": {
		Name:    "csbs",
		Version: "v1",
		Product: "CSBS",
	},
	"vbs": {
		Name:    "vbs",
		Version: "v2",
		Product: "VBS",
	},
	"sdrs": {
		Name:    "sdrs",
		Version: "v1",
		Product: "SDRS",
	},

	// ******* catalog for network ******
	"vpc": {
		Name:             "vpc",
		Version:          "v1",
		WithOutProjectID: true,
		Product:          "VPC",
	},
	"networkv2": {
		Name:             "vpc",
		Version:          "v2.0",
		WithOutProjectID: true,
		Product:          "VPC",
	},
	"vpcv3": {
		Name:    "vpc",
		Version: "v3",
		Product: "VPC",
	},
	"nat": {
		Name:             "nat",
		Version:          "v2.0",
		WithOutProjectID: true,
		Product:          "NAT",
	},
	"elbv2": {
		Name:             "vpc",
		Version:          "v2.0",
		WithOutProjectID: true,
		Product:          "ELB",
	},
	"elbv3": {
		Name:    "vpc",
		Version: "v3",
		Product: "ELB",
	},
	"elb": {
		Name:    "vpc",
		Version: "v2",
		Product: "ELB",
	},
	"fwv2": {
		Name:             "vpc",
		Version:          "v2.0",
		WithOutProjectID: true,
		Product:          "VPC",
	},
	"vpcep": {
		Name:    "vpcep",
		Version: "v1",
		Product: "VPCEP",
	},
	"dns": {
		Name:             "clouddns",
		Version:          "v2",
		WithOutProjectID: true,
		Product:          "DNS",
	},
	"dns_region": {
		Name:             "clouddns",
		Version:          "v2",
		WithOutProjectID: true,
		Product:          "DNS",
	},
	"workspace": {
		Name:    "workspace",
		Version: "v2",
		Product: "Workspace",
	},
	"er": {
		Name:    "er",
		Version: "v3",
		Product: "ER",
	},
	"vpn": {
		Name:    "vpn",
		Version: "v5",
		Product: "VPN",
	},
	"ga": {
		Name:             "ga",
		Version:          "v1",
		WithOutProjectID: true,
		Scope:            "global",
		Product:          "GA",
	},
	"dc": {
		Name:    "dcaas",
		Version: "v3",
		Product: "DC",
	},
	"cfw": {
		Name:    "cfw",
		Version: "v1",
		Product: "CFW",
	},

	// catalog for database
	"rdsv1": {
		Name:    "rds",
		Version: "rds/v1",
		Product: "RDS",
	},
	"rds": {
		Name:    "rds",
		Version: "v3",
		Product: "RDS",
	},
	"ram": {
		Name:             "ram",
		Version:          "v1",
		WithOutProjectID: true,
		Scope:            "global",
		Product:          "RAM",
	},
	"dds": {
		Name:    "dds",
		Version: "v3",
		Product: "DDS",
	},
	"geminidb": {
		Name:    "gaussdb-nosql",
		Version: "v3",
		Product: "GaussDBforNoSQL",
	},
	"geminidbv31": {
		Name:    "gaussdb-nosql",
		Version: "v3.1",
		Product: "GaussDBforNoSQL",
	},
	"gaussdb": {
		Name:    "gaussdb",
		Version: "v3",
		Product: "GaussDB",
	},
	"opengauss": {
		Name:    "gaussdb-opengauss",
		Version: "v3",
		Product: "GaussDBforopenGauss",
	},
	"opengaussv31": {
		Name:    "gaussdb",
		Version: "v3.1",
		Product: "GaussDBforopenGauss",
	},
	"drs": {
		Name:    "drs",
		Version: "v3",
		Product: "DRS",
	},

	// catalog for management service
	"ces": {
		Name:    "ces",
		Version: "V1.0",
		Product: "CES",
	},
	"cesv2": {
		Name:    "ces",
		Version: "v2",
		Product: "CES",
	},
	"cts": {
		Name:    "cts",
		Version: "v1.0",
		Product: "CTS",
	},
	"lts": {
		Name:    "lts",
		Version: "v2",
		Product: "LTS",
	},
	"apm": {
		Name:             "apm2",
		Version:          "v1",
		Product:          "APM",
		WithOutProjectID: true,
	},
	"smn": {
		Name:         "smn",
		Version:      "v2",
		ResourceBase: "notifications",
		Product:      "SMN",
	},
	"smn-tag": {
		Name:    "smn",
		Version: "v2",
		Product: "SMN",
	},
	"sms": {
		Name:             "sms.ap-southeast-1",
		Version:          "v3",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "SMS",
	},
	"tms": {
		Name:             "tms",
		Version:          "v1.0",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "TMS",
	},
	"rms": {
		Name:             "rms",
		Scope:            "global",
		Version:          "v1",
		WithOutProjectID: true,
		Product:          "RMS",
	},
	"organizations": {
		Name:             "organizations",
		Version:          "v1",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "Organizations",
	},
	// catalog for Meeting service, only used for API scan
	"meeting": {
		Name:             "api.meeting",
		Version:          "v1",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "Meeting",
	},

	// catalog for Security service
	"aad": {
		Name:             "aad",
		Version:          "v1",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "AAD",
	},
	"anti-ddos": {
		Name:    "antiddos",
		Version: "v1",
		Product: "Anti-DDoS",
	},
	"kms": {
		Name:             "kms",
		Version:          "v1.0",
		WithOutProjectID: true,
		Product:          "DEW",
	},
	"kmsv1": {
		Name:    "kms",
		Version: "v1",
		Product: "DEW",
	},
	"kmsv3": {
		Name:    "kms",
		Version: "v3",
		Product: "DEW",
	},
	"csms": {
		Name:    "csms",
		Version: "v1",
		Product: "DEW",
	},
	"waf": {
		Name:         "waf",
		Version:      "v1",
		ResourceBase: "waf",
		Product:      "WAF",
	},
	"waf-dedicated": {
		Name:         "waf",
		Version:      "v1",
		ResourceBase: "premium-waf",
		Product:      "WAF",
	},
	"dbss": {
		Name:    "dbss",
		Version: "v2",
		Product: "DBSS",
	},
	"hss": {
		Name:    "hss",
		Version: "v5",
		Product: "HSS",
	},
	"secmaster": {
		Name:    "secmaster",
		Version: "v1",
		Product: "SecMaster",
	},

	// catalog for Enterprise Intelligence
	"mrs": {
		Name:    "mrs",
		Version: "v1.1",
		Product: "MRS",
	},
	"mrsv2": {
		Name:    "mrs",
		Version: "v2",
		Product: "MRS",
	},
	"modelarts": {
		Name:    "modelarts",
		Version: "v1",
		Product: "ModelArts",
	},
	"modelartsv2": {
		Name:    "modelarts",
		Version: "v2",
		Product: "ModelArts",
	},
	"dataarts": {
		Name:    "dayu",
		Version: "v1",
		Product: "DataArtsStudio",
	},
	"dws": {
		Name:    "dws",
		Version: "v1.0",
		Product: "DWS",
	},
	"dwsv2": {
		Name:    "dws",
		Version: "v2",
		Product: "DWS",
	},
	"dli": {
		Name:    "dli",
		Version: "v1.0",
		Product: "DLI",
	},
	"dliv2": {
		Name:    "dli",
		Version: "v2.0",
		Product: "DLI",
	},
	"dis": {
		Name:    "dis",
		Version: "v2",
		Product: "DIS",
	},
	"disv3": {
		Name:    "dis",
		Version: "v3",
		Product: "DIS",
	},
	"css": {

		Name:    "css",
		Version: "v1.0",
		Product: "CSS",
	},
	"cs": {
		Name:    "cs",
		Version: "v1.0",
		Product: "CloudStream",
	},
	"ges": {
		Name:    "ges",
		Version: "v1.0",
		Product: "GES",
	},
	"cloudtable": {
		Name:    "cloudtable",
		Version: "v2",
		Product: "CloudTable",
	},
	"cdm": {
		Name:    "cdm",
		Version: "v1.1",
		Product: "CDM",
	},

	// catalog for Application
	"apig": {
		Name:             "apig",
		Version:          "v1.0",
		ResourceBase:     "apigw",
		WithOutProjectID: true,
		Product:          "APIG",
	},
	"apigv2": {
		Name:         "apig",
		Version:      "v2",
		ResourceBase: "apigw",
		Product:      "APIG",
	},
	"bcs": {
		Name:    "bcs",
		Version: "v2",
		Product: "BCS",
	},
	"cse": {
		Name:    "cse",
		Version: "v2",
		Product: "CSE",
	},
	"dcsv1": {
		Name:             "dcs",
		Version:          "v1.0",
		WithOutProjectID: true,
		Product:          "DCS",
	},
	"dcs": {
		Name:             "dcs",
		Version:          "v2",
		WithOutProjectID: true,
		Product:          "DCS",
	},
	"dms": {
		Name:             "dms",
		Version:          "v1.0",
		WithOutProjectID: true,
		Product:          "DMS",
	},
	"dmsv2": {
		Name:             "dms",
		Version:          "v2",
		WithOutProjectID: true,
		Product:          "DMS",
	},
	"servicestage": {
		Name:    "servicestage",
		Version: "v1",
		Product: "ServiceStage",
	},
	"servicestagev2": {
		Name:    "servicestage",
		Version: "v2",
		Product: "ServiceStage",
	},
	"roma": {
		Name:    "roma",
		Version: "v1",
		Product: "ROMA Connect",
	},
	"romav2": {
		Name:    "roma",
		Version: "v2",
		Product: "ROMA Connect",
	},

	// catalog for IEC which is a global service
	"iec": {
		Name:             "iecs",
		Version:          "v1",
		Scope:            "global",
		Admin:            true,
		WithOutProjectID: true,
		Product:          "IEC",
	},

	// catalog for Others
	"rts": {
		Name:    "rts",
		Version: "v1",
		Product: "RTS",
	},
	"oms": {
		Name:    "oms",
		Version: "v1",
		Product: "OMS",
	},
	"scm": {
		Name:             "scm",
		Version:          "v3",
		WithOutProjectID: true,
		Product:          "SCM",
	},

	// catalog for cc
	"cc": {
		Name:             "cc",
		Version:          "v3",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "CC",
	},

	"cpts": {
		Name:    "cpts",
		Version: "v1",
		Product: "CPTS",
	},

	"live": {
		Name:    "live",
		Version: "v1",
		Product: "Live",
	},

	"mpc": {
		Name:    "mpc",
		Version: "v1",
		Product: "MPC",
	},

	"iotda": {
		Name:    "iotda",
		Version: "v5",
		Product: "IoTDA",
	},

	"vod": {
		Name:    "vod",
		Version: "v1",
		Product: "VOD",
	},

	"cmdb": {
		Name:    "cmdb",
		Version: "v1",
		Scope:   "global",
		Product: "AOM",
	},

	"ddm": {
		Name:             "ddm",
		WithOutProjectID: true,
		Product:          "DDM",
	},

	// catalog for Developer Services
	"codehub": {
		Name:    "codehub-ext",
		Product: "CodeHub",
	},

	"projectman": {
		Name:    "projectman-ext",
		Version: "v4",
		Product: "ProjectMan",
	},

	// catalog for Data Security Center
	"dsc": {
		Name:    "sdg",
		Product: "DSC",
	},

	// catalog for Cloud Phone
	"cph": {
		Name:    "cph",
		Product: "CPH",
	},

	// catalog for Joint-Operation Cloud only
	// it should be at the end of this map, and no necessary to put the key into allServiceCatalog
	"mls": {
		Name:    "mls",
		Version: "v1.0",
		Product: "MLS",
	},
}

// GetServiceEndpoint try to get the endpoint from customizing map
func GetServiceEndpoint(c *HcsConfig, srv, region string) string {
	if endpoint, ok := c.Endpoints[srv]; ok {
		return endpoint
	}

	// get the endpoint from build-in catalog
	catalog, ok := allServiceCatalog[srv]
	if !ok {
		return ""
	}

	var ep string
	if catalog.Scope == "global" && !c.RegionClient {
		ep = fmt.Sprintf("https://%s.%s/", catalog.Name, c.Cloud)
	} else {
		ep = fmt.Sprintf("https://%s.%s.%s/", catalog.Name, region, c.Cloud)
	}
	return ep
}

// GetServiceCatalog returns the catalog object of a service
func GetServiceCatalog(service string) *ServiceCatalog {
	if catalog, ok := allServiceCatalog[service]; ok {
		return &catalog
	}
	return nil
}
