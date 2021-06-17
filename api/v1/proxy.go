package v1

type Proxy struct {
	Enabled    *bool    `json:"enabled,omitempty"`
	ConfigRef  string   `json:"configRef,omitempty"`
	HttpProxy  []string `json:"httpProxy,omitempty"`
	HttpsProxy []string `json:"httpsProxy,omitempty"`
	NoProxy    []string `json:"noProxy,omitempty"`
}

var defaultProxy = Proxy{
	Enabled:    &defaultFalse,
	ConfigRef:  "cp-proxy",
	HttpProxy:  nil,
	HttpsProxy: nil,
	NoProxy:    []string{".svc.cluster.local"},
}
