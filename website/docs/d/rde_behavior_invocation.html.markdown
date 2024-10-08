---
layout: "vcd"
page_title: "VMware Cloud Director: vcd_rde_behavior_invocation"
sidebar_current: "docs-vcd-data-source-rde-behavior-invocation"
description: |-
   Provides the capability of invoking an existing Runtime Defined Entity Behavior in VMware Cloud Director.
---

# vcd\_rde\_behavior\_invocation

~> This feature is **experimental** and may change in future

Provides the capability of invoking an existing [RDE Interface Behavior](/providers/vmware/vcd/latest/docs/resources/rde_interface_behavior)
or [RDE Type Behavior](/providers/vmware/vcd/latest/docs/resources/rde_type_behavior) in VMware Cloud Director.

Supported in provider *v3.11+*

## Example Usage

```hcl
resource "vcd_rde_interface" "interface" {
  nss     = "example"
  version = "1.2.3"
  vendor  = "vmware"
  name    = "MyInterface"
}

resource "vcd_rde_interface_behavior" "behavior" {
  rde_interface_id = vcd_rde_interface.interface.id
  name             = "MyInterfaceBehavior"
  description      = "A cool behavior example"
  execution = {
    "id" : "Noop"
    "type" : "noop"
  }
}

resource "vcd_rde_type" "type" {
  nss           = "example"
  version       = "1.1.0"
  vendor        = "vmware"
  name          = "MyType"
  description   = "A cool type example"
  interface_ids = [vcd_rde_interface.interface.id]
  schema        = file("../my-rde-type.json")

  # Behaviors can't be created after the RDE Interface is used by a RDE Type
  # so we need to depend on the Behaviors to wait for them to be created first.
  depends_on = [vcd_rde_interface_behavior.behavior]
}

resource "vcd_rde" "rde" {
  org          = "System"
  rde_type_id  = vcd_rde_type.type.id
  name         = "MyRde"
  resolve      = true
  input_entity = file("../my-rde.json")
}

# Required Access Levels to invoke Behaviors
resource "vcd_rde_type_behavior_acl" "interface_acl" {
  rde_type_id      = vcd_rde_type.type.id
  behavior_id      = vcd_rde_interface_behavior.behavior.id
  access_level_ids = ["urn:vcloud:accessLevel:FullControl"]
}

data "vcd_rde_behavior_invocation" "invoke" {
  rde_id            = vcd_rde.rde.id
  behavior_id       = vcd_rde_interface_behavior.behavior.id
  invoke_on_refresh = true
}

output "rde_behavior_invocation_output" {
  value = data.vcd_rde_behavior_invocation.invoke.result
}
```

## Argument Reference

The following arguments are supported:

* `rde_id` - (Required) The ID of the [RDE](/providers/vmware/vcd/latest/docs/resources/rde) which Behavior will be invoked
* `behavior_id` - (Required) The ID of the [RDE Interface Behavior](/providers/vmware/vcd/latest/docs/resources/rde_interface_behavior) or
  the [RDE Interface Behavior](/providers/vmware/vcd/latest/docs/resources/rde_type_behavior) to invoke
* `arguments` - (*Deprecated*; Optional) A map with the arguments of the invocation. Deprecated in favor of `arguments_json`
* `arguments_json` - (Optional) A JSON string with the arguments of the invocation. The built-in function `jsonencode` can
  be used to set it.
* `metadata` - (*Deprecated*; Optional) A map with the metadata of the invocation. Deprecated in favor of `metadata_json`
* `metadata_json` - (Optional) A JSON string with the arguments of the invocation. The built-in function `jsonencode` can
  be used to set it.
* `invoke_on_refresh` - (Optional) Defaults to `true`, indicates whether the Behavior should be invoked on every refresh.
  E.g. to suppress the invocation on the refresh step of `terraform destroy` operations, or to stop invoking a given Behavior temporarily.

## Attribute Reference

The following attributes are supported:

* `result` - The invocation result in plain text
