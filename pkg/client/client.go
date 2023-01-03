package client

import (
	"flag"
	"log"
	"net"
)

type Client struct {
}

func (c *Client) Start() {

	// go c.startTunnelProxy()
	// c.startForwardProxy()
	c.start()

	// wait := make(chan int)
	// <-wait
}

var clientListenAddr = flag.String("client-addr", "127.0.0.1:9990", "default: 127.0.0.1:9990")
var proxyDomainAddr = flag.String("ccproxy-proxy-domain", "ccproxy.guftall.ir", "default: ccproxy.guftall.ir")

func (c *Client) start() {

	log.Println("Starting client on", *clientListenAddr)

	addr, err := net.ResolveTCPAddr("tcp", *clientListenAddr)
	if err != nil {
		panic(err)
	}

	tcpListener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal("ListenTCP:", err)
	}

	for {

		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			log.Fatal("AcceptTCP:", err)
		}
		go c.serveTcp(tcpConn)
	}
}

func (c *Client) serveTcp(tcpConn *net.TCPConn) {
	conn := NewConnection(*proxyDomainAddr)
	conn.begin()
	conn.serve(tcpConn)
}

// func (c *Client) startTunnelProxy() {

// 	var addr = "127.0.0.1:9990"

// 	proxy := &tunnelProxy{}

// 	log.Println("Starting tunnel proxy server on", addr)
// 	if err := http.ListenAndServe(addr, proxy); err != nil {
// 		log.Fatal("ListenAndServe:", err)
// 	}
// }

// func (c *Client) startForwardProxy() {

// 	var addr = "127.0.0.1:9991"

// 	proxy := &forwardProxy{}

// 	log.Println("Starting forward proxy server on", addr)
// 	if err := http.ListenAndServe(addr, proxy); err != nil {
// 		log.Fatal("ListenAndServe:", err)
// 	}
// }

// func (c *Client) Resolve(domain string) []dns.Answer {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	dc := doh.Use(doh.CloudflareProvider, doh.GoogleProvider, doh.Quad9Provider)
// 	defer dc.Close()

// 	rsp, err := dc.Query(ctx, dns.Domain(domain), dns.TypeTXT)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("resolved by %s\n", rsp.Provider)
// 	return rsp.Answer
// }
