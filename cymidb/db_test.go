package cymidb

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateDBFile(t *testing.T) {
	db1, err := CreateDBFile(":memory:", "tmp1", "")
	require.NoError(t, err)
	defer db1.Close()
	db2, err := CreateDBFile(":memory:", "tmp2", "")
	require.NoError(t, err)
	defer db2.Close()
	require.NotEqual(t, db1, db2)
}

func TestOpenDBFile(t *testing.T) {
	f, err := ioutil.TempFile("/tmp", "db1")
	require.NoError(t, err)
	f.Close()
	defer os.Remove(f.Name())
	db1, err := CreateDBFile(f.Name(), "tmp", "http://")
	require.NoError(t, err)
	db1.Close()

	db1, err = OpenDBFile(f.Name())
	require.NoError(t, err)
	require.Equal(t, "tmp", db1.Device.Name)
	db1.Close()
}

func TestDB_GetLatest(t *testing.T) {
	f, err := ioutil.TempFile("/tmp", "db")
	require.NoError(t, err)
	f.Close()
	defer os.Remove(f.Name())
	db, err := CreateDBFile(f.Name(), "tmp", "http://")
	require.NoError(t, err)

	db.Device.Name = "temp"
	require.NoError(t, db.SaveNode(db.Device))
	db.Close()

	db, err = OpenDBFile(f.Name())
	require.NoError(t, err)
	require.Equal(t, "temp", db.Device.Name)
	db.Close()
}

func TestDB_Links(t *testing.T) {
	db, err := CreateDBFile(":memory:", "tmp1", "")
	require.NoError(t, err)
	defer db.Close()

	td := NewDataType("blue.gasser/test/blob")

	n1 := NewNode(NodeBlob, Data{td, []byte("blob1")})
	n2 := NewNode(NodeBlob, Data{td, []byte("blob2")})
	require.NoError(t, db.SaveNode(n1))
	require.NoError(t, db.SaveNode(n2))
	require.NoError(t, db.AddLink(n1.NodeID, n2.NodeID))

	children, err := db.GetChildrenNodes(n1.NodeID)
	require.NoError(t, err)
	require.Equal(t, 1, len(children))
	require.NoError(t, n2.Equals(children[0]))

	children, err = db.GetChildrenNodes(n2.NodeID)
	require.NoError(t, err)
	require.Equal(t, 0, len(children))

	ancestors, err := db.GetAncestorsNodes(n2.NodeID)
	require.NoError(t, err)
	require.Equal(t, 1, len(ancestors))
	require.NoError(t, n1.Equals(ancestors[0]))

	ancestors, err = db.GetAncestorsNodes(n1.NodeID)
	require.NoError(t, err)
	require.Equal(t, 0, len(ancestors))
}
