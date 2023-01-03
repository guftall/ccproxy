package client

import (
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

		writer, err := c.client.NextWriter(websocket.BinaryMessage)
		if err != nil {
			log.Println("NextWriter:", err)
			break
		}
		log.Println("writing to ws proxy")

		count, err = writer.Write(buffer[:count])
		if err != nil {
			log.Println("write to ws proxy failed:", err)
			break
		}
		log.Printf("wrote %d bytes to ws proxy", count)
		writer.Close()
	}

	tcpConn.Close()
	c.client.Close()
}

func (c *connection) startReceivingFromProxy(tcpConn *net.TCPConn) {

	for {
		_, reader, err := c.client.NextReader()
		if err != nil {
			log.Println("NextReader:", err)
			break
		}

		buffer := make([]byte, BUFFER_LEN)

		for {
			count, err := reader.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("read from reader:", err)
				break
			}
			if count > 0 {

				tcpConn.Write(buffer[:count])
			}

			if count < BUFFER_LEN {
				break
			}
		}
	}
}
