package files

import (
	"crypto/md5"
	"encoding/hex"
	"io"

	conn "foundry/cli/connection"
	connMsg "foundry/cli/connection/msg"
	"foundry/cli/logger"

	"foundry/cli/zip"

	"github.com/gobwas/glob"
)

var (
	lastArchiveChecksum = ""
)

func Upload(c *conn.Connection, rootDir, serviceAccPath string, promptNotifCh chan<- string, ignore ...glob.Glob) {
	// Zip the project in the memory and send the file in chunks
	buf, err := zip.ArchiveDir(rootDir, serviceAccPath, ignore)
	if err != nil {
		logger.FdebuglnFatal("ArchiveDir error:", err)
		logger.FatalLogln("Error uploading files:", err)
	}

	// err = ioutil.WriteFile("./source.zip", buf.Bytes(), 0644)
	// logger.FatalLogln("Written", err)

	archiveChecksum := checksum(buf.Bytes())

	if lastArchiveChecksum == archiveChecksum {
		promptNotifCh <- "No change in the code detected. Make change to upload the code."
		return
	}
	lastArchiveChecksum = archiveChecksum

	bufferSize := 1024 // 1024B, size of a single chunk
	buffer := make([]byte, bufferSize)
	chunkCount := (buf.Len() / bufferSize) + 1

	checksum := [md5.Size]byte{}
	previousChecksum := [md5.Size]byte{}

	for i := 0; i < chunkCount; i++ {
		bytesread, err := buf.Read(buffer)
		// TODO: Why did this work without err != io.EOF?
		if err != nil && err != io.EOF {
			logger.FdebuglnFatal("Error reading chunk from buffer:", err)
			logger.FatalLogln("Error reading chunk from buffer:", err)
		}

		previousChecksum = checksum
		bytes := buffer[:bytesread]
		checksum = md5.Sum(bytes)

		checkStr := hex.EncodeToString(checksum[:])
		prevCheckStr := hex.EncodeToString(previousChecksum[:])

		lastChunk := i == chunkCount-1

		chunk := connMsg.NewChunkMsg(bytes, checkStr, prevCheckStr, lastChunk)
		if err = c.Send(chunk); err != nil {
			logger.FdebuglnFatal("Error sending chunk", err)
			logger.FatalLogln("Error sending chunk", err)
		}
	}
}

func checksum(data []byte) string {
	hashInBytes := md5.Sum(data)
	return hex.EncodeToString(hashInBytes[:])
}
