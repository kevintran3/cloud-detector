package clouddetector

import (
	"encoding/json"
	"net/http"

	"golang.org/x/sys/unix"
)

type IPInfo struct {
	IP       string
	Org      string
	ISP      string
	Location string
}

type ProviderInfo struct {
	Name    string
	Product string
	Region  string
	Others  map[string]string
}

type Provider interface {
	Identify() string
	GetInfo() ProviderInfo
}

const (
	ProviderNameAmazon       = "Amazon"
	ProviderNameAzure        = "Azure"
	ProviderNameGoogle       = "Google"
	ProviderNameOracle       = "Oracle"
	ProviderNameDigitalOcean = "DigitalOcean"
)

func IdentifyProvider() ProviderInfo {
	providers := []Provider{
		&ProviderAmazon{},
		&ProviderGoogle{},
		&ProviderOracle{},
	}
	foundProv := ProviderInfo{}
	for _, prov := range providers {
		if len(prov.Identify()) > 0 {
			return prov.GetInfo()
		}
	}
	return foundProv
}

func GetHostInfo() map[string]string {
	h := map[string]string{}

	providerInfo := IdentifyProvider()
	h["Cloud"] = providerInfo.Name
	h["CloudProduct"] = providerInfo.Product
	h["CloudRegion"] = providerInfo.Region

	ipInfo := getHostPublicIP()
	h["IP"] = ipInfo.IP
	h["IPISP"] = ipInfo.ISP + " - " + ipInfo.Org
	h["IPLocation"] = ipInfo.Location

	u := unix.Utsname{}
	if err := unix.Uname(&u); err == nil {
		h["Machine"] = string(u.Machine[:])
		h["Nodename"] = string(u.Nodename[:])
		h["Release"] = string(u.Release[:])
		h["Sysname"] = string(u.Sysname[:])
		h["Version"] = string(u.Version[:])
	}

	return h
}

func getHostPublicIP() IPInfo {
	// Getting Public IP detail
	info := IPInfo{}
	req, _ := http.NewRequest("GET", "http://ip-api.com/json/", nil)
	resp, _ := http.DefaultClient.Do(req)
	var data map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&data)
	info.IP = data["query"]
	info.Org = data["org"]
	info.ISP = data["isp"]
	info.Location = data["city"] + ", " + data["regionName"] + ", " + data["country"]
	return info
}
