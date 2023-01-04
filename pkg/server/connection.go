package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

type connection struct {
	client *websocket.Conn
	tunnel net.Conn
}

func NewConnection() *connection {
	return &connection{}
}

func (c *connection) begin(w http.ResponseWriter, r *http.Request) {

	log.Println("upgrading ws connection")
	client, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("we have our ws connection!")
	c.client = client
	c.initTunnelConnection()
}

func (c *connection) initTunnelConnection() {
	tcpConn, err := net.Dial("tcp", "127.0.0.1:9993")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[server] successfuly connected to tunnel")
	c.tunnel = tcpConn
}

func (c *connection) serve() {

	go c.readFromProxy()
	go c.writeToProxy()
}

const BUFFER_LEN = 1024

func (c *connection) readFromProxy() {

	for {
		_, msg, err := c.client.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		reader := bytes.NewReader(msg)

		io.Copy(c.tunnel, reader)
	}
}

func (c *connection) writeToProxy() {

	for {
		buffer := make([]byte, BUFFER_LEN)

		count, err := c.tunnel.Read(buffer)
		if err != nil {
			log.Println("read from tunnel", err)
			break
		}
		log.Printf("read %d bytes from tunnel\n", count)

		err = c.client.WriteMessage(websocket.BinaryMessage, buffer[:count])
		if err != nil {
			log.Println("write to proxy:", err)
			break
		}
		log.Printf("wrote %d bytes to proxy\n", count)
	}
}
