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

package main

import (
	"math/big"

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/env"
	"github.com/ethereum/go-ethereum/common"
)

var (
	// totalSupplyWei is the amount of IMX that is pre-funded to the bridge EOA
	totalSupplyWei = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2e9))

	// immutableDevPremines is a list of EOAs that are pre-funded with a large amount of ETH
	// on dev networks. This is used to allow developers to deploy contracts without having
	// to request funds from a faucet.
	immutableDevPremineAddresses = []string{
		"0x340bC2c77514ede2a23Fd4F42F411A8e351d8eE6",
		"0xebbf4C07a63986204C37cc5A188AaBF53564C583",
		"0xdEAdC0de8a3B037925a895843f96b0c525FBC31f",
		"0xeFE12952541356Ffc969A343A81D1cE7D2806179",
		"0x4AEdf28A437b94749037cC39f83F4422469CF2F7",
		"0x8318a871CC140d9f77a1999f84875AC36EeCC04E",
		"0xCc5C8CEa877f2F351F38c190867BbD31FaFadD22",
		"0x000000000013B7b1B08B3c8EFE02E866F746bD38",
		"0xa6C368164Eb270C31592c1830Ed25c2bf5D34BAE",
		"0xC606830D8341bc9F5F5Dd7615E9313d2655B505D",
		"0x784578949A4A50DeA641Fb15dd2B11C72E76919a",
		"0xEac347177DbA4a190B632C7d9b8da2AbfF57c772",
		"0xD509997AB62fDA51c32E64E69Fb090DF8894105e",
		"0xF6372939CE2d14A68A629B8E4785E9dCB4EdA0cf",
		"0x9C1634bebC88653D2Aebf4c14a3031f62092b1D9",
		"0xb3343666188A694120C18c4985C57e4C0913A6F0",
		"0x2E969d22e6654e064F461cf8B1314Cc0864a4914",
		"0xd9275Eb8276E14b9e28d5f9B12e90dDAAF3586Ef",
		"0x3e290FE8F2A5dB60A81cb47EA296e0299048Dd71",
		"0x4A73506a31DB769AC442b17ca9A1679f44757Bbf",
		"0x7C3E6CE6fd293Fc66d9d73d49fd546CCE1e19F0e",
		"0xEB7FFb9fb0c80437120f6F97EdE60aB59055EAE0",
		"0xe567Ea84e1eB3fFdc8F5aA420BF14A16eeE6A809",
		"0xC8714F989cE817e5d21349888077Aa5Db4A9BCf6",
		"0x0CCB0a3fc5Ca38fcd9FfD8a667Cb83e3194250d7",
		"0x5ABFc3E307b037325BFC6988Ae265dcB211Ec533",
		"0x7442eD1e3c9FD421F47d12A2742AfF5DaFBf43f8",
		"0x4A73506a31DB769AC442b17ca9A1679f44757Bbf",
		"0x7C3E6CE6fd293Fc66d9d73d49fd546CCE1e19F0e",
		"0xEB7FFb9fb0c80437120f6F97EdE60aB59055EAE0",
		"0xed557863FFD4C87537BA8264098B22483c6145f2",
		"0x7924BF4cBb25f7bA2aB1335e293afe6a7E78235a",
	}
)

func immutablePremines(envr env.Environment, bridgeEOA common.Address) []Premine {
	// Always premine bridge EOA
	premines := []Premine{
		{
			Address: bridgeEOA,
			Wei:     totalSupplyWei,
		},
	}
	// Only premine devnet
	if envr != env.Devnet {
		return premines
	}
	for _, address := range immutableDevPremineAddresses {
		premines = append(premines, Premine{
			Address: common.HexToAddress(address),
			Wei:     totalSupplyWei,
		})
	}
	return premines
}
