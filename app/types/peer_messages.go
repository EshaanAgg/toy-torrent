package types

import "fmt"

// PerformHandshake sends a handshake to the peer and waits for a response.
// Returns the received handshake or an error.
func (p *Peer) PerformHandshake(infoHash []byte) (*Handshake, error) {
	handshake := Handshake{
		PeerID:             SERVER_PEER_ID,
		InfoHash:           infoHash,
		SupportsExtensions: true,
	}

	_, err := p.conn.Write(handshake.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error sending handshake: %w", err)
	}

	response, err := p.readExactBytes(68) // 68 bytes for the handshake response
	if err != nil {
		return nil, fmt.Errorf("error reading handshake response: %w", err)
	}
	recievedHandshake := NewHandshakeFromBytes(response)

	return recievedHandshake, nil
}

func (p *Peer) PerformExtensionHandshake() (*ExtensionHandshake, error) {
	// Send handshake message
	handshake := NewExtensionHandshake()
	_, err := p.conn.Write(handshake.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error sending extension handshake: %w", err)
	}

	// Recieve handshake response
	data, err := p.RecieveMessage()
	if err != nil {
		return nil, fmt.Errorf("error receiving extension handshake response: %w", err)
	}
	response, err := NewExtensionHandshakeFromBytes(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing extension handshake response: %w", err)
	}
	return response, nil
}

// SendInterested sends an "interested" message to the peer indicating that
// we want to download data from them. It also waits for the unchoke message
// indicating that the peer is ready to send us data.
func (p *Peer) SendInterested() error {
	interested := []byte{INTERESTED_MESSAGE_ID}
	err := p.SendMessage(interested)
	if err != nil {
		return fmt.Errorf("error sending interested message: %w", err)
	}

	err = p.blockTillUnchokeMessage()
	if err != nil {
		return fmt.Errorf("error while wating for unchoke message: %w", err)
	}
	return nil
}

// blockTillBitFieldMessage blocks the peer until a bitfield message (ID = 5) is received.
// It skips over other messages and discards their payloads.
func (p *Peer) BlockTillBitFieldMessage() error {
	for {
		msg, err := p.RecieveMessage()
		if err != nil {
			return fmt.Errorf("error receiving message: %w", err)
		}

		if msg[0] == BITFIELD_MESSAGE_ID {
			// TODO: Parse the response to get all the pieces present
			return nil

		}

		p.Log("Recieved message bytes: %q while waiting for bitfield message", msg)
	}
}

// blockTillUnchokeMessage blocks the peer until an unchoke message (ID = 1) is received.
// It skips over other messages and discards their payloads.
func (p *Peer) blockTillUnchokeMessage() error {
	for {
		msg, err := p.RecieveMessage()
		if err != nil {
			return fmt.Errorf("error receiving message: %w", err)
		}

		if msg[0] == UNCHOKE_MESSAGE_ID {
			return nil
		}

		p.Log("Recieved message bytes: %q while waiting for unchoke message", msg)
	}
}
