package files

import (
	"crypto/md5"
	"encoding/hex"
	"log"

	"foundry/cli/zip"

	"github.com/gorilla/websocket"
)

var (
	ignore = []string{"node_modules", ".git"}
	lastArchiveChecksum = ""
)

func Upload(c *websocket.Conn, rootDir string) {
  // Zip the project
  buf, err := zip.ArchiveDir(rootDir, ignore)
  if err != nil {
    log.Fatal(err)
  }

  // Get the 16 bytes hash
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
    if err != nil {
      log.Fatal(err)
    }

    previousChecksum = checksum
    bytes := buffer[:bytesread]
    checksum = md5.Sum(bytes)

    checkStr := hex.EncodeToString(checksum[:])
    prevCheckStr := hex.EncodeToString(previousChecksum[:])

    lastChunk := i == chunkCount - 1

    if err = sendChunk(
      c,
      bytes,
      checkStr,
      prevCheckStr,
      lastChunk); err != nil {
      log.Fatal(err)
    }
  }
}

func sendChunk(c *websocket.Conn, b []byte, checksum string, prevChecksum string, last bool) error {
  msg := struct {
    Data              string   `json:"data"`
    PreviousChecksum  string   `json:"previousChecksum"`
    Checksum          string   `json:"checksum"`
    IsLast            bool     `json:"isLast"`
    RunAll            bool     `json:"runAll"`
    Run               []string `json:"run"`
  }{hex.EncodeToString(b),
    prevChecksum,
    checksum,
    last,
    true,
    []string{},
  }
  err := c.WriteJSON(msg)
  if err != nil {
    return err
  }
  return nil
}

func checksum(data []byte) string {
  hashInBytes := md5.Sum(data)
  return hex.EncodeToString(hashInBytes[:])
}
