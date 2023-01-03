package server

import (
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
		_, reader, err := c.client.NextReader()
		if err != nil {
			log.Println(err)
			return
		}

		buffer := make([]byte, BUFFER_LEN)

		for {
			count, err := reader.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("read from proxy failed:", err)
				break
			}

			log.Printf("read %d bytes from proxy\n", count)

			if count > 0 {
				count, err = c.tunnel.Write(buffer[:count])
				if err != nil {
					log.Println("write to tunnel failed", err)
					break
				}
				log.Printf("wrote %d bytes to tunnel \n", count)
			}

			if count < BUFFER_LEN {
				break
			}
		}
	}
}

func (c *connection) writeToProxy() {

	for {
		writer, err := c.client.NextWriter(websocket.BinaryMessage)
		if err != nil {
			log.Println("NextWriter:", err)
			break
		}

		buffer := make([]byte, BUFFER_LEN)

		count, err := c.tunnel.Read(buffer)
		if err != nil {
			log.Println("read from tunnel", err)
			break
		}
		log.Printf("read %d bytes from tunnel\n", count)

		count, err = writer.Write(buffer[:count])
		if err != nil {
			log.Println("write to proxy:", err)
			break
		}
		log.Printf("wrote %d bytes to proxy\n", count)
		writer.Close()
	}
}
