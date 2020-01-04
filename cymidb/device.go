package cymidb

import (
	"errors"
	"fmt"
	"time"
)

// Device holds all information for a new device.
type Device struct {
	Node
	Name string
	URL  string
}

// DataTypeDeviceName is used in a device node to represent the name of the device.
var DataTypeDeviceName = NewDataType("blue.gasser/cybermind/device/name")

// DataTypeDeviceURL is used in a device node to represent the URL of the device.
var DataTypeDeviceURL = NewDataType("blue.gasser/cybermind/device/url")

// NewDeviceFromNode takes a node and returns a device. If the node is not of the correct type,
// or if the name is not present, an error will be returned.
func NewDeviceFromNode(n Node) (dev Device, err error) {
	if n.Type != NodeDev {
		return dev, errors.New("node is not of type device")
	}
	dev.Node = n
	md, err := n.GetData(DataTypeDeviceName)
	if err != nil {
		return dev, fmt.Errorf("couldn't get data: %+v", err)
	}
	dev.Name = string(md)
	return
}

// CreateNewNode returns a node that
func NewDevice(name string) (dev Device) {
	dev.Node = Node{
		NodeID:  RandomNodeID(),
		Type:    NodeDev,
		Date:    time.Now().Unix(),
		Version: 0,
	}
	dev.Name = name
	dev.UpdateNode()
	return
}

// UpdateNode copies the name and url to the node structure.
func (dev *Device) UpdateNode() {
	dev.Node.datas = []Data{
		{Type: DataTypeDeviceName, Data: []byte(dev.Name)},
		{Type: DataTypeDeviceURL, Data: []byte(dev.URL)},
	}
}
