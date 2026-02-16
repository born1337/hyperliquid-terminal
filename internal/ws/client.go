package ws

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	url           string
	conn          *websocket.Conn
	mu            sync.Mutex
	msgCh         chan Message
	done          chan struct{}
	subscriptions []SubRequest
	connected     bool
}

func NewClient(url string, msgCh chan Message) *Client {
	return &Client{
		url:   url,
		msgCh: msgCh,
		done:  make(chan struct{}),
	}
}

func (c *Client) Connected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connected
}

func (c *Client) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return err
	}
	c.mu.Lock()
	c.conn = conn
	c.connected = true
	c.mu.Unlock()

	// Resubscribe
	for _, sub := range c.subscriptions {
		c.send(sub)
	}

	go c.readLoop()
	return nil
}

func (c *Client) Subscribe(sub SubRequest) {
	c.mu.Lock()
	// Deduplicate: check if this subscription already exists
	data, _ := json.Marshal(sub.Subscription)
	key := string(data)
	found := false
	for _, existing := range c.subscriptions {
		ed, _ := json.Marshal(existing.Subscription)
		if string(ed) == key {
			found = true
			break
		}
	}
	if !found {
		c.subscriptions = append(c.subscriptions, sub)
	}
	c.mu.Unlock()
	c.send(sub)
}

func (c *Client) send(msg interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *Client) readLoop() {
	defer func() {
		c.mu.Lock()
		c.connected = false
		if c.conn != nil {
			c.conn.Close()
		}
		c.mu.Unlock()
		c.reconnect()
	}()

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		select {
		case c.msgCh <- msg:
		case <-c.done:
			return
		default:
			log.Printf("ws: dropped message on channel %s (buffer full)", msg.Channel)
		}
	}
}

func (c *Client) reconnect() {
	backoff := time.Second
	maxBackoff := 30 * time.Second

	for {
		select {
		case <-c.done:
			return
		case <-time.After(backoff):
			if err := c.Connect(); err != nil {
				log.Printf("ws reconnect failed: %v", err)
				backoff *= 2
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
				continue
			}
			return
		}
	}
}

func (c *Client) Close() {
	close(c.done)
	c.mu.Lock()
	if c.conn != nil {
		c.conn.Close()
	}
	c.mu.Unlock()
}
