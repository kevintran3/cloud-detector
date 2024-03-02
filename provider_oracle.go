package clouddetector

import (
	"os"
	"strings"
)

type ProviderOracle struct{}

const (
	ProductOracleOKE     = "Container Engine for Kubernetes"
	ProductOracleUnknown = ""
)

func (p *ProviderOracle) Identify() string {
	data, _ := os.ReadFile("/sys/class/dmi/id/chassis_asset_tag")
	if strings.Contains(string(data), "OracleCloud") {
		return ProviderNameOracle
	}
	return ""
}

func (p *ProviderOracle) GetInfo() ProviderInfo {
	info := ProviderInfo{
		Name:    ProviderNameOracle,
		Product: ProductOracleUnknown,
	}

	IMDSv2Url := "http://169.254.169.254/opc/v2/instance/"
	IMDSv2Headers := map[string]string{"Authorization": "Bearer Oracle"}

	info.Region = getMetadata("GET", IMDSv2Url+"regionInfo", IMDSv2Headers)["regionIdentifier"]
	if _, ok := getMetadata("GET", IMDSv2Url+"metadata", IMDSv2Headers)["oke-cluster-id"]; ok {
		info.Product = ProductOracleOKE
	}

	return info
}
