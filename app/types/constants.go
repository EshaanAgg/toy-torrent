package types

import "github.com/EshaanAgg/toy-bittorrent/app/utils"

const BLOCK_SIZE uint32 = 16384 // 16 KB

const UNCHOKE_MESSAGE_ID = 1
const INTERESTED_MESSAGE_ID = 2
const BITFIELD_MESSAGE_ID = 5
const REQUEST_MESSAGE_ID = 6
const PIECE_MESSAGE_ID = 7

const EXTENSION_HANDSHAKE_HEADER_MESSAGE_ID = 20
const EXTENSION_HANDSHAKE_PAYLOAD_MESSAGE_ID = 0

// We assume that our server always uses the ID 1
// for the "ut_metadata" extension.
const UT_METADATA_EXTENSION_ID = 1

var SERVER_PEER_ID = utils.GetRandomPeerID()

var peerAddressToIDMap = make(map[string]uint32)

// getPeerID returns a unique ID for the peer based on its address.
// This makes the logs from the peers more readable and helps in debugging.
func getPeerID(address string) uint32 {
	id, ok := peerAddressToIDMap[address]
	if !ok {
		id = uint32(len(peerAddressToIDMap)) + 1
		peerAddressToIDMap[address] = id
	}
	return id
}
