package main

import (
	"fmt"
	"net/http"

	"entronica.co.th/SBC/httpControl"
	"github.com/julienschmidt/httprouter"
)

func main() {

	router := httprouter.New()
	router.GET("/", httpControl.Index)
	router.GET("/hello/:name", httpControl.Hello)

	//	router.GET("/user/:uid", httpControl.getuser)
	//	router.POST("/adduser/:uid", httpControl.adduser)
	//	router.DELETE("/deluser/:uid", httpControl.deleteuser)
	//	router.PUT("/moduser/:uid", httpControl.modifyuser)

	fmt.Println(http.ListenAndServe(":8080", router))

}
