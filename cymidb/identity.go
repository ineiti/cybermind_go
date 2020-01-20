package cymidb

import (
	"bytes"
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
	if err := n.DecodeNodeType(NodeIdentity, DataTypeIdentity, &ident); err != nil {
		return ident, fmt.Errorf("couldn't decode node: %+v", err)
	}
	ident.node = n
	return
}

// NewIdentity returns a new Identity
func NewIdentity(a string, emails []string) (ident Identity, err error) {
	ident.node = NewNode(NodeIdentity)
	ident.Alias = a
	ident.Emails = emails
	return
}

// GetNode returns the updated node
func (ident Identity) GetNode() (Node, error) {
	err := ident.node.EncodeData(DataTypeIdentity, &ident)
	if err != nil {
		return ident.node, fmt.Errorf("couldn't update with data: %+v", err)
	}
	return ident.node, nil
}

// CompareTo returns nil if the two identities are equal, or an error otherwise.
func (ident Identity) Equals(other Identity) error {
	if bytes.Compare(ident.node.NodeID, other.node.NodeID) != 0 {
		return errors.New("not the same NodeID")
	}
	if ident.node.Version != other.node.Version {
		return errors.New("not the same version")
	}
	if ident.Alias != other.Alias {
		return errors.New("not the same alias")
	}
	if len(ident.Emails) != len(other.Emails) {
		return errors.New("not the same amount of emails")
	}
	for e := range ident.Emails {
		if ident.Emails[e] != other.Emails[e] {
			return errors.New("different email")
		}
	}
	return nil
}
