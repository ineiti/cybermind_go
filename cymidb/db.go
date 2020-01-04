package cymidb

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DB represents one CyMiDB.
type DB struct {
	gdb    *gorm.DB
	Device Device
}

// NewDBFile opens the DB with the given file and autoMigrates for MemoryLaneEntry and NodeVersions.
func NewDBFile(file string) (db DB, err error) {
	db.gdb, err = gorm.Open("sqlite3", file)
	if err != nil {
		return db, fmt.Errorf("coulnd't open sqlite3")
	}
	//db.gdb.LogMode(true)
	db.gdb.AutoMigrate(&Node{}, &Link{})
	return
}

// CreateDBFile creates a new DB in a given file and returns an initialized DB containing only the given device.
func CreateDBFile(file string, name, url string) (db DB, err error) {
	db, err = NewDBFile(file)
	if err != nil {
		return db, err
	}

	db.Device = NewDevice(name)
	err = db.SaveNode(db.Device)
	if err != nil {
		return db, fmt.Errorf("couldn't create new node: %+v", err)
	}
	return
}

// OpenDBFile returns a db initialised with a file. It returns either the db, if successful,
// or an error. If the db did not exist previously, the method will return an error.
func OpenDBFile(file string) (db DB, err error) {
	db, err = NewDBFile(file)
	if err != nil {
		return db, err
	}

	n := Node{}
	err = db.gdb.First(&n).Error
	if err != nil {
		return db, fmt.Errorf("couldn't get first node: %+v", err)
	}
	node, err := db.GetLatest(n.NodeID)
	if err != nil {
		return db, fmt.Errorf("couldn't get latest node version: %+v", err)
	}
	db.Device, err = NewDeviceFromNode(node)
	if err != nil {
		return db, fmt.Errorf("couldn't get device from node: %+v", err)
	}
	return
}

// Closes the connection to the database. No further action is possible after this call.
func (db DB) Close() error {
	return db.gdb.Close()
}

// NewNode takes a node and inserts it in the DB.
func (db DB) SaveNode(n Noder) error {
	node, err := n.GetNode()
	if err != nil {
		return fmt.Errorf("couldn't get node: %+v", err)
	}
	var exist Node
	db.gdb.Last(&exist, &Node{NodeID: node.NodeID})
	if bytes.Compare(exist.NodeID, node.NodeID) == 0 {
		node.Version = exist.Version + 1
	}
	err = db.gdb.Save(&node).Error
	if err != nil {
		return fmt.Errorf("couldn't create new node: %+v", err)
	}
	return nil
}

// AddLink creates a new link between two nodes.
func (db DB) AddLink(from, to NodeID) error {
	return db.gdb.Save(&Link{from, to}).Error
}

// GetNodes returns all nodes given by the ids.
func (db DB) GetNodes(ids []NodeID) (nodes []Node, err error) {
	for _, l := range ids {
		var n Node
		err = db.gdb.Last(&n, &Node{NodeID: l}).Error
		if err != nil {
			return nil, fmt.Errorf("couldn't get node %x: %+v", l, err)
		}
		nodes = append(nodes, n)
	}
	return
}

// GetChildren searches for nodes that have the given node as ancestor and returns their ids.
func (db DB) GetChildren(from NodeID) (children []NodeID, err error) {
	var links []Link
	err = db.gdb.Find(&links, &Link{From: from}).Error
	for _, l := range links {
		children = append(children, l.To)
	}
	return
}

// GetChildrenNodes searches for nodes that have the given node as ancestor and returns the nodes.
func (db DB) GetChildrenNodes(from NodeID) (children []Node, err error) {
	ids, err := db.GetChildren(from)
	if err != nil {
		return nil, fmt.Errorf("couldn't get ids: %+v", ids)
	}
	return db.GetNodes(ids)
}

// GetAncestors searches for nodes that have the given node as child, and returns their ids.
func (db DB) GetAncestors(to NodeID) (ancestors []NodeID, err error) {
	var links []Link
	err = db.gdb.Find(&links, &Link{To: to}).Error
	for _, l := range links {
		ancestors = append(ancestors, l.From)
	}
	return
}

// GetAncestorsNodes searches for nodes that have the given node as child, and returns their ids.
func (db DB) GetAncestorsNodes(to NodeID) (ancestors []Node, err error) {
	ids, err := db.GetAncestors(to)
	if err != nil {
		return nil, fmt.Errorf("couldn't get ancestor links: %+v", err)
	}
	return db.GetNodes(ids)
}

// GetNodeVersions gets the NodeVersions entry and also fetches all related MemoryLaneEntries.
func (db DB) GetNodeVersions(id NodeID) (nodes []Node, err error) {
	err = db.gdb.Find(&nodes, &Node{NodeID: id}).Error
	if err != nil {
		return nodes, fmt.Errorf("couldn't get NodeVersions: %+v", err)
	}
	return
}

// GetLatest returns the latest version of the node with the given id.
func (db DB) GetLatest(id NodeID) (n Node, err error) {
	nodes, err := db.GetNodeVersions(id)
	if err != nil {
		return n, fmt.Errorf("couldn't get latest node: %+v", err)
	}
	if len(nodes) == 0 {
		return n, errors.New("no node with this id")
	}
	n = nodes[len(nodes)-1]
	return
}
