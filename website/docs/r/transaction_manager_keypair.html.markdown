---
layout: "quorum"
page_title: "Quorum: quorum_transaction_manager_keypair"
sidebar_current: "docs-quorum-transaction-manager-keypair"
description: |-
   Use this resource to create a key pair used in a transaction manager.
   
   This key pair provides attributes which are useful when building the configuration for a transaction manager.
---

# quorum_transaction_manager_keypair

Use this resource to create a key pair used in a transaction manager.

This key pair provides attributes which are useful when building the configuration for a transaction manager.

## Example Usage

```hcl
resource "quorum_transaction_manager_keypair" "test" {
  password = "foo"
}
```

## Argument Reference

- `config` - (Optional) Key generation config

    Each `config` supports the following

    - `iterations` -(Optional) Number of iterations to cycle through
    - `memory` -(Optional) Memory limit
    - `parallelism` -(Optional) Number of threads to use
    - `variant` -(Optional) Algorithm to use when hashing. Allowed values are `id` or `i`

- `password` - (Optional) A password to protect the keypair

## Attributes Reference

- `key_data` - Key Data in JSON format to be used by Private Transaction Manager
- `private_key_json` - Private key in JSON representation
- `public_key_b64` - Public key in standard base64 encoding
