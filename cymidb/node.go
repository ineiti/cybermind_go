package cymidb

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// NodeType points to one of the general node types. Most of the types have sub-types
type NodeType uint64

const (
	NodeDev = iota * 0x1000
	NodeIdentity
	NodeHook
	NodeACL
	NodeBlob
	NodeLink
	NodeTag
)

func RandomNodeID() []byte {
	nid := make([]byte, 32)
	_, err := rand.Read(nid)
	if err != nil {
		panic("couldn't read random value: " + err.Error())
	}
	return nid
}

type NodeID []byte

// Node is the basic type in the DB. Every node can have 0 or more fields that are either data, or point to other nodes.
type Node struct {
	ID      []byte
	Type    NodeType
	Version uint64
	Date    int64
	Links   []Link
	Datas   []Data
}

type LinkType uint64

// Link is either pointing to another node, or has data in it.
type Link struct {
	Type LinkType
	Link []byte
}

func NewDataType(url string) DataType {
	sha := sha256.Sum256([]byte(url))
	return DataType(binary.LittleEndian.Uint64(sha[0:8]))
}

type DataType uint64

type Data struct {
	Type DataType
	Data []byte
}

func (n Node) GetLinks(t LinkType) ([]NodeID, error) {
	return nil, nil
}

func (n Node) GetData(t DataType) (d []byte, err error) {
	for _, d := range n.Datas {
		if d.Type == t {
			return d.Data, nil
		}
	}
	return nil, fmt.Errorf("couldn't find dataType: %d", t)
}
