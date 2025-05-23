package bencode

import "fmt"

type DataType string

const (
	StringType     DataType = "String"
	IntegerType    DataType = "Integer"
	ListType       DataType = "List"
	DictionaryType DataType = "Dictionary"
)

type ValueInterface interface {
	String() string // Returns the string representation of the value
	Encode() []byte // Returns the encoded bytes according to bencode format
}

type BencodeData struct {
	Type  DataType
	Value ValueInterface // Stores the pointers to the underlying values
}

func NewBencodeData(s []byte) (*BencodeData, error) {
	d := newDecoder(s)

	v, err := d.parse()
	if err != nil {
		return nil, err
	}

	// Check that we are at the last byte of the decoder
	if (d.idx) != d.dataLen {
		return nil, fmt.Errorf("data bytes left unprocessed from index %d even after parsing: '%s'", d.idx, d.data[d.idx:])
	}

	return v, nil
}

// NewPartialBencodeData is used to parse a partial bencode data.
// It does not check if the entire data has been parsed.
// It also returns the remaining bytes after parsing.
func NewPartialBencodeData(s []byte) (*BencodeData, []byte, error) {
	d := newDecoder(s)

	v, err := d.parse()
	if err != nil {
		return nil, nil, err
	}

	leftData := d.data[d.idx:]
	return v, leftData, nil
}

func NewDataString(s string) *BencodeData {
	return &BencodeData{
		Type: StringType,
		Value: &BencodeString{
			Length: len(s),
			Value:  []byte(s),
		},
	}
}

func NewDataInteger(i int) *BencodeData {
	return &BencodeData{
		Type: IntegerType,
		Value: &BencodeInteger{
			Value: i,
		},
	}
}

func NewDataList(l []*BencodeData) *BencodeData {
	return &BencodeData{
		Type: ListType,
		Value: &BencodeList{
			Array:  l,
			Length: len(l),
		},
	}
}

func NewDataDictionary(m map[string]*BencodeData) *BencodeData {
	return &BencodeData{
		Type: DictionaryType,
		Value: &BencodeDictionary{
			Map:    m,
			Length: len(m),
		},
	}
}
