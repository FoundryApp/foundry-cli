package connection

import (
	"fmt"
	"time"

	"foundry/cli/logger"
	"foundry/cli/connection/msg"

	"github.com/gorilla/websocket"
)

type ListenCallback func(data []byte, err error) ()

type Connection struct {
	token string
	wsconn 	*websocket.Conn
}

type ConnectionMessage interface {
	Body() interface{}
}

const (
	baseURL = "127.0.0.1:8000" // autorun
	// baseURL = "127.0.0.1:3500" // podm
  // baseURL = "ide.foundryapp.co"

  wsScheme = "ws"
  // wsScheme = "wss"

  pingScheme = "http"
  // pingScheme = "https"
)


func New(token string) (*Connection, error) {
	logger.Debugln("<Connection> WS dialing")
	url := WebSocketURL(token)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	logger.Debugln("<Connection> WS connected")

	return &Connection{token, c}, nil
}

func (c *Connection) Close() {
	logger.Debugln("<Connection> WS closing")
	c.wsconn.Close()
}

func (c *Connection) Listen(cb ListenCallback) {
	logger.Debugln("<Connection> WS listening")

	for {
		_, msg, err := c.wsconn.ReadMessage()
		cb(msg, err)
	}
}

// Sends WS message
func (c *Connection) Send(cm ConnectionMessage) error {
	b := cm.Body()
	err := c.wsconn.WriteJSON(b)
  if err != nil {
    return err
	}
	return nil
}

// Pings server so the WS connection stays open
func (c *Connection) Ping(pm *msg.PingMsg, ticker *time.Ticker, stop <-chan struct{}) {
	logger.Debugln("<Ping> Pinging")
	for {
		select {
		case <-ticker.C:
			if err := pm.Send(); err != nil {
				logger.Debugln("<Ping> Failed to ping server", err)
			}
		case <-stop:
			logger.Debugln("<Ping> Stopping ping")
			ticker.Stop()
			return
		}
	}
}

func WebSocketURL(token string) string {
	return fmt.Sprintf("%s://%s/ws/%s", wsScheme, baseURL, token)
}

func PingURL() string {
	return fmt.Sprintf("%s://%s/ping", pingScheme, baseURL)
}