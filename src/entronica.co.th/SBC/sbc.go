package main

import (
	"fmt"
	"net"

	"github.com/pixelbender/go-stun/stun"
)

func main() {
	srv := stun.NewServer(nil)
	l, err := net.ListenPacket("udp", ":6060")
	if err != nil {
		fmt.Print("listen error", err)
	}
	defer l.Close()
	srv.ServePacket(l)

}
