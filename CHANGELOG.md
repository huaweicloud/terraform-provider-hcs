# CHANGELOG

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
