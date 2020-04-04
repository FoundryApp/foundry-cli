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

var (
	buf = new(bytes.Buffer)
)

// Recursively zips the directory
func ArchiveDir(dir string, ignore []glob.Glob) (*bytes.Buffer, error) {
	buf.Reset()
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
		logger.Fdebugln("")
		logger.Fdebugln("path in zip", path)

		// Prepend path with the "./" so the prefix
		// is same as the ignore array in the config
		// file.
		// TODO: Should the prefix be foundryConf.RootDir?
		// path = "." + string(os.PathSeparator) + path

		if err != nil {
			// TODO: If path is in the ignored array, should we ignore the error?
			if ignored(path, ignore) {
				return nil
			}
			logger.FdebuglnError("Zip walk error - path", path)
			logger.FdebuglnError("Zip walk error - error", err)
			return err
		}

		if ignored(path, ignore) {
			// If it's a directory, skip the whole directory
			if info.IsDir() {
				logger.Fdebugln("\t- Skipping dir")
				return filepath.SkipDir
			}
			// If it's a file, skip the file by returning nil
			logger.Fdebugln("\t- Skipping file")
			return nil
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

func ignored(s string, globs []glob.Glob) bool {
	logger.Fdebugln("string to match:", s)
	for _, g := range globs {
		logger.Fdebugln("\t- glob:", g)
		logger.Fdebugln("\t- match:", g.Match(s))
		if g.Match(s) {
			return true
		}
	}
	return false
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
