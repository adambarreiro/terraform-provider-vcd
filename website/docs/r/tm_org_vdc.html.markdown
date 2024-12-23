---
layout: "vcd"
page_title: "VMware Cloud Director: vcd_tm_org_vdc"
sidebar_current: "docs-vcd-resource-tm-org-vdc"
description: |-
  Provides a resource to manage VMware Cloud Foundation Tenant Manager Organization VDC.
---

# vcd\_tm\_org\_vdc

Provides a resource to manage VMware Cloud Foundation Tenant Manager Organization VDC.

This resource is exclusive to **VMware Cloud Foundation Tenant Manager**. Supported in provider
*v4.0+*

## Example Usage

```hcl
data "vcd_vcenter" "vc" {
  name = "my-vcenter"
}

data "vcd_tm_supervisor" "supervisor" {
  name       = "my-supervisor"
  vcenter_id = vcd_vcenter.vc.id
}

data "vcd_tm_region" "one" {
  name = "region-one"
}

data "vcd_tm_region_zone" "one" {
  region_id = data.vcd_tm_region.one.id
  name      = "my-zone"
}

resource "vcd_tm_org_vdc" "first" {
  name           = "my-org-vdc"
  org_id         = vcd_tm_org.test.id
  region_id      = data.vcd_tm_region.one.id
  supervisor_ids = [data.vcd_tm_supervisor.test.id]
  zone_resource_allocations {
    region_zone_id         = data.vcd_tm_region_zone.one.id
    cpu_limit_mhz          = 2000
    cpu_reservation_mhz    = 100
    memory_limit_mib       = 1024
    memory_reservation_mib = 512
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name for the Org VDC
* `org_id` - (Required) An Org ID for this Org VDC to be assigned to
* `region_id` - (Required) A Region ID that this Org VDC should be backed by
* `supervisor_ids` - (Required) A set of Supervisor IDs that back this Org VDC. Can be looked up
  using [`vcd_tm_supervisor`](/providers/vmware/vcd/latest/docs/data-sources/tm_supervisor) data source
* `zone_resource_allocations` - (Required) A set of Zone Resource Allocation definitions. See [Zone Resource Allocations](#zone-resource-allocations-block)

<a id="zone-resource-allocations-block"></a>
## Zone Resource Allocations

* `region_zone_id` - (Required) Can be looked up using
  [`vcd_tm_region_zone`](/providers/vmware/vcd/latest/docs/data-sources/tm_region_zone) data source
* `cpu_limit_mhz` - (Required) Maximum CPU consumption limit in MHz
* `cpu_reservation_mhz` - (Required) Defines reserved CPU capacity in MHz
* `memory_limit_mib` - (Required) Maximum memory consumption limit in MiB
* `memory_reservation_mib` - (Required) Defines reserved memory capacity in MiB

A computed attribute `region_zone_name` will be set in each `zone_resource_allocations` block.


## Attribute Reference

The following attributes are exported on this resource:

* `status` - The creation status of the Org VDC. Possible values are `READY`, `NOT_READY`, `ERROR`,
  `FAILED`

## Importing

~> **Note:** The current implementation of Terraform import can only import resources into the
state. It does not generate configuration. However, an experimental feature in Terraform 1.5+ allows
also code generation. See [Importing resources][importing-resources] for more information.

An existing Org VDC configuration can be [imported][docs-import] into this resource
via supplying path for it. An example is
below:

[docs-import]: https://www.terraform.io/docs/import/

```
terraform import vcd_tm_org_vdc.imported my-org-name.my-org-vdc-name
```

The above would import the `my-org-vdc-name` Tenant Org VDC that belongs to my-org-name Organization.