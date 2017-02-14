package stunSBC

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

type Sbc struct {
	Conn *net.UDPConn
}

func NewSBCServer(port string) *Sbc {

	sbc := Open(port)
	log.Println("sbc: ", sbc)
	log.Println("sbc.Conn: ", sbc.Conn.LocalAddr().String())
	go sbc.UDPServer(sbc.Conn)
	log.Println("OUT sbc: ", sbc)
	log.Println("OUT sbc.Conn: ", sbc.Conn)
	return &sbc

}

func Open(port string) Sbc {
	UDPPort, _ := strconv.Atoi(port)

	addr := net.UDPAddr{
		Port: UDPPort,
		IP:   net.ParseIP(""),
		// IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		// return nil
	}
	log.Println("ser: ", ser)
	// log.Println("sbc: ", sbc)
	// log.Println("sbc.Conn: ", sbc.Conn.LocalAddr().String())
	return Sbc{Conn: ser}
	// sbc.UDPServer(ser)
}

func SecondSBCServer(sbc *Sbc, port string) *Sbc {
	sbc2 := Open(port)
	log.Println("Nsbc: ", &sbc)
	log.Println("Nsbc.Conn: ", sbc.Conn.LocalAddr().String())
	log.Println("Nsbc2.Conn: ", sbc2.Conn.LocalAddr().String())
	go sbc.UDPServer(sbc2.Conn)
	log.Println("NOUT sbc2: ", sbc2)
	log.Println("NOUT sbc2.Conn: ", sbc2.Conn)
	return &sbc2

}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP([]byte("From server: Hello I got your mesage "), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}

func (sbc *Sbc) UDPServer(ser *net.UDPConn) {
	p := make([]byte, 2048)
	// UDPPort, _ := strconv.Atoi(port)

	// addr := net.UDPAddr{
	// 	Port: UDPPort,
	// 	IP:   net.ParseIP(""),
	// 	// IP:   net.ParseIP("127.0.0.1"),
	// }
	// ser, err := net.ListenUDP("udp", &addr)
	// if err != nil {
	// 	fmt.Printf("Some error %v\n", err)
	// 	return
	// }
	// sbc.Conn = ser
	for {
		_, remoteaddr, err := ser.ReadFromUDP(p)
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		go sendResponse(ser, remoteaddr)
	}
}

func (sbc *Sbc) Running(port string) {
	p := make([]byte, 2048)
	UDPPort, _ := strconv.Atoi(port)

	addr := net.UDPAddr{
		Port: UDPPort,
		IP:   net.ParseIP(""),
		// IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	sbc.Conn = ser
	for {
		_, remoteaddr, err := ser.ReadFromUDP(p)
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		go sendResponse(ser, remoteaddr)
	}
}

func Handler(port string) {

	// if port == "6060"{
	// 	srv := NewSBCServer(port)
	// 	log.Println("Port as : ", srv)
	// }else{

	// }
	// NewSBCServer(port)
	// log.Println("Port : ", srv.port1)
	// log.Println("srv.lcon1.LocalAddr().Network() : ", srv.lcon1.LocalAddr().Network())
	// log.Println("srv.lcon1.LocalAddr().String() : ", srv.lcon1.LocalAddr().String())
	// go NewSBCServer(port)
	// log.Println("Port : ", port)
	// for {
	// 	log.Println("Port : ", port)
	// 	time.Sleep(time.Second * 2)
	// }

	// portOpen := make(chan string, 1)
	// for {
	// 	go NewSBCServer(port)
	// }
}
