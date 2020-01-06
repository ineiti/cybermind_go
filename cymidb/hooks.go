package cymidb

import (
	"errors"
	"fmt"
)

// Hook represents a place where external services can listen and send requests to add new blobs and other pieces in
// the DB.
type Hook struct {
	// Name is the name of the hook and must be unique
	Name string
	node Node
}

// DataTypeHookName is the unique ID for the name
var DataTypeHookName = NewDataType("blue.gasser/cybermind/hook/name")

func NewHookFromNode(n Node) (h Hook, err error) {
	if n.Type != NodeHook {
		return h, errors.New("node is not of type hook")
	}
	h.node = n
	name, err := n.GetData(DataTypeHookName)
	if err != nil {
		return h, fmt.Errorf("couldn't get name of hook: %+v", err)
	}
	h.Name = string(name)
	return
}

// NewHook returns a hook defined on a new node with the name set.
func NewHook(name string) (h Hook) {
	h.node = NewNode(NodeHook)
	h.Name = name
	return
}

// GetNode is used to implement Noders.
func (h Hook) GetNode() (Node, error) {
	err := h.node.SetDatas(Data{Type: DataTypeHookName, Data: []byte(h.Name)})
	return h.node, err
}
