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

func NewEncodeString(s string) *BencodeData {
	return &BencodeData{
		Type: StringType,
		Value: &BencodeString{
			Length: len(s),
			Value:  []byte(s),
		},
	}
}

func NewEncodeInteger(i int) *BencodeData {
	return &BencodeData{
		Type: IntegerType,
		Value: &BencodeInteger{
			Value: i,
		},
	}
}

func NewEncodeList(l []*BencodeData) *BencodeData {
	return &BencodeData{
		Type: ListType,
		Value: &BencodeList{
			Array:  l,
			Length: len(l),
		},
	}
}

func NewEncodeDictionary(m map[string]*BencodeData) *BencodeData {
	return &BencodeData{
		Type: DictionaryType,
		Value: &BencodeDictionary{
			Map:    m,
			Length: len(m),
		},
	}
}
