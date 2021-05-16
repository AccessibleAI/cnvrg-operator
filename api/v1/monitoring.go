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
	Enabled             *bool             `json:"enabled,omitempty"`
	Image               string            `json:"image,omitempty"`
	BasicAuthProxyImage string            `json:"basicAuthProxyImage,omitempty"`
	Requests            Requests          `json:"requests,omitempty"`
	SvcName             string            `json:"svcName,omitempty"`
	Port                int               `json:"port,omitempty"`
	NodePort            int               `json:"nodePort,omitempty"`
	StorageSize         string            `json:"storageSize,omitempty"`
	StorageClass        string            `json:"storageClass,omitempty"`
	CredsRef            string            `json:"credsRef,omitempty"`
	UpstreamRef         string            `json:"upstreamRef,omitempty"`
	NodeSelector        map[string]string `json:"nodeSelector,omitempty"`
}

type NodeExporter struct {
	Enabled *bool             `json:"enabled,omitempty"`
	Image   string            `json:"image,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

type KubeStateMetrics struct {
	Enabled *bool             `json:"enabled,omitempty"`
	Image   string            `json:"image,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
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
	Enabled *bool             `json:"enabled,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

type DcgmExporter struct {
	Enabled *bool  `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type CnvrgIdleMetricsExporter struct {
	Enabled *bool             `json:"enabled,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

type MetricsServer struct {
	Enabled *bool  `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type CnvrgInfraMonitoring struct {
	PrometheusOperator       PrometheusOperator       `json:"prometheusOperator,omitempty"`
	Prometheus               Prometheus               `json:"prometheus,omitempty"`
	NodeExporter             NodeExporter             `json:"nodeExporter,omitempty"`
	KubeStateMetrics         KubeStateMetrics         `json:"kubeStateMetrics,omitempty"`
	Grafana                  Grafana                  `json:"grafana,omitempty"`
	DcgmExporter             DcgmExporter             `json:"dcgmExporter,omitempty"`
	DefaultServiceMonitors   DefaultServiceMonitors   `json:"defaultServiceMonitors,omitempty"`
	CnvrgIdleMetricsExporter CnvrgIdleMetricsExporter `json:"cnvrgIdleMetricsExporter,omitempty"`
}

type CnvrgAppMonitoring struct {
	Prometheus               Prometheus               `json:"prometheus,omitempty"`
	Grafana                  Grafana                  `json:"grafana,omitempty"`
	CnvrgIdleMetricsExporter CnvrgIdleMetricsExporter `json:"cnvrgIdleMetricsExporter,omitempty"`
}

var grafanaInfraDefault = Grafana{
	Enabled:    &defaultEnabled,
	Image:      "grafana/grafana:7.3.4",
	SvcName:    "grafana",
	Port:       8080,
	NodePort:   30012,
	OauthProxy: OauthProxyServiceConf{SkipAuthRegex: []string{`\/api\/health`}},
}

var grafanaAppDefault = Grafana{
	Enabled:    &defaultEnabled,
	Image:      "grafana/grafana:7.3.4",
	SvcName:    "grafana",
	Port:       8080,
	NodePort:   30014,
	OauthProxy: OauthProxyServiceConf{SkipAuthRegex: []string{`\/api\/health`}},
}

var prometheusInfraDefault = Prometheus{
	Enabled:             &defaultEnabled,
	Image:               "quay.io/prometheus/prometheus:v2.22.1",
	BasicAuthProxyImage: "docker.io/nginx:1.20",
	Requests: Requests{
		Cpu:    "200m",
		Memory: "500Mi",
	},
	SvcName:      "prometheus",
	Port:         9091, // basic auth nginx proxy is enabled by default
	NodePort:     30910,
	StorageSize:  "50Gi",
	StorageClass: "",
	CredsRef:     "prom-creds",
	UpstreamRef:  "upstream-prom-static-config",
	NodeSelector: nil,
}

var prometheusAppDefault = Prometheus{
	Enabled:             &defaultEnabled,
	Image:               "quay.io/prometheus/prometheus:v2.22.1",
	BasicAuthProxyImage: "docker.io/nginx:1.20",
	Requests: Requests{
		Cpu:    "200m",
		Memory: "500Mi",
	},
	SvcName:      "prometheus",
	Port:         9091, // basic auth nginx proxy is enabled by default
	NodePort:     30909,
	StorageSize:  "50Gi",
	StorageClass: "",
	CredsRef:     "prom-creds",
	UpstreamRef:  "upstream-prom-static-config",
	NodeSelector: nil,
}

var cnvrgAppMonitoringDefault = CnvrgAppMonitoring{
	Prometheus: prometheusAppDefault,
	Grafana:    grafanaAppDefault,
	CnvrgIdleMetricsExporter: CnvrgIdleMetricsExporter{
		Enabled: &defaultEnabled,
	},
}

var infraMonitoringDefault = CnvrgInfraMonitoring{
	Prometheus: prometheusInfraDefault,
	Grafana:    grafanaInfraDefault,
	PrometheusOperator: PrometheusOperator{
		Enabled: &defaultEnabled,
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
	DefaultServiceMonitors: DefaultServiceMonitors{
		Enabled: &defaultEnabled,
	},
	CnvrgIdleMetricsExporter: CnvrgIdleMetricsExporter{
		Enabled: &defaultEnabled,
	},
}
