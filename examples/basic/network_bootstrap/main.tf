variable "network_name" {
  type        = string
  description = "Name of the network"
}

variable "nodes_config" {
  type        = list
  description = "A complex structure configures nodes in a network"
}

variable "node_dir_prefix" {
  default = "node-"
}

locals {
  number_of_nodes = length(var.nodes_config)
  genesis_file    = format("%s/%s", quorum_bootstrap_network.this.network_dir_abs, "genesis.json")
  accounts_alloc  = <<-EOT
  %{for a in flatten(quorum_bootstrap_keystore.node.*.account)~}
    "${a.address}" : {
      "balance": "${a.balance}"
    },
  %{endfor~}
  EOT
}

resource "random_integer" "network_id" {
  max = 3000
  min = 1400
}

resource "quorum_bootstrap_network" "this" {
  name = var.network_name
}

resource "quorum_bootstrap_keystore" "node" {
  count        = local.number_of_nodes
  keystore_dir = format("%s/%s%s/%s", quorum_bootstrap_network.this.network_dir_abs, var.node_dir_prefix, count.index, "keystore")

  dynamic "account" {
    for_each = var.nodes_config[count.index].accounts
    content {
      passphrase = lookup(account.value, "passphrase", "")
      balance    = lookup(account.value, "balance", "1000000000000000000000000000")
    }
  }
}

resource "quorum_bootstrap_data_dir" "node" {
  count    = local.number_of_nodes
  data_dir = format("%s/%s%s", quorum_bootstrap_network.this.network_dir_abs, var.node_dir_prefix, count.index)
  genesis  = local_file.genesis.content
}

resource "local_file" "genesis" {
  filename = local.genesis_file
  content  = <<-JSON
  {
    "alloc": {
      ${substr(local.accounts_alloc, 0, length(local.accounts_alloc) - 2)}
    },
    "coinbase": "0x0000000000000000000000000000000000000000",
    "config": {
      "homesteadBlock": 0,
      "byzantiumBlock": 0,
      "constantinopleBlock":0,
      "chainId": ${random_integer.network_id.result},
      "eip150Block": 0,
      "eip155Block": 0,
      "eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "eip158Block": 0,
      "isQuorum": true,
      "maxCodeSize": 50
    },
    "difficulty": "0x0",
    "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "gasLimit": "0xFFFFFF00",
    "mixhash": "0x00000000000000000000000000000000000000647572616c65787365646c6578",
    "nonce": "0x0",
    "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "timestamp": "0x00"
  }
  JSON
}