package cymidb

import (
	"fmt"
)

// Device holds all information for a new device.
type Device struct {
	Name string
	URL  string
	node Node
}

// NewDeviceFromNode takes a node and returns a device. If the node is not of the correct type,
// or if the name is not present, an error will be returned.
func NewDeviceFromNode(n Node) (dev Device, err error) {
	err = n.DecodeNodeType(NodeDev, &dev)
	if err != nil {
		return dev, fmt.Errorf("couldn't decode device node: %v", err)
	}
	dev.node = n
	return
}

// CreateNewNode returns a node that
func NewDevice(name string) (dev Device) {
	dev.node = NewNode(NodeDev)
	dev.Name = name
	return
}

// GetNode makes sure that the dataBuf of the node is updated and returns the updated node.
func (dev Device) GetNode() (Node, error) {
	err := dev.node.EncodeData(&dev)
	return dev.node, err
}
