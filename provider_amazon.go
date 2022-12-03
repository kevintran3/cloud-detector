package clouddetector

import (
	"os"
	"strings"
)

type ProviderAmazon struct{}

const (
	ProductAmazonLambda  = "Lambda"
	ProductAmazonUnknown = ""
)

func (p *ProviderAmazon) Identify() string {
	data, _ := os.ReadFile("/sys/class/dmi/id/sys_vendor")
	if strings.Contains(string(data), "amazon") || len(os.Getenv("AWS_LAMBDA_FUNCTION_NAME")) > 0 {
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
	data := getMetadata("GET", "http://169.254.169.254/latest/dynamic/instance-identity/document", nil)
	info.Region = data["region"]
	info.Others = data
	return info
}
