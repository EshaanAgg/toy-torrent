package bencode

import (
	"errors"
	"fmt"
)

type decoder struct {
	data    string
	idx     int
	dataLen int
}

func newDecoder(data string) *decoder {
	return &decoder{
		data:    data,
		dataLen: len(data),
		idx:     0,
	}
}

func (d *decoder) peek() *byte {
	if d.idx >= len(d.data) {
		return nil
	}
	b := d.data[d.idx]
	return &b
}

func (d *decoder) next() *byte {
	if d.idx >= d.dataLen {
		return nil
	}
	b := d.data[d.idx]
	d.idx++
	return &b
}

func (d *decoder) expect(b byte) error {
	ch := d.peek()
	if ch == nil || *ch != b {
		return fmt.Errorf("expected character at index %d to be '%q': got '%q'", d.idx, b, *ch)
	}

	// Consume the character
	d.next()
	return nil
}

func (d *decoder) getNBytes(l int) *string {
	if d.idx+l-1 >= d.dataLen {
		return nil
	}

	v := d.data[d.idx : d.idx+l]
	d.idx += l
	return &v
}

func (d *decoder) readInteger() (int, error) {
	nxt := d.peek()
	if nxt == nil {
		return 0, errors.New("not enough bytes")
	}
	if !isDigit(nxt) {
		return 0, fmt.Errorf("expected a digit, recieved '%q'", *nxt)
	}

	// Consume the first byte to read the number
	v := int(*nxt) - int('0')
	d.next()

	// Continue parsing
	for isDigit(d.peek()) {
		n := int(*d.next()) - int('0')
		v = v*10 + n
	}

	return v, nil
}

func (d *decoder) parse() (*BencodeData, error) {
	ch := d.peek()

	if ch == nil {
		return nil, fmt.Errorf("unexpected end of data")
	}

	// Parse string if it starts with a numeric character
	if isDigit(ch) {
		s, err := d.parseString()
		if err != nil {
			return nil, fmt.Errorf("error parsing string: %w", err)
		}
		return &BencodeData{
			Type:  StringType,
			Value: s,
		}, nil
	}

	// Parse integer if it starts with byte 'i'
	if *ch == 'i' {
		n, err := d.parseInteger()
		if err != nil {
			return nil, fmt.Errorf("error parsing integer: %w", err)
		}
		return &BencodeData{
			Type:  IntegerType,
			Value: n,
		}, nil
	}

	// Parse list if it starts with byte 'l'
	if *ch == 'l' {
		l, err := d.parseList()
		if err != nil {
			return nil, fmt.Errorf("error parsing list: %w", err)
		}
		return &BencodeData{
			Type:  ListType,
			Value: l,
		}, nil
	}

	return nil, fmt.Errorf("unrecognized character to start parsing: %q", *ch)
}
