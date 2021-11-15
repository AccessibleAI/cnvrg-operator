package v1

type ServicesCheck struct {
	Enabled bool `json:"enabled,omitempty"`
}

var servicesCheckDefault = ServicesCheck{
	Enabled: true,
}
