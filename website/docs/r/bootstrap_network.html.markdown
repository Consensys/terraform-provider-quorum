---
layout: "quorum"
page_title: "Quorum: quorum_bootstrap_network"
sidebar_current: "docs-quorum-bootstrap-network"
description: |-
   Use this resource to create a new directory that represents a new Quorum network.
   
   Bootstraping data will be kept in this folder.
---

# quorum_bootstrap_network

Use this resource to create a new directory that represents a new Quorum network.

Bootstraping data will be kept in this folder.

## Example Usage

```hcl
resource "quorum_bootstrap_network" "test" {
  name       = "test-network"
  target_dir = "%s"
}
```

## Argument Reference

- `name` - (Required) Name of a new network. Directory name restriction applied
- `target_dir` - (Optional) File system path to the directory on which new directory will be created. Default is current working directory

## Attributes Reference

- `network_dir_abs` - Absolute path to a directory representing this new network
