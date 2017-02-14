package httpControl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	// conf "sbc/conf"

	// "encoding/json"
	// "fmt"
	// "github.com/jmcvetta/napping"
	// "gopkg.in/jmcvetta/napping.v3"
	// "log"
	"net/http"
)

func RequestHTTTP(host, data string) string {
	/*url := "http://restapi3.apiary.io/notes"
	fmt.Println("URL:>", url)*/

	var jsonStr = []byte(data)
	req, err := http.NewRequest("POST", "http://"+host, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return string(body)
}
