package cymidb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNewFS sets up a simple directory and links all nodes together.
func TestNewFS(t *testing.T) {
	db, err := CreateDBFile(":memory:", "tmp", "")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, db.Close())
	}()

	rootDir := NewDir("/", 0777)
	docDir := NewDir("Documents", 0777)
	require.NoError(t, rootDir.AddSubdir(db, docDir))
	emptyDir := NewDir("Empty", 0777)
	require.NoError(t, rootDir.AddSubdir(db, emptyDir))

	gitConfig := NewFile(".gitignore", 0777)
	require.NoError(t, rootDir.AddFile(db, gitConfig))
	gitConfigData := NewFileData([]byte("[alias]"))
	require.NoError(t, gitConfig.AddData(db, gitConfigData))

	todo := NewFile("TODO.md", 0777)
	require.NoError(t, docDir.AddFile(db, todo))
	todoData := NewFileData([]byte("Finish Project"))
	require.NoError(t, todo.AddData(db, todoData))

	require.NoError(t, db.SaveNode(rootDir, docDir, emptyDir,
		gitConfig, gitConfigData, todo, todoData))

	subdirs, err := rootDir.GetDirs(db)
	require.NoError(t, err)
	require.Equal(t, 2, len(subdirs))
	for _, sd := range subdirs {
		switch sd.Name {
		case docDir.Name:
			require.NoError(t, docDir.node.CompareTo(sd.node))
		case emptyDir.Name:
			require.NoError(t, emptyDir.node.CompareTo(sd.node))
		default:
			require.Fail(t, "found unknown directory")
		}
	}

	subdirs, err = docDir.GetDirs(db)
	require.NoError(t, err)
	require.Equal(t, 0, len(subdirs))

	subdirs, err = emptyDir.GetDirs(db)
	require.NoError(t, err)
	require.Equal(t, 0, len(subdirs))

	files, err := rootDir.GetFiles(db)
	require.NoError(t, err)
	require.Equal(t, 1, len(files))
	require.NoError(t, gitConfig.node.CompareTo(files[0].node))

	files, err = docDir.GetFiles(db)
	require.NoError(t, err)
	require.Equal(t, 1, len(files))
	require.NoError(t, todo.node.CompareTo(files[0].node))

	files, err = emptyDir.GetFiles(db)
	require.NoError(t, err)
	require.Equal(t, 0, len(files))
}
