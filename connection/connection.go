package connection

import (
	"fmt"
	"strconv"
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

// TODO: Use channels so the Connection struct is thread safe.
// Gorilla's websocket.Conn can be accessed only from a single
// goroutine.

func New(token string, admin bool) (*Connection, error) {
	logger.Fdebugln("WS dialing")
	url := WebSocketURL(token, admin)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	logger.Fdebugln("WS connected")

	return &Connection{token, c}, nil
}

func (c *Connection) Close() {
	if c.wsconn == nil {
		return
	}
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

func WebSocketURL(token string, admin bool) string {
	return fmt.Sprintf("%s://%s/ws/%s?admin=%s", endpoint.WebSocketScheme, endpoint.WebSocketURL, token, strconv.FormatBool(admin))
}

func PingURL() string {
	return fmt.Sprintf("%s://%s/ping", endpoint.PingScheme, endpoint.PingURL)
}
