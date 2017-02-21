package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	// "io/ioutil"
	b64 "encoding/base64"
	"net/http"
)

type Profile struct {
	XSession        string `json:"x-session"`
	SDP             string `json:"SDP"`
	CallbackAddress string `json:"Callback-Address"`
	CallbackSession string `json:"Callback-Session"`
	TestAllocate    string `json:"Test_Allocate"`
}

func main() {

	prof := Profile{CallbackAddress: "202.129.207.231:3000"}

	prof.XSession = "1111111111"
	prof.SDP = "dj0wDQpvPTExMTExMTExMTEgNzg2IDMyMDQgSU4gSVA0IDIwMi4xMzkuMjA3LjIzMQ0Kcz1UYWxrDQpjPUlOIElQNCAyMDIuMTM5LjIwNy4yMzENCnQ9MCAwDQphPXJ0Y3AteHI6cmN2ci1ydHQ9YWxsOjEwMDAwIHN0YXQtc3VtbWFyeT1sb3NzLGR1cCxqaXR0LFRUTCB2b2lwLW1ldHJpY3MNCm09YXVkaW8gNzA3OCBSVFAvQVZQIDAgOCAxMDENCmE9cnRwbWFwOjEwMSB0ZWxlcGhvbmUtZXZlbnQvODAwMA0KYT1ydGNwLWZiOiogdHJyLWludCA1MDAwDQptPXZpZGVvIDkwNzggUlRQL0FWUCA5Ng0KYT1ydHBtYXA6OTYgVlA4LzkwMDAwDQphPXJ0Y3AtZmI6KiB0cnItaW50IDUwMDANCmE9cnRjcC1mYjo5NiBuYWNrIHBsaQ0KYT1ydGNwLWZiOjk2IG5hY2sgc2xpDQphPXJ0Y3AtZmI6OTYgYWNrIHJwc2kNCmE9cnRjcC1mYjo5NiBjY20gZmlyDQo="
	prof.CallbackSession = "MO:uID2x0Xnpj"
	prof.TestAllocate = "Allocate_1"

	sDec0, _ := b64.StdEncoding.DecodeString(prof.SDP)
	fmt.Println(string(sDec0))
	fmt.Println()

	profResp := send(prof)

	fmt.Println("Profile.XSession", profResp.XSession)
	fmt.Println("Profile.CallbackAddress", profResp.CallbackAddress)
	fmt.Println("Profile.CallbackSession", profResp.CallbackSession)
	fmt.Println("Profile.TestAllocate", profResp.TestAllocate)
	fmt.Println("Profile.SDP", profResp.SDP)

	sDec, _ := b64.StdEncoding.DecodeString(profResp.SDP)
	fmt.Println(string(sDec))
	fmt.Println()

	profa := reuse("2222222222", profResp.CallbackAddress, "MT:uID2x0Xnpj", profResp.TestAllocate, profResp.SDP)
	fmt.Println(profa)
	profResp1 := send(profa)
	fmt.Println("1 Profile.XSession", profResp1.XSession)
	fmt.Println("1 Profile.CallbackAddress", profResp1.CallbackAddress)
	fmt.Println("1 Profile.CallbackSession", profResp1.CallbackSession)
	fmt.Println("1 Profile.TestAllocate", profResp1.TestAllocate)
	fmt.Println("1 Profile.SDP", profResp1.SDP)

	sDec1, _ := b64.StdEncoding.DecodeString(profResp1.SDP)
	fmt.Println(string(sDec1))
	fmt.Println()

	profResp1.SDP = "dj0wDQpvPTIyMjIyMjIyMjIgNzg2IDMyMDQgSU4gSVA0IDIwMi4xMzkuMjA3LjIzMQ0Kcz1UYWxrDQpjPUlOIElQNCAyMDIuMTM5LjIwNy4yMzENCnQ9MCAwDQphPXJ0Y3AteHI6cmN2ci1ydHQ9YWxsOjEwMDAwIHN0YXQtc3VtbWFyeT1sb3NzLGR1cCxqaXR0LFRUTCB2b2lwLW1ldHJpY3MNCm09YXVkaW8gODA4OCBSVFAvQVZQIDAgOCAxMDENCmE9cnRwbWFwOjEwMSB0ZWxlcGhvbmUtZXZlbnQvODAwMA0KYT1ydGNwLWZiOiogdHJyLWludCA1MDAwDQptPXZpZGVvIDkwNzggUlRQL0FWUCA5Ng0KYT1ydHBtYXA6OTYgVlA4LzkwMDAwDQphPXJ0Y3AtZmI6KiB0cnItaW50IDUwMDANCmE9cnRjcC1mYjo5NiBuYWNrIHBsaQ0KYT1ydGNwLWZiOjk2IG5hY2sgc2xpDQphPXJ0Y3AtZmI6OTYgYWNrIHJwc2kNCmE9cnRjcC1mYjo5NiBjY20gZmlyDQ=="

	profb := reuse("2222222222", profResp1.CallbackAddress, "MT:uID2x0Xnpj", profResp1.TestAllocate, profResp1.SDP)
	fmt.Println(profb)
	profResp2 := send(profb)
	fmt.Println("2 Profile.XSession", profResp2.XSession)
	fmt.Println("2 Profile.CallbackAddress", profResp2.CallbackAddress)
	fmt.Println("2 Profile.CallbackSession", profResp2.CallbackSession)
	fmt.Println("2 Profile.TestAllocate", profResp2.TestAllocate)
	fmt.Println("2 Profile.SDP", profResp2.SDP)

	sDec2, _ := b64.StdEncoding.DecodeString(profResp2.SDP)
	fmt.Println(string(sDec2))
	fmt.Println()

	profc := reuse("1111111111", profResp2.CallbackAddress, "MO:uID2x0Xnpj", profResp2.TestAllocate, profResp2.SDP)
	fmt.Println(profc)
	profResp3 := send(profc)
	fmt.Println("3 Profile.XSession", profResp3.XSession)
	fmt.Println("3 Profile.CallbackAddress", profResp3.CallbackAddress)
	fmt.Println("3 Profile.CallbackSession", profResp3.CallbackSession)
	fmt.Println("3 Profile.TestAllocate", profResp3.TestAllocate)
	fmt.Println("3 Profile.SDP", profResp3.SDP)

	sDec3, _ := b64.StdEncoding.DecodeString(profResp3.SDP)
	fmt.Println(string(sDec3))
	fmt.Println()
}

func reuse(xSession, callA, callS, allocate, sdp string) Profile {
	fmt.Println("Update")
	return Profile{XSession: xSession, CallbackAddress: callA, CallbackSession: callS, SDP: sdp, TestAllocate: allocate}
}

func send(prof Profile) *Profile {
	domain := "http://202.129.207.231:8080"
	url := domain + "/p-SSF/1.0.0/SBC/ResourceAllocate1/" + prof.XSession + "?"
	fmt.Println("URL:>", url)

	js, err := json.Marshal(prof)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	fmt.Println(string(js))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(js))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	// defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var profResp Profile
	errr := decoder.Decode(&profResp)
	if errr != nil {
		panic(errr)
	}

	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	return &profResp
}
