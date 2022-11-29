package clouddetector

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
	ProviderNameUnknown      = "Unknown"
)

func IdentifyProvider() ProviderInfo {
	providers := []Provider{
		&ProviderAmazon{},
		&ProviderGoogle{},
		&ProviderOracle{},
	}
	foundProv := ProviderInfo{
		Name: ProviderNameUnknown,
	}
	for _, prov := range providers {
		if len(prov.Identify()) > 0 {
			return prov.GetInfo()
		}
	}
	return foundProv
}
