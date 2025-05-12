package deconz

import (
	"deconz-homekit/internal/client"
)

type Configuration struct {
	ApiVersion          string  `json:"apiversion"`
	BridgeId            string  `json:"bridgeid"`
	DeviceName          string  `json:"devicename"`
	DHCP                bool    `json:"dhcp"`
	ZigbeeFirmware      string  `json:"fwversion"`
	NetworkGateway      string  `json:"gateway"`
	IpAddress           string  `json:"ipaddress"`
	LinkEnabled         bool    `json:"linkbutton"`
	Time                string  `json:"localtime"`
	MacAddress          string  `json:"mac"`
	ModelId             string  `json:"modelid"`
	Name                string  `json:"name"`
	Netmask             string  `json:"netmask"`
	NetworkOpenDuration uint16  `json:"networkopenduration"`
	NTP                 *string `json:"ntp"`
	PanId               uint16  `json:"panid"`
	PortalServices      bool    `json:"portalservices"`
	RfConnected         bool    `json:"rfconnected"`
	SwVersion           string  `json:"swversion"`
	TimeFormat          string  `json:"timeformat"`
	TimeZone            string  `json:"timezone"`
	UTC                 string  `json:"UTC"`
	UUID                string  `json:"uuid"`
	WebsocketNotifyAll  bool    `json:"websocketnotifyall"`
	WebsocketPort       int     `json:"websocketport"`
	ZigbeeChannel       int     `json:"zigbeechannel"`
}

func (ac *ApiClient) GetConfiguration() (*Configuration, error) {
	return client.Get[Configuration](ac.buildUrl("/config"))
}

type GatewayState struct {
}

func (ac *ApiClient) GetState() (*GatewayState, error) {
	return client.Get[GatewayState](ac.baseUrl + "/api/" + ac.apiKey)
}
