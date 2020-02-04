package cymidb

import (
	"fmt"
)

// Hook represents a place where external services can listen and send requests to add new blobs and other pieces in
// the DB.
type Hook struct {
	// Name is the name of the hook and must be unique
	Name string
	// Devices holds all IDs of devices this hook is active on
	Devices []NodeID
	// Types is an array of types which will create events when they change
	Types []NodeType
	node  Node
	db    DB
}

func NewHookFromNode(db DB, n Node) (h Hook, err error) {
	err = n.DecodeNodeType(NodeHook, &h)
	if err != nil {
		return h, fmt.Errorf("couldn't decode hook node: %v", err)
	}
	h.node = n
	h.db = db
	return
}

// NewHook returns a hook defined on a new node with the name set.
func NewHook(name string, devices []NodeID, types []NodeType) (h Hook) {
	h.node = NewNode(NodeHook)
	h.Name = name
	h.Devices = devices
	h.Types = types
	return
}

// GetNode is used to implement Noders.
func (h Hook) GetNode() (Node, error) {
	err := h.node.EncodeData(&h)
	return h.node, err
}
