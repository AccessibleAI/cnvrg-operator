package v1

type Tenancy struct {
	Enabled        string `json:"enabled"`
	DedicatedNodes string `json:"dedicatedNodes"`
	Cnvrg          Cnvrg  `json:"cnvrg"`
}

type Cnvrg struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var tenancyDefault = Tenancy{
	Enabled:        "false",
	DedicatedNodes: "false",
	Cnvrg: Cnvrg{
		Key:   "cnvrg-taint",
		Value: "true",
	},
}
