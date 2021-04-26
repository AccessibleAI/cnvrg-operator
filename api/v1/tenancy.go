package v1

type Tenancy struct {
	NamespaceTenancy bool   `json:"namespaceTenancy,omitempty"`
	Enabled          bool   `json:"enabled,omitempty"`
	DedicatedNodes   bool   `json:"dedicatedNodes,omitempty"`
	Key              string `json:"key,omitempty"`
	Value            string `json:"value,omitempty"`
}

var defaultTenancy = Tenancy{
	NamespaceTenancy: false,
	Enabled:          false,
	DedicatedNodes:   false,
	Key:              "purpose",
	Value:            "cnvrg-ccp",
}
