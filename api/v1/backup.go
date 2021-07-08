package v1

type Backup struct {
	Enabled   *bool  `json:"enabled,omitempty"`
	BucketRef string `json:"bucketRef,omitempty"`
	CredsRef  string `json:"credsRef,omitempty"`
	Rotation  int    `json:"rotation,omitempty"`
	Period    int    `json:"period,omitempty"`
}
