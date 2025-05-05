package bencode

import (
	"fmt"
	"strings"
)

type BencodeList struct {
	Array  []*BencodeData
	Length int
}

// Parses a list encoded in the bencode format.
// Format: l<item><item>...e
func (d *decoder) parseList() (*BencodeList, error) {
	err := d.expect('l')
	if err != nil {
		return nil, err
	}

	var items []*BencodeData
	for {
		ch := d.peek()
		if ch != nil && *ch == 'e' {
			break
		}

		item, err := d.parse()
		if err != nil {
			return nil, fmt.Errorf("parsing list item: %w", err)
		}
		items = append(items, item)
	}

	err = d.expect('e')
	if err != nil {
		return nil, err
	}

	return &BencodeList{
		Array:  items,
		Length: len(items),
	}, nil
}

func (s *BencodeList) String() string {
	elements := make([]string, 0)
	for _, item := range s.Array {
		elements = append(elements, item.Value.String())
	}

	return fmt.Sprintf("[%s]", strings.Join(elements, ","))
}

func (s *BencodeList) Encode() string {
	str := "l"
	for _, item := range s.Array {
		str += item.Value.Encode()
	}
	str += "e"

	return str
}
