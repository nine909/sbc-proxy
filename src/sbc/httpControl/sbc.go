package httpControl

import (
	"fmt"
	"log"
	"net"
	"sbc/logs"
	"strconv"
	// "strings"
	// "time"

	conf "sbc/conf"
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
	logs.Logger.Debug("Server Listener...", ser.LocalAddr().String())
	return ser, nil
}

func (sbc *Sbc) StartServer(rAddr, port string) error {
	logs.Logger.Debug("UDP Server Starting...")
	logs.Logger.Debug("PortOpen:", port)
	logs.Logger.Debug("backTrack:", rAddr)

	logs.Logger.Debug("sbc.clients:", sbc.clients)
	client := Client{remoteAddr: rAddr}
	if sbc.handler == "MT" && len(sbc.clients) >= 1 {
		logs.Logger.Debug("Connect to MO!!!")
		client.isclient = true
		sbc.clients[port] = client
		sbc.ConnectToMO(port)
		go sbc.UDPServer(port)
	} else {

		logs.Logger.Debug("PortOpen:", port)
		logs.Logger.Debug("backTrack:", rAddr)

		conn, err := sbc.Open(port)
		if err != nil {
			fmt.Printf("Some error %v\n", err)
			return err
		}
		client.OpenConn = conn
		sbc.clients[port] = client
		logs.Logger.Debug("UDP Server Started!!!")
		go sbc.UDPServer(port)
	}
	// if _, ok := sbc.clients[port]; !ok {
	// 	client.OpenConn = conn
	// }
	logs.Logger.Debug("sbc.handler", sbc.handler)
	logs.Logger.Debug("len(sbc.clients)", len(sbc.clients))
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
	logs.Logger.Debug("UDPDetail:", &conn)
	// defer conn.Close()
	for {
		p := make([]byte, 2048)
		logs.Logger.Debug("Waiting Incoming...", port)
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
		logs.Logger.Debug(c.sAddr)
		logs.Logger.Debug("Remote Address Reader:", remoteaddr)
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

		logs.Logger.Debug("Read a message from %v %s \n", remoteaddr, p)
		sbc.clients[port] = c
		go sbc.sendTo(port, p)
		// go sbc.sendTo(conn, p)
	}

}

func (sbc *Sbc) findSender(port string) {
	logs.Logger.Debug("Starting Finding Sender:", port)

	connCurrent := sbc.clients[port]

	logs.Logger.Debug("connCurrent", connCurrent)
	logs.Logger.Debug("connCurrent.sAddr", connCurrent.sAddr)
	logs.Logger.Debug("connCurrent.localAddr", connCurrent.localAddr)
	logs.Logger.Debug("connCurrent.remoteAddr", connCurrent.remoteAddr)
	logs.Logger.Debug("connCurrent.connForward", connCurrent.connForward)
	logs.Logger.Debug("connCurrent.OpenConn", connCurrent.OpenConn)

	logs.Logger.Debug("Client range :", len(sbc.clients))
	logs.Logger.Debug(sbc.clients)
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
	logs.Logger.Debug("c", c)
	logs.Logger.Debug("c.sAddr", c.sAddr)
	logs.Logger.Debug("c.localAddr", c.localAddr)
	logs.Logger.Debug("c.remoteAddr", c.remoteAddr)
	logs.Logger.Debug("c.connForward", c.connForward)
	logs.Logger.Debug("c.OpenConn", c.OpenConn)

	logs.Logger.Debug("c.LocalAddr:", c.localAddr)
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
	logs.Logger.Debug("ConnectToMO with Local:", port)

	connCurrent := sbc.clients[port]

	logs.Logger.Debug("connCurrent", connCurrent)
	logs.Logger.Debug("connCurrent.sAddr", connCurrent.sAddr)
	logs.Logger.Debug("connCurrent.localAddr", connCurrent.localAddr)
	logs.Logger.Debug("connCurrent.remoteAddr", connCurrent.remoteAddr)
	logs.Logger.Debug("connCurrent.connForward", connCurrent.connForward)
	logs.Logger.Debug("connCurrent.OpenConn", connCurrent.OpenConn)

	logs.Logger.Debug("Client range :", len(sbc.clients))
	logs.Logger.Debug(sbc.clients)
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
	logs.Logger.Debug("c", c)
	logs.Logger.Debug("c.sAddr", c.sAddr)
	logs.Logger.Debug("c.localAddr", c.localAddr)
	logs.Logger.Debug("c.remoteAddr", c.remoteAddr)
	logs.Logger.Debug("c.connForward", c.connForward)
	logs.Logger.Debug("c.OpenConn", c.OpenConn)

	logs.Logger.Debug("c.LocalAddr:", c.localAddr)
	if c.localAddr == nil {
		return
	}

	localIp := conf.Conf.Localip
	logs.Logger.Debug("localIp", localIp)
	LocalAddr, err := net.ResolveUDPAddr("udp", localIp+":"+port)
	if err != nil {
		logs.Logger.Debug("Error: ", err)
	}

	ServerAddr, err := net.ResolveUDPAddr("udp", c.remoteAddr)
	if err != nil {
		logs.Logger.Debug("Error: ", err)
	}

	ConnRemote, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	if err != nil {
		logs.Logger.Debug("Error: ", err)
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
	logs.Logger.Debug("conn", conn)
	logs.Logger.Debug("conn.sAddr", conn.sAddr)
	logs.Logger.Debug("conn.localAddr", conn.localAddr)
	logs.Logger.Debug("conn.remoteAddr", conn.remoteAddr)
	logs.Logger.Debug("conn.connForward", conn.connForward)
	logs.Logger.Debug("conn.OpenConn", conn.OpenConn)
	logs.Logger.Debug("conn.OpenConn.LocalAddr()", conn.OpenConn.LocalAddr())
	logs.Logger.Debug("conn.OpenConn.RemoteAddr()", conn.OpenConn.RemoteAddr())

	var c Client
	// var lport string
	for key, val := range sbc.clients {
		if key != port {
			c = val
			// lport = key
			break
		}
	}
	logs.Logger.Debug("c", c)
	logs.Logger.Debug("c.sAddr", c.sAddr)
	logs.Logger.Debug("c.localAddr", c.localAddr)
	logs.Logger.Debug("c.remoteAddr", c.remoteAddr)
	logs.Logger.Debug("c.connForward", c.connForward)
	logs.Logger.Debug("c.OpenConn", c.OpenConn)

	if c.OpenConn == nil {
		return
	}

	logs.Logger.Debug("c.OpenConn.LocalAddr()", c.OpenConn.LocalAddr())
	logs.Logger.Debug("c.OpenConn.RemoteAddr()", c.OpenConn.RemoteAddr())

	logs.Logger.Debug("sbc.handler", sbc.handler)
	if sbc.handler == "MO" && c.sAddr != nil {
		// defer c.OpenConn.Close()
		n, err := c.OpenConn.WriteTo(p, c.sAddr)
		if err != nil {
			logs.Logger.Debug("Couldn't send response", c.sAddr, err)
		}
		// log.Println("Send to : %s --> %s", c.OpenConn.LocalAddr().String(), c.OpenConn.RemoteAddr().String())
		logs.Logger.Debug(n, err)
	} else if sbc.handler == "MT" && c.sAddr != nil {
		// defer c.OpenConn.Close()
		if c.isclient {
			// defer c.OpenConn.Close()
			n, err := c.OpenConn.Write(p)
			if err != nil {
				log.Println("Couldn't send response", c.sAddr, err)
			}
			// log.Println("Send to : %s --> %s", c.OpenConn.LocalAddr().String(), c.OpenConn.RemoteAddr().String())
			logs.Logger.Debug(n, err)
		} else {
			n, err := c.OpenConn.WriteTo(p, c.sAddr)
			if err != nil {
				log.Println("Couldn't send response", c.sAddr, err)
			}
			// log.Println("Send to : %s --> %s", c.OpenConn.LocalAddr().String(), c.OpenConn.RemoteAddr().String())
			logs.Logger.Debug(n, err)
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
