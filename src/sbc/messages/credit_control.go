package messages

import (
	"encoding/json"
	"log"
	"strconv"
	"time"
)

var seq int = 0

type Credit_Control struct {
	SessionId         string               `json:"Session-Id"`
	AuthApplicationId string               `json:"Auth-Application-Id"` //0 , 1 ,2 ,3 4
	ServiceContextId  string               `json:"Service-Context-Id"`
	CCRequestType     string               `json:"CC-Request-Type"`
	CCRequestNumber   string               `json:"CC-Request-Number"`
	EventTimestamp    string               `json:"Event-Timestamp"`
	ServiceIdentifier string               `json:"Service-Identifier"`
	RouteRecord       string               `json:"Route-Record"`
	SubscriptId       SubscriptionId       `json:"Subscription-Id"`
	ReqSernit         RequestedServiceUnit `json:"Requested-Service-Unit"`
	UsedSerUnit       UsedServiceUnit      `json:"Used-Service-Unit"`
	ServiceInfo       ServiceInformation   `json:"Service-Information"`
}
type SubscriptionId struct {
	Type string `json:"Subscription-Id-Type"`
	Data string `json:"Subscription-Id-Data"`
}
type RequestedServiceUnit struct {
	CCTime string `json:"CC-Time"`
}
type UsedServiceUnit struct {
	CCTime string `json:"CC-Time"`
}
type ServiceInformation struct {
	InInfo INInformation `json:"IN-Information"`
}
type ResourceAllocateResponse struct {
	resultcode       string
	developermessage string
	SDP              string
}
type INInformation struct {
	ChargeFlowType             string `json:"Charge-Flow-Type"`
	SSPTime                    string `json:"SSP-Time"`
	TimeZone                   string `json:"Time-Zone"`
	CallingPartyAddressNature  string `json:"Calling-Party-Address-Nature"`
	CalledPartyAddressNature   string `json:"Called-Party-Address-Nature"`
	calledPartyBCDNumberNature string `json:"called-Party-BCDNumber-Nature"`
	EventTypeBCSM              string `json:"EventType-BCSM"`
}

func ConstructCCR_I(sesstion string) string {
	ccri := &Credit_Control{}
	ccri.SessionId = sesstion
	ccri.AuthApplicationId = ""
	ccri.ServiceContextId = ""
	ccri.CCRequestType = "1"
	ccri.CCRequestNumber = strconv.Itoa(seq + 1) //update to seq
	ccri.EventTimestamp = time.Now().String()
	ccri.ServiceIdentifier = ""
	ccri.RouteRecord = ""
	ccri.SubscriptId = SubscriptionId{"dd", "ddd"}
	ccri.ReqSernit.CCTime = ""
	ccri.UsedSerUnit.CCTime = "180"

	inInfo1 := INInformation{}
	inInfo1.ChargeFlowType = ""
	inInfo1.SSPTime = ""
	inInfo1.TimeZone = ""
	inInfo1.CallingPartyAddressNature = ""
	inInfo1.CalledPartyAddressNature = ""
	inInfo1.calledPartyBCDNumberNature = ""
	inInfo1.EventTypeBCSM = ""

	ccri.ServiceInfo.InInfo = inInfo1

	ccriJson, _ := json.Marshal(ccri)
	log.Println(string(ccriJson))
	return string(ccriJson)
}
