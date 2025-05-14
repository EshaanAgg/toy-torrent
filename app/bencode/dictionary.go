package bencode

import (
	"fmt"
	"sort"
	"strings"
)

type BencodeDictionary struct {
	Map    map[string]*BencodeData
	Length int
}

func NewBencodeDictionary() *BencodeDictionary {
	return &BencodeDictionary{
		Map:    make(map[string]*BencodeData),
		Length: 0,
	}
}

func (d *BencodeDictionary) Add(key string, value *BencodeData) {
	d.Map[key] = value
	d.Length++
}

// Returns the sorted keys of the dictionary. This is used to ensure that the
// dictionary is always encoded in a consistent order, which is important for
// hashing and comparison purposes.
func getSortedKeys(m map[string]*BencodeData) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
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

		keyStr := string(key.GetString().Value)
		itemsMap[keyStr] = item
	}

	err = d.expect('e')
	if err != nil {
		return nil, err
	}

	return &BencodeDictionary{
		Map:    itemsMap,
		Length: len(itemsMap),
	}, nil
}

func (s *BencodeDictionary) String() string {
	elements := make([]string, 0)
	keys := getSortedKeys(s.Map)

	for _, key := range keys {
		elements = append(elements, fmt.Sprintf("\"%s\":%s", key, s.Map[key].Value.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(elements, ","))
}

func (s *BencodeDictionary) Encode() []byte {
	res := []byte{'d'}
	keys := getSortedKeys(s.Map)

	for _, key := range keys {
		res = append(res, fmt.Sprintf("%d:%s", len(key), key)...) // Encode key as string
		res = append(res, s.Map[key].Value.Encode()...)
	}
	res = append(res, 'e')

	return res
}
