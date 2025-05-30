---
subcategory : "Additional IP"
---

# ovh_ip_service (Data Source)

Use this data source to retrieve information about an IP service.

## Example Usage

```terraform
data "ovh_ip_service" "my_ip" {
  service_name  = "XXXXXX"
}
```

## Argument Reference

* `service_name` - The service name

## Attributes Reference

`id` is set to ip service ip block. In addition, the following attributes are exported.

* `can_be_terminated` - can be terminated
* `country` - country
* `description` - Custom description on your ip
* `ip` - ip block
* `organisation_id` - IP block organisation Id
* `routed_to` - Routage information
  * `service_name` - Service where ip is routed to
* `type` - Possible values for ip type ( "cdn", "cloud", "dedicated", "failover", "hosted_ssl", "housing", "loadBalancing", "mail", "overthebox", "pcc", "pci", "private", "vpn", "vps", "vrack", "xdsl")
