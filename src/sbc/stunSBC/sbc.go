package stunSBC

import (
	"fmt"
	"net"

	"github.com/pixelbender/go-stun/stun"
)

func ServerListener(port string) {
	srv := stun.NewServer(nil)
	p := ":" + port
	fmt.Println(p)
	protocol := "udp"
	l, err := net.ListenPacket(protocol, p)
	if err != nil {
		fmt.Print("listen error", err)
	}

	defer l.Close()
	srv.ServePacket(l)

}
