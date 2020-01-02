package cymidb

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDBFile(t *testing.T) {
	db1, err := CreateDBFile(":memory:", "tmp1", "")
	require.NoError(t, err)
	defer db1.Close()
	db2, err := CreateDBFile(":memory:", "tmp2", "")
	require.NoError(t, err)
	defer db2.Close()
	require.NotEqual(t, db1, db2)
}

func TestNewDBFile2(t *testing.T) {
	f, err := ioutil.TempFile("/tmp", "db1")
	require.NoError(t, err)
	f.Close()
	db1, err := CreateDBFile(f.Name(), "tmp", "http://")
	require.NoError(t, err)
	db1.Close()

	db1, err = OpenDBFile(f.Name())
	require.NoError(t, err)
	require.Equal(t, "tmp", db1.Device.Name)
	db1.Close()
}
