package httpControl

import (
	"fmt"
	"log"
	"net"
)

type service interface {
	Open(port string) *net.UDPConn
	StartServer(port string)
	DeletePort()
}

type Sbc struct {
	portOpened map[string]*net.UDPConn
	clients    map[*net.UDPConn]Client
}

type Client struct {
	addr *net.UDPAddr
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

func (sbc *Sbc) StartServer(port string) error {
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
	log.Println("UDP Server Started!!!")
	go sbc.UDPServer(conn)
	return nil
}

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
		_, remoteaddr, err := conn.ReadFromUDP(p)

		if err != nil {
			fmt.Printf("Some error  %v", err)
			// continue
			return err
		}

		log.Println("Remote Address:", remoteaddr)
		client := Client{addr: remoteaddr}
		if val, ok := sbc.clients[conn]; !ok {
			log.Println("Added New Client:", client.addr)
		} else {
			log.Println("Client is Existed:", val.addr)
			log.Println("Update Client:", client.addr)
		}
		sbc.clients[conn] = client

		log.Println("Client UDPAddr:", sbc.clients[conn])
		// sbc.MapUDPAddrs[ser] = remoteaddr
		// log.Println("MapUDPAddress:", sbc.MapUDPAddrs[ser])
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)

		go sbc.sendResponse(conn, p)
	}
}
