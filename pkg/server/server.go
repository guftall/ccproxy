package server

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/guftall/ccproxy/pkg/proxy"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var serverListenAddr = flag.String("server-addr", "127.0.0.1:9991", "default: 127.0.0.1:9991")
var tunnelListenAddr = flag.String("tunnel-addr", "127.0.0.1:9993", "default: 127.0.0.1:9993")

type Server struct {
}

func (s *Server) Start() {
	go s.startHttpServer()
	go s.startTunnelServer()

	ch := make(chan int)

	<-ch
}

func (s *Server) startHttpServer() {

	http.HandleFunc("/ccproxy/main", s.handle)

	log.Println("starting server on", *serverListenAddr)
	if err := http.ListenAndServe(*serverListenAddr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func (s *Server) startTunnelServer() {

	proxy := &proxy.TunnelProxy{}

	log.Println("Starting tunnel server on", *tunnelListenAddr)
	if err := http.ListenAndServe(*tunnelListenAddr, proxy); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	c := NewConnection()
	c.begin(w, r)
	c.serve()
}
