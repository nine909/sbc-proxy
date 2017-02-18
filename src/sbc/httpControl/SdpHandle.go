package httpControl

import (
	"log"
	"net"

	conf "sbc/conf"

	//	"io/ioutil"

	"github.com/ernado/sdp"
	//	"os"
	//	"github.com/julienschmidt/httprouter"
)

func SdpParser(sdpByte []byte, aport, vport string) OriginSdp {
	var (
		s   sdp.Session
		err error
	)

	if s, err = sdp.DecodeSession(sdpByte, s); err != nil {
		log.Fatal("err:", err)
	}
	// for k, v := range s {
	// 	fmt.Println(k, v)
	// }
	d := sdp.NewDecoder(s)
	m := new(sdp.Message)
	if err = d.Decode(m); err != nil {
		log.Fatal("err:", err)
	}
	log.Println("Decoded session", m.Name)
	log.Println("Info:", m.Info)
	log.Println("Origin:", m.Origin)
	log.Println("IP 4: ", m.Origin.Address)
	log.Println("IP 4: ", m.Timing)
	log.Println("NetworkType: ", m.Connection.NetworkType)
	log.Println("AddressType: ", m.Connection.AddressType)
	log.Println("IP: ", m.Connection.IP)
	log.Println("TTL: ", m.Connection.TTL)
	log.Println("Addresses: ", m.Connection.Addresses)

	orig := OriginSdp{}
	medias1 := sdp.Medias{}

	orig.ip = m.Origin.Address
	isMultiMediaAudio := false
	isMultiMediaVideo := false
	// var isRemove []int
	for i, media := range m.Medias {
		log.Println("=======================")
		log.Println("Type: ", media.Description.Type)
		log.Println("Port: ", media.Description.Port)
		log.Println("PortsNumber: ", media.Description.PortsNumber)
		log.Println("Protocol: ", media.Description.Protocol)
		log.Println("Format: ", media.Description.Format)
		log.Println("Medias Connection: ", media.Connection)
		log.Println("Medias Attributes: ", media.Attributes)
		log.Println("Medias Encryption: ", media.Encryption)
		log.Println("Medias Bandwidths: ", media.Bandwidths)

		switch media.Description.Type {
		case "audio":
			if !isMultiMediaAudio {
				orig.audio.port = media.Description.Port
				orig.audio.portsNumber = media.Description.PortsNumber
				orig.audio.protocal = media.Description.Protocol

				isMultiMediaAudio = true
				medias1 = append(medias1, m.Medias[i])
			}
		case "video":
			if !isMultiMediaVideo {
				orig.video.port = media.Description.Port
				orig.video.portsNumber = media.Description.PortsNumber
				orig.video.protocal = media.Description.Protocol

				isMultiMediaVideo = true
				medias1 = append(medias1, m.Medias[i])
			}
		}

	}
	// defining medias
	m.Medias = medias1
	//replace ip
	//	m.Origin.Address = conf.Conf.Localip
	m.Connection.IP = net.ParseIP(conf.Conf.Localip)
	//encode base64
	//	log.Print("dddddddddddddddddd", orig.audio.port)
	return orig

}
func ConstructSdp(me *sdp.Message) string {
	var (
		s sdp.Session
		b []byte
	)

	// defining message
	m := me
	//	m.AddFlag("nortpproxy:yes")

	// appending message to session
	s = m.Append(s)

	// appending session to byte buffer
	b = s.AppendTo(b)
	log.Println("SDP Construct :", string(b)+"\n")
	return string(b) + "\n"
}

type OriginSdp struct {
	ip    string
	audio struct {
		port        int
		portsNumber int
		protocal    string
	}
	video struct {
		port        int
		portsNumber int
		protocal    string
	}
}
