---
layout: "quorum"
page_title: "Quorum: quorum_bootstrap_data_dir"
sidebar_current: "docs-quorum-bootstrap-data-dir"
description: |-
   Use this resource to create a data dir locally. This equivalent to execute `geth init`.
---

# quorum_bootstrap_data_dir

Use this resource to create a data dir locally. This equivalent to execute `geth init`.

## Example Usage

```hcl
resource "quorum_bootstrap_data_dir" "test" {
  data_dir = "%s"
  genesis  = <<EOF
{
  "alloc": { },
  "coinbase": "0x0000000000000000000000000000000000000000",
  "config": {
    "byzantiumBlock": 1,
    "chainId": 10,
    "eip150Block": 1,
    "eip155Block": 0,
    "eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "eip158Block": 1,
    "isQuorum": true
  },
  "difficulty": "0x00",
  "extraData": "0x00",
  "gasLimit": "0xE0000000",
  "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "nonce": "0x00",
  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "timestamp": "0x00",
  "number": "0x00",
  "gasUsed": "0x00"
}
EOF
}
```

## Argument Reference

- `data_dir` - (Required) Directory to intialize a genesis block
- `genesis` - (Required) Genesis file content in JSON format
- `instance_name` - (Optional) The instance name of the node. This must be the same as the value in geth node config. Default is `geth`

## Attributes Reference

- `data_dir_abs` - Absolute path to the data dir
