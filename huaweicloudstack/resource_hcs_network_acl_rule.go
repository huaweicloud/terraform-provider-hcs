/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2023-2023. All rights reserved.
 */

package huaweicloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v2/extensions/fwaas_v2/policies"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v2/extensions/fwaas_v2/rules"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/pagination"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/logp"
)

func ResourceNetworkACLRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkACLRuleCreate,
		Read:   resourceNetworkACLRuleRead,
		Update: resourceNetworkACLRuleUpdate,
		Delete: resourceNetworkACLRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"tcp", "udp", "icmp", "any",
				}, true),
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"allow", "deny",
				}, true),
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  4,
			},
			"source_ip_address": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use source_ip_address instead",
			},
			"destination_ip_address": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use destination_ip_address instead",
			},
			"source_ip_addresses": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"destination_ip_addresses": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source_port": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use source_ports instead",
			},
			"destination_port": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use destination_ports instead",
			},
			"source_ports": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"destination_ports": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceNetworkACLRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := config.GetHcsConfig(meta)
	fwClient, err := config.FwV2Client(common.GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloudStack fw client: %s", err)
	}

	enabled := d.Get("enabled").(bool)
	ipVersion := normalizeNetworkACLRuleIPVersion(d.Get("ip_version").(int))
	protocol := normalizeNetworkACLRuleProtocol(d.Get("protocol").(string))

	ruleConfiguration := rules.CreateOpts{
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		Action:               d.Get("action").(string),
		IPVersion:            ipVersion,
		Protocol:             protocol,
		SourceIPAddress:      d.Get("source_ip_address").(string),
		DestinationIPAddress: d.Get("destination_ip_address").(string),
		SourcePort:           d.Get("source_port").(string),
		DestinationPort:      d.Get("destination_port").(string),
		Enabled:              &enabled,
	}
	sourceIPAddresses := d.Get("source_ip_addresses").(*schema.Set).List()
	if len(sourceIPAddresses) > 0 {
		ruleConfiguration.SourceIPAddresses = make([]string, len(sourceIPAddresses))
		for i, r := range sourceIPAddresses {
			ruleConfiguration.SourceIPAddresses[i] = r.(string)
		}
	}
	destinationIPAddresses := d.Get("destination_ip_addresses").(*schema.Set).List()
	if len(destinationIPAddresses) > 0 {
		ruleConfiguration.DestinationIPAddresses = make([]string, len(destinationIPAddresses))
		for i, r := range destinationIPAddresses {
			ruleConfiguration.DestinationIPAddresses[i] = r.(string)
		}
	}
	sourcePorts := d.Get("source_ports").(*schema.Set).List()
	if len(sourcePorts) > 0 {
		ruleConfiguration.SourcePorts = make([]string, len(sourcePorts))
		for i, r := range sourcePorts {
			ruleConfiguration.SourcePorts[i] = r.(string)
		}
	}
	destinationPorts := d.Get("destination_ports").(*schema.Set).List()
	if len(destinationPorts) > 0 {
		ruleConfiguration.DestinationPorts = make([]string, len(destinationPorts))
		for i, r := range destinationPorts {
			ruleConfiguration.DestinationPorts[i] = r.(string)
		}
	}

	logp.Printf("[DEBUG] Create Network ACL rule: %#v", ruleConfiguration)
	rule, err := rules.Create(fwClient, ruleConfiguration).Extract()
	if err != nil {
		return err
	}

	logp.Printf("[DEBUG] Network ACL rule with id %s", rule.ID)
	d.SetId(rule.ID)

	return resourceNetworkACLRuleRead(d, meta)
}

func resourceNetworkACLRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := config.GetHcsConfig(meta)
	fwClient, err := config.FwV2Client(common.GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloudStack fw client: %s", err)
	}

	rule, err := rules.Get(fwClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "Network ACL rule")
	}

	logp.Printf("[DEBUG] Retrieve HuaweiCloudStack Network ACL rule %s: %#v", d.Id(), rule)

	d.Set("action", rule.Action)
	d.Set("name", rule.Name)
	d.Set("description", rule.Description)
	d.Set("ip_version", rule.IPVersion)
	d.Set("source_ip_address", rule.SourceIPAddress)
	d.Set("source_ip_addresses", rule.SourceIPAddresses)
	d.Set("destination_ip_address", rule.DestinationIPAddress)
	d.Set("destination_ip_addresses", rule.DestinationIPAddresses)
	d.Set("source_port", rule.SourcePort)
	d.Set("destination_port", rule.DestinationPort)
	d.Set("source_ports", rule.SourcePorts)
	d.Set("destination_ports", rule.DestinationPorts)
	d.Set("enabled", rule.Enabled)

	if rule.Protocol == "" {
		d.Set("protocol", "any")
	} else {
		d.Set("protocol", rule.Protocol)
	}

	return nil
}

func resourceNetworkACLRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := config.GetHcsConfig(meta)
	fwClient, err := config.FwV2Client(common.GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloudStack fw client: %s", err)
	}

	var updateOpts rules.UpdateOpts
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("protocol") {
		protocol := d.Get("protocol").(string)
		if protocol == "any" {
			updateOpts.Protocol = nil
		} else {
			updateOpts.Protocol = &protocol
		}
	}
	if d.HasChange("action") {
		action := d.Get("action").(string)
		updateOpts.Action = &action
	}
	if d.HasChange("ip_version") {
		ipVersion := normalizeNetworkACLRuleIPVersion(d.Get("ip_version").(int))
		updateOpts.IPVersion = &ipVersion
	}
	if d.HasChange("source_ip_address") {
		sourceIPAddress := d.Get("source_ip_address").(string)
		updateOpts.SourceIPAddress = &sourceIPAddress
	}
	if d.HasChange("source_ip_addresses") {

		sourceIPAddresses := d.Get("source_ip_addresses").(*schema.Set).List()
		addresses := make([]string, len(sourceIPAddresses))
		if len(sourceIPAddresses) > 0 {
			for i, r := range sourceIPAddresses {
				addresses[i] = r.(string)
			}
		}
		updateOpts.SourceIPAddresses = &addresses
	}
	if d.HasChange("source_port") {
		sourcePort := d.Get("source_port").(string)
		updateOpts.SourcePort = &sourcePort
	}
	if d.HasChange("source_ports") {
		sourcePorts := d.Get("source_ports").(*schema.Set).List()
		ports := make([]string, len(sourcePorts))
		if len(sourcePorts) > 0 {
			for i, r := range sourcePorts {
				ports[i] = r.(string)
			}
		}
		updateOpts.SourcePorts = &ports
	}
	if d.HasChange("destination_ip_address") {
		destinationIPAddress := d.Get("destination_ip_address").(string)
		updateOpts.DestinationIPAddress = &destinationIPAddress
	}
	if d.HasChange("destination_ip_addresses") {

		destinationIPAddresses := d.Get("destination_ip_addresses").(*schema.Set).List()
		addresses := make([]string, len(destinationIPAddresses))
		if len(destinationIPAddresses) > 0 {
			for i, r := range destinationIPAddresses {
				addresses[i] = r.(string)
			}
		}
		updateOpts.DestinationIPAddresses = &addresses
	}
	if d.HasChange("destination_port") {
		destinationPort := d.Get("destination_port").(string)
		updateOpts.DestinationPort = &destinationPort
	}
	if d.HasChange("destination_ports") {
		destinationPorts := d.Get("destination_ports").(*schema.Set).List()
		ports := make([]string, len(destinationPorts))
		if len(destinationPorts) > 0 {
			for i, r := range destinationPorts {
				ports[i] = r.(string)
			}
		}
		updateOpts.DestinationPorts = &ports
	}
	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		updateOpts.Enabled = &enabled
	}

	logp.Printf("[DEBUG] Updating Network ACL rule %s: %#v", d.Id(), updateOpts)
	err = rules.Update(fwClient, d.Id(), updateOpts).Err
	if err != nil {
		return err
	}

	return resourceNetworkACLRuleRead(d, meta)
}

func resourceNetworkACLRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := config.GetHcsConfig(meta)
	fwClient, err := config.FwV2Client(common.GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloudStack fw client: %s", err)
	}

	rule, err := rules.Get(fwClient, d.Id()).Extract()
	if err != nil {
		return err
	}

	policyID, err := assignedPolicyID(fwClient, rule.ID)
	if err != nil {
		return err
	}
	if policyID != "" {
		_, err := policies.RemoveRule(fwClient, policyID, rule.ID).Extract()
		if err != nil {
			return err
		}
	}

	logp.Printf("[DEBUG] Destroy Network ACL rule: %s", d.Id())
	return rules.Delete(fwClient, d.Id()).Err
}

func normalizeNetworkACLRuleIPVersion(ipv int) golangsdk.IPVersion {
	// Determine the IP Version
	var ipVersion golangsdk.IPVersion
	switch ipv {
	case 4:
		ipVersion = golangsdk.IPv4
	case 6:
		ipVersion = golangsdk.IPv6
	}

	return ipVersion
}

func normalizeNetworkACLRuleProtocol(p string) rules.Protocol {
	var protocol rules.Protocol
	switch p {
	case "any":
		protocol = rules.ProtocolAny
	case "icmp":
		protocol = rules.ProtocolICMP
	case "tcp":
		protocol = rules.ProtocolTCP
	case "udp":
		protocol = rules.ProtocolUDP
	}

	return protocol
}

func assignedPolicyID(fwClient *golangsdk.ServiceClient, ruleID string) (string, error) {
	pager := policies.List(fwClient, policies.ListOpts{})
	policyID := ""
	err := pager.EachPage(func(page pagination.Page) (b bool, err error) {
		policyList, err := policies.ExtractPolicies(page)
		if err != nil {
			return false, err
		}
		for _, policy := range policyList {
			for _, rule := range policy.Rules {
				if rule == ruleID {
					policyID = policy.ID
					return false, nil
				}
			}
		}
		return true, nil
	})
	if err != nil {
		return "", err
	}
	return policyID, nil
}
