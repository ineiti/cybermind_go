package cymidb

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFileHook(t *testing.T) {
	root, err := ioutil.TempDir("", "hooks")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(root))
	}()

	require.NoError(t, ioutil.WriteFile(path.Join(root, ".gitconfig"), []byte("[remote]"), 0770))
	require.NoError(t, os.Mkdir(path.Join(root, "Documents"), 0770))
	require.NoError(t, ioutil.WriteFile(path.Join(root, "Documents", "README.md"), []byte("very important"), 0770))
	require.NoError(t, os.Mkdir(path.Join(root, "Empty"), 0770))

	db, err := CreateDBFile(":memory:", "tmp", "")
	require.NoError(t, err)
	defer db.Close()

	rootDir := NewDir("/", 0777)
	require.NoError(t, db.SaveNode(rootDir))

	hook := NewHook("filer", nil, []NodeType{NodeTypeFile, NodeTypeDir, NodeTypeFileData})
	require.NoError(t, db.SaveNode(hook))
	require.NoError(t, db.AddLink(rootDir, hook))

	pfh, err := NewPluginFileHook(hook, root, PFHOUpdateDB|PFHOUpdateFS)
	require.NoError(t, err)
	log.Print(pfh)
}
