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
	"crypto/ecdsa"
	"fmt"
	"net"
	"os"

	"github.com/ethereum/go-ethereum/accounts/immutable"
	"github.com/ethereum/go-ethereum/cmd/geth/immutable/keys"
	"github.com/ethereum/go-ethereum/cmd/geth/immutable/node"
	"github.com/ethereum/go-ethereum/cmd/geth/immutable/role"
	"github.com/ethereum/go-ethereum/cmd/immutable/remote/aws"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/urfave/cli/v2"
)

func runBootNodeCommand(c *cli.Context) error {
	// Retrieve the boot node's p2p key
	var p2pKey *ecdsa.PrivateKey

	pwFilepath := os.Getenv("GETH_FLAG_PASSWORD_FILEPATH")
	if pwFilepath != "" {
		// Retrieve the key from local keystore (for local development only)
		log.Info("retrieving P2P key from local keystore")
		keyStoreDirpath := c.String(utils.KeyStoreDirFlag.Name)
		store, err := immutable.NewKeystore(keyStoreDirpath, pwFilepath)
		if err != nil {
			return err
		}
		p2pKey, err = store.GetPrivateKey(c.Context)
		if err != nil {
			return err
		}
	} else {
		// Retrieve the key from AWS
		log.Info("retrieving P2P key from Secrets Manager")
		awsRegion := os.Getenv("GETH_FLAG_IMMUTABLE_AWS_REGION")
		if awsRegion == "" {
			return fmt.Errorf("GETH_FLAG_IMMUTABLE_AWS_REGION is required")
		}
		store, err := aws.NewSecretsManager(awsRegion)
		if err != nil {
			return fmt.Errorf("failed to create AWS SecretsManager Store: %v ", err)
		}
		podOrdinal, err := ordinalFromPodNameEnvVar()
		if err != nil {
			return err
		}
		secretIDTemplate, err := keys.SecretIDTemplate()
		if err != nil {
			return err
		}
		secretID, err := keys.SecretID(secretIDTemplate, node.CanonicalNodeName(role.Boot, podOrdinal))
		if err != nil {
			return err
		}
		p2pKey, err = keys.Retrieve(
			c.Context,
			store,
			secretID,
		)
		if err != nil {
			return fmt.Errorf("failed to retrieve P2P key: %v", err)
		}
	}

	// Run the boot node
	return runBootNode(p2pKey, c.Int(utils.ListenPortFlag.Name))
}

func runBootNode(p2pKey *ecdsa.PrivateKey, port int) error {
	// Log
	glogger := log.NewGlogHandler(log.NewTerminalHandler(os.Stderr, false))
	glogger.Verbosity(4)
	glogger.Vmodule("")
	log.SetDefault(log.NewLogger(glogger))

	// UDP session
	listenAddr := fmt.Sprintf(":%d", port)
	addr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %w", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen UDP: %w", err)
	}
	defer conn.Close()

	// Enode DB
	db, _ := enode.OpenDB("")
	ln := enode.NewLocalNode(db, p2pKey)

	// Notice
	listenerAddr := conn.LocalAddr().(*net.UDPAddr)
	printNotice(&p2pKey.PublicKey, *listenerAddr)

	// Run
	cfg := discover.Config{
		PrivateKey:  p2pKey,
		NetRestrict: nil,
	}
	if _, err := discover.ListenUDP(conn, ln, cfg); err != nil {
		return fmt.Errorf("failed to listen UDP: %w", err)
	}

	// Block
	select {}
}

func printNotice(nodeKey *ecdsa.PublicKey, addr net.UDPAddr) {
	if addr.IP.IsUnspecified() {
		addr.IP = net.IP{127, 0, 0, 1}
	}
	n := enode.NewV4(nodeKey, addr.IP, 0, addr.Port)
	fmt.Println(n.URLv4())
}
