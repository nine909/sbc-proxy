package main

import (
	"fmt"
	"net/http"

	"sbc/httpControl"

	"github.com/julienschmidt/httprouter"
)

func main() {

	router := httprouter.New()
	router.GET("/", httpControl.Index)
	router.GET("/hello/:name", httpControl.Hello)

	router.GET("/user/:uid", httpControl.Getuser)
	router.POST("/adduser/:uid", httpControl.Adduser)
	//	router.DELETE("/deluser/:uid", httpControl.deleteuser)
	//	router.PUT("/moduser/:uid", httpControl.modifyuser)

	fmt.Println(http.ListenAndServe(":8080", router))

}
