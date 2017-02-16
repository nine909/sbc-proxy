package httpControl

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"

	conf "sbc/conf"

	b64 "encoding/base64"
	msg "sbc/messages"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"

	//	"io"
	//	"io/ioutil"

	"github.com/ernado/sdp"
	//	"os"
	"strconv"
	"sync"
	//	"github.com/julienschmidt/httprouter"
	"github.com/tideland/golib/scene"
)

// var OrigSdp OriginSdp

type sdpMessages struct {
	SDP             string
	XSession        string `json:"x-session"`
	CallbackAddr    string `json:"Callback-Address"`
	CallbackSession string `json:"Callback-Session"`

	Resultcode        string `json:"resultcode"`
	Deverlopermessage string `json:"deverlopermessage"`
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
	log.Println("SBC Clients:", sbc.clients)
	for key := range sbc.clients {
		log.Println("sbc.Client: ", sbc.clients[key].addr)
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

	// log.Println("SBC Clients:", sbc.portOpened)
	// for key := range sbc.clients {
	// 	log.Println("sbc.Client: ", sbc.clients[key].addr)
	// }

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

	log.Println("SBC Clients:", sbc.clients)
	for key := range sbc.clients {
		log.Println("sbc.Client: ", sbc.clients[key].addr)
	}

}

func TestClient(w http.ResponseWriter, r *http.Request, session sessions.Session, ps martini.Params) {
	// uid := r.FormValue("uid")
	//	uid := ps.ByName("uid")
	//	fmt.Println(r.))
	//	fmt.Fprintf(w, "you are add user %s", uid)

	decoder := json.NewDecoder(r.Body)
	var sdp sdpMessages
	err := decoder.Decode(&sdp)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	log.Println("SDP Encode: ", sdp.SDP)
	log.Println("SDP X-Session: ", sdp.XSession)
	log.Println("SDP result code: ", sdp.Resultcode)
	log.Println("SDP dev msg code: ", sdp.Deverlopermessage)

	log.Println("Send request To P-WRTC ", sdp.CallbackAddr)
	ccri := msg.ConstructCCR_I(sdp.CallbackSession)

	an, errhttp := RequestHTTTP(sdp.CallbackAddr+"/CCR-I/"+sdp.CallbackSession+"?", ccri)
	if errhttp != nil {
		fmt.Fprintf(w, "error")
		return
	}

	log.Println("Recieve response from P-WRTC")
	log.Println("Data: ", an)

	var aport string
	var vport string
	for {
		ap := rand.Int()
		vp := rand.Int()
		aport = strconv.Itoa(int(ap))[:5]
		vport = strconv.Itoa(int(vp))[:5]
		_, aerr := net.ResolveUDPAddr("udp", ":"+aport)
		_, verr := net.ResolveUDPAddr("udp", ":"+vport)
		if aerr == nil && verr == nil {
			log.Println("new audio port:", aport)
			log.Println("new video port:", vport)
			break
		}
	}

	//decode base64
	desdp, _ := b64.StdEncoding.DecodeString(sdp.SDP)
	log.Println("SDP Decode: ", string(desdp))
	mediaDesc, _ := sbpParser(desdp, aport, vport)
	log.Println(mediaDesc.ip)

	log.Println("Old port: ", mediaDesc.audio.port)
	log.Println("New port: ", aport)
	s := string(desdp)
	oip := mediaDesc.ip
	oldport := strconv.Itoa(mediaDesc.audio.port)
	newSdp := strings.Replace(s, oldport, aport, -1)
	newSdp = strings.Replace(newSdp, "c=IN IP4 "+oip, "c=IN IP4 "+conf.Conf.Localip, -1)
	//	c=IN IP4 192.168.0.32

	log.Println("New sdp: ", newSdp)
	//start RTP
	log.SetFlags(log.Lshortfile)
	value, err := scn.Fetch(sdp.CallbackAddr + sdp.CallbackSession)
	if err != nil {
		log.Println(err)
		value = NewSBCServer()
	}
	sbc := value.(*Sbc)
	log.Println(sbc)
	var wg sync.WaitGroup
	// var sbc *Sbc

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		sbc.StartServer(aport)
		sbc.StartServer(vport)

	}(&wg)
	wg.Wait()
	log.Println("UPDServer : ", sbc)
	log.Println("SBC Clients:", sbc.clients)
	for key := range sbc.clients {
		log.Println("sbc.Client: ", sbc.clients[key].addr)
	}

	errStore := scn.Store(sdp.CallbackAddr+sdp.CallbackSession, sbc)
	if errStore != nil {
		log.Println(errStore)
	}
	//end rtp

	//encode base64
	sEnc := b64.StdEncoding.EncodeToString([]byte(newSdp))
	fmt.Printf(sEnc)
	sdpRes := &sdpMessages{
		SDP:               sEnc,
		XSession:          sdp.XSession,
		CallbackAddr:      sdp.CallbackAddr,
		CallbackSession:   sdp.CallbackSession,
		Resultcode:        "200",
		Deverlopermessage: "OK"}

	res2B, _ := json.Marshal(sdpRes)
	fmt.Println(string(res2B))
	fmt.Fprintf(w, string(res2B))
}

func TestClient2(w http.ResponseWriter, r *http.Request, session sessions.Session, ps martini.Params) {
	// uid := r.FormValue("uid")
	//	uid := ps.ByName("uid")
	//	fmt.Println(r.))
	//	fmt.Fprintf(w, "you are add user %s", uid)
	decoder := json.NewDecoder(r.Body)
	var sdp sdpMessages
	err := decoder.Decode(&sdp)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	log.Println("SDP Encode: ", sdp.SDP)
	log.Println("SDP X-Session: ", sdp.XSession)

	log.Println("Send request To P-WRTC ", sdp.CallbackAddr)
	//	ccri := msg.ConstructCCR_I()

	//	an := RequestHTTTP(sdp.CallbackAddr, ccri)

	//	log.Println("Recieve response from P-WRTC")
	//	log.Println("Data: ", an)

	var aport string
	var vport string
	for {
		ap := rand.Int()
		vp := rand.Int()
		aport = strconv.Itoa(int(ap))[:5]
		vport = strconv.Itoa(int(vp))[:5]
		_, aerr := net.ResolveUDPAddr("udp", ":"+aport)
		_, verr := net.ResolveUDPAddr("udp", ":"+vport)
		if aerr == nil && verr == nil {
			log.Println("new audio port:", aport)
			log.Println("new video port:", vport)
			break
		}
	}

	//decode base64
	desdp, _ := b64.StdEncoding.DecodeString(sdp.SDP)
	log.Println("SDP Decode: ", string(desdp))
	mediaDesc, _ := sbpParser(desdp, aport, vport)

	log.Println("Old port: ", mediaDesc.audio.port)
	log.Println("New port: ", aport)
	s := string(desdp)
	oip := mediaDesc.ip
	oldport := strconv.Itoa(mediaDesc.audio.port)
	newSdp := strings.Replace(s, oldport, aport, -1)
	newSdp = strings.Replace(newSdp, "c=IN IP4 "+oip, "c=IN IP4 "+conf.Conf.Localip, -1)

	log.Println(mediaDesc.ip)
	log.Println(newSdp)

	//start RTP
	log.SetFlags(log.Lshortfile)
	value, err := scn.Fetch(sdp.CallbackAddr + sdp.CallbackSession)
	if err != nil {
		log.Println(err)
		value = NewSBCServer()
	}
	sbc := value.(*Sbc)
	log.Println(sbc)
	var wg sync.WaitGroup
	// var sbc *Sbc

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		sbc.StartServer(aport)
		sbc.StartServer(vport)

	}(&wg)
	wg.Wait()
	log.Println("UPDServer : ", sbc)
	log.Println("SBC Clients:", sbc.clients)
	for key := range sbc.clients {
		log.Println("sbc.Client: ", sbc.clients[key].addr)
	}

	errStore := scn.Store(keyStore, sbc)
	if errStore != nil {
		log.Println(errStore)
	}
	//end rtp

	//encode base64
	sEnc := b64.StdEncoding.EncodeToString([]byte(newSdp))
	// fmt.Fprintf(w, sEnc)
	sdpRes := &sdpMessages{
		SDP:               sEnc,
		XSession:          sdp.XSession,
		CallbackAddr:      sdp.CallbackAddr,
		CallbackSession:   sdp.CallbackSession,
		Resultcode:        "200",
		Deverlopermessage: "OK"}
	res2B, _ := json.Marshal(sdpRes)
	fmt.Println(string(res2B))
	fmt.Fprintf(w, string(res2B))
}
func sbpParser(sdpByte []byte, aport, vport string) (OriginSdp, string) {
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
	log.Println("Decoded session", m.Name)
	log.Println("Info:", m.Info)
	log.Println("Origin:", m.Origin)
	log.Println("IP 4: ", m.Origin.Address)
	log.Println("IP 4: ", m.Timing)
	log.Println("NetworkType: ", m.Connection.NetworkType)
	log.Println("AddressType: ", m.Connection.AddressType)
	log.Println("IP: ", m.Connection.IP)
	log.Println("TTL: ", m.Connection.TTL)
	log.Println("Addresses: ", m.Connection.Addresses)

	orig := OriginSdp{}
	medias1 := sdp.Medias{}

	orig.ip = m.Origin.Address
	isMultiMediaAudio := false
	isMultiMediaVideo := false
	// var isRemove []int
	for i, media := range m.Medias {
		log.Println("=======================")
		log.Println("Type: ", media.Description.Type)
		log.Println("Port: ", media.Description.Port)
		log.Println("PortsNumber: ", media.Description.PortsNumber)
		log.Println("Protocol: ", media.Description.Protocol)
		log.Println("Format: ", media.Description.Format)
		log.Println("Medias Connection: ", media.Connection)
		log.Println("Medias Attributes: ", media.Attributes)
		log.Println("Medias Encryption: ", media.Encryption)
		log.Println("Medias Bandwidths: ", media.Bandwidths)

		a, _ := strconv.Atoi(aport)
		v, _ := strconv.Atoi(vport)
		switch media.Description.Type {
		case "audio":
			if !isMultiMediaAudio {
				orig.audio.port = media.Description.Port
				orig.audio.portsNumber = media.Description.PortsNumber
				orig.audio.protocal = media.Description.Protocol
				m.Medias[i].Description.Port = a
				isMultiMediaAudio = true
				medias1 = append(medias1, m.Medias[i])
			}
		case "video":
			if !isMultiMediaVideo {
				orig.video.port = media.Description.Port
				orig.video.portsNumber = media.Description.PortsNumber
				orig.video.protocal = media.Description.Protocol
				m.Medias[i].Description.Port = v
				isMultiMediaVideo = true
				medias1 = append(medias1, m.Medias[i])
			}
		}

	}
	// defining medias
	m.Medias = medias1
	//replace ip
	//	m.Origin.Address = conf.Conf.Localip
	m.Connection.IP = net.ParseIP(conf.Conf.Localip)

	newSdp := constructSdp(m)

	//encode base64
	//	log.Print("dddddddddddddddddd", orig.audio.port)
	return orig, newSdp

}

func constructSdp(me *sdp.Message) string {
	var (
		s sdp.Session
		b []byte
	)

	// defining message
	m := me
	//	m.AddFlag("nortpproxy:yes")

	// appending message to session
	s = m.Append(s)

	// appending session to byte buffer
	b = s.AppendTo(b)
	log.Println("SDP Construct :", string(b)+"\n")
	return string(b) + "\n"
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
