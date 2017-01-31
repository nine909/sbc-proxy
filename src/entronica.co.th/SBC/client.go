package main

import (
	"fmt"

	"github.com/pixelbender/go-stun/stun"
)

func main() {
	addr, err := stun.Lookup("stun:127.0.0.1:6060", "username", "password")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(addr)
	}
}
