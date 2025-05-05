package bencode

import (
	"fmt"
	"slices"
	"strings"
)

type BencodeDictionary struct {
	Map    map[string]*BencodeData
	Length int
}

// sortedMapByKey sorts the map by its keys and returns a new map with the sorted keys.
// This is necessary because Go maps do not maintain order.
func sortedMapByKey(m map[string]*BencodeData) map[string]*BencodeData {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	sorted := make(map[string]*BencodeData)
	for _, k := range keys {
		sorted[k] = m[k]
	}
	return sorted
}

// Parses a dictionary encoded in the bencode format.
// Format: d<key1><item1>...e
func (d *decoder) parseDictionary() (*BencodeDictionary, error) {
	err := d.expect('d')
	if err != nil {
		return nil, err
	}

	itemsMap := make(map[string]*BencodeData)
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

		itemsMap[key.GetString().Value] = item
	}

	err = d.expect('e')
	if err != nil {
		return nil, err
	}

	itemsMap = sortedMapByKey(itemsMap)
	return &BencodeDictionary{
		Map:    itemsMap,
		Length: len(itemsMap),
	}, nil
}

func (s *BencodeDictionary) String() string {
	elements := make([]string, 0)
	s.Map = sortedMapByKey(s.Map)
	for key, item := range s.Map {
		elements = append(elements, fmt.Sprintf("\"%s\":%s", key, item.Value.String()))
	}

	return fmt.Sprintf("{%s}", strings.Join(elements, ","))
}

func (s *BencodeDictionary) Encode() string {
	str := "d"
	s.Map = sortedMapByKey(s.Map)
	for key, item := range s.Map {
		str += fmt.Sprintf("%d:%s", len(key), key) // Encode the key as string
		str += item.Value.Encode()
	}
	str += "e"

	return str
}
