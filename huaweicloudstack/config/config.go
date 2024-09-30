package config

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	huaweiConfig "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/helper/mutexkv"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/identity/v3/domains"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/identity/v3/projects"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/identity/v3/users"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/obs"
)

const (
	providerUserAgent string = "terraform-provider-iac"
)

// MutexKV is a global lock on all resources, it can lock the specified shared string (such as resource ID, resource
// Name, port, etc.) to prevent other resources from using it, for concurrency control.
// Usage: MutexKV.Lock({resource ID}) and MutexKV.Unlock({resource ID})
var MutexKV = mutexkv.NewMutexKV()

type Config = huaweiConfig.Config

type HcsConfig struct {
	Config
	HcsHwClient     *golangsdk.ProviderClient
	HcsDomainClient *golangsdk.ProviderClient
}

func GetHcsConfig(meta interface{}) *HcsConfig {
	conf, ok := meta.(*Config)
	if ok {
		if hcsConf, ok := conf.Metadata.(*HcsConfig); ok {
			return hcsConf
		}
	}

	if hcsConf, ok := meta.(*HcsConfig); ok {
		return hcsConf
	}

	return nil
}

func (c *HcsConfig) LoadAndValidate() error {
	if c.MaxRetries < 0 {
		return fmt.Errorf("max_retries should be a positive value")
	}

	err := buildClient(c)
	if err != nil {
		return err
	}

	if c.Region == "" {
		return fmt.Errorf("region should be provided")
	}

	if c.TenantName == "" && c.TenantID == "" {
		return fmt.Errorf("project_name or project_id should be provided")
	}

	if c.Cloud == "" && len(c.Endpoints) <= 0 {
		return fmt.Errorf("cloud or endpoints should be provided")
	}

	// Assume role
	if c.AssumeRoleAgency != "" {
		err = buildClientByAgency(c)
		if err != nil {
			return err
		}
	}

	if c.HwClient != nil && c.HwClient.ProjectID != "" {
		c.RegionProjectIDMap[c.Region] = c.HwClient.ProjectID
	}
	log.Printf("[DEBUG] init region and project map: %#v", c.RegionProjectIDMap)

	// set DomainID for IAM resource
	if c.DomainID == "" {
		if domainID, err := c.getDomainID(); err == nil {
			c.DomainID = domainID

			// update DomainClient.AKSKAuthOptions
			if c.DomainClient.AKSKAuthOptions.AccessKey != "" {
				c.DomainClient.AKSKAuthOptions.DomainID = c.DomainID
			}
		} else {
			log.Printf("[WARN] get domain id failed: %s", err)
		}
	}

	if c.UserID == "" && c.Username != "" {
		if userID, err := c.getUserIDbyName(c.Username); err == nil {
			c.UserID = userID
		} else {
			log.Printf("[WARN] get user id failed: %s", err)
		}
	}

	return nil
}

func retryBackoffFunc(ctx context.Context, respErr *golangsdk.ErrUnexpectedResponseCode, e error, retries uint) error {
	minutes := int(math.Pow(2, float64(retries)))
	if minutes > 30 { // won't wait more than 30 minutes
		minutes = 30
	}

	log.Printf("[WARN] Received StatusTooManyRequests response code, try to sleep %d minutes", minutes)
	sleep := time.Duration(minutes) * time.Minute

	if ctx != nil {
		select {
		case <-time.After(sleep):
		case <-ctx.Done():
			return e
		}
	} else {
		//lintignore:R018
		time.Sleep(sleep)
	}

	return nil
}

func getObsEndpoint(c *HcsConfig, region string) string {
	if endpoint, ok := c.Endpoints["obs"]; ok {
		// replace the region in customizing OBS endpoint
		subparts := strings.Split(endpoint, ".")
		if len(subparts) >= 3 && subparts[1] != region {
			subparts[1] = region
			return strings.Join(subparts, ".")
		}
		return endpoint
	}
	return fmt.Sprintf("https://obs.%s.%s/", region, c.Cloud)
}

func (c *HcsConfig) ObjectStorageClientWithSignature(region string) (*obs.ObsClient, error) {
	if c.AccessKey == "" || c.SecretKey == "" {
		return nil, fmt.Errorf("missing credentials for OBS, need access_key and secret_key values for provider")
	}

	clientConfigure := obs.WithHttpClient(&c.DomainClient.HTTPClient)
	userAgentConfigure := obs.WithUserAgent(buildObsUserAgent())
	obsEndpoint := getObsEndpoint(c, region)
	if c.SecurityToken != "" {
		return obs.New(c.AccessKey, c.SecretKey, obsEndpoint,
			obs.WithSignature("OBS"), obs.WithSecurityToken(c.SecurityToken), clientConfigure,
			userAgentConfigure)
	}
	return obs.New(c.AccessKey, c.SecretKey, obsEndpoint, obs.WithSignature("OBS"), clientConfigure, userAgentConfigure)
}

func (c *HcsConfig) ObjectStorageClient(region string) (*obs.ObsClient, error) {
	if c.AccessKey == "" || c.SecretKey == "" {
		return nil, fmt.Errorf("missing credentials for OBS, need access_key and secret_key values for provider")
	}

	if !c.SecurityKeyExpiresAt.IsZero() {
		c.SecurityKeyLock.Lock()
		defer c.SecurityKeyLock.Unlock()
		timeNow := time.Now().Unix()
		expairesAtInt := c.SecurityKeyExpiresAt.Unix()
		if timeNow+keyExpiresDuration > expairesAtInt {
			err := c.reloadSecurityKey()
			if err != nil {
				return nil, err
			}
		}
	}

	clientConfigure := obs.WithHttpClient(&c.DomainClient.HTTPClient)
	userAgentConfigure := obs.WithUserAgent(buildObsUserAgent())
	obsEndpoint := getObsEndpoint(c, region)
	if c.SecurityToken != "" {
		return obs.New(c.AccessKey, c.SecretKey, obsEndpoint, obs.WithSecurityToken(c.SecurityToken), clientConfigure,
			userAgentConfigure)
	}
	return obs.New(c.AccessKey, c.SecretKey, obsEndpoint, clientConfigure, userAgentConfigure)
}

func buildObsUserAgent() string {
	var agent string = providerUserAgent
	if customUserAgent := os.Getenv("HCS_TF_CUSTOM_UA"); customUserAgent != "" {
		agent = fmt.Sprintf("%s %s", customUserAgent, providerUserAgent)
	}

	return agent
}

// NewServiceClient create a ServiceClient which was assembled from ServiceCatalog.
// If you want to add new ServiceClient, please make sure the catalog was already in allServiceCatalog.
// the endpoint likes https://{Name}.{Region}.myhuaweicloud.com/{Version}/{project_id}/{ResourceBase}
func (c *HcsConfig) NewServiceClient(srv, region string) (*golangsdk.ServiceClient, error) {
	serviceCatalog, ok := allServiceCatalog[srv]
	if !ok {
		return nil, fmt.Errorf("service type %s is invalid or not supportted", srv)
	}

	if !c.SecurityKeyExpiresAt.IsZero() {
		c.SecurityKeyLock.Lock()
		defer c.SecurityKeyLock.Unlock()
		timeNow := time.Now().Unix()
		expairesAtInt := c.SecurityKeyExpiresAt.Unix()
		if timeNow+keyExpiresDuration > expairesAtInt {
			err := c.reloadSecurityKey()
			if err != nil {
				return nil, err
			}
		}
	}

	client := c.HcsHwClient
	if serviceCatalog.Admin {
		client = c.HcsDomainClient
	}

	if endpoint, ok := c.Endpoints[srv]; ok {
		return c.newServiceClientByEndpoint(client, srv, endpoint)
	}
	return c.newServiceClientByName(client, serviceCatalog, region)
}

func (c *HcsConfig) newServiceClientByName(client *golangsdk.ProviderClient, catalog ServiceCatalog, region string) (*golangsdk.ServiceClient, error) {
	if catalog.Name == "" {
		return nil, fmt.Errorf("must specify the service name")
	}

	// Custom Resource-level region only supports AK/SK authentication.
	// If set it when using non AK/SK authentication, then it must be the same as Provider-level region.
	if region != c.Region && (c.AccessKey == "" || c.SecretKey == "") {
		return nil, fmt.Errorf("Resource-level region must be the same as Provider-level region when using non AK/SK authentication if Resource-level region set")
	}

	c.RPLock.Lock()
	defer c.RPLock.Unlock()
	projectID, ok := c.RegionProjectIDMap[region]
	if !ok {
		// Not find in the map, then try to query and store.
		err := c.loadUserProjects(client, region)
		if err != nil {
			return nil, err
		}
		projectID = c.RegionProjectIDMap[region]
	}

	// update ProjectID and region in ProviderClient
	clone := new(golangsdk.ProviderClient)
	*clone = *client
	clone.ProjectID = projectID
	clone.AKSKAuthOptions.ProjectId = projectID
	clone.AKSKAuthOptions.Region = region

	sc := &golangsdk.ServiceClient{
		ProviderClient: clone,
	}

	if catalog.Scope == "global" && !c.RegionClient {
		sc.Endpoint = fmt.Sprintf("https://%s.%s/", catalog.Name, c.Cloud)
	} else {
		sc.Endpoint = fmt.Sprintf("https://%s.%s.%s/", catalog.Name, region, c.Cloud)
	}

	sc.ResourceBase = sc.Endpoint
	if catalog.Version != "" {
		sc.ResourceBase = sc.ResourceBase + catalog.Version + "/"
	}
	if !catalog.WithOutProjectID {
		sc.ResourceBase = sc.ResourceBase + projectID + "/"
	}
	if catalog.ResourceBase != "" {
		sc.ResourceBase = sc.ResourceBase + catalog.ResourceBase + "/"
	}

	return sc, nil
}

// newServiceClientByEndpoint returns a ServiceClient which the endpoint was initialized by customer
// the format of customer endpoint likes https://{Name}.{Region}.xxxx.com
func (c *HcsConfig) newServiceClientByEndpoint(client *golangsdk.ProviderClient, srv, endpoint string) (*golangsdk.ServiceClient, error) {
	catalog, ok := allServiceCatalog[srv]
	if !ok {
		return nil, fmt.Errorf("service type %s is invalid or not supportted", srv)
	}

	sc := &golangsdk.ServiceClient{
		ProviderClient: client,
		Endpoint:       endpoint,
	}

	sc.ResourceBase = sc.Endpoint
	if catalog.Version != "" {
		sc.ResourceBase = sc.ResourceBase + catalog.Version + "/"
	}
	if !catalog.WithOutProjectID {
		sc.ResourceBase = sc.ResourceBase + client.ProjectID + "/"
	}
	if catalog.ResourceBase != "" {
		sc.ResourceBase = sc.ResourceBase + catalog.ResourceBase + "/"
	}
	return sc, nil
}

func (c *HcsConfig) getDomainID() (string, error) {
	identityClient, err := c.IdentityV3Client(c.Region)
	if err != nil {
		return "", fmt.Errorf("Error creating IAM client: %s", err)
	}

	// the List request does not support query options
	allPages, err := domains.List(identityClient, nil).AllPages()
	if err != nil {
		return "", fmt.Errorf("List domains failed, err=%s", err)
	}

	all, err := domains.ExtractDomains(allPages)
	if err != nil {
		return "", fmt.Errorf("Extract domains failed, err=%s", err)
	}

	if len(all) == 0 {
		return "", fmt.Errorf("domain was not found")
	}

	if c.DomainName != "" && c.DomainName != all[0].Name {
		return "", fmt.Errorf("domain %s was not found, got %s", c.DomainName, all[0].Name)
	}

	return all[0].ID, nil
}

func (c *HcsConfig) getUserIDbyName(name string) (string, error) {
	identityClient, err := c.IdentityV3Client(c.Region)
	if err != nil {
		return "", fmt.Errorf("Error creating IAM client: %s", err)
	}

	opts := users.ListOpts{
		Name: name,
	}
	allPages, err := users.List(identityClient, opts).AllPages()
	if err != nil {
		return "", fmt.Errorf("query IAM user %s failed, err=%s", name, err)
	}

	all, err := users.ExtractUsers(allPages)
	if err != nil {
		return "", fmt.Errorf("Extract users failed, err=%s", err)
	}

	if len(all) == 0 {
		return "", fmt.Errorf("IAM user %s was not found", name)
	}

	if name != "" && name != all[0].Name {
		return "", fmt.Errorf("IAM user %s was not found, got %s", name, all[0].Name)
	}

	return all[0].ID, nil
}

// loadUserProjects will query the region-projectId pair and store it into RegionProjectIDMap
func (c *HcsConfig) loadUserProjects(client *golangsdk.ProviderClient, region string) error {

	log.Printf("[DEBUG] Load project ID for region: %s", region)
	domainID := client.DomainID
	opts := projects.ListOpts{
		DomainID: domainID,
		Name:     region,
	}
	sc := new(golangsdk.ServiceClient)
	sc.Endpoint = c.IdentityEndpoint + "/"
	sc.ProviderClient = client
	allPages, err := projects.List(sc, &opts).AllPages()
	if err != nil {
		return fmt.Errorf("List projects failed, err=%s", err)
	}

	all, err := projects.ExtractProjects(allPages)
	if err != nil {
		return fmt.Errorf("Extract projects failed, err=%s", err)
	}

	if len(all) == 0 {
		return fmt.Errorf("Wrong name or no access to the region: %s", region)
	}

	for _, item := range all {
		log.Printf("[DEBUG] add %s/%s to region and project map", item.Name, item.ID)
		c.RegionProjectIDMap[item.Name] = item.ID
	}
	return nil
}

// GetProjectID is used to get the project ID for services
func (c *HcsConfig) GetProjectID(region string) string {
	c.RPLock.Lock()
	defer c.RPLock.Unlock()

	projectID, ok := c.RegionProjectIDMap[region]
	if !ok {
		// Not find in the map, then try to query and store.
		if err := c.loadUserProjects(c.HcsDomainClient, region); err != nil {
			log.Printf("[WARN] can not find the project ID of %s: %s", region, err)
			return ""
		}
		projectID = c.RegionProjectIDMap[region]
	}

	return projectID
}

// GetRegion returns the region that was specified in the resource. If a
// region was not set, the provider-level region is checked. The provider-level
// region can either be set by the region argument or by HCS_REGION_NAME.
func (c *HcsConfig) GetRegion(d *schema.ResourceData) string {
	if v, ok := d.GetOk("region"); ok {
		return v.(string)
	}

	return c.Region
}

// GetEnterpriseProjectID returns the enterprise_project_id that was specified in the resource.
// If it was not set, the provider-level value is checked. The provider-level value can
// either be set by the `enterprise_project_id` argument or by HCS_ENTERPRISE_PROJECT_ID.
func (c *HcsConfig) GetEnterpriseProjectID(d *schema.ResourceData) string {
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		return v.(string)
	}

	return c.EnterpriseProjectID
}

// DataGetEnterpriseProjectID returns the enterprise_project_id that was specified in the data source.
// If it was not set, the provider-level value is checked. The provider-level value can
// either be set by the `enterprise_project_id` argument or by HCS_ENTERPRISE_PROJECT_ID.
// If the provider-level value is also not set, `all_granted_eps` will be returned.
func (c *HcsConfig) DataGetEnterpriseProjectID(d *schema.ResourceData) string {
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		return v.(string)
	}
	if c.EnterpriseProjectID != "" {
		return c.EnterpriseProjectID
	}
	return "all_granted_eps"
}

// ********** client for Global Service **********
func (c *HcsConfig) IAMV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("iam", region)
}

func (c *HcsConfig) IdentityV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("identity", region)
}

func (c *HcsConfig) IAMNoVersionClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("iam_no_version", region)
}

func (c *HcsConfig) CdnV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cdn", region)
}

func (c *HcsConfig) EnterpriseProjectClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("eps", region)
}

// ********** client for Compute **********
func (c *HcsConfig) ComputeV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("ecs", region)
}

func (c *HcsConfig) ComputeV11Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("ecsv11", region)
}

func (c *HcsConfig) ComputeV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("ecsv21", region)
}

func (c *HcsConfig) AutoscalingV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("autoscaling", region)
}

func (c *HcsConfig) ImageV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("imsv1", region)
}

func (c *HcsConfig) ImageV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("ims", region)
}

func (c *HcsConfig) CceV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("ccev1", region)
}

func (c *HcsConfig) CceV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cce", region)
}

func (c *HcsConfig) CceAddonV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cce_addon", region)
}

func (c *HcsConfig) AomV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("aom", region)
}

func (c *HcsConfig) CciV1BetaClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cciv1_bata", region)
}

func (c *HcsConfig) CciV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cci", region)
}

func (c *HcsConfig) FgsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("fgs", region)
}

func (c *HcsConfig) SwrV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("swr", region)
}

func (c *HcsConfig) BmsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("bms", region)
}

func (c *HcsConfig) AosV1Client(region string) (*golangsdk.ServiceClient, error) {
	client, err := c.NewServiceClient("aos", region)
	if err != nil {
		return nil, err
	}
	u, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	client.MoreHeaders = map[string]string{
		"Content-Type":      "application/json",
		"X-Language":        "en-us",
		"Client-Request-Id": u,
	}
	return client, nil
}

// ********** client for Storage **********
func (c *HcsConfig) BlockStorageV21Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("evsv21", region)
}

func (c *HcsConfig) BlockStorageV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("evs", region)
}

func (c *HcsConfig) SfsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("sfs", region)
}

func (c *HcsConfig) SfsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("sfs-turbo", region)
}

func (c *HcsConfig) CbrV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cbr", region)
}

func (c *HcsConfig) CsbsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("csbs", region)
}

func (c *HcsConfig) VbsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("vbs", region)
}

func (c *HcsConfig) SdrsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("sdrs", region)
}

// ********** client for Network **********
func (c *HcsConfig) NetworkingV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("vpc", region)
}

// NetworkingV2Client returns a ServiceClient for neutron APIs
// the endpoint likes: https://vpc.{region}.myhuaweicloud.com/v2.0/
func (c *HcsConfig) NetworkingV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("networkv2", region)
}

func (c *HcsConfig) NetworkingV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("vpcv3", region)
}

// VPCEPClient returns a ServiceClient for VPC Endpoint APIs
// the endpoint likes: https://vpcep.{region}.myhuaweicloud.com/v1/{project_id}/
func (c *HcsConfig) VPCEPClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("vpcep", region)
}

func (c *HcsConfig) NatGatewayClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("nat", region)
}

// ElbV2Client is the client for elb v2.0 (openstack) api
func (c *HcsConfig) ElbV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("elbv2", region)
}

// ElbV3Client is the client for elb v3 api
func (c *HcsConfig) ElbV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("elbv3", region)
}

// LoadBalancerClient is the client for elb v2 api
func (c *HcsConfig) LoadBalancerClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("elb", region)
}

func (c *HcsConfig) FwV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("networkv2", region)
}

func (c *HcsConfig) DnsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dns", region)
}

func (c *HcsConfig) DnsWithRegionClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dns_region", region)
}

func (c *HcsConfig) ErV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("er", region)
}

func (c *HcsConfig) DcV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dc", region)
}

// ********** client for Management **********
func (c *HcsConfig) CtsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cts", region)
}

func (c *HcsConfig) CesV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("ces", region)
}

func (c *HcsConfig) CesV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cesv2", region)
}

func (c *HcsConfig) LtsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("lts", region)
}

func (c *HcsConfig) SmnV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("smn", region)
}

func (c *HcsConfig) SmnV2TagClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("smn-tag", region)
}

func (c *HcsConfig) RmsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("rms", region)
}

// ********** client for Security **********
func (c *HcsConfig) AntiDDosV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("anti-ddos", region)
}

func (c *HcsConfig) AadV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("aad", region)
}

func (c *HcsConfig) KmsKeyV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("kms", region)
}

func (c *HcsConfig) KmsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("kmsv1", region)
}

func (c *HcsConfig) KmsV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("kmsv3", region)
}

// WafV1Client is not avaliable in HuaweiCloudStack, will be imported by other clouds
func (c *HcsConfig) WafV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("waf", region)
}

func (c *HcsConfig) WafDedicatedV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("waf-dedicated", region)
}

// ********** client for Enterprise Intelligence **********
func (c *HcsConfig) MrsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("mrs", region)
}

func (c *HcsConfig) MrsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("mrsv2", region)
}

func (c *HcsConfig) DwsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dws", region)
}

func (c *HcsConfig) DwsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dwsv2", region)
}

func (c *HcsConfig) DliV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dli", region)
}

func (c *HcsConfig) DliV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dliv2", region)
}

func (c *HcsConfig) DisV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dis", region)
}

func (c *HcsConfig) DisV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("disv3", region)
}

func (c *HcsConfig) CssV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("css", region)
}

func (c *HcsConfig) CloudStreamV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cs", region)
}

func (c *HcsConfig) CloudtableV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cloudtable", region)
}

func (c *HcsConfig) CdmV11Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cdm", region)
}

func (c *HcsConfig) GesV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("ges", region)
}

func (c *HcsConfig) ModelArtsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("modelarts", region)
}

func (c *HcsConfig) ModelArtsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("modelartsv2", region)
}

func (c *HcsConfig) DataArtsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dataarts", region)
}

func (c *HcsConfig) WorkspaceV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("workspace", region)
}

// ********** client for Application **********
func (c *HcsConfig) ApiGatewayV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("apig", region)
}

func (c *HcsConfig) ApigV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("apigv2", region)
}

func (c *HcsConfig) BcsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("bcs", region)
}

func (c *HcsConfig) CseV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cse", region)
}

func (c *HcsConfig) DcsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dcsv1", region)
}

func (c *HcsConfig) DcsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dcs", region)
}

func (c *HcsConfig) DmsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dms", region)
}

func (c *HcsConfig) DmsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dmsv2", region)
}

func (c *HcsConfig) ServiceStageV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("servicestage", region)
}

func (c *HcsConfig) ServiceStageV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("servicestagev2", region)
}

// ********** client for Database **********
func (c *HcsConfig) RdsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("rdsv1", region)
}

func (c *HcsConfig) RdsV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("rds", region)
}

func (c *HcsConfig) DdsV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dds", region)
}

func (c *HcsConfig) GeminiDBV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("geminidb", region)
}

func (c *HcsConfig) GeminiDBV31Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("geminidbv31", region)
}

func (c *HcsConfig) OpenGaussV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("opengauss", region)
}

func (c *HcsConfig) GaussdbV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("gaussdb", region)
}

func (c *HcsConfig) DrsV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("drs", region)
}

// ********** client for edge / IoT **********

// IECV1Client returns a ServiceClient for IEC Endpoint APIs
func (c *HcsConfig) IECV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("iec", region)
}

// ********** client for Others **********
func (c *HcsConfig) BssV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("bss", region)
}

func (c *HcsConfig) BssV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("bssv2", region)
}

func (c *HcsConfig) MaasV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("oms", region)
}

func (c *HcsConfig) SmsV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("sms", region)
}

func (c *HcsConfig) ScmV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("scm", region)
}

func (c *HcsConfig) RomaConnectV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("roma", region)
}

func (c *HcsConfig) RomaConnectV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("romav2", region)
}

// the following clients are used for Joint-Operation Cloud only

// MlsV1Client has the endpoint: https://mls.{{region}}/{{cloud}}/v1.0/{{project_id}}
func (c *HcsConfig) MlsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("mls", region)
}

// NatV2Client has the endpoint: https://nat.{{region}}/{{cloud}}/v2.0/
func (c *HcsConfig) NatV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("natv2", region)
}
