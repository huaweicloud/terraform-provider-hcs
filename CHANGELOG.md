# CHANGELOG

## 2.4.2 (June 29, 2024)

* The new resources and data sources added in the new version are as follows:

* **New Resources:**
  + `hcs_dcs_backup`
  + `hcs_dcs_instance`
  + `hcs_swr_image_retention_policy`
  + `hcs_swr_image_trigger`
  + `hcs_swr_organization`
  + `hcs_swr_repository_sharing`

* **New Data Sources:**
  + `hcs_dcs_flavors`
  + `hcs_dcs_instances`
  + `hcs_dcs_template_detail`
  + `hcs_dcs_templates`

## 2.4.1 (March 29, 2024)

* Fix the bug that need create BMS instance without security_groups.

## 2.4.0 (March 5, 2024)

* Release a version for Huawei Cloud Stack v8.3.1.
  The new resources and data sources added in the new version are as follows.

  * **New Resources:**
    + `hcs_dws_cluster`
    + `hcs_dws_alarm_subscription`
    + `hcs_dws_event_subscription`
    + `hcs_dws_ext_data_source`
    + `hcs_dws_snapshot`
    + `hcs_dws_snapshot_policy`
    + `hcs_gaussdb_opengauss_instance`
    + `hcs_mrs_cluster`
    + `hcs_mrs_job`
    + `hcs_waf_address_group`
    + `hcs_waf_certificate`
    + `hcs_waf_dedicated_domain`
    + `hcs_waf_dedicated_instance`
    + `hcs_waf_policy`
    + `hcs_waf_reference_table`
    + `hcs_waf_rule_blacklist`
    + `hcs_waf_rule_cc_protection`
    + `hcs_waf_rule_data_masking`
    + `hcs_waf_rule_geolocation_access_control`
    + `hcs_waf_rule_known_attack_source`
    + `hcs_waf_rule_global_protection_whitelist`
    + `hcs_waf_rule_information_leakage_prevention`
    + `hcs_waf_rule_precise_protection`
    + `hcs_waf_rule_web_tamper_protection`

  * **New Data Sources:**
    + `hcs_dws_flavors`
    + `hcs_gaussdb_opengauss_instance`
    + `hcs_gaussdb_opengauss_instances`
    + `hcs_mrs_versions`
    + `hcs_mrs_clusters`
    + `hcs_waf_certificate`
    + `hcs_waf_dedicated_instances`
    + `hcs_waf_policies`
    + `hcs_waf_reference_tables`

## 2.3.4 (January 27, 2024)

* Add more unit tests.
* Fix some bugs.

* **New Data Sources:**
  + `hcs_volume_types`

## 2.3.1 (September 27, 2023)

* Add more unit tests and more examples.
* Fix some bugs.
* Refine the docs.

## 2.3.0 (September 11, 2023)

* **New Resources:**
    + `hcs_cce_addon`
    + `hcs_cce_cluster`
    + `hcs_cce_namespace`
    + `hcs_cce_node`
    + `hcs_cce_node_attach`
    + `hcs_cce_node_pool`
    + `hcs_cce_pvc`
    + `hcs_cfw_address_group`
    + `hcs_cfw_address_group_member`
    + `hcs_cfw_black_white_list`
    + `hcs_cfw_eip_protection`
    + `hcs_cfw_protection_rule`
    + `hcs_cfw_service_group_member`
    + `hcs_cfw_service_group`
    + `hcs_enterprise_project`
    + `hcs_dns_recordset`
    + `hcs_dns_zone`
    + `hcs_vpcep_approval`
    + `hcs_vpcep_endpoint`
    + `hcs_vpcep_service`
    + `hcs_vpc_bandwidth`
    + `hcs_vpc_eip`
    + `hcs_vpc_eip_associate`
    + `hcs_vpc_bandwidth_associate`
    + `hcs_vpc_bandwidth_v2`
    + `hcs_vpc_eip_v1`
    + `hcs_evs_volume`
    + `hcs_evs_snapshot`
    + `hcs_elb_certificate`
    + `hcs_elb_l7policy`
    + `hcs_elb_l7rule`
    + `hcs_elb_listener`
    + `hcs_elb_loadbalancer`
    + `hcs_elb_member`
    + `hcs_elb_monitor`
    + `hcs_elb_pool`
    + `hcs_elb_security_policy`
    + `hcs_ecs_compute_volume_attach`
    + `hcs_ecs_compute_server_group`
    + `hcs_ecs_compute_interface_attach`
    + `hcs_ecs_compute_instance`
    + `hcs_ecs_compute_keypair`
    + `hcs_ecs_compute_eip_associate`
    + `hcs_networking_eip_associate`
    + `hcs_nat_gateway`
    + `hcs_nat_snat_rule`
    + `hcs_nat_dnat_rule`
    + `hcs_smn_topic`
    + `hcs_smn_subscription`
    + `hcs_smn_message_template`
    + `hcs_smn_topic_v2`
    + `hcs_smn_subscription_v2`
    + `hcs_as_bandwidth_policy`
    + `hcs_as_configuration`
    + `hcs_as_group`
    + `hcs_as_instance_attach`
    + `hcs_as_lifecycle_hook`
    + `hcs_as_notification`
    + `hcs_as_policy`
    + `hcs_ims_image`
    + `hcs_ims_image_share`
    + `hcs_ims_image_share_accepter`
    + `hcs_vpc`
    + `hcs_vpc_subnet`
    + `hcs_vpc_route`
    + `hcs_vpc_route_v2`
    + `hcs_vpc_v1`
    + `hcs_vpc_subnet_v1`
    + `hcs_networking_vip`
    + `hcs_networking_vip_associate`
    + `hcs_vpc_peering_connection`
    + `hcs_vpc_peering_connection_accepter`
    + `hcs_network_acl`
    + `hcs_network_acl_rule`
    + `hcs_networking_secgroup`
    + `hcs_networking_secgroup_rule`
    + `hcs_bms_instance`

* **New Data Sources:**
    + `hcs_enterprise_project`
    + `hcs_vpcep_public_services`
    + `hcs_cfw_firewalls`
    + `hcs_vpc_bandwidth`
    + `hcs_vpc_eip`
    + `hcs_vpc_eips`
    + `hcs_evs_volumes`
    + `hcs_elb_certificate`
    + `hcs_elb_pools`
    + `hcs_nat_gateway`
    + `hcs_smn_topics`
    + `hcs_ims_images`
    + `hcs_vpc`
    + `hcs_vpc_subnet`
    + `hcs_vpc_subnet_v1`
    + `hcs_vpc_subnet_ids`
    + `hcs_vpc_subnet_ids_v1`
    + `hcs_vpcs`
    + `hcs_vpc_subnets`
    + `hcs_vpc_peering_connection`
    + `hcs_networking_port`
    + `hcs_networking_secgroup`
    + `hcs_networking_secgroups`
    + `hcs_as_configurations`
    + `hcs_as_groups`
    + `hcs_bms_flavors`
    + `hcs_cce_cluster`
    + `hcs_cce_clusters`
    + `hcs_cce_addon_template`
    + `hcs_cce_node_pool`
    + `hcs_cce_node`
    + `hcs_cce_nodes`
    + `hcs_availability_zones`
    + `hcs_ecs_compute_flavors`
    + `hcs_ecs_compute_instance`
    + `hcs_ecs_compute_instances`
    + `hcs_ecs_compute_servergroups`
