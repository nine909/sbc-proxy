package httpControl

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	//"strings"

	conf "sbc/conf"

	b64 "encoding/base64"
	msg "sbc/messages"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"

	//	"io"
	//	"io/ioutil"

	"github.com/ernado/sdp"
	//	"os"

	"sync"
	//	"github.com/julienschmidt/httprouter"
	"github.com/tideland/golib/scene"
)

// var OrigSdp OriginSdp

type sdpMessage struct {
	SDP             string
	XSession        string `json:"x-session"`
	CallbackAddr    string `json:"Callback-Address"`
	CallbackSession string `json:"Callback-Session"`
}

func Index(w http.ResponseWriter, r *http.Request, session sessions.Session) {

	fmt.Fprintln(w, "Welcomes!\n")
	//	go stunSBC.ServerListener("6006")
	//	fmt.Fprintln(w, r)
}

var scn = scene.Start()
var keyStore = "sbc001"

func Hello(w http.ResponseWriter, r *http.Request, session sessions.Session, ps martini.Params) {

	log.SetFlags(log.Lshortfile)
	fmt.Fprintf(w, "hello, %s!\n", ps["portgu"])

	value, err := scn.Fetch(keyStore)
	if err != nil {
		log.Println(err)
		value = NewSBCServer()
	}
	sbc := value.(*Sbc)
	log.Println(sbc)
	var wg sync.WaitGroup
	// var sbc *Sbc
	port := ps["portgu"]

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		sbc.StartServer(port)

	}(&wg)
	wg.Wait()
	log.Println("UPDServer : ", sbc)
	for _, conn := range sbc.Conn {
		log.Println("sbc.Conn: ", conn.LocalAddr().String())
	}

	errStore := scn.Store(keyStore, sbc)
	if errStore != nil {
		log.Println(errStore)
	}

}

func DeleteWTF(w http.ResponseWriter, r *http.Request, session sessions.Session, ps martini.Params) {

	fmt.Fprintf(w, "delete, %s!\n", ps["delete"])
	fmt.Println("Delete Handler")
	value, err := scn.Fetch(keyStore)
	if err != nil {
		log.Println(err)
	}
	sbc := value.(*Sbc)
	fmt.Println(keyStore, sbc)

	for _, conn := range sbc.Conn {
		log.Println("Port Num: ", conn.LocalAddr().String())
	}

	sbc.DeletePort()
}

func Lists(w http.ResponseWriter, r *http.Request, session sessions.Session, ps martini.Params) {

	fmt.Fprintf(w, "Show list PortUDP already\n")
	value, err := scn.Fetch(keyStore)
	if err != nil {
		log.Println(err)
	}
	sbc := value.(*Sbc)
	fmt.Println(keyStore, sbc)

	for _, conn := range sbc.Conn {
		log.Println("Port Num: ", conn.LocalAddr().String())
	}

}

func TestClient(w http.ResponseWriter, r *http.Request, session sessions.Session, ps martini.Params) {
	// uid := r.FormValue("uid")
	//	uid := ps.ByName("uid")
	//	fmt.Println(r.))
	//	fmt.Fprintf(w, "you are add user %s", uid)
	decoder := json.NewDecoder(r.Body)
	var sdp sdpMessage
	err := decoder.Decode(&sdp)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	log.Println("SDP Encode: ", sdp.SDP)
	log.Println("SDP X-Session: ", sdp.XSession)

	log.Println("Send request To P-WRTC ", sdp.CallbackAddr)
	ccri := msg.ConstructCCR_I()

	an := RequestHTTTP(sdp.CallbackAddr, ccri)

	log.Println("Recieve response from P-WRTC")
	log.Println("Data: ", an)

	//decode base64
	desdp, _ := b64.StdEncoding.DecodeString(sdp.SDP)
	log.Println("SDP Decode: ", string(desdp))
	mediaDesc, newSdp := sbpParser(desdp)
	fmt.Println(mediaDesc.ip)
	fmt.Println(newSdp)

	//encode base64
	sEnc := b64.StdEncoding.EncodeToString([]byte(newSdp))
	// fmt.Fprintf(w, sEnc)

	sdpRes := &sdpMessage{
		SDP:             sEnc,
		XSession:        sdp.XSession,
		CallbackAddr:    sdp.CallbackAddr,
		CallbackSession: sdp.CallbackSession}
	res2B, _ := json.Marshal(sdpRes)
	fmt.Println(string(res2B))
	fmt.Fprintf(w, string(res2B))

}

func sbpParser(sdpByte []byte) (OriginSdp, string) {
	var (
		s   sdp.Session
		err error
	)

	if s, err = sdp.DecodeSession(sdpByte, s); err != nil {
		log.Fatal("err:", err)
	}
	// for k, v := range s {
	// 	fmt.Println(k, v)
	// }
	d := sdp.NewDecoder(s)
	m := new(sdp.Message)
	if err = d.Decode(m); err != nil {
		log.Fatal("err:", err)
	}
	fmt.Println("Decoded session", m.Name)
	fmt.Println("Info:", m.Info)
	fmt.Println("Origin:", m.Origin)
	fmt.Println("IP 4: ", m.Origin.Address)
	fmt.Println("IP 4: ", m.Timing)
	fmt.Println("NetworkType: ", m.Connection.NetworkType)
	fmt.Println("AddressType: ", m.Connection.AddressType)
	fmt.Println("IP: ", m.Connection.IP)
	fmt.Println("TTL: ", m.Connection.TTL)
	fmt.Println("Addresses: ", m.Connection.Addresses)

	orig := OriginSdp{}
	medias1 := sdp.Medias{}

	orig.ip = m.Origin.Address
	isMultiMediaAudio := false
	isMultiMediaVideo := false
	// var isRemove []int
	for i, media := range m.Medias {
		fmt.Println("=======================")
		fmt.Println("Type: ", media.Description.Type)
		fmt.Println("Port: ", media.Description.Port)
		fmt.Println("PortsNumber: ", media.Description.PortsNumber)
		fmt.Println("Protocol: ", media.Description.Protocol)
		fmt.Println("Format: ", media.Description.Format)
		fmt.Println("Medias Connection: ", media.Connection)
		fmt.Println("Medias Attributes: ", media.Attributes)
		fmt.Println("Medias Encryption: ", media.Encryption)
		fmt.Println("Medias Bandwidths: ", media.Bandwidths)

		switch media.Description.Type {
		case "audio":
			if !isMultiMediaAudio {
				orig.audio.port = media.Description.Port
				orig.audio.portsNumber = media.Description.PortsNumber
				orig.audio.protocal = media.Description.Protocol
				m.Medias[i].Description.Port = 11111
				isMultiMediaAudio = true
				medias1 = append(medias1, m.Medias[i])
			}
		case "video":
			if !isMultiMediaVideo {
				orig.video.port = media.Description.Port
				orig.video.portsNumber = media.Description.PortsNumber
				orig.video.protocal = media.Description.Protocol
				m.Medias[i].Description.Port = 555
				isMultiMediaVideo = true
				medias1 = append(medias1, m.Medias[i])
			}
		}

	}
	// defining medias
	m.Medias = medias1
	//replace ip
	m.Origin.Address = conf.Conf.Localip
	m.Connection.IP = net.ParseIP(conf.Conf.Localip)

	newSdp := constructSdp(m)

	//encode base64

	return orig, newSdp

}

func constructSdp(me *sdp.Message) string {
	var (
		s sdp.Session
		b []byte
	)

	// defining message
	m := me
	// m.AddFlag("recvonly")

	// appending message to session
	s = m.Append(s)

	// appending session to byte buffer
	b = s.AppendTo(b)
	fmt.Println("SDP Construct :", string(b))
	return string(b)
}

type OriginSdp struct {
	ip    string
	audio struct {
		port        int
		portsNumber int
		protocal    string
	}
	video struct {
		port        int
		portsNumber int
		protocal    string
	}
}
