---
subcategory: "Gateway"
---

# ovh_cloud_project_gateway

Create a new Gateway for existing subnet in the specified public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_project_network_private" "mypriv" {
  service_name = "xxxxxxxxxx"
  vlan_id      = "0"
  name         = "mypriv"
  regions      = ["GRA9"]
}

resource "ovh_cloud_project_network_private_subnet" "myprivsub" {
  service_name = ovh_cloud_project_network_private.mypriv.service_name
  network_id   = ovh_cloud_project_network_private.mypriv.id
  region       = "GRA9"
  start        = "10.0.0.2"
  end          = "10.0.255.254"
  network      = "10.0.0.0/16"
  dhcp         = true
}

resource "ovh_cloud_project_gateway" "gateway" {
  service_name = ovh_cloud_project_network_private.mypriv.service_name
  name         = "my-gateway"
  model        = "s"
  region       = ovh_cloud_project_network_private_subnet.myprivsub.region
  network_id   = tolist(ovh_cloud_project_network_private.mypriv.regions_attributes[*].openstackid)[0]
  subnet_id    = ovh_cloud_project_network_private_subnet.myprivsub.id
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) ID of the private network.
* `name` - (Required) Name of the gateway.
* `region` - (Required) Region of the gateway.
* `model` - (Required) Model of the gateway.
* `network_id` - (Required) ID of the private network.
* `subnet_id` - (Required) ID of the subnet.

## Attributes Reference

The following attributes are exported:

* `id` - Identifier of the gateway.
* `service_name` - See Argument Reference above.
* `name` - See Argument Reference above.
* `region` - Region of the gateway.
* `model` - See Argument Reference above.
* `network_id` - See Argument Reference above.
* `subnet_id` - See Argument Reference above.
* `external_information` - List of External Information of the gateway.
  * `network_id` - External network ID of the gateway.
  * `ips` - List of external ips of the gateway.
    * `ip` - External IP of the gateway.
    * `subnet_id` - Subnet ID of the ip.
* `interfaces` - Interfaces list of the gateway.
  * `id` - ID of the interface.
  * `ip` - IP of the interface.
  * `network_id` - Network ID of the interface.
  * `subnet_id` - Subnet ID of the interface.
* `status` - Status of the gateway.

## Import

A gateway can be imported using the `service_name`, `region` and `id` (identifier of the gateway) properties, separated by a `/`.
However, please note that in the case of a gateway import, `network_id` and `subnet_id` values used at gateway creation are not injected back in the state.
If you want to define these properties on your imported resource, you have to add an "ignore_changes" lifecycle argument in order not to trigger a recreation, as suggested in the following example. 

```terraform
resource "ovh_cloud_project_gateway" "imported_gateway" {
  service_name = ovh_cloud_project_network_private.mypriv.service_name
  name         = "<my-imported-gateway>"
  model        = "<my-model>"
  region       = "<my-region>"
  network_id   = "<my-imported-gateway-network-id>"
  subnet_id    = "<my-imported-gateway-subnet-id>"
  lifecycle {
    ignore_changes = [network_id, subnet_id]
  }
}

import {
  id = "<service-name>/<region>/<gateway-id>"
  to = ovh_cloud_project_gateway.imported_gateway
}
```

```bash
$ terraform import ovh_cloud_project_gateway.gateway service_name/region/id
```
