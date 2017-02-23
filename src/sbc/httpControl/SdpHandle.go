package httpControl

import (
	"net"

	conf "sbc/conf"
	"sbc/logs"

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
		logs.Logger.Debug("err:", err)
	}
	// for k, v := range s {
	// 	fmt.Println(k, v)
	// }
	d := sdp.NewDecoder(s)
	m := new(sdp.Message)
	if err = d.Decode(m); err != nil {
		logs.Logger.Debug("err:", err)
	}
	logs.Logger.Debug("Decoded session", m.Name)
	logs.Logger.Debug("Info:", m.Info)
	logs.Logger.Debug("Origin:", m.Origin)
	logs.Logger.Debug("IP 4: ", m.Origin.Address)
	logs.Logger.Debug("IP 4: ", m.Timing)
	logs.Logger.Debug("NetworkType: ", m.Connection.NetworkType)
	logs.Logger.Debug("AddressType: ", m.Connection.AddressType)
	logs.Logger.Debug("IP: ", m.Connection.IP)
	logs.Logger.Debug("TTL: ", m.Connection.TTL)
	logs.Logger.Debug("Addresses: ", m.Connection.Addresses)

	orig := OriginSdp{}
	medias1 := sdp.Medias{}

	orig.ip = m.Connection.IP.String()
	isMultiMediaAudio := false
	isMultiMediaVideo := false
	// var isRemove []int
	for i, media := range m.Medias {
		logs.Logger.Debug("=======================")
		logs.Logger.Debug("Type: ", media.Description.Type)
		logs.Logger.Debug("Port: ", media.Description.Port)
		logs.Logger.Debug("PortsNumber: ", media.Description.PortsNumber)
		logs.Logger.Debug("Protocol: ", media.Description.Protocol)
		logs.Logger.Debug("Format: ", media.Description.Format)
		logs.Logger.Debug("Medias Connection: ", media.Connection)
		logs.Logger.Debug("Medias Attributes: ", media.Attributes)
		logs.Logger.Debug("Medias Encryption: ", media.Encryption)
		logs.Logger.Debug("Medias Bandwidths: ", media.Bandwidths)

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
	logs.Logger.Debug("SDP Construct :", string(b)+"\n")
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
