package main

import (
	"fmt"
	"net/http"
	"sbc/logs"
	// "os"
	"sbc/conf"
	"sbc/httpControl"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
)

func main() {
	logs.Logger.Info("Start server at:")
	// log.WithFields(log.Fields{}).Info("A group of walrus emerges from the ocean")

	// log.WithFields(log.Fields{}).Warn("The group's number increased tremendously!")

	// log.WithFields(log.Fields{}).Fatal("The ice breaks!")

	config := conf.ReadConfig()

	m := martini.Classic()

	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("my_session", store))

	m.Get("/", httpControl.Index)
	m.Get("/hello/:portgu", httpControl.Hello)
	m.Get("/delete", httpControl.DeleteWTF)
	m.Get("/list", httpControl.Lists)
	m.Get("/testflow", httpControl.TestFlow)

	m.Post("/p-SSF/1.0.0/SBC/ResourceAllocate1/:uid", httpControl.TestClient)
	m.Post("/p-SSF/1.0.0/SBC/ResourceAllocate2/:uid", httpControl.TestClient2)
	m.Post("/p-SSF/1.0.0/SBC/ResourceUnAllocate/:uid", httpControl.UnResoreceAllocate1)
	fmt.Println("Base URL: localhost:" + config.HttpPort)
	fmt.Println(http.ListenAndServe(":"+config.HttpPort, m))

}

func Hello(w http.ResponseWriter, r *http.Request, session *sessions.Session, ps martini.Params) {

	fmt.Fprintf(w, "hello, %s!\n", ps["portgu"])

	//	go stunSBC.ServerListener(ps.ByName("portgu"))

}
