package main

import (
	"fmt"
	"net/http"

	"sbc/httpControl"

	//	"github.com/julienschmidt/httprouter"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
)

func main() {

	m := martini.Classic()

	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("my_session", store))

	//	router := httprouter.New()
	//	m.GET("/", httpControl.Index)
	m.Get("/hello/:portgu", httpControl.Hello)
	m.Get("/second/:portsec", httpControl.Second)
	m.Get("/delete", httpControl.DeleteWTF)
	m.Get("/list", httpControl.Lists)

	//	m.GET("/user/:uid", httpControl.Getuser)
	//	m.POST("/adduser/:uid", httpControl.Adduser)
	//	router.DELETE("/deluser/:uid", httpControl.deleteuser)
	//	router.PUT("/moduser/:uid", httpControl.modifyuser)

	fmt.Println(http.ListenAndServe(":8080", m))

}

func Hello(w http.ResponseWriter, r *http.Request, session *sessions.Session, ps martini.Params) {

	fmt.Fprintf(w, "hello, %s!\n", ps["portgu"])

	//	go stunSBC.ServerListener(ps.ByName("portgu"))

}
