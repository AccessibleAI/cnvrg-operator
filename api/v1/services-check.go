package v1

type IngressCheck struct {
	Enabled bool `json:"enabled,omitempty"`
}

var ingressCheckDefault = IngressCheck{
	Enabled: true,
}
