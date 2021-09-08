package v1

type Tenancy struct {
	Enabled bool   `json:"enabled,omitempty"`
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
}

var tenancyDefault = Tenancy{
	Enabled: false,
	Key:     "purpose",
	Value:   "cnvrg-control-plane",
}
