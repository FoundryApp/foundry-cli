package zip

import (
	"archive/zip"
	"bytes"
	"foundry/cli/logger"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gobwas/glob"
)

// Recursively zips the directory
func ArchiveDir(dir string, ignore []glob.Glob) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	defer zw.Close()

	fs, err := walk(dir, ignore)
	if err != nil {
		return nil, err
	}

	for _, f := range fs {
		err = addToZip(f, zw)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func walk(start string, ignore []glob.Glob) ([]string, error) {
	var arr []string

	walkfn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fname := info.Name()
		logger.Fdebugln("")
		logger.Fdebugln("fname in zip:", fname)
		logger.Fdebugln("ignore in zip:", ignore)
		for _, g := range ignore {
			logger.Fdebugln("\t- glob:", g)
			logger.Fdebugln("\t- match:", g.Match(fname))
			// Check if the file name matches the glob pattern
			if g.Match(fname) {
				// If it's a directory, skip the whole directory
				if info.IsDir() {
					logger.Fdebugln("\t- Skipping dir")
					return filepath.SkipDir
				}
				// If it's a file, skip the file by returning nil
				logger.Fdebugln("\t- Skipping file")
				return nil
			}
		}

		// Dirs aren't zipped - zip file creates a folder structure
		// automatically if we later specify full paths for files
		if !info.IsDir() {
			arr = append(arr, path)
		}

		return nil
	}

	err := filepath.Walk(start, walkfn)

	return arr, err
}

func addToZip(fname string, zw *zip.Writer) error {
	f, err := os.Open(fname)
	// fi, err := f.Stat()
	// log.Println("add", fi.Size())
	if err != nil {
		return err
	}
	defer f.Close()

	// Get the file information
	info, err := f.Stat()
	if err != nil {
		return err
	}

	// file -> info header -> edit header -> create hader in the zip using zip writer

	h, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	h.Name = fname

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	h.Method = zip.Deflate

	// Reset time values so they don't influence
	// the checksum of the created zip file
	h.Modified = time.Time{}
	h.ModifiedTime = uint16(0)
	h.ModifiedDate = uint16(0)

	w, err := zw.CreateHeader(h)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, f)
	return err
}
