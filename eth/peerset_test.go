package eth

import (
	"testing"

	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/p2p"
)

func TestAllPeers_EmptyPeerSet_ZeroPeers(t *testing.T) {
	t.Parallel()
	testPeerSet := newPeerSet()
	peers := testPeerSet.allPeers()
	if len(peers) != 0 {
		t.Fatalf("expected 0 peers but got %d", len(peers))
	}
}

func TestAllPeers_NonEmptyPeerset_NonZeroPeers(t *testing.T) {
	t.Parallel()
	testPeerSet := newPeerSet()

	// create peers
	expectedNumPeers := 10
	for i := 0; i < expectedNumPeers; i++ {
		p2pPeer := p2p.NewPeer([32]byte{byte(i)}, "", []p2p.Cap{})
		peer := eth.NewPeer(1, p2pPeer, nil, nil)
		err := testPeerSet.registerPeer(peer, nil)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	peers := testPeerSet.allPeers()
	if len(peers) != expectedNumPeers {
		t.Fatalf("expected %d peers but got %d", expectedNumPeers, len(peers))
	}
}
