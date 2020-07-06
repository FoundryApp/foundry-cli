package session

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"foundry/cli/logger"

	"github.com/gorilla/websocket"
)

type ListenCallback func(data []byte, err error)

func (sess *Session) Connect(token string) error {
	c, _, err := websocket.DefaultDialer.Dial(sess.webSocketURL(token), nil)
	if err != nil {
		return err
	}
	sess.wsconn = c
	return nil
}

func (sess *Session) Close() error {
	if sess.wsconn == nil {
		return errors.New("No active WebSocket connection")
	}
	return sess.wsconn.Close()
}

func (sess *Session) Listen(data chan<- []byte, err chan<- error) {
	if sess.wsconn == nil {
		err <- errors.New("No active WebSocket connection")
	}

	for {
		_, msg, readErr := sess.wsconn.ReadMessage()
		if readErr != nil {
			err <- readErr
		} else {
			data <- msg
		}
	}
}

func (sess *Session) SendData(buf *bytes.Buffer) error {
	// 1024B is the size of a single chunk
	chunkSize := 1024
	chunkBuffer := make([]byte, chunkSize)
	chunkCount := (buf.Len() / chunkSize) + 1

	chsum := ""
	prevchsum := ""

	for i := 0; i < chunkCount; i++ {
		bytesread, err := buf.Read(chunkBuffer)
		if err != nil && err != io.EOF {
			logger.FdebuglnError("send data buffer read error:", err)
			return err
		}

		bytes := chunkBuffer[:bytesread]
		prevchsum = chsum
		chsum = checksum(bytes)

		isLast := i == chunkCount-1
		chunk := NewChunkMsg(bytes, chsum, prevchsum, isLast)
		if err := sess.write(chunk); err != nil {
			logger.FdebuglnError("send data write chunk error:", err)
			return err
		}
	}

	return nil
}

func (sess *Session) write(ch *ChunkMsg) error {
	b := ch.Body()
	if err := sess.wsconn.WriteJSON(b); err != nil {
		return err
	}
	return nil
}

func (sess *Session) webSocketURL(token string) string {
	return fmt.Sprintf("%s://%s%s", webSocketHost, sess.WebSocketHost, token)
}

func checksum(data []byte) string {
	hashInBytes := md5.Sum(data)
	return hex.EncodeToString(hashInBytes[:])
}
