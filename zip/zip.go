package zip

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Recursively zips the directory
func ArchiveDir(dir string, ignore []string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	defer zw.Close()

	fs, err := walk(dir, ignore)
	if err != nil {
		return nil, err
	}

	// log.Println(fs)

	for _, f := range fs {
		err = addToZip(f, zw)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func walk(start string, ignore []string) ([]string, error) {
	var arr []string

	err := filepath.Walk(start, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fname := info.Name()
		if isInArr(fname, ignore) {
			return filepath.SkipDir
		}

		// Dirs aren't zipped - zip file creates a folder structure
		// automatically if we later specify full paths for files
		if !info.IsDir() {
			arr = append(arr, path)
		}

		// if fname == "main.go" {
		// 	arr = append(arr, path)
		// }

		return nil
	})

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

func isInArr(str string, arr []string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}
