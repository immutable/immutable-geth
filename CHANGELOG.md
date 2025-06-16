# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.0.0-beta.16]

- Return `413` status code instead of `500` when RPC request body exceeds size limit 

## [v1.0.0-beta.15]

- Reduce noisy logs

## [v1.0.0-beta.14]

- Limit the number of bytes read by NR RPC middleware
- Add fix [#30014](https://github.com/ethereum/go-ethereum/pull/30014) and [#30430](https://github.com/ethereum/go-ethereum/pull/30430) from upstream
- Enforce pricelimit (10 gwei) on rewards values returned by `eth_feeHistory`

## [v1.0.0-beta.13]

- Remove deployer allowlist
- Add `GETH_FLAG_IMMUTABLE_LONG_RANGE_SYNC` flag to allow snap sync from genesis

## [v1.0.0-beta.12]

#### Cancun

This release enables Cancun fork on all Immutable zkEVM networks.
| Network | Unix Timestamp | Date |
| -------- | ------- | ------- |
| Devnet | 1724796000 | Tue Aug 27 22:00:00 UTC 2024 |
| Testnet | 1727128800 | Mon Sep 23 22:00:00 UTC 2024 |
| Mainnet | 1728338400 | Mon Oct 7 22:00:00 UTC 2024 |

- Enable `ExcessBlobGas`, `BlobGasUsed`, `ParentBeaconRoot`
  - All values are `0x0`
- Enable `WithdrawalsHash`, and `Withdrawals` headers
  - `Withdrawals` are empty and `WithdrawalsHash` is the corresponding digest
- Enable `TSTORE`, `TLOAD`, and `MCOPY` op codes
- Enable `BLOBHASH` op code
- Enable Point Evaluation precompile
- Update `SELFDESTRUCT` op code
- Update clique stack to support Cancun
- Disable blob transactions
- Re-order logic that sets blob block headers and beacon root header to after `engine.Prepare` in order to interoperate with clique which produces variable blocktimes (rather than fixed slots)

#### Other

- Add `--gossipdefault` flag to toggle default geth tx gossiping
- Add GETH_FLAG_P2P_SUBNET env var to limit inbound messages based on subnet
- Add specific peer message handling in `eth/handler.go` to disable ingestion of state
- Add `--rpcproxy` flag to toggle RPC proxy forwarding
- Add RPC proxy forwarding to Immutable zkEVM
- Add `--disabletxpoolgossip` flag to disable tx gossiping
- Add `--gossipdefault` flag to toggle default geth tx gossiping
- Add GETH_FLAG_P2P_SUBNET env var to limit inbound messages based on subnet

## [v1.0.0-beta.11]

- Added partner public role to init container logic

## [v1.0.0-beta.10]

- Add log to correlate block period with block number and hash

## [v1.0.0-beta.9]

- Reject forkids that do not contain prevrandao fork
- Disable more RPC namespaces
- Correct the embedded mainnet.toml's price limit value to 10 gwei

## [v1.0.0-beta.8]

- Add prevrandao fork to forkid
- Pull v1.13.15 from upstream: [ethereum/go-ethereum](https://github.com/ethereum/go-ethereum)
- Log peer fullnames rather than abbreviated
- Do not rate limit peer connections that match on the supplied networks from the net restrict configuration
- Reduce p2p discv4 default refresh interval

## [v1.0.0-beta.7]

- Create Prevrandao fork
- Fix issues syncing with testnet and mainnet relating to Prevrandao fork

## [v1.0.0-beta.6]

- Update geth version logic to be based on immutable/go-ethereum releases
- Add `downloader/sync` metric

## [v1.0.0-beta.5]

- Pull v1.13.14 from upstream: [ethereum/go-ethereum](https://github.com/ethereum/go-ethereum)
- Make `DefaultBaseFeeChangeDenominator` consistent with upstream (8) and mutate it based on chain configuration
  - If chain configuration matches Immutable zkEVM network ID and clique settings, set `DefaultBaseFeeChangeDenominator` to 50
- Validate chain configuration based on expected values if Immutable zkEVM network ID is specified
- Use minimum price limit instead of last price to represent the priority fee of each fetched empty block inside suggested tip cap endpoint.
- Add block period, block propagation, and suggested tip cap metrics

## [v1.0.0-beta.4]

- Fix PREVRANDAO opcode on shanghai network by setting random to mixhash

## [v1.0.0-beta.3]

- Testing release flow

## [v1.0.0-beta.2]

- Enforce minimum price limit on suggested tip cap endpoint

## [v1.0.0-beta.1]

This release enables Shanghai fork on all Immutable zkEVM networks.
| Network | Unix Timestamp | Date |
| -------- | ------- | ------- |
| Devnet | 1709067600 | Tue Feb 27 21:00:00 UTC 2024 |
| Testnet | 1710280800 | Tue Mar 12 22:00:00 UTC 2024 |
| Mainnet | 1711490400 | Tue Mar 26 22:00:00 UTC 2024 |

The following changes were made to support the forks:

- Update clique to
  - allow for Shanghai to be enabled; and
  - log more detail around existence of withdrawals and withdrawals hashes
- Update geth node initialization from genesis
  - Add command to be used instead of `geth init`
  ```sh
  geth immutable bootstrap rpc --zkevm testnet --datadir "..."
  ```
  - genesis.json has been removed from Docker image
- Update geth node run configuration
  - Add `--zkevm [testnet|mainnet]` flag to automatically configure genesis and fork overrides

## [v0.0.16]

- Add `geth immutable run boot` command for running boot node on cluster
- Move boot node p2p keys to Secrets Manager
- Update logs around allowlists
- Added more logging for ACL initialisation and a log when ACL rejects a TX from pool

## [v0.0.15]

- Re-enabled single sequencer invariants

## [v0.0.14]

- Change pod initialization process to use Put instead of Push

## [v0.0.13]

- Implement Immutable wallet backend
  - Add AWS Secrets Manager backend store implementation for private key access via AWS
  - Add local keystore backend store implementation for testing purposes
  - Add `GETH_FLAG_IMMUTABLE_AWS_REGION` and `POD_NAMESPACE` env var support for configuring aws wallet backend
  - Add `GETH_FLAG_PASSWORD_FILEPATH` env var to existing password filepath flag
- Added flag to disable Clique endpoints on RPC server
- Added Clique Client
- Added CLI for Clique Voting

## [v0.0.12]

- Revert dialHistoryExpiration change back to 35s

## [v0.0.11]

- Revert change to always drop peers when find node errors are greater than the threshold
- Reduce dialHistoryExpiration from 35 seconds to 5

## [v0.0.10]

- Decrease `seedMaxAge` to 1 second to ensure that nodes always query boot node for peer IPs
- Always drop peers when find node errors are greater than the threshold

## [v0.0.9]

- Automation / CD updates only. Releasing with latest prod workflows.

## [v0.0.8]

- Added max block range for logs queried using eth_getLogs (#346)
- Default mainnet configuration (`/etc/geth/mainnet.toml`) and genesis (`/etc/geth/mainnet.json`) added to Docker image

## [v0.0.7]

- Automation / CD updates only: immutable init command

## [v0.0.6]

- Automation / CD updates only

## [v0.0.5]

- Reduce price limit 100->10 Gwei

## [v0.0.4]

### Added

- Default testnet configuration (`/etc/geth/testnet.toml`) and genesis (`/etc/geth/testnet.json`) added to Docker image

## [v0.0.3]

### Added

- Add NewRelic agent to geth runtime for application metrics

### Fixed

- Price limit being 100 Gwei rather than 100 Gwei - 1 Wei
- Fixed bug in ACLs that caused transactions to be blocked when they should not have been

## [v0.0.2]

### Changed

- Change to add access control layer to geth legacypool, which safeguards txs
  from entering txpool if the sender is part of a collective SDN list (#17)
- Mempool Rebroadcasting: Legacy pool now rebroadcasts pending transactions
  whenever a pool reorg is triggered. Pool reorgs are triggered when a new block
  is received - (#103)
- Transaction Broadcasting: Transactions are broadcast to all peers, not a
  square root. All peers get the transaction, regardless of whether they may
  have previously received it. All peers get the full transaction payload, i.e.
  no more announcements. (#103)
- Changed `maxUncleDist` from 7 to 0. This means old blocks are rejected,
  reducing chance of reorgs. (#156)
- Changed `maxQueueDist` from 32 to 256. This means more blocks in the future
  can be received without being rejected. This is to account for our increased
  block production due to reduced block time. (#156)
- Changed `blockLimit` from 64 to 256. Since we only have a small amount of
  peers, more will come from the same set of peers, so we need to allow for
  this. (#156)
- Added invariants to prevents reorgs (#175)
- Added check for NoLocals on legacytxpool to reject underpriced TXs
- Changed the `DefaultBaseFeeChangeDenominator` from 8 to 50. This makes the max
  base fee rate of change 2%, instead of 12.5%. With a block time of 2 seconds,
  we match Ethereum in that it would take 72 seconds for the base fee to double.
  (#178)
- Mine a tx even if its effective tipcap (tipcap-basefee) is less than miner's minimum tipcap (which is pricelimit, or 100 gwei) (#258)

### Added

- CLI commands for bootstrapping local and remote Immutable chains
- Added CLI parameter options for the geth immutable bootstraper, which can
  potentially takes in a list of filepaths for blocklists and allowlists (#83)
- Added flags to disable Admin/Txpool/Engine/Debug endpoints on RPC server
