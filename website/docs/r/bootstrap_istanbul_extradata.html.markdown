---
layout: "quorum"
page_title: "Quorum: quorum_bootstrap_istanbul_extradata"
sidebar_current: "docs-quorum-bootstrap-istanbul-extradata"
description: |-
   Use this resource to construct `extradata` field used in the genesis file.
   
   `istanbul_address` can be referenced from `quorum_bootstrap_node_key` data source or newly created from `quorum_bootstrap_node_key` resources.
---

# quorum_bootstrap_istanbul_extradata

Use this resource to construct `extradata` field used in the genesis file.

`istanbul_address` can be referenced from `quorum_bootstrap_node_key` data source or newly created from `quorum_bootstrap_node_key` resources.

## Example Usage

```hcl
resource "quorum_bootstrap_node_key" "test" {
  count = 3
}

resource "quorum_bootstrap_istanbul_extradata" "test" {
  istanbul_addresses = quorum_bootstrap_node_key.test.*.istanbul_address
}
```

## Argument Reference

- `ibft2_mode` - (Optional) generate extradata using RLP encoding mode used by Hyperledger Besu for IBFT2
- `istanbul_addresses` - (Required) list of Istanbul address to construct extradata
- `vanity` - (Optional) Vanity Hex Value to be included in the extradata

## Attributes Reference

- `extradata` - Computed value which can be used in genesis file
