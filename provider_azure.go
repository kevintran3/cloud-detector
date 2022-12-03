package clouddetector

import (
	"os"
	"strings"
)

type ProviderAzure struct{}

const (
	ProductAzureContainerApps = "Container Apps"
	ProductAzureUnknown       = ""
)

func (p *ProviderAzure) Identify() string {
	data, _ := os.ReadFile("/sys/class/dmi/id/chassis_asset_tag")
	if strings.Contains(string(data), "Microsoft Corporation") {
		return ProviderNameAzure
	}
	return ""
}

func (p *ProviderAzure) GetInfo() ProviderInfo {
	info := ProviderInfo{
		Name:    ProviderNameAzure,
		Product: ProductAzureUnknown,
	}
	data := getMetadata("GET", "http://169.254.169.254/metadata/instance?api-version=2021-02-01", map[string]string{"Metadata": "true"})
	info.Region = data["location"]
	info.Others = data
	return info
}
