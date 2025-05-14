package types

const BLOCK_SIZE uint32 = 16384 // 16 KB

const UNCHOKE_MESSAGE_ID = 1
const INTERESTED_MESSAGE_ID = 2
const BITFIELD_MESSAGE_ID = 5
const REQUEST_MESSAGE_ID = 6
const PIECE_MESSAGE_ID = 7

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
