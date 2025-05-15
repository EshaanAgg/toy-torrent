package types

import (
	"encoding/binary"
	"fmt"

	"github.com/EshaanAgg/toy-bittorrent/app/bencode"
)

type Handshake struct {
	SupportsExtensions bool
	PeerID             []byte
	InfoHash           []byte
}

func (h *Handshake) Bytes() []byte {
	var b []byte

	b = append(b, 19) // Protocol length
	b = append(b, "BitTorrent protocol"...)

	for range 8 {
		b = append(b, 0)
	}

	// To signal support for extensions, set the 20th bit
	// from right to 1 (out of the 64 reserved bits)
	if h.SupportsExtensions {
		b[25] |= 1 << 4
	}

	b = append(b, h.InfoHash...)
	b = append(b, h.PeerID...)
	return b
}

func NewHandshakeFromBytes(data []byte) *Handshake {
	if len(data) < 68 {
		return nil
	}

	h := &Handshake{
		SupportsExtensions: (data[25] & (1 << 4)) != 0,
		InfoHash:           data[28:48],
		PeerID:             data[48:68],
	}

	return h
}

type ExtensionHandshake struct {
	ExtensionMap map[string]int
}

func NewExtensionHandshake() *ExtensionHandshake {
	ext := &ExtensionHandshake{
		ExtensionMap: make(map[string]int),
	}
	// Add the "ut_metadata" extension
	ext.ExtensionMap["ut_metadata"] = UT_METADATA_EXTENSION_ID
	return ext
}

func (e *ExtensionHandshake) getDictionaryBytes() []byte {
	// Create a dictionary with the extension map
	dict := make(map[string]*bencode.BencodeData)
	for k, v := range e.ExtensionMap {
		dict[k] = bencode.NewDataInteger(v)
	}

	// Base dictionary contains the extension map
	// under the key "m"
	bd := bencode.NewBencodeDictionary()
	bd.Add("m", bencode.NewDataDictionary(dict))
	return bd.Encode()
}

func (e *ExtensionHandshake) Bytes() []byte {
	// Get the body bytes
	body := make([]byte, 0)
	body = append(body, EXTENSION_HANDSHAKE_PAYLOAD_MESSAGE_ID)
	body = append(body, e.getDictionaryBytes()...)

	// Create the message
	msg := make([]byte, 0)
	l := len(body) + 1 // +1 for the message ID
	msg = binary.BigEndian.AppendUint32(msg, uint32(l))
	msg = append(msg, EXTENSION_HANDSHAKE_HEADER_MESSAGE_ID)
	msg = append(msg, body...)

	return msg
}

func NewExtensionHandshakeFromBytes(data []byte) (*ExtensionHandshake, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("data too short for extension handshake: %d", len(data))
	}

	// Parse the message IDs from the header and payload
	if data[0] != EXTENSION_HANDSHAKE_HEADER_MESSAGE_ID {
		return nil, fmt.Errorf("invalid message ID: expected %d, got %d", EXTENSION_HANDSHAKE_HEADER_MESSAGE_ID, data[0])
	}
	if data[1] != EXTENSION_HANDSHAKE_PAYLOAD_MESSAGE_ID {
		return nil, fmt.Errorf("invalid message ID: expected %d, got %d", EXTENSION_HANDSHAKE_PAYLOAD_MESSAGE_ID, data[0])
	}
	data = data[2:]

	// Make a new extension handshake object to store the parsed data
	handshake := &ExtensionHandshake{
		ExtensionMap: make(map[string]int),
	}

	// Parse the dictionary
	dict, err := bencode.NewBencodeData(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing dictionary: %w", err)
	}
	if dict.Type != bencode.DictionaryType {
		return nil, fmt.Errorf("expected dictionary type, got %s", dict.Type)
	}
	if extensionMap, ok := dict.GetDictionary().Map["m"]; ok {
		if extensionMap.Type != bencode.DictionaryType {
			return nil, fmt.Errorf("expected dictionary type for extensions, got %s", extensionMap.Type)
		}

		for k, v := range extensionMap.GetDictionary().Map {
			if v.Type != bencode.IntegerType {
				return nil, fmt.Errorf("expected integer type for extension value, got %s", v.Type)
			}
			handshake.ExtensionMap[k] = v.GetInteger().Value
		}
	}

	return handshake, nil
}
