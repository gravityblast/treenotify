package treenotify

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func New() (Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &watcher{
		w:      fsw,
		stop:   make(chan struct{}, 1),
		events: make(chan fsnotify.Event),
	}

	return w, nil
}

type Watcher interface {
	Watch(string) (chan fsnotify.Event, error)
	Close()
}

type watcher struct {
	w      *fsnotify.Watcher
	stop   chan struct{}
	events chan fsnotify.Event
}

func (w *watcher) Watch(root string) (chan fsnotify.Event, error) {
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		w.w.Add(path)
		return err
	})

	if err != nil {
		return nil, err
	}

	go w.watch()

	return w.events, nil
}

func (w *watcher) watch() {
	working := true
	for working {
		select {
		case <-w.stop:
			working = false
		case event := <-w.w.Events:
			if event.Op == fsnotify.Create && isDir(event.Name) {
				w.w.Add(event.Name)
			}
			w.events <- event
		}
	}
}

func (w *watcher) Close() {
	w.stop <- struct{}{}
	w.w.Close()
}

func (w *watcher) Events() chan fsnotify.Event {
	return w.events
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fi.Mode().IsDir()
}
