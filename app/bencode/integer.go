package bencode

import (
	"fmt"
)

type BencodeInteger struct {
	Value int
}

// Parses a integer encoded in the bencode format.
// Format: i<number>e
// Examples: i52e, i-52e
func (d *decoder) parseInteger() (*BencodeInteger, error) {
	err := d.expect('i')
	if err != nil {
		return nil, err
	}

	isNeg := false
	nxt := d.peek()
	if nxt != nil && *nxt == '-' {
		isNeg = true
		d.next()
	}

	n, err := d.readInteger()
	if err != nil {
		return nil, err
	}

	err = d.expect('e')
	if err != nil {
		return nil, fmt.Errorf("expected 'e' to terminate integer: %w", err)
	}

	if isNeg {
		n = n * -1
	}
	return &BencodeInteger{
		Value: n,
	}, nil
}

func (i *BencodeInteger) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *BencodeInteger) Encode() string {
	return fmt.Sprintf("i%de", i.Value)
}
