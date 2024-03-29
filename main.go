package clouddetector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

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
	ProviderNameAmazon = "Amazon"
	ProviderNameAzure  = "Azure"
	ProviderNameGoogle = "Google"
	ProviderNameOracle = "Oracle"
)

func IdentifyProvider() ProviderInfo {
	providers := []Provider{
		&ProviderAmazon{},
		&ProviderAzure{},
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
	if len(providerInfo.Name) > 0 {
		h["Provider"] = fmt.Sprintf("%s %s (%s)", providerInfo.Name, providerInfo.Product, providerInfo.Region)
	}

	u := unix.Utsname{}
	if err := unix.Uname(&u); err == nil {
		h["OS"] = fmt.Sprintf("%s %s %s (%s) %s",
			string(u.Sysname[:]), string(u.Machine[:]), string(u.Release[:]), string(u.Version[:]), runtime.Version())
	}

	return h
}

func getMetadata(method string, url string, headers map[string]string) map[string]string {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil
	}
	for h, v := range headers {
		req.Header.Set(h, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	var data map[string]string
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil
	}
	return data
}
