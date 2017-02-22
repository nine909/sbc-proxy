package httpControl

import (
	"fmt"
	"log"
	"net"
	"strconv"
	// "strings"
	// "time"
)

type service interface {
	Open(port string) *net.UDPConn
	StartServer(port string)
	DeletePort()
}

type Sbc struct {
	// portOpened map[string]
	handler string
	clients map[string]Client
}

type Client struct {
	addr        *net.UDPAddr
	sAddr       *net.UDPAddr
	connForward *net.UDPConn
	remoteAddr  string
	localAddr   *net.UDPAddr
	OpenConn    *net.UDPConn
	isclient    bool
}

func NewSBCServer() *Sbc {
	return &Sbc{clients: make(map[string]Client)}
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
	log.Println("PortOpen:", port)
	log.Println("backTrack:", rAddr)

	log.Println("sbc.clients:", sbc.clients)
	client := Client{remoteAddr: rAddr}
	if sbc.handler == "MT" && len(sbc.clients) >= 1 {
		log.Println("Connect to MO!!!")
		client.isclient = true
		sbc.clients[port] = client
		sbc.ConnectToMO(port)
		go sbc.UDPServer(port)
	} else {

		log.Println("PortOpen:", port)
		log.Println("backTrack:", rAddr)

		conn, err := sbc.Open(port)
		if err != nil {
			fmt.Printf("Some error %v\n", err)
			return err
		}
		client.OpenConn = conn
		sbc.clients[port] = client
		log.Println("UDP Server Started!!!")
		go sbc.UDPServer(port)
	}
	// if _, ok := sbc.clients[port]; !ok {
	// 	client.OpenConn = conn
	// }
	log.Println("sbc.handler", sbc.handler)
	log.Println("len(sbc.clients)", len(sbc.clients))
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

	// log.Println("SBC Port Opened:", sbc.portOpened)
	// for key, value := range sbc.portOpened {
	// 	log.Println("sbc. Port Closed : ", value.LocalAddr().String())
	// 	value.Close()
	// 	delete(sbc.portOpened, key)
	// }
	// log.Println("SBC Port Opened:", sbc.portOpened)
}

// func (sbc *Sbc) sendResponse(conn *net.UDPConn, p []byte) {

// 	log.Println(sbc.clients)
// 	for key := range sbc.clients {
// 		log.Println("KeyName:", key)
// 		if key != conn {
// 			// n, err := key.WriteToUDP([]byte("From server: Hello I got your mesage \n"), sbc.clients[key].addr)
// 			log.Println("Send to ", key, sbc.clients[key].addr)
// 			n, err := key.WriteToUDP(p, sbc.clients[key].addr)
// 			if err != nil {
// 				fmt.Println("Couldn't send response %v", err)
// 			}
// 			fmt.Println(n, err)

// 		}

// 	}
// }

func (sbc *Sbc) UDPServer(port string) error {

	conn := sbc.clients[port].OpenConn
	c := sbc.clients[port]
	log.Println("UDPDetail:", &conn)
	// defer conn.Close()
	for {
		p := make([]byte, 2048)
		log.Println("Waiting Incoming...", port)
		// sbc.findSender(conn)
		n, remoteaddr, err := conn.ReadFromUDP(p)
		p = p[:n]
		if err != nil {
			fmt.Printf("Some error  %v", err)
			// continue
			return err
		}
		sbc.findSender(port)
		c.sAddr = remoteaddr
		log.Println(c.sAddr)
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

		log.Printf("Read a message from %v %s \n", remoteaddr, p)
		sbc.clients[port] = c
		go sbc.sendTo(port, p)
		// go sbc.sendTo(conn, p)
	}

}

func (sbc *Sbc) findSender(port string) {
	log.Println("Starting Finding Sender:", port)

	connCurrent := sbc.clients[port]

	log.Println("connCurrent", connCurrent)
	log.Println("connCurrent.sAddr", connCurrent.sAddr)
	log.Println("connCurrent.localAddr", connCurrent.localAddr)
	log.Println("connCurrent.remoteAddr", connCurrent.remoteAddr)
	log.Println("connCurrent.connForward", connCurrent.connForward)
	log.Println("connCurrent.OpenConn", connCurrent.OpenConn)

	log.Println("Client range :", len(sbc.clients))
	log.Println(sbc.clients)
	var c Client
	for key, val := range sbc.clients {
		if key != port {
			c = val
			udpPort, _ := strconv.Atoi(key)
			addr := net.UDPAddr{
				Port: udpPort,
				IP:   net.ParseIP(""),
			}

			c.localAddr = &addr
			sbc.clients[key] = c
			break
		}
	}
	log.Println("c", c)
	log.Println("c.sAddr", c.sAddr)
	log.Println("c.localAddr", c.localAddr)
	log.Println("c.remoteAddr", c.remoteAddr)
	log.Println("c.connForward", c.connForward)
	log.Println("c.OpenConn", c.OpenConn)

	log.Println("c.LocalAddr:", c.localAddr)
	if c.localAddr == nil {
		return
	}

	// conToRemote, err := net.DialUDP("udp", c.localAddr, connCurrent.remoteAddr)
	// if err != nil {
	// 	fmt.Printf("Some error %v", err)
	// 	return
	// }
	// c.connForward = conToRemote
	// log.Println(c)
	// log.Println("Conn By Local", c.connForward.LocalAddr().String())
	// log.Println("Conn By Remote", c.connForward.RemoteAddr().String())
	// log.Println("Conn By ConnForward", c.connForward)
}

// func (c *Client) AddRemoteAddress(rAddr string) {
// 	log.Println("Remote Address:", rAddr)
// 	if c.remoteAddr == nil {
// 		remoteAddr, err := net.ResolveUDPAddr("udp", rAddr)
// 		if err != nil {
// 			fmt.Println("Error: ", err)
// 		}
// 		c.remoteAddr = remoteAddr
// 	}
// 	log.Println("Remote UDPAddr:", c.remoteAddr.String())
// }

func (sbc *Sbc) ConnectToMO(port string) {
	log.Println("ConnectToMO with Local:", port)

	connCurrent := sbc.clients[port]

	log.Println("connCurrent", connCurrent)
	log.Println("connCurrent.sAddr", connCurrent.sAddr)
	log.Println("connCurrent.localAddr", connCurrent.localAddr)
	log.Println("connCurrent.remoteAddr", connCurrent.remoteAddr)
	log.Println("connCurrent.connForward", connCurrent.connForward)
	log.Println("connCurrent.OpenConn", connCurrent.OpenConn)

	log.Println("Client range :", len(sbc.clients))
	log.Println(sbc.clients)
	var c Client
	// var lport string
	for key, val := range sbc.clients {
		if key != port {
			c = val
			// lport = key
			udpPort, _ := strconv.Atoi(key)
			addr := net.UDPAddr{
				Port: udpPort,
				IP:   net.ParseIP(""),
			}

			c.localAddr = &addr
			sbc.clients[key] = c
			break
		}
	}
	log.Println("c", c)
	log.Println("c.sAddr", c.sAddr)
	log.Println("c.localAddr", c.localAddr)
	log.Println("c.remoteAddr", c.remoteAddr)
	log.Println("c.connForward", c.connForward)
	log.Println("c.OpenConn", c.OpenConn)

	log.Println("c.LocalAddr:", c.localAddr)
	if c.localAddr == nil {
		return
	}

	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+port)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	ServerAddr, err := net.ResolveUDPAddr("udp", c.remoteAddr)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	ConnRemote, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	connCurrent.OpenConn = ConnRemote
	connCurrent.sAddr = ServerAddr
	sbc.clients[port] = connCurrent
	// defer ConnRemote.Close()
	// for {
	// 	n, err := ConnRemote.Write(p)
	// 	// n, err := c.OpenConn.WriteToUDP(p, conn.remoteAddr)
	// 	if err != nil {
	// 		log.Println("Couldn't send response %v", err)
	// 	}
	// }
}

func (sbc *Sbc) sendTo(port string, p []byte) {
	// sbc.findSender(conn)
	conn := sbc.clients[port]
	log.Println("conn", conn)
	log.Println("conn.sAddr", conn.sAddr)
	log.Println("conn.localAddr", conn.localAddr)
	log.Println("conn.remoteAddr", conn.remoteAddr)
	log.Println("conn.connForward", conn.connForward)
	log.Println("conn.OpenConn", conn.OpenConn)
	log.Println("conn.OpenConn.LocalAddr()", conn.OpenConn.LocalAddr())
	log.Println("conn.OpenConn.RemoteAddr()", conn.OpenConn.RemoteAddr())

	var c Client
	// var lport string
	for key, val := range sbc.clients {
		if key != port {
			c = val
			// lport = key
			break
		}
	}
	log.Println("c", c)
	log.Println("c.sAddr", c.sAddr)
	log.Println("c.localAddr", c.localAddr)
	log.Println("c.remoteAddr", c.remoteAddr)
	log.Println("c.connForward", c.connForward)
	log.Println("c.OpenConn", c.OpenConn)

	if c.OpenConn == nil {
		return
	}

	log.Println("c.OpenConn.LocalAddr()", c.OpenConn.LocalAddr())
	log.Println("c.OpenConn.RemoteAddr()", c.OpenConn.RemoteAddr())

	log.Println("sbc.handler", sbc.handler)
	if sbc.handler == "MO" && c.sAddr != nil {
		// defer c.OpenConn.Close()
		n, err := c.OpenConn.WriteTo(p, c.sAddr)
		if err != nil {
			log.Println("Couldn't send response", c.sAddr, err)
		}
		// log.Println("Send to : %s --> %s", c.OpenConn.LocalAddr().String(), c.OpenConn.RemoteAddr().String())
		log.Println(n, err)
	} else if sbc.handler == "MT" && c.sAddr != nil {
		// defer c.OpenConn.Close()
		if c.isclient {
			// defer c.OpenConn.Close()
			n, err := c.OpenConn.Write(p)
			if err != nil {
				log.Println("Couldn't send response", c.sAddr, err)
			}
			// log.Println("Send to : %s --> %s", c.OpenConn.LocalAddr().String(), c.OpenConn.RemoteAddr().String())
			log.Println(n, err)
		} else {
			n, err := c.OpenConn.WriteTo(p, c.sAddr)
			if err != nil {
				log.Println("Couldn't send response", c.sAddr, err)
			}
			// log.Println("Send to : %s --> %s", c.OpenConn.LocalAddr().String(), c.OpenConn.RemoteAddr().String())
			log.Println(n, err)
		}

	}

	// else {
	// 	log.Printf("Send a message to %v %s \n", conn.remoteAddr, p)

	// 	s := strings.Split(conn.remoteAddr, ":")
	// 	// ip, port := s[0], s[1]
	// 	var sport string
	// 	sport = s[1]
	// 	log.Println(lport)
	// 	// LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+lport)
	// 	// if err != nil {
	// 	// 	fmt.Println("Error: ", err)
	// 	// }

	// 	ServerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+sport)
	// 	if err != nil {
	// 		fmt.Println("Error: ", err)
	// 	}
	// 	for i := 0; i < 3; i++ {
	// 		log.Println("Close Conn", c.OpenConn)
	// 		c.OpenConn.Close()
	// 		time.Sleep(2 * time.Second)
	// 	}
	// 	LocalAddr := c.localAddr
	// 	// ServerAddr = conn.remoteAddr
	// 	log.Println("LocalAddr", LocalAddr)
	// 	log.Println("ServerAddr", ServerAddr)
	// 	log.Println("LocalAddr.String()", LocalAddr.String())
	// 	log.Println("ServerAddr.String()", ServerAddr.String())
	// 	ConnRemote, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	// 	if err != nil {
	// 		fmt.Println("Error: ", err)
	// 	}

	// 	// defer ConnRemote.Close()
	// 	// for i := 0; i < 3; i++ {
	// 	n, err := ConnRemote.Write(p)
	// 	// n, err := c.OpenConn.WriteToUDP(p, conn.remoteAddr)
	// 	if err != nil {
	// 		log.Println("Couldn't send response %v", err)
	// 	}
	// 	log.Println(n, err)
	// 	c.OpenConn = ConnRemote
	// 	c.sAddr = ServerAddr
	// 	sbc.clients[lport] = c
	// 	// }
	// }

}
