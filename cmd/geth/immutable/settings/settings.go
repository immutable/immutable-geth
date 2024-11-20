// Copyright 2024 The Immutable go-ethereum Authors
// This file is part of the Immutable go-ethereum library.
//
// The Immutable go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Immutable go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Immutable go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package settings

const (
	// PriceLimit is the minimum gas price that RPC nodes will accept for transactions.
	PriceLimit = 10 * 1e9
	// SecondsPerBlock is the amount of time between blocks.
	SecondsPerBlock = 2
	// BaseFeeChangeDenominator is set to a value which leads to a smaller
	// max base fee rate of change (12.5% -> 2%). This accounts for our short block time
	// of 2 seconds, preventing large price fluctuations. At 50, we match Ethereum in that
	// it would take 72 seconds for the baseFee to double.
	BaseFeeChangeDenominator = 50

	// MainnetNetworkID is the network ID for the mainnet.
	MainnetNetworkID = 13371
	// TestnetNetworkID is the network ID for the testnet.
	TestnetNetworkID = 13473
	// DevnetNetworkID is the network ID for the devnet.
	DevnetNetworkID = 15003

	// MainnetRPC is the mainnet RPC endpoint.
	MainnetRPC = "https://rpc.immutable.com"
	// TestnetRPC is the testnet RPC endpoint.
	TestnetRPC = "https://rpc.testnet.immutable.com"
	// DevnetRPC is the devnet RPC endpoint.
	DevnetRPC = "https://rpc.dev.immutable.com"
)

var (
	// DevnetShanghaiFork is the timestamp of the Shanghai devnet fork.
	DevnetShanghaiFork = newFork("Tue Feb 27 21:00:00 UTC 2024")
	// TestnetShanghaiFork is the timestamp of the Shanghai testnet fork.
	TestnetShanghaiFork = newFork("Tue Mar 12 22:00:00 UTC 2024")
	// MainnetShanghaiFork is the timestamp of the Shanghai mainnet fork.
	MainnetShanghaiFork = newFork("Tue Mar 26 22:00:00 UTC 2024")

	// DevnetPrevrandaoFork is the timestamp of the Prevrandao devnet fork.
	DevnetPrevrandaoFork = DevnetShanghaiFork
	// TestnetPrevrandaoFork is the timestamp of the Prevrandao testnet fork.
	TestnetPrevrandaoFork = TestnetShanghaiFork
	// MainnetPrevrandaoFork is the timestamp of the Prevrandao mainnet fork.
	// Only mainnet has a Prevrandao fork separate from the Shanghai fork.
	// Block 0x41DDE4.
	// TZ=UTC gdate --date @1710899402
	MainnetPrevrandaoFork = newFork("Wed Mar 20 01:50:02 UTC 2024")

	// DevnetCancunFork is the timestamp of the Cancun devnet fork.
	DevnetCancunFork = newFork("Tue Aug 27 22:00:00 UTC 2024")
	// TestnetCancunFork is the timestamp of the Cancun testnet fork.
	TestnetCancunFork = newFork("Mon Sep 23 22:00:00 UTC 2024")
	// MainnetCancunFork is the timestamp of the Cancun mainnet fork.
	// Unixdate format expects an extra whitespace if the day is a single digit.
	// e.g. Mon Jan _2 15:04:05 MST 2006
	MainnetCancunFork = newFork("Mon Oct  7 22:00:00 UTC 2024")
)
