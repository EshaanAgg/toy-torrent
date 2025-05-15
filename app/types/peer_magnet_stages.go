package types

import "fmt"

func (p *Peer) PerformMagnetHandshake(m *MagnetURI, logIDs bool) (*Handshake, error) {
	handshake, err := p.PerformHandshake(m.InfoHash)
	if err != nil {
		return nil, fmt.Errorf("error performing handshake: %w", err)
	}
	if logIDs {
		fmt.Printf("Peer ID: %x\n", handshake.PeerID)
	}

	// TODO: Send a bitfield message to the peer
	err = p.BlockTillBitFieldMessage()
	if err != nil {
		return nil, fmt.Errorf("error waiting for bitfield message: %w", err)
	}

	if handshake.SupportsExtensions {
		extHandshake, err := p.PerformExtensionHandshake()
		if err != nil {
			return nil, fmt.Errorf("error performing extension handshake: %w", err)
		}
		if mtExtensionId, ok := extHandshake.ExtensionMap["ut_metadata"]; ok {
			p.ExtensionMessageID = mtExtensionId
			if logIDs {
				fmt.Printf("Peer Metadata Extension ID: %d\n", mtExtensionId)
			}
		}
	}

	return handshake, nil
}

func (p *Peer) MagnetHandshakeAndInfoFile(m *MagnetURI) (*TorrentFileInfo, error) {
	_, err := p.PerformMagnetHandshake(m, false)
	if err != nil {
		return nil, fmt.Errorf("error performing magnet handshake: %w", err)
	}

	// Get the torrent file info from the magnet link
	infoFile, err := p.GetInfoFile(m)
	if err != nil {
		return nil, fmt.Errorf("error getting info file: %w", err)
	}

	return infoFile, nil
}

// PrepareToGetPieceData is to be used for stages with magnet links.
func (p *Peer) PrepareToGetPieceData_Magnet(m *MagnetURI) (*TorrentFileInfo, error) {

	// Get the torrent file info from the magnet link
	infoFile, err := p.GetInfoFile(m)
	if err != nil {
		return nil, fmt.Errorf("error getting info file: %w", err)
	}

	err = p.SendInterested()
	if err != nil {
		return nil, fmt.Errorf("error while sending interested message: %w", err)
	}

	p.Log("completed initialization. ready to download pieces")
	return infoFile, nil
}
