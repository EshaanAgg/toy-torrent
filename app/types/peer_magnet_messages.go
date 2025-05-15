package types

import (
	"fmt"

	"github.com/EshaanAgg/toy-bittorrent/app/bencode"
)

func (p *Peer) GetInfoFile(m *MagnetURI) (*TorrentFileInfo, error) {
	// Send the magnet request message
	err := p.SendMagnetRequestMessage(0)
	if err != nil {
		return nil, fmt.Errorf("error sending magnet request message: %w", err)
	}

	// Receive the magnet data message
	infoDict, err := p.GetMagnetDataMessage(m, 0)
	if err != nil {
		return nil, fmt.Errorf("error receiving magnet data message: %w", err)
	}

	return infoDict, nil
}

func (p *Peer) SendMagnetRequestMessage(pieceIndex int) error {
	// Create the body of the message
	body := make([]byte, 0)
	body = append(body, EXTENSION_HANDSHAKE_HEADER_MESSAGE_ID)
	body = append(body, byte(p.ExtensionMessageID))

	payload := bencode.NewBencodeDictionary()
	payload.Add("msg_type", bencode.NewDataInteger(0))
	payload.Add("piece", bencode.NewDataInteger(pieceIndex))
	body = append(body, payload.Encode()...)

	// Send the message
	err := p.SendMessage(body)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}
	return nil
}

func (p *Peer) GetMagnetDataMessage(m *MagnetURI, pieceIndex int) (*TorrentFileInfo, error) {
	msg, err := p.RecieveMessage()
	if err != nil {
		return nil, fmt.Errorf("error receiving message: %w", err)
	}

	infoDict, err := parseDataMessage(pieceIndex, msg)
	if err != nil {
		return nil, fmt.Errorf("error parsing data message: %w", err)
	}

	return NewTorrentFileInfoFromMagnet(m, infoDict)
}

func parseDataMessage(pieceIdx int, msg []byte) (*bencode.BencodeDictionary, error) {
	if len(msg) < 2 {
		return nil, fmt.Errorf("message too short: %d", len(msg))
	}

	if msg[0] != EXTENSION_HANDSHAKE_HEADER_MESSAGE_ID {
		return nil, fmt.Errorf("invalid message ID: expected %d, got %d", EXTENSION_HANDSHAKE_HEADER_MESSAGE_ID, msg[0])
	}

	if msg[1] != UT_METADATA_EXTENSION_ID {
		return nil, fmt.Errorf("invalid extension ID: expected %d, got %d", UT_METADATA_EXTENSION_ID, msg[1])
	}

	return getDataMessageDict(pieceIdx, msg[2:])
}

func getDataMessageDict(pieceIdx int, data []byte) (*bencode.BencodeDictionary, error) {
	d, leftData, err := bencode.NewPartialBencodeData(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing dictionary: %w", err)
	}

	if d.Type != bencode.DictionaryType {
		return nil, fmt.Errorf("expected dictionary type, got %s", d.Type)
	}
	dict := d.GetDictionary()

	// Valid the different keys
	if msg_type, err := dict.GetInteger("msg_type"); err != nil || msg_type != 1 {
		return nil, fmt.Errorf("invalid message type: %w", err)
	}
	if piece_index, err := dict.GetInteger("piece"); err != nil || piece_index != pieceIdx {
		return nil, fmt.Errorf("invalid piece index: %w", err)
	}
	if total_size, err := dict.GetInteger("total_size"); err != nil || total_size != len(leftData) {
		return nil, fmt.Errorf("invalid total size: %w", err)
	}

	// Decode the remaining data
	decodedData, err := bencode.NewBencodeData(leftData)
	if decodedData.Type != bencode.DictionaryType {
		return nil, fmt.Errorf("expected dictionary type, got %s", decodedData.Type)
	}
	return decodedData.GetDictionary(), nil
}
