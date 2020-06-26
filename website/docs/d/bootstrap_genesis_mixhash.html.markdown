---
layout: "quorum"
page_title: "Quorum: quorum_bootstrap_genesis_mixhash"
sidebar_current: "docs-quorum-bootstrap-genesis-mixhash"
description: |-
   Use this data source to obtain `MixHash` value being used in the genesis file.
   Especially when using Istanbul consensus algorithm which defines a constant digest value.
---

# quorum_bootstrap_genesis_mixhash

Use this data source to obtain `MixHash` value being used in the genesis file.
Especially when using Istanbul consensus algorithm which defines a constant digest value.

## Example Usage

```hcl
data "quorum_bootstrap_genesis_mixhash" "test" {
}

output "istanbul_mix_hash" {
  value = data.quorum_bootstrap_genesis_mixhash.test.istanbul
}
```

## Argument Reference


## Attributes Reference

- `istanbul` - Mix hash that is used to identify whether the block is from Istanbul consensus engine
