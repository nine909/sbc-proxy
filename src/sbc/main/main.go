package main

import (
	"fmt"
	"net/http"
	"sbc/conf"

	"sbc/httpControl"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
)

func main() {
	config := conf.ReadConfig()
	fmt.Println("Base URL: ", config.Baseurl)

	m := martini.Classic()

	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("my_session", store))

	m.Get("/", httpControl.Index)
	m.Get("/hello/:portgu", httpControl.Hello)

	m.Post("/Test/:uid", httpControl.TestClient)

	fmt.Println(http.ListenAndServe(":8080", m))

}

func Hello(w http.ResponseWriter, r *http.Request, session *sessions.Session, ps martini.Params) {

	fmt.Fprintf(w, "hello, %s!\n", ps["portgu"])

	//	go stunSBC.ServerListener(ps.ByName("portgu"))

}
