package client

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/url"

	"github.com/gorilla/websocket"
)

type connection struct {
	client *websocket.Conn
	url    url.URL
}

func NewConnection(addr string) *connection {
	u := url.URL{Scheme: "wss", Host: addr, Path: "/ccproxy/main"}
	return &connection{
		url: u,
	}
}

func (c *connection) begin() {

	log.Println("connecting to proxy server", c.url.String())
	client, _, err := websocket.DefaultDialer.Dial(c.url.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	log.Println("connected to proxy")

	c.client = client
}

func (c *connection) serve(tcpConn *net.TCPConn) {

	go c.startReceivingFromProxy(tcpConn)
	go c.startSendingToProxy(tcpConn)
}

const BUFFER_LEN = 1024

func (c *connection) startSendingToProxy(tcpConn *net.TCPConn) {

	for {
		buffer := make([]byte, BUFFER_LEN)
		count, err := tcpConn.Read(buffer)
		if err != nil {
			log.Println("read from tcp connection failed:", err)
			break
		}
		log.Printf("read %d bytes from tcp client", count)

		log.Println("writing to ws proxy")
		err = c.client.WriteMessage(websocket.BinaryMessage, buffer[:count])
		if err != nil {
			log.Println("write to ws proxy failed:", err)
			break
		}

		log.Printf("wrote %d bytes to ws proxy", count)
	}

	tcpConn.Close()
	c.client.Close()
}

func (c *connection) startReceivingFromProxy(tcpConn *net.TCPConn) {

	for {
		_, msg, err := c.client.ReadMessage()
		if err != nil {
			log.Println("NextReader:", err)
			break
		}

		log.Printf("received %d bytes from proxy\n", len(msg))
		reader := bytes.NewReader(msg)
		io.Copy(tcpConn, reader)
		log.Printf("wrote %d bytes to client\n", len(msg))
	}
}
