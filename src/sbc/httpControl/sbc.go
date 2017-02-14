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
	clients map[*net.UDPConn]Client
	// Conn    []*net.UDPConn
	// MapUDPAddrs map[*net.UDPConn]*net.UDPAddr
	// addrForward map[*net.UDPConn][]*net.UDPConn
}

type Client struct {
	addr *net.UDPAddr
}

func NewSBCServer() *Sbc {
	return &Sbc{clients: make(map[*net.UDPConn]Client)}
}

func (sbc *Sbc) Open(port string) (*net.UDPConn, error) {
	UDPPort, _ := strconv.Atoi(port)

	addr := net.UDPAddr{
		Port: UDPPort,
		IP:   net.ParseIP(""),
		// IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
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
	client := Client{}
	sbc.clients = append(sbc.clients, client)
	// sbc.Conn = append(sbc.Conn, conn)
	log.Println("UDP Server Started!!!")
	go sbc.UDPServer(conn)
	return nil
}

func (sbc *Sbc) DeletePort() {
	log.Println("ConnList : ", len(sbc.Conn))
	for _, conn := range sbc.Conn {
		conn.Close()
		log.Println("Port Closed: ", conn.LocalAddr().String())
	}
	sbc.Conn = sbc.Conn[:0]
}

func (sbc *Sbc) sendResponse(conn *net.UDPConn, p []byte) {

	remoteAddr := sbc.addrForward[conn]
	// _, err := conn.WriteToUDP([]byte("From server: Hello I got your mesage "), addr)

	log.Println(sbc.MapUDPAddrs)
	for _, addr := range remoteAddr {
		log.Println(addr)
		log.Println(sbc.MapUDPAddrs[addr])
		_, err := conn.WriteToUDP([]byte("From server: Hello I got your mesage "), sbc.MapUDPAddrs[addr])
		if err != nil {
			fmt.Printf("Couldn't send response %v", err)
		}
	}

	// _, err := conn.WriteToUDP(p, addr)
	// if err != nil {
	// 	fmt.Printf("Couldn't send response %v", err)
	// }
}

func (sbc *Sbc) UDPServer(ser *net.UDPConn) error {

	// if len(sbc.Conn) <= 2 {
	// 	return
	// }
	// addrForwardMaps := make(map[*net.UDPConn][]*net.UDPConn)
	if len(sbc.Conn) >= 2 {
		for i, conn := range sbc.Conn {
			if ser != conn {
				log.Println(sbc.Conn)
				s1 := sbc.Conn[:i]
				s2 := sbc.Conn[i+1]
				log.Println(conn)
				log.Println(s1)
				log.Println(s2)
				if len(sbc.Conn) == 2 {
					sbc.addrForward[conn] = []*net.UDPConn{s2}
					sbc.addrForward[s2] = []*net.UDPConn{conn}
					break
				}
			}
		}
		// log.Println(addrForwardMaps)
		// sbc.addrForward = addrForwardMaps
		log.Println(sbc.addrForward)
	}
	log.Println("UDPDetail:", &ser)
	log.Println(&sbc.addrForward)
	p := make([]byte, 2048)
	for {
		log.Println("Waiting Incoming...")
		_, remoteaddr, err := ser.ReadFromUDP(p)
		log.Println("Remote Address:", remoteaddr)
		sbc.MapUDPAddrs[ser] = remoteaddr
		log.Println("MapUDPAddress:", sbc.MapUDPAddrs[ser])
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			// continue
			return err
		}
		go sbc.sendResponse(ser, p)
	}
}
