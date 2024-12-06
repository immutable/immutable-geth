# Go-Ethereum Internal Audit

## Authors and Contributors

**Author**: 
- Peter Robinson
- Ryan Teoh

**Contributors**:
- Michael Baker
- Mark Art

## Purpose and Scope

This audit reviews Immutable's fork of [go-ethereum](https://github.com/immutable/geth). It analyses changes relative to the upstream version (https://github.com/ethereum/go-ethereum).

The scope of the audit is all file changes in Immutable's repository at GitHash [`f05c916773b30d40b02db377dd9d152ff218aee1`](https://github.com/immutable/immutable-geth/tree/f05c916773b30d40b02db377dd9d152ff218aee1) (equivalent to tag: audit.1, committed November 20, 2024) relative to Ethereum Foundation's go-ethereum repository at GitHash [`c5ba367eb6232e3eddd7d6226bfd374449c63164`](https://github.com/ethereum/go-ethereum/tree/c5ba367eb6232e3eddd7d6226bfd374449c63164) (equivalent to tag: v1.13.15, committed April 17, 2024). 

Additionally, the audit will check that all relevant fixes from go-ethereum releases v1.14.00 to v1.14.12 (released November 20, 2024) have been integrated into the immutable-geth repo


## Background

### Diff File

The souce code differences between the target of the audit and go-ethereum version 1.13.15 are available in this diff file: [forkdiff](https://immutable.github.io/immutable-geth/audit/202411-internal-audit/forkdiff.html). This file was generated using the commands:

```
mkdir temp
go env -w GOBIN=<absolute path>/temp  
go install github.com/protolambda/forkdiff@latest

git clone https://github.com/ethereum/go-ethereum.git
git clone https://github.com/immutable/immutable-geth.

cd temp
./forkdiff -fork ../immutable-geth/diff/fork.yaml -repo ../immutable-geth -upstream-repo ../go-ethereum -out forkdiff.html
```

### Description of Changes

Core changes directly modify key logic in Go-Ethereum that could result in soft and hard forks, as well as non-protocol changes to areas such as the mempool

Access Control: 

* Change to add access control layer to geth legacypool, which safeguards txs from entering txpool if the sender is part of a collective SDN list (#17)

* Added Allowlist for contract deployment

Key Management:

* Add AWS Secrets Manager backend store implementation for private key access via AWS

Block Time:

* Changed the DefaultBaseFeeChangeDenominator from 8 to 50. This makes the max base fee rate of change 2%, instead of 12.5%. With a block time of 2 seconds, we match Ethereum in that it would take 72 seconds for the base fee to double. (#178)

RPC:

* Max Filter range for logs queried using Get Logs
* Suggest Tip Calculation changes to handle empty blocks
* Minimum Suggested Price Limit of X gwei
* Added flags to disable Admin/Txpool/Engine/Debug endpoints on RPC server

Clique: 

* Updated Clique to allow for Shanghai and Cancun

Reorgs/Finality: 

* Changed maxUncleDist from 7 to 0. This means old blocks are rejected, reducing chance of reorgs. (#156)
* Added invariants to prevent reorgs (#175)
* Max Recommit Interval set to maximum to avoid reorgs

P2P: 

* Decrease seedMaxAge to 1 second to ensure that nodes always query boot node for peer IPs
* Mempool Rebroadcasting, as our network is smaller, there is less mempool and peer redundancy. Source changes made to support this.
* Transaction Broadcasting: Transactions are broadcast to all peers, not a square root. All peers get the transaction, regardless of whether they may have previously received it. All peers get the full transaction payload, i.e. no more announcements. (#103)
* Changed blockLimit from 64 to 256. Since we only have a small amount of peers, more will come from the same set of peers, so we need to allow for this. (#156)
* Changed maxQueueDist from 32 to 256. This means more blocks in the future can be received without being rejected. This is to account for our increased block production due to reduced block time. (#156)

Misc:

* Disabled Miner Tip Enforcement, allowing for transactions accepted by RPCs to be mined rather than being stuck.
* Prevrandao Hard Fork
* Shanghai Fork overrides

Supporting Changes:

_These are mostly additive changes that do not affect the runtime of Go-Ethereum, e.g. Local network bootstrapping_

* Automation
  * Deployments setup on Merges and Releases
  * CLI commands for bootstrapping local and remote Immutable chains
  * Invariants to ensure Single Sequencer
  * CLI for Clique Voting
  * Flags to run Immutable zkEVM with testnet/mainnet configurations
* Observability
  * Add NewRelic agent to geth runtime for application metrics
  * Dashboards 
  * Alerts
  * Custom Metrics
  * Block Propagation Time
  * Tip Cap
  * Sync Duration Timer
  * TxPoolOverflow Counter
  * Geth worker loop
  * Block construction time


## Attack Surfaces and Perceived Attacks

### EVM

The Ethereum Virtual Machine (EVM) is the execution environment in which contracts operate. 

Perceived attacks:

* Double spend / accounting errors
* Opcode changes cause unexpected / breaking behaviour in smart contracts
* Precompile changes allow attackers to bypass signature checks
* Bugs relating to Shanghai, Prevrandao, Cancun changes in clique stack and other parts of the client stack
* User is able to cause client to enter into a code path that should be impossible based on our invariants. Could cause client to crash loop.

Analysis of whether changes could be exploited given the perceived attacks:

* No changes have been made to the EVM or any precompiles.
* No changes have been made to geth’s state transition logic beyond the block header values and opcodes pertaining to Shanghai and Cancun upgrades.
* Note the following opcodes always return 0x0:
  * PREVRANDAO 
  * BLOBBASEFEE

### RPC

The Remote Procedure Call (RPC) interface is the external API used to submit transactions, execute queries (view calls), and interrogate blockchain state. Some RPC APIs have administration capabilities.


Perceived attacks:

* DDOS RPC nodes through specific requests
* Change configuration via admin RPC methods.

Analysis of whether changes could be exploited given the perceived attacks:

* NewRelic integration added for metrics pertaining to RPC requests.
* Admin and deprecated RPC methods are disabled.

### P2P

The Peer to Peer (P2P) protocol is the means by which blockchain nodes communicate with each other. It is used for discovering other blockchain nodes, consensus protocols, communication of pending transactions and minted blocks.

Perceived attacks:

* DDOS public P2P nodes
* Penetrate mempool through public P2P nodes

Analysis of whether changes could be exploited given the perceived attacks:

* P2P network is exposed through a dedicated node in each network. Node is configured to only allow inbound P2P traffic pertaining to specific subnets. This means that external connections should not be able to write to the mempool. This involved changes to geth.


### Consensus

The consensus protocol is the means by which the next block is determined. The consensus algorithm used for Immutable zkEVM must ensure no forking, as the games that rely on the blockchain are not designed to handle blockchain reorganisation.

Perceived attacks:

* Liveness: Do something such that the consensus algorithm is unable to produce the next block.
* Forking: Do something such that the consensus algorithm produces two blocks at the same block height.

Analysis of whether changes could be exploited given the perceived attacks:

* The Clique consensus protocol has not been modified, beyond allowing for the Shanghai and Cancun hard forks. 
  Fork prevention is addressed in the core [blockchain code](https://github.com/immutable/immutable-geth/blob/f05c916773b30d40b02db377dd9d152ff218aee1/core/blockchain.go#L2187).


## Recommendations

### LOW: Docker Base Container

#### Issue

The Docker container used as the Ethereum client executable to run the Immutable chain uses a  [Dockerfile](https://github.com/immutable/immutable-geth/blob/f05c916773b30d40b02db377dd9d152ff218aee1/Dockerfile). This Dockerfile sets the base Docker container used for deployment as [alpine](https://hub.docker.com/_/alpine) 3.18.4. This should be updated to the latest stable version in the 3.18 series: 3.18.9.

#### Team Response
Accepted and [fixed](https://github.com/immutable/immutable-geth/pull/8)


### LOW: Review of Recent Go-Ethereum Changes 

#### Issue

This section reviews changes made to go-ethereum in releases v1.14.00 to v1.14.12 that should be included 
in the Immutable Geth repo that haven't been included thus far. The two recommended bug fixes that should
be included are:

* [Pull Request 30430](https://github.com/ethereum/go-ethereum/pull/30430): Fix potential out-of-bound issue in mempool.
* [Pull Request 30014](https://github.com/ethereum/go-ethereum/pull/30014): Fix out of bounds access in json unmarshalling.

Resolving these has been deemed low as no code was provided demonstrating a practical exploit of these bugs.

#### Team Response
Accepted and [fixed](https://github.com/immutable/immutable-geth/pull/9)

### LOW: Clique Panics

#### Issue

The implementation of [Clique includes golang panics](https://github.com/immutable/immutable-geth/blob/f05c916773b30d40b02db377dd9d152ff218aee1/consensus/clique/clique.go#L824) when checking values of Withdrawals, Beacon Chain Parent Roots, Blob Gas, and Excess Blob Gas. If an attacker was able to submit a Blob transaction, it would cause the `panic` instruction to execute, which would crash the blockchain client. 

Blob instructions are rejected when constructing blocks from [transactions in the transaction pool](https://github.com/immutable/immutable-geth/blob/f05c916773b30d40b02db377dd9d152ff218aee1/miner/worker.go#L808). As such, the consensus algorithm will not see a blob instruction which could cause a panic.

Recommendation:
- Team should evaluate whether converting these panics into logged errors and returned error values would improve system reliability and maintainability.

- Blob instructions currently risk getting stuck in the transaction pool, consuming space until they are eventually dropped. Consider rejecting these transactions upon entry to the pool to prevent inefficiencies and avoid requiring users to submit substitute transactions with the same nonce.


#### Team Response
Accepted and pending fixes.


### LOW: Unlocked account on RPC Node

#### Issue

There’s an unlocked eth account on our mainnet RPC endpoint 0x8a68f4cd3726d39414fc965f6cfa841464a0b670 , this is likely related to [this](https://imtbl.slack.com/archives/C050AE6CHHV/p1719813817671079), any user can send transactions on behalf of the address using RPC method eth_sendTransaction , no security implication identified yet, but we should clean this up.

Curl Command
```
curl -X POST -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}' \
     https://rpc.immutable.com
```

Recommendation: 
- Identify and lock the account using the personal_lockAccount method.
- Review node configurations to ensure no accounts are unlocked by default.
- Regularly audit exposed RPC methods.


#### Team Response
Accepted, fix pending.

### LOW: Local File Inclusion via txpool.blocklistfilepaths

#### Issue

When initialising a new access controller using New(), the accesscontrol.load() function prints file content without validation. Since the file path is supplied via `txpool.blocklistfilepaths` cmd flag, user can misuse this function to print file content. Considering that local access is required, will rate this as low.

```
func load(filePath string) (map[common.Address]struct{}, error) {
	addresses := make(map[common.Address]struct{})

	log.Info("Loading ACL file", "filepath", filePath)
	byteValue, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	log.Info("Loaded ACL file", "filepath", filePath, "content", string(byteValue))
[REDACTED]
}
```
[Code reference](https://github.com/immutable/immutable-geth/blob/f05c916773b30d40b02db377dd9d152ff218aee1/core/txpool/immutable/accesscontrol/sdn_provider.go#L48)


#### Team Response
Accepted, fix pending.

### LOW: New Relic sensitive logging

#### Issue

The newRelicMiddleware wrapper function currently logs certain headers, including x-api-key, x-imx-eth-address, and x-forwarded-for. While this logging behavior can assist in debugging and monitoring, it raises potential concerns regarding the unintentional exposure of sensitive information. Logging headers such as x-api-key and x-imx-eth-address could inadvertently capture confidential data, potentially violating privacy or security standards. Similarly, logging the x-forwarded-for header, which often contains client IP addresses, might lead to the collection of personally identifiable information (PII) without adequate safeguards. To mitigate these risks, it is essential to assess whether logging these headers is strictly necessary and, if so, ensure that sensitive information is appropriately masked, encrypted, or omitted to align with best practices for secure logging and compliance.

```
_, handler := newrelic.WrapHandleFunc(nrApp, "/", func(w http.ResponseWriter, req *http.Request) {
	txn := newrelic.FromContext(req.Context())
	txn.AddAttribute("x-api-key", req.Header.Get("x-api-key"))

	// Capture ethAddress from headers in order to monitor usage and errors
	txn.AddAttribute("x-imx-eth-address", req.Header.Get("x-imx-eth-address"))
	// Capture SDK version from headers in order to monitor usage and errors
	txn.AddAttribute("x-sdk-version", req.Header.Get("x-sdk-version"))
	txn.AddAttribute("x-forwarded-for", req.Header.Get("x-forwarded-for"))
	txn.AddAttribute("x-zkevm-rpc-sticky", req.Header.Get("x-zkevm-rpc-sticky"))
	txn.AddAttribute("k6-load-test-id", req.Header.Get("k6-load-test-id"))
	// the W3C trace context header entries.
	// see: https://github.com/w3c/trace-context/blob/main/spec/20-http_request_header_format.md
	// RFC: https://www.w3.org/TR/trace-context/
	// NR, k6 and any other services following the W3C standard can generate these header entries.
	txn.AddAttribute("traceparent", req.Header.Get("traceparent"))
	txn.AddAttribute("tracestate", req.Header.Get("tracestate"))
```

[Code reference](https://github.com/immutable/immutable-geth/blob/f05c916773b30d40b02db377dd9d152ff218aee1/node/immutable_newrelic.go#L48)


#### Team Response
Accepted and [fixed](https://github.com/immutable/immutable-geth/pull/7/files). The logging platform access is strictly controlled with just-in-time role assignment, any access to production logs require internal approvals.


### LOW: New Relic WS and RPC Handler

#### Issue

The newRelicMiddleware wrapper is designed to support WebSocket (WS) and Remote Procedure Call (RPC) protocols. However, it appears to function exclusively with HTTP requests. Since only HTTP servers are currently enabled in the implementation, the handlers are not expected to process requests using other protocols. Given this limitation, the associated security risk is assessed as low due to the restricted attack surface and the unavailability of WS or RPC endpoints in the current setup. Nonetheless, it is advisable to confirm that this behavior aligns with the intended design and that additional protocols are appropriately disabled if not required.

```
func newRelicMiddleware(nrApp *newrelic.Application, next http.Handler) http.Handler { //nolint:unused
	// CHANGE(immutable) add NR agent
	if nrApp == nil {
		log.Error("Failed to initialise New Relic middlware: nrApp is nil")
		return next
	}
	_, handler := newrelic.WrapHandleFunc(nrApp, "/", func(w http.ResponseWriter, req *http.Request) {
		txn := newrelic.FromContext(req.Context())
    [REDACTED]
  }
}
```

[Code reference](https://github.com/immutable/immutable-geth/blob/f05c916773b30d40b02db377dd9d152ff218aee1/node/immutable_newrelic.go#L40)

#### Team Response
Accepted, fix pending.

## Conclusion

This document provides a comprehensive review of the audit findings and analysis for the internal fork of the Go Ethereum client. The findings highlight areas of strength, such as thorough reviews of the cmd, Bootstrap, and Configuration components, while identifying opportunities for improvement, including addressing go-lang panic commands and changes introduced from go-ethereum v1.14.00 to v1.14.12. By leveraging insights from industry leading Audit Framework and team expertise, the audit aligns with best practices, offering actionable recommendations to enhance the security and functionality of the fork. Addressing these recommendations within the timeline will strengthen the fork's robustness and ensure its alignment with the project's objectives.
