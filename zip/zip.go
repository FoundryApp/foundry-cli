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
func ArchiveDir(rootDir, serviceAccPath string, ignore []glob.Glob) (*bytes.Buffer, error) {
	buf.Reset()
	zw := zip.NewWriter(buf)
	defer zw.Close()

	// Walk all dirs inside the dir and
	// return all file paths (also walks
	// all subdirs)
	fPaths, err := walk(rootDir, ignore)
	if err != nil {
		return nil, err
	}

	for _, fPath := range fPaths {
		err = addToZip(rootDir, fPath, zw)
		if err != nil {
			return nil, err
		}
	}

	// Zip service account - service account
	// might be in a completely different dir.
	// That's why we are adding it separately
	if serviceAccPath != "" {
		err = addServiceAccToZip(serviceAccPath, zw)
		if err != nil {
			return nil, err
		}
	}
	return buf, nil
}

func walk(start string, ignore []glob.Glob) ([]string, error) {
	var filePaths []string

	walkfn := func(path string, info os.FileInfo, err error) error {
		logger.Fdebugln("")
		logger.Fdebugln("path in zip", path)

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
			filePaths = append(filePaths, path)
		}

		return nil
	}

	err := filepath.Walk(start, walkfn)

	return filePaths, err
}

func addToZip(rootDir, fPath string, zw *zip.Writer) error {
	fileToZip, err := os.Open(fPath)
	// fi, err := f.Stat()
	// log.Println("add", fi.Size())
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	// file -> info header -> edit header -> create hader in the zip using zip writer

	h, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we want to get a relative path of the fPath
	// to the current working directory
	relativeFilePath, err := filepath.Rel(rootDir, fPath)
	if err != nil {
		return err
	}
	h.Name = relativeFilePath

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	h.Method = zip.Deflate

	// Reset time values so they don't influence
	// the checksum of the created zip file
	h.Modified = time.Time{}
	h.ModifiedTime = uint16(0)
	h.ModifiedDate = uint16(0)

	headerWriter, err := zw.CreateHeader(h)
	if err != nil {
		return err
	}

	_, err = io.Copy(headerWriter, fileToZip)
	return err
}

// Everything is same as addToZip besides  preserving the
// serviceAcc's file path structure. We want to get only
// the last part of the serviceAcc's path so it's in the
// root of the zip file
func addServiceAccToZip(fPath string, zw *zip.Writer) error {
	fileToZip, err := os.Open(fPath)
	// fi, err := f.Stat()
	// log.Println("add", fi.Size())
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	// file -> info header -> edit header -> create hader in the zip using zip writer

	h, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// We want to add service account into the root of the zip file
	// Therefore we take only the last part (= file name) of its path
	_, fName := filepath.Split(fPath)
	h.Name = fName

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	h.Method = zip.Deflate

	// Reset time values so they don't influence
	// the checksum of the created zip file
	h.Modified = time.Time{}
	h.ModifiedTime = uint16(0)
	h.ModifiedDate = uint16(0)

	headerWriter, err := zw.CreateHeader(h)
	if err != nil {
		return err
	}

	_, err = io.Copy(headerWriter, fileToZip)
	return err
	return nil
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
