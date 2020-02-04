package cymidb

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// NodeType points to one of the general node types. Most of the types have sub-types
type NodeType uint64

const (
	NodeDev = NodeType(iota * (1 << 56))
	NodeIdentity
	NodeHook
	NodeACL
	NodeBlob
	NodeLink
	NodeTag
)

func (nt NodeType) SubType(url string) NodeType {
	sha := sha256.Sum256([]byte(url))
	sub := binary.LittleEndian.Uint64(sha[:]) % (1 << 56)
	return NodeType(uint64(nt) + sub)
}

// RandomNodeID returns a random 256-bit ID.
func RandomNodeID() NodeID {
	nid := make([]byte, NodeIDLen)
	_, err := rand.Read(nid)
	if err != nil {
		panic("couldn't read random value: " + err.Error())
	}
	return nid
}

// NodeID
type NodeID []byte

// NodeIDLen is the length of the nodeID
const NodeIDLen = 32

// Node is the basic type in the DB. Every node can have 0 or more fields that are either Data, or point to other nodes.
type Node struct {
	gorm.Model
	NodeID  NodeID
	Type    NodeType
	Version uint64
	Date    int64
	Data    []byte
}

// Noder can be used for inherited types that need to be stored,
// so they can prepare eventual cached Data and write it to the node before storing.
type Noder interface {
	GetNode() (Node, error)
}

// Link is used to link a parent to a child node, or a child to an ancestor.
type Link struct {
	From NodeID
	To   NodeID
}

// NewNode creates a node and sets up all internal structures accordingly.
// The caller can add any number of Data in the arguments, including 0.
func NewNode(t NodeType) Node {
	n := Node{
		NodeID: RandomNodeID(),
		Type:   t,
		Date:   time.Now().Unix(),
	}
	return n
}

// CompareTo tests if the two nodes have the equal content, not if they are the same object in the DB.
// Two nodes can be different objects in the DB (having different dates, IDs), but nevertheless be the same node.
func (n Node) CompareTo(o Node) error {
	if bytes.Compare(n.NodeID, o.NodeID) != 0 {
		return errors.New("NodeID differs")
	}
	if n.Type != o.Type {
		return errors.New("type differs")
	}
	if n.Version != o.Version {
		return errors.New("version differs")
	}
	if n.Date != o.Date {
		return errors.New("date differs")
	}
	if bytes.Compare(n.Data, o.Data) != 0 {
		return errors.New("dataBuf differs")
	}
	return nil
}

func (n Node) GetNode() (Node, error) {
	return n, nil
}

func (n Node) DecodeNodeType(t NodeType, i interface{}) error {
	if n.Type != t {
		return errors.New("node is not of correct type")
	}
	dec := gob.NewDecoder(bytes.NewBuffer(n.Data))
	err := dec.Decode(i)
	if err != nil {
		return fmt.Errorf("couldn't decode Data: %v", err)
	}
	return nil
}

func (n *Node) EncodeData(i interface{}) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(i)
	if err != nil {
		return fmt.Errorf("couldn't encode Data: %v", err)
	}
	n.Data = buf.Bytes()
	return nil
}

func NoderCompare(noder1, noder2 Noder) error {
	node1, err := noder1.GetNode()
	if err != nil {
		return fmt.Errorf("couldn't get Node of noder1: %v", err)
	}
	node2, err := noder1.GetNode()
	if err != nil {
		return fmt.Errorf("couldn't get Node of noder2: %v", err)
	}
	return node1.CompareTo(node2)
}
