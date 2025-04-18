---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "maas_network_interface_tag Resource - terraform-provider-maas"
subcategory: ""
description: |-
  Provides a resource to manage tags as strings on a network interface.
---

# maas_network_interface_tag (Resource)

Provides a resource to manage tags as strings on a network interface.

## Example Usage

```terraform
resource "maas_network_interface_tag" "test" {
  machine      = "abc123"
  interface_id = 12
  tags = [
    "tag1",
    "tag2",
  ]
}

resource "maas_network_interface_tag" "test2" {
  device       = "cheerful-owl"
  interface_id = 13
  tags = [
    "tag3",
    "tag4",
  ]
}

resource "maas_network_interface_tag" "test3" {
  device       = "def456"
  interface_id = 14
  tags = [
    "tag3",
    "tag4",
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `interface_id` (Number) The network interface ID to tag.
- `tags` (Set of String) The tags to assign to the network interface. Tag names should be short and without spaces.

### Optional

- `device` (String) The identifier (system ID, hostname, or FQDN) of the device with the network interface. Either `machine` or `device` must be provided.
- `machine` (String) The identifier (system ID, hostname, or FQDN) of the machine with the network interface. Either `machine` or `device` must be provided.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# A network interface tag can be imported using the machine or device system id and the network interface id. e.g.
$ terraform import maas_network_interface_tag.test abc123/12
```
