package conf

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

var Conf Config

// Info from config file
type Config struct {
	Baseurl   string
	Title     string
	Templates string
	Posts     string
	Public    string
	Admin     string
	Metadata  string
	Index     string
	Host      string
	Url       string
	Localip   string
	HttpPort  string
}

// Reads info from config file
func ReadConfig() Config {
	var config Config
	if _, err := toml.DecodeFile(os.Getenv("GOPATH")+"/conf/Configure.conf", &config); err != nil {
		log.Fatal(err)
	}
	Conf = config
	return config
}
