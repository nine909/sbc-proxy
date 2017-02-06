package main

import (
	"fmt"
	"net/http"

	"sbc/httpControl"

	//	"github.com/julienschmidt/httprouter"
	"github.com/martini-contrib/sessions"
)

func main() {

	m := martini.Classic()

	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("my_session", store))

	//	router := httprouter.New()
	m.GET("/", httpControl.Index)
	m.GET("/hello/:portgu", httpControl.Hello)

	m.GET("/user/:uid", httpControl.Getuser)
	m.POST("/adduser/:uid", httpControl.Adduser)
	//	router.DELETE("/deluser/:uid", httpControl.deleteuser)
	//	router.PUT("/moduser/:uid", httpControl.modifyuser)

	fmt.Println(http.ListenAndServe(":8080", m))

}
