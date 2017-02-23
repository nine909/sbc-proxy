package logs

import (
	"fmt"
	seelog "github.com/cihub/seelog"
	"io/ioutil"
)

var Logger seelog.LoggerInterface
var pathcof string = "../logs/LogConf.xml"

func loadAppConfig() {

	xmlFile, err := ioutil.ReadFile(pathcof)
	check(err)
	appConfig := string(xmlFile)
	logger, err := seelog.LoggerFromConfigAsBytes([]byte(appConfig))
	if err != nil {
		fmt.Println(err)
		return
	}
	UseLogger(logger)
}

func init() {
	DisableLog()
	loadAppConfig()

}

// DisableLog disables all library log output
func DisableLog() {
	Logger = seelog.Disabled
}

// UseLogger uses a specified seelog.LoggerInterface to output library log.
// Use this func if you are using Seelog logging system in your app.
func UseLogger(newLogger seelog.LoggerInterface) {
	Logger = newLogger
}
func check(e error) {
	if e != nil {
		fmt.Println(pathcof + " Dose not Exist")
		panic(e)
	}
}
