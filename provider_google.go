package clouddetector

import (
	"io"
	"net/http"
	"os"
	"strings"
)

type ProviderGoogle struct{}

const (
	ProductGoogleAppEngine     = "App Engine"
	ProductGoogleCloudRun      = "Cloud Run"
	ProductGoogleCloudFunction = "Cloud Function"
	ProductGoogleUnknown       = ""
)

func (p *ProviderGoogle) Identify() string {
	data, _ := os.ReadFile("/sys/class/dmi/id/product_name")
	if strings.Contains(string(data), "Google") || len(os.Getenv("K_SERVICE")) > 0 || len(os.Getenv("GAE_ENV")) > 0 {
		return ProviderNameGoogle
	}
	return ""
}

func (p *ProviderGoogle) GetInfo() ProviderInfo {
	info := ProviderInfo{
		Name:    ProviderNameGoogle,
		Product: ProductGoogleUnknown,
	}
	// Getting Google Cloud Product using Env
	if len(os.Getenv("FUNCTION_REGION")) > 0 {
		info.Product = ProductGoogleCloudFunction
		info.Region = os.Getenv("FUNCTION_REGION")
		return info
	}
	if len(os.Getenv("K_SERVICE")) > 0 {
		info.Product = ProductGoogleCloudRun
	}
	if len(os.Getenv("GAE_ENV")) > 0 {
		info.Product = ProductGoogleAppEngine
	}
	// Format projects/PROJECT-NUMBER/regions/REGION
	GoogleRegion := getGoogleMetaData("/computeMetadata/v1/instance/region")
	info.Region = strings.Split(GoogleRegion, "/")[3]
	return info
}

func getGoogleMetaData(endpoint string) string {
	// Getting Google Cloud Metadata for Region
	req, err := http.NewRequest("GET", "http://metadata.google.internal"+endpoint, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Metadata-Flavor", "Google")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(resBody)
}
