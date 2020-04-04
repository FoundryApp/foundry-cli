package connection

import (
	"fmt"
	"time"

	"foundry/cli/connection/endpoint"
	"foundry/cli/connection/msg"
	"foundry/cli/logger"

	"github.com/gorilla/websocket"
)

type ListenCallback func(data []byte, err error)

type Connection struct {
	token  string
	wsconn *websocket.Conn
}

type ConnectionMessage interface {
	Body() interface{}
}

func New(token string) (*Connection, error) {
	logger.Fdebugln("WS dialing")
	url := WebSocketURL(token)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	logger.Fdebugln("WS connected")

	return &Connection{token, c}, nil
}

func (c *Connection) Close() {
	logger.Fdebugln("WS closing")
	c.wsconn.Close()
}

func (c *Connection) Listen(cb ListenCallback) {
	logger.Fdebugln(" WS listening")
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
	logger.Fdebugln("Ping")
	for {
		select {
		case <-ticker.C:
			if err := pm.Send(); err != nil {
				logger.FdebuglnFatal("Failed to ping server", err)
				logger.FatalLogln("Failed to ping server", err)
			}
		case <-stop:
			logger.Fdebugln("Stop pinging")
			ticker.Stop()
			return
		}
	}
}

func WebSocketURL(token string) string {
	return fmt.Sprintf("%s://%s/ws/%s", endpoint.WebSocketScheme, endpoint.WebSocketURL, token)
}

func PingURL() string {
	return fmt.Sprintf("%s://%s/ping", endpoint.PingScheme, endpoint.PingURL)
}
