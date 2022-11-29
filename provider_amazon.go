package clouddetector

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

type ProviderAmazon struct{}

const (
	ProductAmazonLambda  = "Amazon Lambda"
	ProductAmazonUnknown = "Amazon"
)

func (p *ProviderAmazon) Identify() string {
	data, _ := os.ReadFile("/sys/class/dmi/id/product_version")
	if strings.Contains(string(data), "amazon") {
		return ProviderNameAmazon
	}
	return ""
}

func (p *ProviderAmazon) GetInfo() ProviderInfo {
	info := ProviderInfo{
		Name:    ProviderNameAmazon,
		Product: ProductAmazonUnknown,
	}
	// Getting Amazon Product using Env
	if len(os.Getenv("AWS_LAMBDA_FUNCTION_NAME")) > 0 {
		info.Product = ProductAmazonLambda
		info.Region = os.Getenv("AWS_REGION")
		return info
	}
	data := getAmazonMetadata()
	info.Region = data["region"]
	info.Others = data
	return info
}

func getAmazonMetadata() map[string]string {
	// Getting Amazon Metadata
	req, err := http.NewRequest("GET", "http://169.254.169.254/latest/dynamic/instance-identity/document", nil)
	if err != nil {
		return nil
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
