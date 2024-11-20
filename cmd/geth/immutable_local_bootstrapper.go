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
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/cmd/geth/immutable/env"
	"github.com/ethereum/go-ethereum/cmd/geth/immutable/node"
	"github.com/ethereum/go-ethereum/cmd/geth/immutable/role"
	"github.com/ethereum/go-ethereum/cmd/geth/immutable/settings"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	ethnode "github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/params"
	"github.com/urfave/cli/v2"
)

// Premine contains addresses and amounts of premined balances
// in the genesis block
type Premine struct {
	Address common.Address
	Wei     *big.Int
}

// ChainOptions contains all the options for generating a genesis
type ChainOptions struct {
	GasLimit        uint64
	SecondsPerBlock uint64

	Validators []common.Address
	Premines   []Premine

	ChainID int
	Dirpath string
}

// Genesis contains the genesis and bootstrap information about it
type Genesis struct {
	// Genesis is the content of the genesis
	Genesis *core.Genesis
	// Filepath is the path to the genesis file
	Filepath string
	// JSON is a convenience field to avoid having to re-marshal the genesis
	JSON []byte
}

// LocalBootstrapper is used to generate genesis files and validator credentials
type LocalBootstrapper struct {
	rootEnvDirpath string
	boots          []node.Node
	validators     []node.Node
	rpcs           []node.Node
}

type bootstrapOptions struct {
	rootDirpath    string
	validatorCount int
	bootCount      int
	rpcCount       int
	gasLimit       uint64

	blockListFilepath string

	// Set these if you want to bootstrap against a remote network only.
	remoteNetwork        string
	remoteConfigFilepath string

	// NOTE: This type should not be used here (compartmentalization bad).
	// We have not refactored geth to improve the bootstrapper in this regard.
	ctx *cli.Context
}

// NewLocalBootstrapper will bootstrap a set of nodes to be run locally.
// It may be used to start a new network or start a node pointing to an existing, remote network.
func NewLocalBootstrapper(opts *bootstrapOptions) (*LocalBootstrapper, error) {
	// Default to devnet configuration
	networkName := "devnet"

	// If a remote network is specified, use that network's configuration
	if opts.remoteNetwork != "" {
		networkName = opts.remoteNetwork
	}
	network, err := settings.NewNetwork(networkName)
	if err != nil {
		return nil, err
	}

	// Set up the paths we will bootstrap files to
	rootEnvDirpath := filepath.Join(opts.rootDirpath, network.String())
	chainSubdir := fmt.Sprint("chain", "-", network.ID())
	chainDirpath := filepath.Join(rootEnvDirpath, chainSubdir)

	// Config TOMLs are either created in this process or derived from an existing file
	configFilepath := opts.remoteConfigFilepath
	if configFilepath == "" {
		configFilepath = filepath.Join(chainDirpath, "config.toml")
	}

	// Clear chaindir
	if err := os.RemoveAll(chainDirpath); err != nil {
		return nil, fmt.Errorf("failed to clear netdir: %w", err)
	}
	// Create dir
	if err := os.MkdirAll(chainDirpath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create datadir: %w", err)
	}

	// Create the nodes
	nodeOpts := make([]node.Options, opts.bootCount+opts.validatorCount+opts.rpcCount)
	for i := range nodeOpts {
		u, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", node.PortP2P+i))
		if err != nil {
			return nil, fmt.Errorf("failed to parse url: %w", err)
		}
		role, ordinal := genesisNodeRoleAndOrdinal(opts.bootCount, opts.validatorCount, opts.rpcCount, i)
		nodeOpts[i] = node.Options{
			Role:         role,
			Ordinal:      ordinal,
			Password:     nodePassword,
			URL:          u,
			ChainDirpath: chainDirpath,
		}
	}
	nodeConf := node.Config{
		Network:           *network,
		ConfigFilepath:    configFilepath,
		BlockListFilepath: opts.blockListFilepath,
	}
	nodes, err := createNodes(nodeOpts, nodeConf)
	if err != nil {
		return nil, err
	}

	// Store the node passwords to files
	for i := range nodes {
		if err := os.WriteFile(localPasswordFilepath(nodes[i]), []byte(nodes[i].Password()), os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create local pw file: %w", err)
		}
	}

	// Create the remaining chain artifacts
	boots, validators, rpcs := splitGenesisNodes(nodes, opts.bootCount, opts.validatorCount)
	var gen *Genesis
	if opts.remoteNetwork == "" {
		bridgeEOAAddress := "0x02F0d131F1f97aef08aEc6E3291B957d9Efe7105" // Also in cmd/geth/testdata/key.addr
		chainOpts := ChainOptions{
			GasLimit:        opts.gasLimit,
			SecondsPerBlock: settings.SecondsPerBlock,
			Validators:      nodesToAddresses(validators),
			Premines:        immutablePremines(env.Devnet, common.HexToAddress(bridgeEOAAddress)),
			Dirpath:         chainDirpath,
			ChainID:         network.ID(),
		}
		// Generate a clique genesis
		var err error
		gen, err = Clique(chainOpts)
		if err != nil {
			return nil, err
		}
	} else {
		// Derive a genesis from the remote network
		network, err := settings.NewNetwork(opts.remoteNetwork)
		if err != nil {
			return nil, err
		}
		gen = &Genesis{
			Filepath: fmt.Sprintf("./cmd/geth/immutable/genesis/%s.json", network.String()),
			Genesis:  core.ImmutableGenesisBlock(network.String()),
			JSON:     []byte(network.GenesisJSON()),
		}
	}

	// Make sure we run an immutable network
	if !gen.Genesis.Config.IsValidImmutableZKEVM() {
		return nil, fmt.Errorf("invalid genesis config: %+v", *gen.Genesis.Config)
	}

	// Init the geth nodes with genesis
	if err := renderLocalChainState(opts.ctx, gen, nodes); err != nil {
		return nil, err
	}

	// Render the config if we are not connecting to an existing network
	if opts.remoteNetwork == "" {
		if err := renderLocalConfig(boots, configFilepath, network.ID()); err != nil {
			return nil, err
		}
	}

	return &LocalBootstrapper{
		rootEnvDirpath: rootEnvDirpath,
		boots:          boots,
		validators:     validators,
		rpcs:           rpcs,
	}, nil
}

func (b *LocalBootstrapper) Clean() {
	os.RemoveAll(b.rootEnvDirpath)
}

// createNodes will generate EOA credentials and store them in encrypted keystore files
func createNodes(opts []node.Options, conf node.Config) ([]node.Node, error) {
	// Generate validator accounts
	nodeCount := len(opts)
	nodes := make([]node.Node, 0, nodeCount)
	for i := 0; i < nodeCount; i++ {
		// Construct the node
		node, err := node.New(opts[i], conf)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, *node)
	}

	return nodes, nil
}

// Clique will generate a genesis for a clique chain and store it to file
func Clique(opts ChainOptions) (*Genesis, error) {
	gen, err := clique(opts)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(gen.Filepath, gen.JSON, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to write genesis file: %w", err)
	}
	return gen, nil
}

func clique(opts ChainOptions) (*Genesis, error) {
	// To encode the signer addresses in extradata,
	// concatenate 32 zero bytes, all signer addresses and 65 further zero bytes.
	// This is based on the clique spec: https://eips.ethereum.org/EIPS/eip-225
	extraDataLength := crypto.DigestLength + common.AddressLength*len(opts.Validators) + crypto.SignatureLength
	extraData := make([]byte, extraDataLength)
	for i := range opts.Validators {
		from := crypto.DigestLength + common.AddressLength*i
		to := from + common.AddressLength
		copy(extraData[from:to], opts.Validators[i].Bytes())
	}

	// Premine
	alloc := types.GenesisAlloc{}
	for i := range opts.Premines {
		alloc[opts.Premines[i].Address] = types.Account{
			Balance: opts.Premines[i].Wei,
		}
	}

	// Marshal the genesis
	gen := core.Genesis{
		Config: &params.ChainConfig{
			ChainID:             new(big.Int).SetUint64(uint64(opts.ChainID)),
			HomesteadBlock:      common.Big0,
			EIP150Block:         common.Big0,
			EIP155Block:         common.Big0,
			EIP158Block:         common.Big0,
			ByzantiumBlock:      common.Big0,
			ConstantinopleBlock: common.Big0,
			PetersburgBlock:     common.Big0,
			IstanbulBlock:       common.Big0,
			MuirGlacierBlock:    common.Big0,
			BerlinBlock:         common.Big0,
			LondonBlock:         common.Big0,
			ArrowGlacierBlock:   common.Big0,
			GrayGlacierBlock:    common.Big0,
			MergeNetsplitBlock:  common.Big0,
			Clique: &params.CliqueConfig{
				Period: opts.SecondsPerBlock,
				Epoch:  30000,
			},
			IsReorgBlocked: true,
			// So as to reflect devnet, testnet, and mainnet, these forks should not be enabled in genesis.
			ShanghaiTime:   nil,
			PrevrandaoTime: nil,
			CancunTime:     nil, // If genesis block has Cancun enabled, you must set the blob-related headers too.
		},

		Difficulty: big.NewInt(1),
		GasLimit:   opts.GasLimit,
		Mixhash:    common.Hash{},
		ExtraData:  extraData,
		Alloc:      alloc,
	}
	json, err := json.MarshalIndent(gen, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal genesis: %w", err)
	}

	// Write genesis file
	path := filepath.Join(opts.Dirpath, "genesis.json")
	return &Genesis{
		Genesis:  &gen,
		Filepath: path,
		JSON:     json,
	}, nil
}

// renderLocalChainState will initialize the chaindata and lightchaindata dbs
func renderLocalChainState(c *cli.Context, gen *Genesis, nodes []node.Node) error {
	// Init node dbs and write genesis to dbs
	for i := range nodes {
		// Instantiate node based on config and dir
		cfg := loadBaseConfig(c)
		cfg.Node.DataDir = nodes[i].Dirpath()
		stack, err := ethnode.New(&cfg.Node)
		if err != nil {
			return fmt.Errorf("failed to create node: %w", err)
		}
		defer stack.Close()

		// Init genesis for chain and lightchain dbs of node
		for _, name := range []string{"chaindata", "lightchaindata"} {
			// Create and open leveldb
			chaindb, err := stack.OpenDatabaseWithFreezer(name, 0, 0, "", "", false)
			if err != nil {
				return fmt.Errorf("failed to open db: %w", err)
			}
			defer chaindb.Close()

			// Create and open trie db
			triedb := utils.MakeTrieDatabaseNoContext(cfg.Eth.StateScheme, chaindb, false, false, false)
			defer triedb.Close()

			// Write genesis block to trie and save to leveldb
			if _, _, err := core.SetupGenesisBlock(chaindb, triedb, gen.Genesis); err != nil {
				return fmt.Errorf("failed to write genesis block to db: %w", err)
			}
		}
	}

	return nil
}

// renderLocalConfig will write the config file for each node.
// This needs to be run after nodekeys are generated by renderLocalChainState.
func renderLocalConfig(boots []node.Node, configFilepath string, chainID int) (err error) {
	enodes := make([]*enode.Node, len(boots))
	for i := range enodes {
		en, err := enodeFromNodeKey(boots[i])
		if err != nil {
			return err
		}
		enodes[i] = en
	}

	f, err := os.OpenFile(configFilepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()
	conf := &gethConfig{
		Eth:  ImmutableEthConfig(chainID),
		Node: nodeConfig(enodes),
	}

	if err := tomlSettings.NewEncoder(f).Encode(conf); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}
	return nil
}

// enodeFromNodeKey reads the nodekey file and parses the enode.
// The nodekey files are available after chain state has been rendered.
func enodeFromNodeKey(node node.Node) (*enode.Node, error) {
	// Derive the public key
	nodeKeyFilepath := filepath.Join(node.Dirpath(), "geth", "nodekey")
	nodeKey, err := crypto.LoadECDSA(nodeKeyFilepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read nodekey file: %w", err)
	}

	// Parse the enode
	enodeURL := fmt.Sprintf("enode://%x@%s:%s", crypto.FromECDSAPub(&nodeKey.PublicKey)[1:], node.URL().Hostname(), node.URL().Port())
	enode, err := enode.Parse(enode.ValidSchemes, enodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse enode: %w", err)
	}
	return enode, nil
}

func nodesToAddresses(nodes []node.Node) []common.Address {
	addrs := make([]common.Address, len(nodes))
	for i := range nodes {
		addrs[i] = nodes[i].Account().Address
	}
	return addrs
}

func genesisNodeRoleAndOrdinal(bootCount, validatorCount, rpcCount, index int) (role.Role, int) {
	// Boots
	if index < bootCount {
		return role.Boot, index
	}
	// Validators
	if index < bootCount+validatorCount {
		return role.Validator, index - bootCount
	}
	// RPCs
	if index < bootCount+validatorCount+rpcCount {
		return role.RPC, index - bootCount - validatorCount
	}
	// Partner
	return role.Partner, index - bootCount - validatorCount - rpcCount
}

func splitGenesisNodes(nodes []node.Node, bootCount, validatorCount int) (boots, validators, rpcs []node.Node) {
	return nodes[:bootCount], nodes[bootCount : bootCount+validatorCount], nodes[bootCount+validatorCount:]
}

func localPasswordFilepath(node node.Node) string {
	return filepath.Join(node.Dirpath(), node.Password())
}
