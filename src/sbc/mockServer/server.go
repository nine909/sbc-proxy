package main

import (
	// "encoding/json"
	"fmt"
	"net/http"
)

// type Profile struct {
//   Name    string
//   Hobbies []string
// }

func main() {
	http.HandleFunc("/", foo)
	http.ListenAndServe(":3000", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {
	// profile := Profile{"Alex", []string{"snowboarding", "programming"}}

	// js, err := json.Marshal(profile)
	// if err != nil {
	//   http.Error(w, err.Error(), http.StatusInternalServerError)
	//   return
	// }

	var jsonStr = []byte(`{"Session-Id":"MO:uID2x0Xnpj","Result-Code":"2001","Origin-Host":"OriginHost","Origin-Realm":"OriginRealm","Auth-Application-Id":"4","CC-Request-Type":"1","CC-Request-Number":"0"}`)
	fmt.Println(string(jsonStr))
	w.Header().Set("Content-Type", "application/json")
	// w.Write(js)
	w.Write(jsonStr)
}
