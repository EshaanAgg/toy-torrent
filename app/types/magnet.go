package types

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
)

type MagnetURI struct {
	FileToDownload string
	TrackerURL     string

	InfoHash    []byte
	InfoHashHex string
}

func (m *MagnetURI) setValue(key, value string) error {
	switch key {
	case "dn":
		m.FileToDownload = value

	case "tr":
		m.TrackerURL = value

	case "xt":
		if !strings.HasPrefix(value, "urn:btih:") {
			return fmt.Errorf("invalid info hash format: %s", value)
		}
		v := strings.TrimPrefix(value, "urn:btih:")
		if len(v) != 40 {
			return fmt.Errorf("invalid info hash length: %s", v)
		}
		m.InfoHashHex = v

		bytes, err := hex.DecodeString(v)
		if err != nil {
			return fmt.Errorf("failed to decode info hash: %v", err)
		}
		m.InfoHash = bytes

	default:
		return fmt.Errorf("unknown key: %s", key)
	}

	return nil
}

func NewMagnetURI(s string) (*MagnetURI, error) {
	s, err := url.PathUnescape(s)
	if err != nil {
		return nil, fmt.Errorf("failed to unescape magnet link: %v", err)
	}

	p := &parser{s, 0}

	dec := p.get(8)
	if dec == nil || *dec != "magnet:?" {
		return nil, fmt.Errorf("expected the URL to start with 'magnet:?'")
	}

	m := &MagnetURI{}

	for !p.isAtEnd() {
		t := p.get(2)
		if t == nil {
			return nil, fmt.Errorf("expected a query parameter, but reached end of string")
		}
		if err := p.expect('='); err != nil {
			return nil, err
		}
		val := p.readUntil('&')

		if err := m.setValue(*t, val); err != nil {
			return nil, err
		}
	}

	return m, nil
}
