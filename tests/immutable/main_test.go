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

package test

import (
	"flag"
	"log"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/settings"
)

// flags constitutes flag inputs
var flags struct {
	PrivKeyFilepath        string
	BlockedPrivKeyFilepath string
	RPCURLStr              string
	ValidatorAdminURLStr   string
	SkipVoting             bool
	NetworkForks           string
}

// config is rendered from env and flags
var config struct {
	rpcURL       *url.URL
	validatorURL *url.URL
	testUser     *Wallet
	blockedUser  *Wallet
	skipVoting   bool
	cancun       *settings.Fork
}

func skipCancun(t *testing.T) bool {
	t.Helper()
	if config.cancun != nil && !config.cancun.IsEnabledAt(time.Now().Unix()) {
		t.Logf("Skipping test because now (%v) is before the Cancun fork (%v)", time.Now().UTC(), config.cancun.Time.UTC())
		return true
	}
	return false
}

func TestMain(m *testing.M) {
	// Read flags
	flag.StringVar(&flags.RPCURLStr, "rpc", "http://localhost:8545", "RPC URL of the node")
	flag.StringVar(&flags.ValidatorAdminURLStr, "validatoradmin", "http://localhost:8550", "RPC URL of the validator admin endpoint")
	flag.StringVar(&flags.PrivKeyFilepath, "privkey", "", "Filepath to private key to use for sending txs")
	flag.StringVar(&flags.BlockedPrivKeyFilepath, "blockedprivkey", "", "Filepath to bloccked private key to use for simulating blocked user")
	flag.BoolVar(&flags.SkipVoting, "skipvoting", true, "If true, voting tests will be skipped")
	flag.StringVar(&flags.NetworkForks, "forks", "", "If (devnet, testnet, mainnet) is set, the test will skip forks before the timestamp of the network")
	flag.Parse()

	// Load the wallet
	var user *Wallet
	var err error
	if flags.PrivKeyFilepath != "" {
		user, err = loadWalletFromFile(flags.PrivKeyFilepath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		privKeyHex := os.Getenv("PRIV_KEY")
		if privKeyHex == "" {
			log.Fatal("-privkey flag and PRIV_KEY env var not set")
		}
		user, err = loadWallet(privKeyHex)
		if err != nil {
			log.Fatal(err)
		}
	}

	config.testUser = user

	if flags.BlockedPrivKeyFilepath != "" {
		blockedUser, err := loadWalletFromFile(flags.BlockedPrivKeyFilepath)
		if err != nil {
			log.Fatal(err)
		}
		config.blockedUser = blockedUser
	}

	// Render the rpc url
	url, err := url.Parse(flags.RPCURLStr)
	if err != nil {
		log.Fatal("failed to parse RPC URL: ", err)
	}
	config.rpcURL = url

	// Render the validator admin url
	valurl, err := url.Parse(flags.ValidatorAdminURLStr)
	if err != nil {
		log.Fatal("failed to parse Validator Admin URL: ", err)
	}
	config.validatorURL = valurl
	config.skipVoting = flags.SkipVoting
	if flags.NetworkForks != "" {
		network, err := settings.NewNetwork(flags.NetworkForks)
		if err != nil {
			log.Fatal(err)
		}
		cancun := network.Cancun()
		config.cancun = &cancun
	}
	// Run the tests
	os.Exit(m.Run())
}
