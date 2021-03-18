package v1

type Fluentbit struct {
	Enabled string `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}
