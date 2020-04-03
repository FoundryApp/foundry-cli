package rwatch

import (
	"foundry/cli/logger"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/gobwas/glob"
)

type Watcher struct {
	Events chan fsnotify.Event
	Errors chan error

	fsnotify *fsnotify.Watcher
	done     chan struct{}

	ignore []glob.Glob
}

// var ignore = []string{".git", "node_modules", ".foundry"}

func New(ignore []glob.Glob) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	w := &Watcher{
		fsnotify: fsw,
		Events:   make(chan fsnotify.Event),
		Errors:   make(chan error),
		done:     make(chan struct{}),
		ignore:   ignore,
	}

	go w.start()
	return w, nil
}

func (w *Watcher) AddRecursive(dir string) error {
	return w.traverse(dir, true)
}

func (w *Watcher) Close() {
	log.Println("Closing rwatch")
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
// We don't care for files, only for directories because we are watching whole dirs
func (w *Watcher) traverse(start string, watch bool) error {
	walkfn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			fname := info.Name()
			logger.Fdebugln("")
			logger.Fdebugln("fname in rwatch:", fname)
			logger.Fdebugln("ignore in rwatch:", w.ignore)

			for _, g := range w.ignore {
				logger.Fdebugln("\t- glob:", g)
				logger.Fdebugln("\t- match:", g.Match(fname))
				if g.Match(fname) {
					logger.Fdebugln("\t- Skipping dir")
					return filepath.SkipDir
				}
			}

			if watch {
				logger.Fdebugln("\t- Adding dir to rwatch")
				return w.fsnotify.Add(path)
			}
			logger.Fdebugln("\t- Removing dir from rwatch")
			return w.fsnotify.Remove(path)
		}
		return nil
	}
	return filepath.Walk(start, walkfn)
}
