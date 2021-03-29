## v0.2.0

**Updated Resources**
- `quorum_bootstrap_istanbul_extradata`: Added a new argument `ibft2_mode` that creates the `extraData` for genesis JSON used for IBFT2 consensus on Hyperledger Besu.

## v0.1.0

*Released on March 24th 2021*

This is an initial release

**New Resources**
- `quorum_bootstrap_account`: Create an Ethereum account in a provided wallet. Only `KeyStore` wallet is currently supported
- `quorum_bootstrap_data_dir`: Perform `geth init` on a directory with the given genesis JSON file
- `quorum_bootstrap_istanbul_extradata`: For istanbul consensus algorithm, create `extraData` value being used in genesis JSON
- `quorum_bootstrap_keystore`: Create a `KeyStore` managing a key storage directory on disk
- `quorum_bootstrap_network`: Create a folder containing all resources required to bootstrap a network
- `quorum_bootstrap_node_key`: Create a node key being used by a Quorum node
- `quorum_transaction_manager_keypair`: Create Private Transaction keypair for privacy

**New Data Sources**
- `quorum_bootstrap_genesis_mixhash`: Provide `mixHash` value constants being used in genesis JSON
- `quorum_bootstrap_node_key`: Read node key hex value

**CI/CD**
- Configure for official [TF registry](https://registry.terraform.io/browse/providers) releases

## START
