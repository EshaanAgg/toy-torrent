package bencode

import (
	"fmt"
)

type BencodeString struct {
	Value  []byte
	Length int
}

// Parses a string encoded in the bencode format.
// Format: <length>:<bytes>
func (d *decoder) parseString() (*BencodeString, error) {
	l, err := d.readInteger()
	if err != nil {
		return nil, fmt.Errorf("length of string: %w", err)
	}

	err = d.expect(':')
	if err != nil {
		return nil, err
	}

	b := d.getNBytes(l)
	if b == nil {
		return nil, fmt.Errorf("not enough bytes: expected to read %d bytes", l)
	}

	return &BencodeString{
		Value:  b,
		Length: len(b),
	}, nil
}

func (s *BencodeString) String() string {
	return fmt.Sprintf("\"%s\"", s.Value)
}

func (s *BencodeString) Encode() []byte {
	return append(fmt.Appendf([]byte{}, "%d:", s.Length), s.Value...)
}
