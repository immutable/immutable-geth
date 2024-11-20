# Immutable E2E Tests

Use this pkg to run E2E tests against local or remote deployments.

You can use the CI scripts to spin up a local network and test:

```sh
.github/scripts/bootstrap_test.sh
```

If you want to run individual tests for example:

```sh
export PRIV_KEY="some throwaway key"
./build/bin/geth immutable bootstrap local --override.prevrandao=$(($(date +%s) + 5)) --override.cancun=$(($(date +%s) + 5)) --override.shanghai=$(($(date +%s) + 5))
go test -v ./tests/immutable/... -run TestImmutable_Cancun_4844TransactionsDisabled
```

That will run the local chain with all required forks 5 seconds after genesis. These forks must be enabled after genesis.

You can use `-privkey` instead of an env var if you wish. For example:

```sh
go test -v ./tests/immutable/... -privkey="/path/to/file" -run TestImmutable_Cancun_4844TransactionsDisabled
```
