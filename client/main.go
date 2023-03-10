package main

import (
	"log"

	"github.com/he-man-david/distributed-rate-limiter-basic/src/client/clients"
)

var (
	start = 8001
	end   = 8003
)

func main() {
	c, err := clients.NewClients(start, end)
	if err != nil {
		log.Fatalf("[Main] failed to create new clients :: ERR: %v", err)
	}
	// close all connections on shutdown
	defer shutdown(c)
	
	log.Println(" *** Testing scenarios *** ")
	sendManyToOneNode(c)
	
	select{}
}

func shutdown(c *clients.Clients) {
	for _, conn := range c.Conns {
		conn.Close()
	}
}

// *** Testing scenarios ***

func sameKeyAllNode(c *clients.Clients) {
	key := 1111
	for port := start; port <= end; port++ {
		c.AllowRequest(port, key)
	}
}

func sendManyToOneNode(c *clients.Clients) {
	key := 2222
	for i := 0; i < 20; i++ {
		c.AllowRequest(int(8001), key)
	}
}