package bencode

import (
	"fmt"
	"strings"
)

type BencodeDictionary struct {
	Map    map[string]*BencodeData
	Length int
}

// Parses a dictionary encoded in the bencode format.
// Format: d<key1><item1>...e
func (d *decoder) parseDictionary() (*BencodeDictionary, error) {
	err := d.expect('d')
	if err != nil {
		return nil, err
	}

	itemMap := make(map[string]*BencodeData)
	for {
		ch := d.peek()
		if ch != nil && *ch == 'e' {
			break
		}

		key, err := d.parse()
		if err != nil {
			return nil, fmt.Errorf("parsing dictionary key: %w", err)
		}
		if key.Type != StringType {
			return nil, fmt.Errorf("expected string key in dictionary, got %s", key.Type)
		}

		item, err := d.parse()
		if err != nil {
			return nil, fmt.Errorf("parsing dictionary item: %w", err)
		}

		itemMap[key.GetString().Value] = item
	}

	err = d.expect('e')
	if err != nil {
		return nil, err
	}

	return &BencodeDictionary{
		Map:    itemMap,
		Length: len(itemMap),
	}, nil
}

func (s *BencodeDictionary) String() string {
	elements := make([]string, 0)
	for key, item := range s.Map {
		elements = append(elements, fmt.Sprintf("\"%s\":%s", key, item.Value.String()))
	}

	return fmt.Sprintf("{%s}", strings.Join(elements, ","))
}

func (s *BencodeDictionary) Encode() string {
	str := "d"
	for key, item := range s.Map {
		str += fmt.Sprintf("%d:%s", len(key), key) // Encode the key as string
		str += item.Value.Encode()
	}
	str += "e"

	return str
}
