package client

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var AllClients = make(map[string]*Client)

//Client struct of clients
type Client struct {
	Username string
	Socket   *websocket.Conn
	Join     chan bool
	Leave    chan bool
	Message  chan Message
}

type Message struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}

func (c *Client) Read() {
	
	for {
		msg := Message{}
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("%v",msg)
		c.Message <- msg
	}
	c.Socket.Close()


}
func (c *Client) Write(msg Message) {
	
		err := c.Socket.WriteJSON(msg)
		if err != nil {
			fmt.Println(err)
			c.Socket.Close()
		}

		
	
}
func AddClient(c *Client) {
	AllClients[c.Username]= c
			fmt.Printf("%v", AllClients)
}

func FindAndSend(msg Message) {
	fmt.Printf("find")
	if val, ok := AllClients[msg.Recipient]; ok {
		fmt.Printf("%v %v", val, ok)
		if err := val.Socket.WriteJSON(msg); err != nil {
			val.Socket.Close()
		}
	}
}

func NewClient(conn *websocket.Conn, username string) *Client {

	return &Client{
		Username: username,
		Socket:   conn,
		Join:     make(chan bool),
		Leave:    make(chan bool),
		Message:  make(chan Message),
	}
}

func (c *Client) Listen() {

	for {
		select {
		case <-c.Join:
			fmt.Printf("Uso")
			AddClient(c)
			
		case <-c.Leave:
			fmt.Printf("Izaso")
			delete(AllClients, c.Username)
		case msg := <-c.Message:
			
			c.Write(msg)
			FindAndSend(msg)
		}
	}
}
