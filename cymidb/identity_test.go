package cymidb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewIdentity(t *testing.T) {
	db, err := CreateDBFile(":memory:", "tmp", "")
	require.NoError(t, err)
	defer db.Close()

	ident, err := NewIdentity("test", []string{"one@test.com", "two@test.com"})
	require.NoError(t, err)
	require.NoError(t, db.SaveNode(ident))
	require.NoError(t, db.AddLink(db.Device.node.NodeID, ident.node.NodeID))

	// Retrieve all identities
	ids, err := db.GetChildren(db.Device.node.NodeID)
	require.NoError(t, err)
	require.Equal(t, 1, len(ids))
	require.Equal(t, ident.node.NodeID, ids[0])

	nodes, err := db.GetNodes(ids)
	ident2, err := NewIdentityFromNode(nodes[0])
	require.NoError(t, err)
	require.NoError(t, ident.Equals(ident2))
}
