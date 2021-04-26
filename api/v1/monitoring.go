package v1

type Images struct {
	OperatorImage                 string `json:"operatorImage,omitempty"`
	ConfigReloaderImage           string `json:"configReloaderImage,omitempty"`
	PrometheusConfigReloaderImage string `json:"prometheusConfigReloaderImage,omitempty"`
	KubeRbacProxyImage            string `json:"kubeRbacProxyImage,omitempty"`
}
type PrometheusOperator struct {
	Enabled *bool  `json:"enabled,omitempty"`
	Images  Images `json:"images,omitempty"`
}
type Prometheus struct {
	Enabled             *bool  `json:"enabled,omitempty"`
	Image               string `json:"image,omitempty"`
	BasicAuthProxyImage string `json:"basicAuthProxyImage,omitempty"`
	CPURequest          string `json:"cpuRequest,omitempty"`
	MemoryRequest       string `json:"memoryRequest,omitempty"`
	SvcName             string `json:"svcName,omitempty"`
	Port                int    `json:"port,omitempty"`
	NodePort            int    `json:"nodePort,omitempty"`
	StorageSize         string `json:"storageSize,omitempty"`
	StorageClass        string `json:"storageClass,omitempty"`
	CredsRef            string `json:"credsRef,omitempty"`
	UpstreamRef         string `json:"upstreamRef"`
}

type NodeExporter struct {
	Enabled *bool  `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}
type KubeStateMetrics struct {
	Enabled *bool  `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}
type Grafana struct {
	Enabled    *bool                 `json:"enabled,omitempty"`
	Image      string                `json:"image,omitempty"`
	SvcName    string                `json:"svcName,omitempty"`
	Port       int                   `json:"port,omitempty"`
	NodePort   int                   `json:"nodePort,omitempty"`
	OauthProxy OauthProxyServiceConf `json:"oauthProxy,omitempty"`
}
type DefaultServiceMonitors struct {
	Enabled *bool `json:"enabled,omitempty"`
}
type SidekiqExporter struct {
	Enabled *bool  `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}
type MinioExporter struct {
	Enabled *bool  `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}
type DcgmExporter struct {
	Enabled *bool  `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}
type IdleMetricsExporter struct {
	Enabled *bool `json:"enabled,omitempty"`
}
type MetricsServer struct {
	Enabled *bool  `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type CnvrgInfraMonitoring struct {
	Enabled               *bool              `json:"enabled,omitempty"`
	PrometheusOperator    PrometheusOperator `json:"prometheusOperator,omitempty"`
	Prometheus            Prometheus         `json:"prometheus,omitempty"`
	KubeletServiceMonitor *bool              `json:"kubeletServiceMonitor,omitempty"`
	NodeExporter          NodeExporter       `json:"nodeExporter,omitempty"`
	KubeStateMetrics      KubeStateMetrics   `json:"kubeStateMetrics,omitempty"`
	Grafana               Grafana            `json:"grafana,omitempty"`
	DcgmExporter          DcgmExporter       `json:"dcgmExporter,omitempty"`
}

type CnvrgAppMonitoring struct {
	Enabled    *bool      `json:"enabled,omitempty"`
	Prometheus Prometheus `json:"prometheus,omitempty"`
	Grafana    Grafana    `json:"grafana,omitempty"`
}

var grafanaDefault = Grafana{
	Enabled:    &defaultEnabled,
	Image:      "grafana/grafana:7.3.4",
	SvcName:    "grafana",
	Port:       8080,
	NodePort:   30012,
	OauthProxy: OauthProxyServiceConf{SkipAuthRegex: []string{`\/api\/health`}},
}

var prometheusDefault = Prometheus{
	Enabled:             &defaultEnabled,
	Image:               "quay.io/prometheus/prometheus:v2.22.1",
	BasicAuthProxyImage: "docker.io/nginx:1.20",
	CPURequest:          "200m",
	MemoryRequest:       "500Mi",
	SvcName:             "prometheus",
	Port:                9091, // basic auth nginx proxy is enabled by default
	NodePort:            30909,
	StorageSize:         "50Gi",
	StorageClass:        "",
	CredsRef:            "prom-creds",
	UpstreamRef:         "upstream-prom-static-config",
}

var cnvrgAppMonitoringDefault = CnvrgAppMonitoring{
	Enabled:    &defaultEnabled,
	Prometheus: prometheusDefault,
	Grafana:    grafanaDefault,
}

var infraMonitoringDefault = CnvrgInfraMonitoring{
	Enabled:    &defaultEnabled,
	Prometheus: prometheusDefault,
	Grafana:    grafanaDefault,
	PrometheusOperator: PrometheusOperator{
		Images: Images{
			OperatorImage:                 "quay.io/prometheus-operator/prometheus-operator:v0.44.1",
			PrometheusConfigReloaderImage: "quay.io/prometheus-operator/prometheus-config-reloader:v0.44.1",
			KubeRbacProxyImage:            "quay.io/brancz/kube-rbac-proxy:v0.8.0",
		},
	},
	KubeStateMetrics: KubeStateMetrics{
		Enabled: &defaultEnabled,
		Image:   "quay.io/coreos/kube-state-metrics:v1.9.7",
	},
	NodeExporter: NodeExporter{
		Enabled: &defaultEnabled,
		Image:   "quay.io/prometheus/node-exporter:v1.0.1",
	},
	DcgmExporter: DcgmExporter{
		Enabled: &defaultEnabled,
		Image:   "nvcr.io/nvidia/k8s/dcgm-exporter:2.1.4-2.3.1-ubuntu18.04",
	},
}
