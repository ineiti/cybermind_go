package cymidb

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
)

// Identity holds all information for an identity.
type Identity struct {
	Alias  string
	Emails []string
	node   Node
}

// DataTypeIdentity is used in a device node to represent the metadata of the identity.
var DataTypeIdentity = NewDataType("blue.gasser/cybermind/data/identity")

// NewIdentityFromNode takes a node and returns an Identity. If the node is not of the correct type,
// or if the name is not present, an error will be returned.
func NewIdentityFromNode(n Node) (ident Identity, err error) {
	if n.Type != NodeIdentity {
		return ident, errors.New("node is not of type Identity")
	}
	ident.node = n
	md, err := n.GetData(DataTypeIdentity)
	if err != nil {
		return ident, fmt.Errorf("couldn't get data: %+v", err)
	}
	dec := json.NewDecoder(bytes.NewBuffer(md))
	err = dec.Decode(&ident)
	if err != nil {
		return ident, fmt.Errorf("couldn't decode JSON data: %+v", err)
	}
	return
}

// NewIdentity returns a new Identity
func NewIdentity(a string, emails []string) (ident Identity, err error) {
	ident.node = NewNode(NodeDev)
	ident.Alias = a
	ident.Emails = emails
	if err != nil {
		return ident, fmt.Errorf("couldn't update node: %+v", err)
	}
	return
}

// updateNode copies the name and url to the node structure.
func (ident *Identity) updateNode() error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(&ident)
	if err != nil {
		return fmt.Errorf("couldn't encode meta: %+v", err)
	}
	return ident.node.SetDatas(Data{Type: DataTypeIdentity, Data: buf.Bytes()})
}

// GetNode returns the updated node
func (ident Identity) GetNode() (Node, error) {
	err := ident.updateNode()
	if err != nil {
		return ident.node, fmt.Errorf("couldn't update with data: %+v", err)
	}
	return ident.node, nil
}
