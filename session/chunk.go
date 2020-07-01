package session

import (
	"encoding/hex"
)

type ChunkContent struct {
	Data             string `json:"data"`
	PreviousChecksum string `json:"previousChecksum"`
	Checksum         string `json:"checksum"`
	IsLast           bool   `json:"isLast"`
}

type ChunkMsg struct {
	Type    string       `json:"type"`
	Content ChunkContent `json:"content"`
}

func NewChunkMsg(b []byte, checksum, checksumPrev string, last bool) *ChunkMsg {
	c := ChunkContent{hex.EncodeToString(b), checksumPrev, checksum, last}
	return &ChunkMsg{"chunk", c}
}

func (cm *ChunkMsg) Body() interface{} {
	return cm
}

// msg := struct {
//   Data              string   `json:"data"`
//   PreviousChecksum  string   `json:"previousChecksum"`
//   Checksum          string   `json:"checksum"`
//   IsLast            bool     `json:"isLast"`
//   RunAll            bool     `json:"runAll"`
//   Run               []string `json:"run"`
// }{hex.EncodeToString(b),
//   prevChecksum,
//   checksum,
//   last,
//   true,
//   []string{},
// }
// err := c.WriteJSON(msg)
// if err != nil {
//   return err
// }
