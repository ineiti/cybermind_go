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
	DataBuf []byte
	datas   []Data
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

// Returns a new dataType, which is the first 64 bits of a sha256 hash of a unique URL. Using the birthday paradox,
// having 2**64 possible Data types, at 2**32 different Data types, there is a 50% chance of a collision,
// which we decide to live with.
func NewDataType(url string) DataType {
	sha := sha256.Sum256([]byte(url))
	return DataType(binary.LittleEndian.Uint64(sha[0:8]))
}

// DataType is used in Node.datas to store different types of Data.
type DataType uint64

// Data represents a Data, including its type.
// TODO: most probably this will disappear
type Data struct {
	Type DataType
	Data []byte
}

// NewNode creates a node and sets up all internal structures accordingly.
// The caller can add any number of Data in the arguments, including 0.
func NewNode(t NodeType, datas ...Data) Node {
	n := Node{
		NodeID: RandomNodeID(),
		Type:   t,
		Date:   time.Now().Unix(),
	}
	if len(datas) > 0 {
		err := n.SetDatas(datas...)
		if err != nil {
			panic(err)
		}
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
	if err := n.updateDataBuf(); err != nil {
		return err
	}
	if err := o.updateDataBuf(); err != nil {
		return err
	}
	if bytes.Compare(n.DataBuf, o.DataBuf) != 0 {
		return errors.New("dataBuf differs")
	}
	return nil
}

// GetDatas returns a slice of all Data of the node. As gorm cannot store structures,
// the Data itself is stored as a binary blob in the DB.
func (n Node) GetDatas() ([]Data, error) {
	if err := n.updateDatas(); err != nil {
		return nil, fmt.Errorf("couldn't update datas: %v", err)
	}
	return n.datas, nil
}

// SetDatas overwrites the current Data slice of the node.
func (n *Node) SetDatas(d ...Data) error {
	n.datas = d
	return n.updateDataBuf()
}

// GetData returns the Data of the given dataType. If there is more than one Data entry of the given dataType,
// only the first is returned.
func (n Node) GetData(t DataType) (d []byte, err error) {
	if err = n.updateDatas(); err != nil {
		return nil, fmt.Errorf("couldn't update datas: %v", err)
	}
	for _, d := range n.datas {
		if d.Type == t {
			return d.Data, nil
		}
	}
	return nil, fmt.Errorf("couldn't find dataType: %d", t)
}

func (n Node) GetNode() (Node, error) {
	if err := n.updateDataBuf(); err != nil {
		return n, fmt.Errorf("couldn't update Data buffer: %v", err)
	}
	return n, nil
}

// DataTypeGobEncoder is used to encode any Data in a gob format.
var DataTypeGobEncoder = NewDataType("blue.gasser/cybermind/Data/gobEncoder")

func (n Node) DecodeNodeType(t NodeType, i interface{}) error {
	if n.Type != t {
		return errors.New("node is not of correct type")
	}
	d, err := n.GetData(DataTypeGobEncoder)
	if err != nil {
		return fmt.Errorf("node doesn't have required Data: %v", err)
	}
	dec := gob.NewDecoder(bytes.NewBuffer(d))
	err = dec.Decode(i)
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
	return n.SetDatas(Data{Type: DataTypeGobEncoder, Data: buf.Bytes()})
}

func (n *Node) updateDataBuf() error {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(n.datas)
	if err != nil {
		return fmt.Errorf("couldn't encode datas: %v", err)
	}
	n.DataBuf = buf.Bytes()
	return nil
}

func (n *Node) updateDatas() error {
	if n.datas == nil && n.DataBuf != nil {
		dec := gob.NewDecoder(bytes.NewBuffer(n.DataBuf))
		err := dec.Decode(&n.datas)
		if err != nil {
			return fmt.Errorf("couldn't decode datas: %v", err)
		}
	}
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
