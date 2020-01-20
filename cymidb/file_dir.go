package cymidb

import (
	"errors"
	"fmt"
)

// File_dir implements the File and Dir blobs that can be used to represent files, either on a real file system,
// or files from other sources like email or chats.

// File represents one file. The data itself is stored in a separate blob,
// which can be a virtual blob in the case of a filesystem stored directly on the device itself.
type File struct {
	Name string
	Mask uint16
	Data NodeID
	node Node
}

var NodeTypeFile = NodeBlob.SubType("blue.gasser/cybermind/file")
var DataTypeFile = NewDataType("blue.gasser/cybermind/file")

// FileData is the data of a file. It can be a virtual blob that does only exist on the file system of the device.
// If it's a virtual blob, the NodeID is all 0s.
type FileData struct {
	node Node
	data []byte
}

var NodeTypeFileData = NodeBlob.SubType("blue.gasser/cybermind/filedata")

// DataTypeFileData represents the data of the file.
var DataTypeFileData = NewDataType("blue.gasser/cybermind/filedata")

// Dir holds multiple files and dirs together.
type Dir struct {
	Name string
	Mask uint16
	node Node
}

var NodeTypeDir = NodeBlob.SubType("blue.gasser/cybermind/dir")

// DataTypeDir represents the data of the file.
var DataTypeDir = NewDataType("blue.gasser/cybermind/dir")

func NewFileFromNode(n Node) (f File, err error) {
	err = n.DecodeNodeType(NodeTypeFile, DataTypeFile, &f)
	if err != nil {
		return f, fmt.Errorf("couldn't decode file: %+v", err)
	}
	f.node = n
	return
}

func NewFile(name string, mask uint16) (f File) {
	f.node = NewNode(NodeTypeFile)
	f.Name = name
	f.Mask = mask
	return
}

func (f File) GetNode() (Node, error) {
	err := f.node.EncodeData(DataTypeFile, &f)
	return f.node, err
}

func (f File) AddData(db DB, fd FileData) error {
	return db.AddLink(f.node.NodeID, fd.node.NodeID)
}

func NewFileDataFromNode(n Node) (f FileData, err error) {
	if n.Type != NodeTypeFileData {
		return f, errors.New("node not of type FileData")
	}
	f.node = n
	return
}

func NewFileData(data []byte) (fd FileData) {
	fd.node = NewNode(NodeTypeFileData)
	fd.data = data
	return
}

func (f FileData) GetNode() (Node, error) {
	if f.data != nil {
		if err := f.node.SetDatas(Data{Type: DataTypeFileData, Data: f.data}); err != nil {
			return f.node, fmt.Errorf("couldn't set node data")
		}
	}
	return f.node, nil
}

func NewDirFromNode(n Node) (d Dir, err error) {
	err = n.DecodeNodeType(NodeTypeDir, DataTypeDir, &d)
	if err != nil {
		return d, fmt.Errorf("couldn't decode dir: %+v", err)
	}
	d.node = n
	return
}

func NewDir(name string, mask uint16) (d Dir) {
	d.node = NewNode(NodeTypeDir)
	d.Name = name
	d.Mask = mask
	return
}

func (d Dir) GetNode() (Node, error) {
	err := d.node.EncodeData(DataTypeDir, &d)
	return d.node, err
}

func (d Dir) AddSubdir(db DB, sd Dir) error {
	return db.AddLink(d.node.NodeID, sd.node.NodeID)
}

func (d Dir) AddFile(db DB, f File) error {
	return db.AddLink(d.node.NodeID, f.node.NodeID)
}

func (d Dir) GetDirs(db DB) (dirs []Dir, err error) {
	children, err := db.GetChildrenNodes(d.node.NodeID)
	if err != nil {
		return nil, fmt.Errorf("couldn't get children: %v", err)
	}
	for _, child := range children {
		switch child.Type {
		case NodeTypeDir:
			dir, err := NewDirFromNode(child)
			if err != nil {
				return nil, fmt.Errorf("couldn't get child as dir: %v", err)
			}
			dirs = append(dirs, dir)
		default:
			// Ignoring non-directories
		}
	}
	return
}

func (d Dir) GetFiles(db DB) (files []File, err error) {
	children, err := db.GetChildrenNodes(d.node.NodeID)
	if err != nil {
		return nil, fmt.Errorf("couldn't get children: %v", err)
	}
	for _, child := range children {
		switch child.Type {
		case NodeTypeFile:
			file, err := NewFileFromNode(child)
			if err != nil {
				return nil, fmt.Errorf("couldn't get child as file: %v", err)
			}
			files = append(files, file)
		default:
			// Ignoring non-files
		}
	}
	return
}
