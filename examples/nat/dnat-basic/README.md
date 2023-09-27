# Basic NAT gateway and DNAT rule

This example provisions:

* a ECS instance:
* a VPC and subnet instance:
  > VPC provides an isolated cloud service environment.
* a basic NAT gateway.
* a EIP.
* a SNAT rule:
  > By binding external subnet IP, SNAT function can realize multiple virtual machines across availability zones
  to share external subnet IP and access external data center or other VPCs.
* two DNAT rule:
  > DNAT rule specifies port 8080 as the port for Tomcat to provide external services, which can be used to build nginx services.
  > DNAT rule specifies port 8022 as the port for SSH, which can be used to login in ECS by ssh command.
  > DNAT function is bound with EIP, and EIP is shared across VPC by binding IP mapping,
  which provides services for the Internet.
