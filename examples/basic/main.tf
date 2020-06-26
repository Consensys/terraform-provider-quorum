variable "network_name" {
  default = "my-network"
}

locals {
  default_account_config = {
    passphrase = ""
    balance    = "1000000000000000000000000000"
  }
}

module "etherum" {
  source       = "./network_bootstrap"
  network_name = var.network_name
  // this also defines how many nodes there are in the network
  nodes_config = [
    // node 1 configuration
    {
      accounts = [
        local.default_account_config,
        {
          passphrase = ""
          balance    = "2000000000000000000000000000"
        },
      ] // create 2 new accounts
    },
    // node 2 configuration
    {
      accounts = [
        local.default_account_config,
        {
          passphrase = ""
          balance    = "2000000000000000000000000000"
        },
      ] // create 2 new accounts
    },
    // node 3 configuration
    {
      accounts = [
        local.default_account_config,
        {
          passphrase = ""
          balance    = "2000000000000000000000000000"
        },
      ] // create 2 new accounts
    }
  ]
}
