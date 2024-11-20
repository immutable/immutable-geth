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
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
)

func TestImmutable_RPCDisableFlag(t *testing.T) {
	// Flags are set under cmd/geth/immutable_local.go
	rpcMethods := []string{
		// Debug
		"debug_startGoTrace",
		"debug_stopGoTrace",
		"debug_blockProfile",
		"debug_setBlockProfileRate",
		"debug_writeBlockProfile",
		"debug_mutexProfile",
		"debug_setMutexProfileFraction",
		"debug_writeMutexProfile",
		"debug_writeMemProfile",
		"debug_traceBlock",
		"debug_traceBlockFromFile",
		"debug_traceBadBlock",
		"debug_standardTraceBadBlockToFile",
		"debug_intermediateRoots",
		"debug_standardTraceBlockToFile",
		"debug_traceBlockByNumber",
		"debug_traceBlockByHash",
		"debug_traceTransaction",
		"debug_traceCall",
		"debug_preimage",
		"debug_getBadBlocks",
		"debug_storageRangeAt",
		"debug_getModifiedAccountsByNumber",
		"debug_getModifiedAccountsByHash",
		"debug_freezeClient",
		"debug_getAccessibleState",
		"debug_dbGet",
		"debug_dbAncient",
		"debug_dbAncients",
		"debug_setTrieFlushInterval",
		"debug_getTrieFlushInterval",
		"debug_accountRange",
		"debug_printBlock",
		"debug_getRawHeader",
		"debug_getRawBlock",
		"debug_getRawReceipts",
		"debug_getRawTransaction",
		"debug_setHead",
		"debug_seedHash",
		"debug_dumpBlock",
		"debug_chaindbProperty",
		"debug_chaindbCompact",
		"debug_verbosity",
		"debug_vmodule",
		"debug_backtraceAt",
		"debug_stacks",
		"debug_freeOSMemory",
		"debug_setGCPercent",
		"debug_memStats",
		"debug_gcStats",
		"debug_cpuProfile",
		"debug_startCPUProfile",
		"debug_stopCPUProfile",
		"debug_goTrace",
		// Tx pool
		"txpool_content",
		"txpool_inspect",
		"txpool_status",
		"txpool_contentFrom",
		// Clique
		"clique_getSnapshot",
		"clique_getSnapshotAtHash",
		"clique_getSigner",
		"clique_getSigners",
		"clique_getSignersAtHash",
		"clique_proposals",
		"clique_propose",
		"clique_discard",
		"clique_status",
		// Miner
		"miner_getHashrate",
		"miner_setExtra",
		"miner_setGasPrice",
		"miner_setRecommitInterval",
		"miner_start",
		"miner_stop",
		"miner_setEtherbase",
		"miner_setGasLimit",
		// Personal
		"personal_listAccounts",
		"personal_deriveAccount",
		"personal_ecRecover",
		"personal_importRawKey",
		"personal_listWallets",
		"personal_newAccount",
		"personal_openWallet",
		"personal_sendTransaction",
		"personal_sign",
		"personal_signTransaction",
		"personal_unlockAccount",
		"personal_lockAccount",
		"personal_unpair",
		"personal_initializeWallet",
		"personal_initializeWallets",
	}

	client, err := rpc.DialContext(context.Background(), config.rpcURL.String())
	require.NoError(t, err)

	for _, rpcMethod := range rpcMethods {
		t.Run(fmt.Sprintf("RPCDisabled_%s", rpcMethod), func(t *testing.T) {
			var res string
			err = client.Call(&res, rpcMethod)
			require.Error(t, err)
			//
			require.Equal(t, fmt.Sprintf("the method %s does not exist/is not available", rpcMethod), err.Error())
		})
	}
}

func TestImmutable_RPCDisableFlagsOmitted(t *testing.T) {
	// Some of these flags are disabled in remote environments but not for local tests but
	// it is to check locally the counter when the disabled flags are omitted
	isLocalhost := strings.Contains(config.rpcURL.String(), "localhost") || strings.Contains(config.rpcURL.String(), "127.0.0.1")
	if !isLocalhost {
		t.Skip("Skipping this test because this test is only applicable to localhost setup")
	}

	fileName := "1"
	// This tests outputs material, defer func to clean it up to avoid duplicate write errors
	defer func() {
		// Tests will be run from tests/immutable, file will be saved at root
		relPath := path.Join("../../", fileName)

		_, err := os.Stat(relPath)
		if err == nil {
			// Rm the test file
			err := os.Remove(relPath)
			if err != nil {
				t.Fatal(err)
			}
		}
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				t.Fatalf("Test file %s does not exist, expected test to write test file", fileName)
			}

			t.Fatal(err)
		}
	}()

	rpcMethods := map[string][]interface{}{
		"admin_exportChain": {fileName},
		"admin_nodeInfo":    {}, // No arguments required
	}

	client, err := rpc.DialContext(context.Background(), config.rpcURL.String())
	require.NoError(t, err)

	for rpcMethod, args := range rpcMethods {
		t.Run(fmt.Sprintf("RPCEnabled_%s", rpcMethod), func(t *testing.T) {
			var res interface{}
			if len(args) > 0 {
				err = client.Call(&res, rpcMethod, args...)
			} else {
				err = client.Call(&res, rpcMethod)
			}
			require.NoError(t, err)
		})
	}
}
