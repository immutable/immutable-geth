## Immutable

This is a fork of go-ethereum that is specific to Immutables deployment and infrastructure.
Several commands have been added to the geth-cli to support Immutable's operations. All subcommands
will be accessible under the `immutable` parent command.

## Run a local network

Build the binary:
```sh
make geth
```

To bootstrap a minimal chain locally, run:
```sh
./build/bin/geth immutable bootstrap local
```

this will set up a network with 1 validator (http://localhost:8545) and 1 RPC (http://localhost:8546) node.

Use --help to learn about the other ways of configuring your local network:
```sh
./build/bin/geth immutable bootstrap --help
```

## Test against a running network

See `.github/scripts/bootstrap_test.sh` on how to run go tests against a running local network.
