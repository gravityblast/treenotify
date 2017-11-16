package treenotify

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	assert "github.com/pilu/miniassert"
)

func mkdirAll(root, path string) string {
	fullPath := filepath.Join(root, path)
	err := os.MkdirAll(fullPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	return fullPath
}

func createFile(path, name string) string {
	fullPath := filepath.Join(path, name)
	file, err := os.Create(fullPath)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	return fullPath
}

func checkEvent(t *testing.T, events chan fsnotify.Event, name string) {
	select {
	case event := <-events:
		assert.Equal(t, name, event.Name)
	case <-time.After(time.Second):
		t.Error("event not fired for %v", name)
		t.Fail()
	}
}

func TestWatcher(t *testing.T) {
	tmpRoot := filepath.Join(os.TempDir(), "golang.watcher.tests")
	defer os.RemoveAll(tmpRoot)

	w, err := New()
	assert.Nil(t, err)
	defer func() {
		w.Close()
	}()

	rootPath := mkdirAll(tmpRoot, "foo")
	subPath := mkdirAll(tmpRoot, "foo/bar")

	events, err := w.Watch(rootPath)
	assert.Nil(t, err)

	// new file inside the root
	file1 := createFile(rootPath, "foo.txt")
	checkEvent(t, events, file1)

	// new file inside sub-folder
	file2 := createFile(subPath, "foo.txt")
	checkEvent(t, events, file2)

	// new folder inside sub-folder
	subSubPath := mkdirAll(tmpRoot, "foo/bar/baz")
	checkEvent(t, events, subSubPath)

	// new file inside the folder created after watching
	file3 := createFile(subSubPath, "foo.txt")
	checkEvent(t, events, file3)
}
