package msg

import (
	"encoding/hex"

	// "foundry/cli/logger"
)

type ChunkBody struct {
	Data              string   `json:"data"`
	PreviousChecksum  string   `json:"previousChecksum"`
	Checksum          string   `json:"checksum"`
	IsLast            bool     `json:"isLast"`
	RunAll            bool     `json:"runAll"`
	Run               []string `json:"run"`
}

type ChunkMsg struct {
	MsgBody	ChunkBody
}

func NewChunkMsg(b []byte, checksum, checksumPrev string, last bool) *ChunkMsg {
	return &ChunkMsg{
		MsgBody: ChunkBody{
			hex.EncodeToString(b),
			checksumPrev,
			checksum,
			last,
			true,
			[]string{},
		},
	}
}

func (cm *ChunkMsg) Body() interface{} {
	return cm.MsgBody
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