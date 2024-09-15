package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
	"strings"

	types "common"
	"server/src/auth"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	tempId int64
	send   chan types.Message
	recv   chan types.Message
}

func (self *Client) Close() {
	close(self.recv)
	// keep send open ig :/
}

func (self *Client) handleMessages() {
	for {
		msg, ok := <-self.recv
		if !ok {
			return
		}

		switch msg.Typ {
		default:
			log.Printf("message type '%s' not handled ('%s')\n", msg.Typ.TSName(), msg.Val)
		}
	}

}

type ClientsCtx struct {
	// string -> *Client
	clients sync.Map
	tempId  int64
	mutex   sync.Mutex
}

func (self *ClientsCtx) addClient(name string) *Client {
	val, ok := self.clients.Load(name)
	if !ok {
		client := &Client{
			send: make(chan types.Message),
			recv: make(chan types.Message),
		}
		self.set(name, client)
		return client
	}
	client, _ := val.(*Client)
	return client
}

func (self *ClientsCtx) set(name string, client *Client) int64 {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	id := atomic.AddInt64(&self.tempId, 1)
	client.tempId = id

	self.clients.Store(name, client)

	return id
}

func (self *ClientsCtx) get(name string) (*Client, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	val, ok := self.clients.Load(name)
	if !ok {
		return nil, fmt.Errorf("Client with username '%s' not found", name)
	}
	client, ok := val.(*Client)
	if !ok {
		return nil, fmt.Errorf("Bad value in map")
	}
	return client, nil
}

// don't yeet clients from memory. push messages in them as they will be sent to user after reconnection
func (self *ClientsCtx) remove(name string, tempId int64) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	val, ok := self.clients.LoadAndDelete(name)
	if !ok {
		return
	}
	client, ok := val.(*Client)
	if !ok {
		return
	}
	if client.tempId != tempId {
		self.clients.Store(name, client)
		return
	}

	client.Close()
}

func AppRoutes(route *gin.Engine) {
	var state ClientsCtx

	handleClient := func(ws *websocket.Conn, ctx context.Context, client *Client) {
		defer ws.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-client.send:
				if !ok {
					return
				}
				log.Println(msg)
				ws.SetWriteDeadline(time.Now().Add(time.Second * 5))
				err := ws.WriteJSON(msg)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}

	handleMessages := func(ws *websocket.Conn, close context.CancelFunc, client *Client) {
		defer ws.Close()

		for {
			var msg types.Message
			err := ws.ReadJSON(&msg)
			if err != nil {
				log.Println(err)
				close()
				return
			}
			client.recv <- msg
		}
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	wsHandler := func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			return
		}
		token := bearerToken[1]

		claims, err := auth.VerifyJWT(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			c.JSON(401, gin.H{"error": "Username not found in token"})
			return
		}

		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(400, gin.H{"error": "could not upgrade websocket"})
			return
		}

		ctx, cancel := context.WithCancel(context.Background())

		client := state.addClient(username)

		log.Println("new conn")
		go client.handleMessages()
		go handleClient(ws, ctx, client)
		handleMessages(ws, cancel, client)
	}

	route.GET("/ws", wsHandler)
}
