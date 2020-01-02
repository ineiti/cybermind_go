package cymidb

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// MemoryLaneEntry holds one node, which includes the version of this node.
type MemoryLaneEntry struct {
	gorm.Model
	NodeBuf        []byte
	NodeVersionsID uint
}

// NodeVersions is stored as the ID of the node and points to all stored versions in the MemoryLane.
type NodeVersions struct {
	gorm.Model
	NodeID   []byte `gorm:"type:binary(32)"`
	Versions []MemoryLaneEntry
}

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
	db.gdb.AutoMigrate(&MemoryLaneEntry{}, &NodeVersions{})
	return
}

// CreateDBFile creates a new DB in a given file and returns an initialized DB containing only the given device.
func CreateDBFile(file string, name, url string) (db DB, err error) {
	db, err = NewDBFile(file)
	if err != nil {
		return db, err
	}

	db.Device = NewDevice(name)
	err = db.NewNode(db.Device.Node)
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

	mle := MemoryLaneEntry{}
	err = db.gdb.First(&mle).Error
	if err != nil {
		return db, fmt.Errorf("couldn't get first ml-entry: %+v", err)
	}
	n, err := mle.Node()
	if err != nil {
		return db, fmt.Errorf("couldn't get node: %+v", err)
	}
	node, err := db.GetLatest(n.Hash)
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
func (db DB) NewNode(n Node) error {
	mle, err := NewMemoryLaneEntry(n)
	if err != nil {
		return fmt.Errorf("couldn't create memoryLaneEntry: %+v", err)
	}
	err = db.gdb.Create(&NodeVersions{
		NodeID:   n.Hash,
		Versions: []MemoryLaneEntry{mle},
	}).Error
	if err != nil {
		return fmt.Errorf("couldn't create new node: %+v", err)
	}
	return nil
}

// GetLatest returns the latest version of the node with the given id.
func (db DB) GetLatest(id []byte) (n Node, err error) {
	ni := NodeVersions{}
	err = db.gdb.Where(&NodeVersions{NodeID: id}).First(&ni).Error
	if err != nil {
		return n, fmt.Errorf("couldn't get NodeVersions: %+v", err)
	}
	err = db.gdb.Model(&ni).Related(&ni.Versions).Error
	if err != nil {
		return n, fmt.Errorf("couldn't get versions: %+v", err)
	}
	if len(ni.Versions) == 0 {
		return n, errors.New("no versions stored in nodeIndex")
	}
	return ni.Versions[len(ni.Versions)-1].Node()
}

// NewMemoryLaneEntry creates a
func NewMemoryLaneEntry(n Node) (mle MemoryLaneEntry, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(&n)
	if err != nil {
		return mle, fmt.Errorf("couldn't encode node: %+v", err)
	}
	mle.NodeBuf = buf.Bytes()
	return
}

// Returns the node stored in the MemoryLaneEntry
func (ml MemoryLaneEntry) Node() (n Node, err error) {
	buf := bytes.NewBuffer(ml.NodeBuf)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&n)
	if err != nil {
		return n, fmt.Errorf("couldn't decode node: %+v", err)
	}
	return
}