package httpControl

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

type service interface {
	Open(port string) *net.UDPConn
	StartServer(port string)
	DeletePort()
}

type Sbc struct {
	portOpened map[string]*net.UDPConn
	clients    map[*net.UDPConn]Client
	remoteAddr string
}

type Client struct {
	addr                  *net.UDPAddr
	connForward           *net.UDPConn
	remoteAddr, localAddr *net.UDPAddr
}

func NewSBCServer() *Sbc {
	return &Sbc{portOpened: make(map[string]*net.UDPConn), clients: make(map[*net.UDPConn]Client)}
}

func (sbc *Sbc) Open(port string) (*net.UDPConn, error) {
	// udpPort, _ := strconv.Atoi(port)
	// addr := net.UDPAddr{
	// 	Port: UDPPort,
	// 	IP:   net.ParseIP(""),
	// 	// IP:   net.ParseIP("127.0.0.1"),
	// }

	udpAddr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return nil, err
	}

	ser, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return nil, err
	}
	log.Println("Server Listener...", ser.LocalAddr().String())
	return ser, nil
}

func (sbc *Sbc) StartServer(rAddr, port string) error {
	log.Println("UDP Server Starting...")

	conn, err := sbc.Open(port)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return err
	}
	if val, ok := sbc.portOpened[port]; !ok {
		sbc.portOpened[port] = conn
	} else {
		conn = val
	}

	client := Client{}
	client.AddRemoteAddress(rAddr)
	sbc.clients[conn] = client
	log.Println("UDP Server Started!!!")
	go sbc.UDPServer(conn)
	return nil
}

// func (c *Client) findSendTo(conn *net.UDPConn) (*net.UDPConn, error) {
// 	raddr, err := net.ResolveUDPAddr("udp", c.raddr)
// 	if err != nil {
// 		fmt.Printf("Some error %v\n", err)
// 		return nil, err
// 	}
// 	laddr, err := net.ResolveUDPAddr("udp", ":"+c.port)
// 	if err != nil {
// 		fmt.Printf("Some error %v\n", err)
// 		return nil, err
// 	}
// 	connOld, err := net.DialUDP("udp", laddr, raddr)
// 	if err != nil {
// 		fmt.Printf("Some error %v", err)
// 		return nil, err
// 	}

// 	return connOld, nil
// }

func (sbc *Sbc) DeletePort() {

	log.Println("SBC Port Opened:", sbc.portOpened)
	for key, value := range sbc.portOpened {
		log.Println("sbc. Port Closed : ", value.LocalAddr().String())
		value.Close()
		delete(sbc.portOpened, key)
	}
	log.Println("SBC Port Opened:", sbc.portOpened)
}

func (sbc *Sbc) sendResponse(conn *net.UDPConn, p []byte) {

	log.Println(sbc.clients)
	for key := range sbc.clients {
		log.Println("KeyName:", key)
		if key != conn {
			// n, err := key.WriteToUDP([]byte("From server: Hello I got your mesage \n"), sbc.clients[key].addr)
			log.Println("Send to ", key, sbc.clients[key].addr)
			n, err := key.WriteToUDP(p, sbc.clients[key].addr)
			if err != nil {
				fmt.Println("Couldn't send response %v", err)
			}
			fmt.Println(n, err)

		}

	}
}

func (sbc *Sbc) UDPServer(conn *net.UDPConn) error {

	log.Println("UDPDetail:", &conn)
	defer conn.Close()
	p := make([]byte, 2048)
	for {
		log.Println("Waiting Incoming...")
		// sbc.findSender(conn)
		_, remoteaddr, err := conn.ReadFromUDP(p)

		if err != nil {
			fmt.Printf("Some error  %v", err)
			// continue
			return err
		}
		log.Println("Remote Address Reader:", remoteaddr)
		// client := Client{addr: remoteaddr}
		// if val, ok := sbc.clients[conn]; !ok {
		// 	log.Println("Added New Client:", client.addr)
		// } else {
		// 	log.Println("Client is Existed:", val.addr)
		// 	log.Println("Update Client:", client.addr)
		// }
		// sbc.clients[conn] = client

		// log.Println("Client UDPAddr:", sbc.clients[conn])

		// sbc.MapUDPAddrs[ser] = remoteaddr
		// log.Println("MapUDPAddress:", sbc.MapUDPAddrs[ser])
		// fmt.Printf("Read a message from %v %s \n", remoteaddr, p)

		// go sbc.sendResponse(conn, p)

		// log.Printf("Read a message from %v %s \n", remoteaddr, p)

		go sbc.sendTo(conn, p)
	}

}

// func (sbc *Sbc) SendTo(lAddr, rAddr string) {
// 	// rAddr = "127.0.0.1:10001"
// 	ServerAddr, err := net.ResolveUDPAddr("udp", rAddr)
// 	if err != nil {
// 		fmt.Println("Error: ", err)
// 	}

// 	LocalAddr = sbc.portOpened[lAddr]

// 	if
// 	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
// 	if err != nil {
// 		fmt.Println("Error: ", err)
// 	}
// 	defer Conn.Close()
// 	for {
// 		_, err := Conn.Write(buf)
// 		if err != nil {
// 			fmt.Println(msg, err)
// 		}
// 	}

// }

func (sbc *Sbc) findSender(conn *net.UDPConn) {
	log.Println("Starting Finding Sender:", conn)
	c := sbc.clients[conn]
	log.Println(c.localAddr)
	log.Println(c.remoteAddr)
	log.Println(c.connForward)
	if c.connForward != nil {
		for key, val := range sbc.portOpened {
			log.Println("KeyName:", key)
			if val != conn {
				udpPort, _ := strconv.Atoi(key)
				addr := net.UDPAddr{
					Port: udpPort,
					IP:   net.ParseIP(""),
				}
				c.localAddr = &addr
				break
			}
		}
		log.Println("c.LocalAddr:", c.localAddr)
		conToRemote, err := net.DialUDP("udp", sbc.clients[conn].localAddr, sbc.clients[conn].remoteAddr)
		if err != nil {
			fmt.Printf("Some error %v", err)
			return
		}
		c.connForward = conToRemote
	}
	log.Println("Conn By Local", c.connForward.LocalAddr().String())
	log.Println("Conn By Remote", c.connForward.RemoteAddr().String())
}

func (c *Client) AddRemoteAddress(rAddr string) {
	log.Println("Remote Address:", rAddr)
	if c.remoteAddr == nil {
		remoteAddr, err := net.ResolveUDPAddr("udp", rAddr)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		c.remoteAddr = remoteAddr
	}
	log.Println("Remote UDPAddr:", c.remoteAddr.String())
}

func (sbc *Sbc) sendTo(conn *net.UDPConn, p []byte) {
	sbc.findSender(conn)
	c := sbc.clients[conn]
	log.Println(c.connForward)
	log.Println(conn)
	if c.connForward == nil {
		return
	}
	log.Println("Recieved By Local", conn.LocalAddr().String())
	log.Println("Recieved By Remote", conn.RemoteAddr().String())
	log.Println("Send From Local", sbc.clients[conn].connForward.LocalAddr().String())
	log.Println("Send to ", c.connForward.RemoteAddr().String())
	n, err := c.connForward.WriteToUDP(p, c.remoteAddr)
	if err != nil {
		log.Println("Couldn't send response %v", err)
	}
	log.Println(n, err)

}
