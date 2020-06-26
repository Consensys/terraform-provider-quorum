---
layout: "quorum"
page_title: "Quorum: quorum_bootstrap_account"
sidebar_current: "docs-quorum-bootstrap-account"
description: |-
   Use this resource to create a new Ethereum account
---

# quorum_bootstrap_account

Use this resource to create a new Ethereum account

## Example Usage

```hcl
resource "quorum_bootstrap_keystore" "test" {
  keystore_dir         = "%s"
  use_light_weight_kdf = false
}
resource "quorum_bootstrap_account" "test" {
  wallet_id  = quorum_bootstrap_keystore.test.id
  passphrase = "test"
}
```

## Argument Reference

- `balance` - (Optional) A place holder to keep account initial balance for referencing
- `passphrase` - (Optional) Passphrase to lock/unlock the account. Default is empty
- `wallet_id` - (Required) ID of a wallet storing the newly created account. For keystore, it's the keystore resource id

## Attributes Reference

- `account_url` - URL of the newly generated account
- `address` - Address of the newly generated account
ted account

## Attributes Reference

- `account_url` - URL of the newly generated account
- `address` - Address of the newly generated account
