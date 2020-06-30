package rwatch

import (
	"foundry/cli/logger"
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
	return w, nil
}

func (w *Watcher) AddRecursive(dir string) error {
	return w.traverse(dir, true)
}

func (w *Watcher) Close() {
	logger.Fdebugln("Closing rwatch")
	w.fsnotify.Close()
	close(w.done)
}

func (w *Watcher) Watch() {
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
		logger.Fdebugln("")
		logger.Fdebugln("path in rwatch", path)

		// Prepend path with the "./" so the prefix
		// is same as the ignore array in the config
		// file.
		// TODO: Should the prefix be foundryConf.RootDir?
		// path = "." + string(os.PathSeparator) + path

		if err != nil {
			// TODO: If path is in the ignored array, should we ignore the error?
			// Note: we can't use info.IsDir() here because of the error - the file
			// might not even exist. Using info.IsDir() would cause panic
			logger.FdebuglnError("rwatch walk error - path", path)
			logger.FdebuglnError("rwatch walk error - error", err)
			if w.ignored(path) {
				return nil
			}
			return err
		}

		isIgnored := w.ignored(path)
		logger.Fdebugln("is ignored?", isIgnored)

		if isIgnored {
			// If it's a directory, skip the whole directory
			if info.IsDir() {
				logger.Fdebugln("\t- Skipping dir")
				// No need to remove watcher on an ignored dir because watcher isn't recursive
				// i.e.: if we have following folder structure:
				// rootDir/
				//	file1
				//	subDir/
				//		file2
				// and when we add rootDir to the watcher, the subDir isn't added
				return filepath.SkipDir
			}

			// Always remove watcher on an ignored file because the file could be in a folder that is watched
			logger.Fdebugln("\t- Skipping file (removing watch)")
			return w.fsnotify.Remove(path)
		}

		if watch && !isIgnored {
			logger.Fdebugln("\t- Adding file/dir to rwatch")
			return w.fsnotify.Add(path)
		} else if !watch {
			logger.Fdebugln("\t- Removing file/dir from rwatch")
			return w.fsnotify.Remove(path)
		}

		return nil
	}
	return filepath.Walk(start, walkfn)
}

func (w *Watcher) ignored(s string) bool {
	logger.Fdebugln("string to match:", s)
	for _, g := range w.ignore {
		logger.Fdebugln("\t- glob:", g)
		logger.Fdebugln("\t- match:", g.Match(s))
		if g.Match(s) {
			return true
		}
	}
	return false
}
