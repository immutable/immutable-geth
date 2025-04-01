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

package node

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/geth/immutable/role"
	"github.com/ethereum/go-ethereum/cmd/geth/immutable/settings"
)

const (
	PortP2P     = 30300
	PortRPC     = 8545
	PortWS      = 8535
	PortMetrics = 6060
	PortAuthRPC = 8550
	PortPprof   = 7070
)

// Options are node-specific values
type Options struct {
	Role         role.Role
	Ordinal      int
	Password     string
	URL          *url.URL
	ChainDirpath string
}

// Config contains the configuration relevant for all nodes
type Config struct {
	Network           settings.Network
	ConfigFilepath    string
	BlockListFilepath string
}

// Node is a runnable geth node that is produced by the LocalBootstrapper.
// It depends on files rendered by the bootstrap process.
type Node struct {
	account  accounts.Account
	dirpath  string
	password string
	u        *url.URL

	args    []string
	environ []string
}

// New creates a new node
func New(opts Options, conf Config) (*Node, error) {
	// Create indexed node directory
	nodeDirpath := filepath.Join(opts.ChainDirpath, CanonicalNodeName(opts.Role, opts.Ordinal))
	if err := os.MkdirAll(nodeDirpath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create node directory: %w", err)
	}

	// Store the encrypted key file to disk
	account, err := keystore.StoreKey(
		keystoreDirpath(nodeDirpath),
		opts.Password,
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create key: %w", err)
	}

	// Store the account address to disk
	if err := os.WriteFile(filepath.Join(nodeDirpath, "address"), []byte(account.Address.Hex()), 0644); err != nil {
		return nil, fmt.Errorf("failed to write address to file: %w", err)
	}

	passwordFilepath := filepath.Join(nodeDirpath, opts.Password)
	environ := []string{
		"GETH_FLAG_PASSWORD_FILEPATH",
		passwordFilepath,
	}
	args := []string{}
	switch opts.Role {
	case role.Boot:
		args = []string{
			"./build/bin/geth",
			"immutable",
			"run",
			"boot",
			"--keystore", filepath.Join(nodeDirpath, "keystore"),
			"--port", opts.URL.Port(),
		}
	case role.Validator:
		args = []string{
			"./build/bin/geth",
			"--datadir", nodeDirpath,
			"--log.debug",
			"--networkid", fmt.Sprint(conf.Network.ID()),
			"--metrics",
			"--metrics.addr", "127.0.0.1",
			"--metrics.port", fmt.Sprint(PortMetrics + opts.Ordinal),
			"--authrpc.port", fmt.Sprint(PortAuthRPC + opts.Ordinal),
			"--verbosity", "4",
			"--port", opts.URL.Port(),
			"--password", passwordFilepath,
			"--rpc.debugdisable",
			"--rpc.txpooldisable",
			"--config", conf.ConfigFilepath,
			"--pprof",
			"--pprof.port", fmt.Sprint(PortPprof + opts.Ordinal),
			"--miner.etherbase", account.Address.Hex(),
			"--mine",
			// NOTE: Add these flags to enable RPC on validator
			"--http",
			"--http.port", fmt.Sprint(PortRPC + opts.Ordinal),
			"--cache", "128",
			"--cache.database", "35",
			"--cache.trie", "35",
			"--cache.gc", "10",
			"--cache.snapshot", "20",
		}
	case role.RPC:
		// TODO: This assumes 1 validator, better to parameterize
		rpcPortOffset := opts.Ordinal + 1 // Avoid port collisions on same host
		args = []string{
			"./build/bin/geth",
			"--datadir", nodeDirpath,
			"--log.debug",
			"--networkid", fmt.Sprint(conf.Network.ID()),
			"--metrics",
			"--metrics.addr", "127.0.0.1",
			"--metrics.port", fmt.Sprint(PortMetrics + rpcPortOffset),
			"--authrpc.port", fmt.Sprint(PortAuthRPC + rpcPortOffset),
			"--http",
			"--http.port", fmt.Sprint(PortRPC + rpcPortOffset),
			"--ws.port", fmt.Sprint(PortWS + rpcPortOffset),
			"--txpool.blocklistfilepaths", conf.BlockListFilepath,
			"--rpc.debugdisable",
			"--rpc.txpooldisable",
			"--rpc.cliquedisable",
			"--rpc.minerdisable",
			"--rpc.personaldisable",
			"--verbosity", "4",
			"--port", opts.URL.Port(),
			"--config", conf.ConfigFilepath,
			"--pprof",
			"--pprof.port", fmt.Sprint(PortPprof + rpcPortOffset),
			"--cache", "128",
			"--cache.database", "40",
			"--cache.trie", "40",
			"--cache.gc", "0",
			"--cache.snapshot", "20",
		}
	}

	return &Node{
		account:  account,
		dirpath:  nodeDirpath,
		password: opts.Password,
		u:        opts.URL,
		args:     args,
		environ:  environ,
	}, nil
}

// URL returns the enode URL of the node
func (n Node) URL() *url.URL {
	return n.u
}

// Dirpath returns the directory path of the node's files
func (n Node) Dirpath() string {
	return n.dirpath
}

// Account returns the account of the node
func (n Node) Account() accounts.Account {
	return n.account
}

// Password returns the password of the node
func (n Node) Password() string {
	return n.password
}

// Run runs the node as a sub-process
func (n Node) Run(optionalArgs []string) error {
	args := append(n.args, optionalArgs...)
	cmd, err := prepareCmd(n.environ, args...)
	if err != nil {
		return err
	}
	return cmd.Run()
}

func keystoreDirpath(baseDirpath string) string {
	return filepath.Join(baseDirpath, "keystore")
}

// CanonicalNodeName returns the canonical name of a node based on role and ordinal
func CanonicalNodeName(role role.Role, ordinal int) string {
	return role.String() + "-" + fmt.Sprint(ordinal)
}

func prepareCmd(environ []string, args ...string) (*exec.Cmd, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("not enough args")
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	environLog := ""
	if len(environ)%2 != 0 {
		return nil, fmt.Errorf("environ must be key-value pairs")
	}
	if len(environ) > 0 {
		environLog = "export"
	}
	cmd.Env = os.Environ()
	for i := 0; i < len(environ); i += 2 {
		environStr := fmt.Sprintf("%s=%s", environ[i], environ[i+1])
		cmd.Env = append(cmd.Env, environStr)
		environLog = fmt.Sprintf("%s %s", environLog, environStr)
	}
	fmt.Printf("<<<<<<<<<<<<\n%s; %s\n>>>>>>>>>>>>>\n", environLog, cmd.String())
	return cmd, nil
}
