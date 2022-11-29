package clouddetector

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

type ProviderOracle struct{}

const (
	ProductOracleOKE     = "Oracle Container Engine for Kubernetes"
	ProductOracleUnknown = "Oracle"
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
	data := getOracleMetadata()
	info.Region = data["canonicalRegionName"]
	info.Others = data
	return info
}

func getOracleMetadata() map[string]string {
	// Getting Oracle Metadata
	req, err := http.NewRequest("GET", "http://169.254.169.254/opc/v2/instance/", nil)
	if err != nil {
		return nil
	}
	req.Header.Set("Authorization", "Bearer Oracle")
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
