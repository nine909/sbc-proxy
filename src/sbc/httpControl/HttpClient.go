package httpControl

import (
	"bytes"
	//	"fmt"
	"io/ioutil"
	// conf "sbc/conf"

	// "encoding/json"
	// "fmt"
	// "github.com/jmcvetta/napping"
	// "gopkg.in/jmcvetta/napping.v3"
	"log"
	"net/http"
)

func RequestHTTTP(host, data string) (string, error) {
	/*url := "http://restapi3.apiary.io/notes"
	fmt.Println("URL:>", url)*/

	var jsonStr = []byte(data)
	req, _ := http.NewRequest("POST", "http://"+host, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, errhttp := client.Do(req)
	if errhttp != nil {
		//		panic(nil, err)
		return "", errhttp
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

	return string(body), nil
}
