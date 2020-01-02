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
	db.Device.UpdateNode()
	require.NoError(t, db.UpdateNode(db.Device.Node))
	db.Close()

	db, err = OpenDBFile(f.Name())
	require.NoError(t, err)
	require.Equal(t, "temp", db.Device.Name)
	db.Close()
}
