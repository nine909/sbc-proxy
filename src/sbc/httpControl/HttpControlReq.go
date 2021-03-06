package httpControl

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"sbc/logs"
	"strings"
	"time"

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

func TestFlow(w http.ResponseWriter, r *http.Request, session sessions.Session, ps martini.Params) {

	fmt.Fprintf(w, "hello, %s!\n", ps["portgu"])

	DoUDP("MO:uID2x0Xnpj", "MO", "127.0.0.1:7078", "61294")
	time.Sleep(2 * time.Second)
	DoUDP("MT:uID2x0Xnpj", "MT", "127.0.0.1:61294", "39165")
	time.Sleep(2 * time.Second)
	DoUDP("MT:uID2x0Xnpj", "MT", "127.0.0.1:8088", "60539")
	time.Sleep(2 * time.Second)
	DoUDP("MO:uID2x0Xnpj", "MO", "127.0.0.1:60539", "19762")

}

func DoUDP(session, handler, oldPort, port string) {
	// log.SetFlags(log.Lshortfile)

	keyStore = session
	value, err := scn.Fetch(keyStore)
	if err != nil {
		logs.Logger.Debug(err)
		value = NewSBCServer()
	}
	sbc := value.(*Sbc)
	logs.Logger.Debug(sbc)
	sbc.handler = handler
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		sbc.StartServer(oldPort, port)

	}(&wg)
	wg.Wait()
	logs.Logger.Debug("UPDServer : ", sbc)
	logs.Logger.Debug("SBC Clients:", sbc.clients)
	for key, val := range sbc.clients {
		logs.Logger.Debug("sbc.Client", key, val)
	}

	errStore := scn.Store(keyStore, sbc)
	if errStore != nil {
		logs.Logger.Debug(errStore)
	}
}

func Hello(w http.ResponseWriter, r *http.Request, session sessions.Session, ps martini.Params) {

	// log.SetFlags(log.Lshortfile)
	fmt.Fprintf(w, "hello, %s!\n", ps["portgu"])

	value, err := scn.Fetch(keyStore)
	if err != nil {
		logs.Logger.Debug(err)
		value = NewSBCServer()
	}
	sbc := value.(*Sbc)
	logs.Logger.Debug(sbc)
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
	logs.Logger.Debug("UPDServer : ", sbc)
	logs.Logger.Debug("SBC Clients:", sbc.clients)
	for key := range sbc.clients {
		logs.Logger.Debug("sbc.Client: ", sbc.clients[key].addr)
	}

	errStore := scn.Store(keyStore, sbc)
	if errStore != nil {
		logs.Logger.Debug(errStore)
	}

}

func DeleteWTF(w http.ResponseWriter, r *http.Request, session sessions.Session, ps martini.Params) {

	fmt.Fprintf(w, "delete, %s!\n", ps["delete"])
	logs.Logger.Debug("Delete Handler")
	value, err := scn.Fetch(keyStore)
	if err != nil {
		logs.Logger.Debug(err)
	}
	sbc := value.(*Sbc)
	logs.Logger.Debug(keyStore, sbc)

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
		logs.Logger.Debug(err)
	}
	sbc := value.(*Sbc)
	logs.Logger.Debug(keyStore, sbc)

	logs.Logger.Debug("SBC Clients:", sbc.clients)
	for key := range sbc.clients {
		logs.Logger.Debug("sbc.Client: ", sbc.clients[key].addr)
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

	logs.Logger.Debug("SDP Encode: ", sdp.SDP)
	logs.Logger.Debug("SDP X-Session: ", sdp.XSession)
	logs.Logger.Debug("SDP result code: ", sdp.Resultcode)
	logs.Logger.Debug("SDP dev msg code: ", sdp.Deverlopermessage)

	logs.Logger.Debug("Send request To P-WRTC ", sdp.CallbackAddr)
	ccri := msg.ConstructCCR_I(sdp.CallbackSession)

	an, errhttp := RequestHTTTP(sdp.CallbackAddr+"/CCR-I/"+sdp.CallbackSession+"?", ccri)
	if errhttp != nil {
		fmt.Fprintf(w, "error")
		return
	}

	logs.Logger.Debug("Recieve response from P-WRTC")
	logs.Logger.Debug("Data: ", an)

	var aport string
	var vport string
	for {
		ap := rand.Int()
		vp := rand.Int()
		aport = strconv.Itoa(int(ap))[:5]
		vport = strconv.Itoa(int(vp))[:5]
		break
		// _, aerr := net.ResolveUDPAddr("udp", ":"+aport)
		// _, verr := net.ResolveUDPAddr("udp", ":"+vport)
		// if aerr == nil && verr == nil {
		// 	log.Println("new audio port:", aport)
		// 	log.Println("new video port:", vport)
		// 	break
		// }

	}

	//decode base64
	desdp, _ := b64.StdEncoding.DecodeString(sdp.SDP)
	logs.Logger.Debug("SDP Decode: ", string(desdp))
	mediaDesc := SdpParser(desdp, aport, vport)
	logs.Logger.Debug(mediaDesc.ip)

	logs.Logger.Debug("Old port: ", mediaDesc.audio.port)
	logs.Logger.Debug("New port: ", aport)
	s := string(desdp)
	oip := mediaDesc.ip
	oldport := strconv.Itoa(mediaDesc.audio.port)
	newSdp := strings.Replace(s, oldport, aport, -1)
	newSdp = strings.Replace(newSdp, "c=IN IP4 "+oip, "c=IN IP4 "+conf.Conf.Localip, -1)
	//	c=IN IP4 192.168.0.32

	logs.Logger.Debug("New sdp: ", newSdp)
	//start RTP

	ss := strings.Split(sdp.CallbackSession, ":")
	// ip, port := s[0], s[1]
	handler := ss[0]

	rtpMapping(sdp.CallbackAddr+sdp.CallbackSession, handler, mediaDesc.ip+":"+oldport, aport)

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

	logs.Logger.Debug("SDP Encode: ", sdp.SDP)
	logs.Logger.Debug("SDP X-Session: ", sdp.XSession)

	logs.Logger.Debug("Send request To P-WRTC ", sdp.CallbackAddr)
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
		break
		// _, aerr := net.ResolveUDPAddr("udp", ":"+aport)
		// _, verr := net.ResolveUDPAddr("udp", ":"+vport)
		// if aerr == nil && verr == nil {
		// 	log.Println("new audio port:", aport)
		// 	log.Println("new video port:", vport)
		// 	break
		// }
	}

	//decode base64
	desdp, _ := b64.StdEncoding.DecodeString(sdp.SDP)
	logs.Logger.Debug("SDP Decode: ", string(desdp))
	mediaDesc := SdpParser(desdp, aport, vport)

	logs.Logger.Debug("Old port: ", mediaDesc.audio.port)
	logs.Logger.Debug("New port: ", aport)
	s := string(desdp)
	oip := mediaDesc.ip
	oldport := strconv.Itoa(mediaDesc.audio.port)
	newSdp := strings.Replace(s, oldport, aport, -1)
	newSdp = strings.Replace(newSdp, "c=IN IP4 "+oip, "c=IN IP4 "+conf.Conf.Localip, -1)

	logs.Logger.Debug(mediaDesc.ip)
	logs.Logger.Debug(newSdp)

	//start RTP

	ss := strings.Split(sdp.CallbackSession, ":")
	// ip, port := s[0], s[1]
	handler := ss[0]
	rtpMapping(sdp.CallbackAddr+sdp.CallbackSession, handler, mediaDesc.ip+":"+oldport, aport)

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

func rtpMapping(session, handler, uri, port string) {

	// log.SetFlags(log.Lshortfile)
	value, err := scn.Fetch(session)
	if err != nil {
		logs.Logger.Debug(err)
		value = NewSBCServer()
	}

	sbc := value.(*Sbc)
	logs.Logger.Debug("Session get ", session, sbc)

	sbc.handler = handler
	for key := range sbc.clients {
		logs.Logger.Debug("KeyName:", key, sbc.clients[key].addr)

	}
	logs.Logger.Debug(sbc)
	var wg sync.WaitGroup
	// var sbc *Sbc

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		sbc.StartServer(uri, port)
		//		sbc.StartServer(vport)

	}(&wg)
	wg.Wait()
	logs.Logger.Debug("UPDServer : ", sbc)
	logs.Logger.Debug("SBC Clients:", sbc.clients)
	for key := range sbc.clients {
		logs.Logger.Debug("sbc.Client: ", sbc.clients[key].addr)
	}

	errStore := scn.Store(session, sbc)
	if errStore != nil {
		logs.Logger.Debug(errStore)
	}

	logs.Logger.Debug("Session increment ", session, sbc)
}
func UnResoreceAllocate1(w http.ResponseWriter, r *http.Request, session sessions.Session, ps martini.Params) {

	decoder := json.NewDecoder(r.Body)
	var sdp sdpMessages
	err := decoder.Decode(&sdp)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	logs.Logger.Debug("SDP Encode: ", sdp.SDP)
	logs.Logger.Debug("SDP X-Session: ", sdp.XSession)
	logs.Logger.Debug("SDP result code: ", sdp.Resultcode)
	logs.Logger.Debug("SDP dev msg code: ", sdp.Deverlopermessage)

	logs.Logger.Debug("Send request To P-WRTC ", sdp.CallbackAddr)
	ccrt := msg.ConstructCCR_T(sdp.CallbackSession)

	an, errhttp := RequestHTTTP(sdp.CallbackAddr+"/CCR-T/"+sdp.CallbackSession+"?", ccrt)
	if errhttp != nil {
		fmt.Fprintf(w, "error")
		return
	}

	logs.Logger.Debug("Recieve response from P-WRTC")
	logs.Logger.Debug("Data: ", an)

	//start RTP
	logs.Logger.Debug("Delete Handler")
	value, err := scn.Fetch(sdp.CallbackAddr + sdp.CallbackSession)
	if err != nil {
		logs.Logger.Debug(err)
	}
	sbc := value.(*Sbc)
	// log.Println("SBC Clients:", sbc.portOpened)
	// for key := range sbc.clients {
	// 	log.Println("sbc.Client: ", sbc.clients[key].addr)
	// }

	sbc.DeletePort()

	//end rtp

	//encode base64
	sdpRes := &sdpMessages{
		SDP:               sdp.SDP,
		XSession:          sdp.XSession,
		CallbackAddr:      sdp.CallbackAddr,
		CallbackSession:   sdp.CallbackSession,
		Resultcode:        "200",
		Deverlopermessage: "OK"}

	res2B, _ := json.Marshal(sdpRes)
	logs.Logger.Debug(string(res2B))
	fmt.Fprintf(w, string(res2B))
}
func UnResoreceAllocate2(w http.ResponseWriter, r *http.Request, session sessions.Session, ps martini.Params) {
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

	logs.Logger.Debug("SDP Encode: ", sdp.SDP)
	logs.Logger.Debug("SDP X-Session: ", sdp.XSession)

	logs.Logger.Debug("Send request To P-WRTC ", sdp.CallbackAddr)
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
			logs.Logger.Debug("new audio port:", aport)
			logs.Logger.Debug("new video port:", vport)
			break
		}
	}

	//decode base64
	desdp, _ := b64.StdEncoding.DecodeString(sdp.SDP)
	logs.Logger.Debug("SDP Decode: ", string(desdp))
	mediaDesc := SdpParser(desdp, aport, vport)

	logs.Logger.Debug("Old port: ", mediaDesc.audio.port)
	logs.Logger.Debug("New port: ", aport)
	s := string(desdp)
	oip := mediaDesc.ip
	oldport := strconv.Itoa(mediaDesc.audio.port)
	newSdp := strings.Replace(s, oldport, aport, -1)
	newSdp = strings.Replace(newSdp, "c=IN IP4 "+oip, "c=IN IP4 "+conf.Conf.Localip, -1)

	logs.Logger.Debug(mediaDesc.ip)
	logs.Logger.Debug(newSdp)

	//start RTP

	// rtpMapping(sdp.CallbackAddr+sdp.CallbackSession, mediaDesc.ip+":"+oldport, aport)

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
	logs.Logger.Debug(string(res2B))
	fmt.Fprintf(w, string(res2B))
}
func Ccru() {
	for {
		time.Sleep(10 * time.Second)
		logs.Logger.Debug("Hello")

	}

}
