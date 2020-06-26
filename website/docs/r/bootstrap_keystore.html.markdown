---
layout: "quorum"
page_title: "Quorum: quorum_bootstrap_keystore"
sidebar_current: "docs-quorum-bootstrap-keystore"
description: |-
   Use this resource to create a keystore which maintains multiple Ethereum accounts.
---

# quorum_bootstrap_keystore

Use this resource to create a keystore which maintains multiple Ethereum accounts.

## Example Usage

```hcl
resource "quorum_bootstrap_keystore" "test" {
  keystore_dir         = "%s"
  use_light_weight_kdf = false
  account {
  }
  account {
    passphrase = "acc1"
  }
}
```

## Argument Reference

- `account` - (Optional) Account being created under this keystore

    Each `account` supports the following

    - `account_url` - Local path to the JSON representation of newly generated account private key
    - `address` - Address of the newly generated account
    - `balance` -(Optional) A place holder to keep account initial balance for referencing
    - `passphrase` -(Optional) Passphrase to lock/unlock the account. Default is empty

- `keystore_dir` - (Required) Directory contains private keys
- `use_light_weight_kdf` - (Optional) True to lower the memory and CPU requirements of the key store scrypt KDF at the expense of security

## Attributes Reference

- `keystore_dir_abs` - Absolute path of the keystore directory
