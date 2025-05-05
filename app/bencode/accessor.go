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
