package cymidb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHook(t *testing.T) {
	db, err := CreateDBFile(":memory:", "tmp", "")
	require.NoError(t, err)
	defer db.Close()

	n, err := db.Device.GetNode()
	require.NoError(t, err)
	_, err = NewHookFromNode(n)
	require.Error(t, err)

	h := NewHook("hook")
	n, err = h.GetNode()
	require.NoError(t, err)
	h2, err := NewHookFromNode(n)
	require.NoError(t, err)
	require.Equal(t, h.Name, h2.Name)
}
