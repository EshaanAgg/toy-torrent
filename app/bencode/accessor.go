package bencode

import "fmt"

// Returns the string value associated with the data node.
// If the underlying type of the node is not StringType, the
// function would panic with an error message conveying the same.
func (bd *BencodeData) GetString() *BencodeString {
	if bd.Type != StringType {
		panic(fmt.Sprintf("called GetString on the BencodeData node %v with type %s", *bd, bd.Type))
	}

	v, ok := bd.Value.(*BencodeString)
	if !ok {
		panic(fmt.Sprintf("%v has type String, but conversion to Bencode string failed", *bd))
	}

	return v
}

// Returns the integer value associated with the data node.
// If the underlying type of the node is not IntegerType, the
// function would panic with an error message conveying the same.
func (bd *BencodeData) GetInteger() *BencodeInteger {
	if bd.Type != IntegerType {
		panic(fmt.Sprintf("called GetInteger on the BencodeData node %v with type %s", *bd, bd.Type))
	}

	v, ok := bd.Value.(*BencodeInteger)
	if !ok {
		panic(fmt.Sprintf("%v has type Integer, but conversion to Bencode integer failed", *bd))
	}

	return v
}

// Returns the list value associated with the data node.
// If the underlying type of the node is not ListType, the
// function would panic with an error message conveying the same.
func (bd *BencodeData) GetList() *BencodeList {
	if bd.Type != ListType {
		panic(fmt.Sprintf("called GetList on the BencodeData node %v with type %s", *bd, bd.Type))
	}

	v, ok := bd.Value.(*BencodeList)
	if !ok {
		panic(fmt.Sprintf("%v has type List, but conversion to Bencode list failed", *bd))
	}

	return v
}

// Returns the dictionary value associated with the data node.
// If the underlying type of the node is not DictionaryType, the
// function would panic with an error message conveying the same.
func (bd *BencodeData) GetDictionary() *BencodeDictionary {
	if bd.Type != DictionaryType {
		panic(fmt.Sprintf("called GetDictionary on the BencodeData node %v with type %s", *bd, bd.Type))
	}

	v, ok := bd.Value.(*BencodeDictionary)
	if !ok {
		panic(fmt.Sprintf("%v has type Dictionary, but conversion to Bencode dictionary failed", *bd))
	}

	return v
}
