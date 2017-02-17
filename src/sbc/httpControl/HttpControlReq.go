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

	//	"os"
	"strconv"
	"sync"
	//	"github.com/julienschmidt/httprouter"
	"github.com/tideland/golib/scene"
)

// var OrigSdp OriginSdp

type sdpMessages struct {
	SDP               string
	XSession          string `json:"x-session"`
	CallbackAddr      string `json:"Callback-Address"`
	CallbackSession   string `json:"Callback-Session"`
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
	oldPort := "127.0.0.1:1234"
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		sbc.StartServer(oldPort, port)

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
	mediaDesc := SdpParser(desdp, aport, vport)
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

	rtpMapping(sdp.CallbackAddr+sdp.CallbackSession, mediaDesc.ip+":"+oldport, aport)

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
	mediaDesc := SdpParser(desdp, aport, vport)

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

	rtpMapping(sdp.CallbackAddr+sdp.CallbackSession, mediaDesc.ip+":"+oldport, aport)

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

func rtpMapping(session, uri, port string) {

	log.SetFlags(log.Lshortfile)
	value, err := scn.Fetch(session)
	if err != nil {
		log.Println(err)
		value = NewSBCServer()
	}

	sbc := value.(*Sbc)
	log.Println("Session get ", session, sbc)

	for key := range sbc.clients {
		log.Println("KeyName:", key, sbc.clients[key].addr)

	}
	log.Println(sbc)
	var wg sync.WaitGroup
	// var sbc *Sbc

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		sbc.StartServer(uri, port)
		//		sbc.StartServer(vport)

	}(&wg)
	wg.Wait()
	log.Println("UPDServer : ", sbc)
	log.Println("SBC Clients:", sbc.clients)
	for key := range sbc.clients {
		log.Println("sbc.Client: ", sbc.clients[key].addr)
	}

	errStore := scn.Store(session, sbc)
	if errStore != nil {
		log.Println(errStore)
	}

	log.Println("Session increment ", session, sbc)
}

func UnResoreceAllocate(w http.ResponseWriter, r *http.Request, session sessions.Session, ps martini.Params) {

}
