package main

/* example response
{
   "software": {
      "softwareVersion": "9000.84.0",
      "updateChannel": "stable-channel",
      "updateNewVersion": "0.0.0.0",
      "updateProgress": 0.0,
      "updateRequired": false,
      "updateStatus": "idle"
   },
   "system": {
      "countryCode": "us",
      "deviceId": "52daed7b-307e-f95e-82e7-76dc8a98bca0",
      "groupRole": "none",
      "hardwareId": "WHIRLWIND D3A-Q2Q-Q8B",
      "lan0Link": true,
      "modelId": "ACdqK",
      "uptime": 1761341
   },
   "wan": {
      "captivePortal": false,
      "ethernetLink": true,
      "gatewayIpAddress": "172.78.186.1",
      "invalidCredentials": false,
      "ipAddress": true,
      "ipMethod": "dhcp",
      "ipPrefixLength": 24,
      "leaseDurationSeconds": 300,
      "localIpAddress": "172.78.186.169",
      "nameServers": [ "8.8.8.8", "8.8.4.4", "172.78.186.1" ],
      "online": true,
      "pppoeDetected": false
   }
}
*/

type OnhubStatus struct {
	Software onhubSoftware
	System   onhubSystem
	Wan      onhubWan
}

type onhubSoftware struct {
	SoftwareVersion  string
	UpdateChannel    string
	UpdateNewVersion string
	UpdateProgress   float64
	UpdateRequired   bool
	UpdateStatus     string
}

type onhubSystem struct {
	CountryCode string
	DeviceId    string
	GroupRole   string
	HardwareId  string
	Lan0Link    bool
	ModelId     string
	Uptime      int64
}

type onhubWan struct {
	CaptivePortal        bool
	EthernetLink         bool
	GatewayIpAddress     string
	InvalidCredentials   bool
	IpAddress            bool
	IpMethod             string
	IpPrefixLength       int
	LeaseDurationSeconds int
	LocalIpAddress       string
	NameServers          []string
	Online               bool
	PPPOEDetected        bool
}
