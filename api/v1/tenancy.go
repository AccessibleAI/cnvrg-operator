package v1

type Tenancy struct {
	Enabled          *bool  `json:"enabled,omitempty"`
	DedicatedNodes   string `json:"dedicatedNodes,omitempty"`
	Key              string `json:"key,omitempty"`
	Value            string `json:"value,omitempty"`
}
