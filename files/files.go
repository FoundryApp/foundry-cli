package files

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"regexp"

	conn "foundry/cli/connection"
	connMsg "foundry/cli/connection/msg"
	"foundry/cli/logger"

	"foundry/cli/zip"
)

var (
	lastArchiveChecksum = ""
)

func Upload(c *conn.Connection, rootDir string, ignore ...*regexp.Regexp) {
	// Zip the project
	buf, err := zip.ArchiveDir(rootDir, ignore)
	if err != nil {
		logger.FdebuglnFatal("ArchiveDir error", err)
		logger.ErrorLoglnFatal("Error archiving the directorye", err)
	}

	archiveChecksum := checksum(buf.Bytes())

	// TODO: REMOVE
	// if lastArchiveChecksum == archiveChecksum { return }
	lastArchiveChecksum = archiveChecksum

	bufferSize := 1024 // 1024B, size of a single chunk
	buffer := make([]byte, bufferSize)
	chunkCount := (buf.Len() / bufferSize) + 1

	checksum := [md5.Size]byte{}
	previousChecksum := [md5.Size]byte{}

	for i := 0; i < chunkCount; i++ {
		bytesread, err := buf.Read(buffer)
		// TODO: HEEEEEEEEEEEEEEEEEEEEEEEEEEEREE EOF
		if err != nil && err != io.EOF {
			logger.FdebuglnFatal("Error reading chunk from buffer", err)
			logger.ErrorLoglnFatal("Error reading chunk from buffer", err)
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
			logger.ErrorLoglnFatal("Error sending chunk", err)
		}
	}
}

func checksum(data []byte) string {
	hashInBytes := md5.Sum(data)
	return hex.EncodeToString(hashInBytes[:])
}
