---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "maas_tag Resource - terraform-provider-maas"
subcategory: ""
description: |-
  Provides a resource to manage a MAAS tag, used to tag machines.
---

# maas_tag (Resource)

Provides a resource to manage a MAAS tag, used to tag machines.

## Example Usage

```terraform
resource "maas_tag" "kvm" {
  name = "kvm"
  machines = [
    maas_machine.virsh_vm1.id,
    maas_machine.virsh_vm2.id,
    maas_vm_host_machine.kvm[0].id,
    maas_vm_host_machine.kvm[1].id,
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The new tag name. Because the name will be used in urls, it should be short.

### Optional

- `comment` (String) A description of what the the tag will be used for in natural language.
- `definition` (String) An XPATH query that is evaluated against the hardware_details stored for all nodes. (i.e. the output of ``lshw -xml``)
- `kernel_opts` (String) Nodes associated with this tag will add this string to their kernel options when booting. The value overrides the global ``kernel_opts`` setting. If more than one tag is associated with a node, command line will be concatenated from all associated tags, in alphabetic tag name order.
- `machines` (Set of String) List of MAAS machines' identifiers (system ID, hostname, or FQDN) that will be tagged with the new tag.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# An existing tag can be imported using its name. e.g.
$ terraform import maas_tag.kvm kvm
```
