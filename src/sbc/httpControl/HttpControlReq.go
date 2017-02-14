package httpControl

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	//	"github.com/julienschmidt/httprouter"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
	"github.com/tideland/golib/scene"
)

type test_struct struct {
	Test string
}

func Index(w http.ResponseWriter, r *http.Request, session sessions.Session) {

	fmt.Fprintln(w, "Welcomess!\n")
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

func Getuser(w http.ResponseWriter, r *http.Request, session sessions.Session, ps martini.Params) {
	uid := ps["uid"]
	fmt.Fprintf(w, "you are get user %s", uid)

}

func modifyuser(w http.ResponseWriter, r *http.Request, session sessions.Session, ps martini.Params) {
	uid := ps["uid"]
	fmt.Fprintf(w, "you are modify user %s", uid)
}

func deleteuser(w http.ResponseWriter, r *http.Request, session sessions.Session, ps martini.Params) {
	uid := ps["uid"]
	fmt.Fprintf(w, "you are delete user %s", uid)
}

func Adduser(w http.ResponseWriter, r *http.Request, session sessions.Session, ps martini.Params) {
	// uid := r.FormValue("uid")
	//	uid := ps.ByName("uid")
	//	fmt.Println(r.))
	//	fmt.Fprintf(w, "you are add user %s", uid)
	decoder := json.NewDecoder(r.Body)
	var t test_struct
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	log.Println(t.Test)
}

//func main() {
//	fmt.Print("asd")
//	router := httprouter.New()
//	router.GET("/", Index)
//	router.GET("/hello/:name", Hello)

//	router.GET("/user/:uid", getuser)
//	router.POST("/adduser/:uid", adduser)
//	router.DELETE("/deluser/:uid", deleteuser)
//	router.PUT("/moduser/:uid", modifyuser)

//	fmt.Println(http.ListenAndServe(":8080", router))

//}
