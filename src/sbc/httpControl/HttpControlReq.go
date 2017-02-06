package httpControl

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"sbc/stunSBC"
	//	"github.com/julienschmidt/httprouter"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
)

type test_struct struct {
	Test string
}

func Index(w http.ResponseWriter, r *http.Request, session sessions.Session) {

	fmt.Fprintln(w, "Welcomess!\n")
	//	stunSBC.ServerListener("6006")
	//	fmt.Fprintln(w, r)
}

func Hello(w http.ResponseWriter, r *http.Request, session sessions.Session, ps martini.Params) {

	fmt.Fprintf(w, "hello, %s!\n", ps["portgu"])

	go stunSBC.ServerListener(ps["portgu"])

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
