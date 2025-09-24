# Immutable Geth

Golang execution layer implementation of the Ethereum protocol. Modified for the purposes of the Immutable zkEVM.

All modifications made by Immutable are either contained in files named `*immutable*.go` or have `CHANGE(immutable):` comments above or inside the modified lines of code.

## Build

With Golang 1.20 installed, run:

```
make geth
```

or, to build the full suite of utilities:

```
make all
```

The built client binary is `./build/bin/geth`

You can run a local network with your built binary via the `immutable bootstrap local` command:
```
now=$(date +%s)
./build/bin/geth immutable bootstrap local \
--override.shanghai="$now" \
--override.prevrandao="$now" \
--override.cancun="$now"
```

You can run the E2E tests against your built binary via:
```
.github/scripts/bootstrap_test.sh
```

## Docker

The client is distributed as the following Docker image:
```
docker pull ghcr.io/immutable/go-ethereum/go-ethereum:latest
```

## Run

If you wish to join the P2P network, please follow the instructions [here](https://docs.immutable.com/learn/platform/nodes/).

### Hardware Requirements

Minimum:

* CPU with 2+ cores
* 4GB RAM
* 1TB free storage space to sync the Mainnet
* 8 MBit/sec download Internet service

Recommended:

* Fast CPU with 4+ cores
* 16GB+ RAM
* High-performance SSD with at least 1TB of free space
* 25+ MBit/sec download Internet service

## Contribution

We welcome any contributions and aim to respond promptly to issues and pull requests.

Please make sure your contributions adhere to our coding guidelines:

 * Code must adhere to the official Go [formatting](https://golang.org/doc/effective_go.html#formatting)
   guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
 * Code must be documented adhering to the official Go [commentary](https://golang.org/doc/effective_go.html#commentary)
   guidelines.
 * Pull requests need to be based on and opened against the `main` branch.

## License

The Immutable go-ethereum library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html),
also included in our repository in the `COPYING.LESSER` file.

The Immutable go-ethereum binaries (i.e. all code inside of the `cmd` directory) are licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also
included in our repository in the `COPYING` file.
