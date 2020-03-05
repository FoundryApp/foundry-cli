package rwatch

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	Events		chan fsnotify.Event
	Errors		chan error

	fsnotify 	*fsnotify.Watcher

	done			chan struct{}
}

var ignore = []string{".git", "node_modules", ".foundry"}

func New() (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	w := &Watcher{
		fsnotify: fsw,
		Events: make(chan fsnotify.Event),
		Errors: make(chan error),
		done:		make(chan struct{}),
	}

	go w.start()
	return w, nil
}

func (w *Watcher) AddRecursive(dir string) error {
	if err := w.traverse(dir, true); err != nil {
		return err
	}
	return nil
}

func (w *Watcher) Close() {
	log.Println("Closing")
	w.fsnotify.Close()
	close(w.done)
}

func (w *Watcher) start() {
	for {
		select {
		case ev := <-w.fsnotify.Events:
			fi, err := os.Stat(ev.Name)
			if err == nil && fi != nil && fi.IsDir() {
				if ev.Op == fsnotify.Create {
					if err = w.traverse(ev.Name, true); err != nil {
						w.Errors <- err
					}
				}
			}

			// os.Stat() can't be used on deleted dir/file
			// Pretend it was a directory (we don't really know)
			// and try to remove it
			if ev.Op == fsnotify.Remove {
				w.fsnotify.Remove(ev.Name)
			}

			if ev.Op != fsnotify.Chmod {
				w.Events <- ev
			}
		case err := <-w.fsnotify.Errors:
			w.Errors <- err

		case <-w.done:
			close(w.Events)
			close(w.Errors)
			return
		}
	}
}

// Traverses the root directory and adds watcher for each directory along the way
func (w *Watcher) traverse(start string, watch bool) error {
	err := filepath.Walk(start, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if dirIgnored(info) {
				return filepath.SkipDir
			}

			if watch {
				if err = w.fsnotify.Add(path); err != nil {
					return err
				}
			} else {
				if err = w.fsnotify.Remove(path); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}


func dirIgnored(fi os.FileInfo) bool {
	n := fi.Name()
	for _, i := range ignore {

		if i == n {
			return true
		}
	}
	return false
}
