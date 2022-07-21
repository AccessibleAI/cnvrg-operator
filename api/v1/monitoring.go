package v1

type PrometheusOperator struct {
	Enabled                       bool   `json:"enabled,omitempty"`
	OperatorImage                 string `json:"operatorImage,omitempty"`
	PrometheusConfigReloaderImage string `json:"prometheusConfigReloaderImage,omitempty"`
	KubeRbacProxyImage            string `json:"kubeRbacProxyImage,omitempty"`
}

type Prometheus struct {
	Enabled             bool              `json:"enabled,omitempty"`
	Image               string            `json:"image,omitempty"`
	Replicas            int               `json:"replicas,omitempty"`
	BasicAuthProxyImage string            `json:"basicAuthProxyImage,omitempty"`
	Requests            Requests          `json:"requests,omitempty"`
	Limits              Limits            `json:"limits,omitempty"`
	SvcName             string            `json:"svcName,omitempty"`
	Port                int               `json:"port,omitempty"`
	NodePort            int               `json:"nodePort,omitempty"`
	Retention           string            `json:"retention,omitempty"`
	StorageSize         string            `json:"storageSize,omitempty"`
	StorageClass        string            `json:"storageClass,omitempty"`
	CredsRef            string            `json:"credsRef,omitempty"`
	UpstreamRef         string            `json:"upstreamRef,omitempty"`
	NodeSelector        map[string]string `json:"nodeSelector,omitempty"`
}

type NodeExporter struct {
	Enabled bool              `json:"enabled,omitempty"`
	Image   string            `json:"image,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

type KubeStateMetrics struct {
	Enabled bool              `json:"enabled,omitempty"`
	Image   string            `json:"image,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

type Grafana struct {
	Enabled    bool                  `json:"enabled,omitempty"`
	Image      string                `json:"image,omitempty"`
	SvcName    string                `json:"svcName,omitempty"`
	Port       int                   `json:"port,omitempty"`
	NodePort   int                   `json:"nodePort,omitempty"`
	OauthProxy OauthProxyServiceConf `json:"oauthProxy,omitempty"`
	CredsRef   string                `json:"credsRef,omitempty"`
}

type DefaultServiceMonitors struct {
	Enabled bool              `json:"enabled,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

type DcgmExporter struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type HabanaExporter struct {
	Enabled   bool   `json:"enabled,omitempty"`
	Image     string `json:"image,omitempty"`
	HlmlImage string `json:"hlmlImage,omitempty"`
}

type CnvrgIdleMetricsExporter struct {
	Enabled bool              `json:"enabled,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

type CnvrgInfraMonitoring struct {
	PrometheusOperator       PrometheusOperator       `json:"prometheusOperator,omitempty"`
	Prometheus               Prometheus               `json:"prometheus,omitempty"`
	NodeExporter             NodeExporter             `json:"nodeExporter,omitempty"`
	KubeStateMetrics         KubeStateMetrics         `json:"kubeStateMetrics,omitempty"`
	Grafana                  Grafana                  `json:"grafana,omitempty"`
	DcgmExporter             DcgmExporter             `json:"dcgmExporter,omitempty"`
	HabanaExporter           HabanaExporter           `json:"habanaExporter,omitempty"`
	DefaultServiceMonitors   DefaultServiceMonitors   `json:"defaultServiceMonitors,omitempty"`
	CnvrgIdleMetricsExporter CnvrgIdleMetricsExporter `json:"cnvrgIdleMetricsExporter,omitempty"`
}

type CnvrgAppMonitoring struct {
	Prometheus               Prometheus               `json:"prometheus,omitempty"`
	Grafana                  Grafana                  `json:"grafana,omitempty"`
	CnvrgIdleMetricsExporter CnvrgIdleMetricsExporter `json:"cnvrgIdleMetricsExporter,omitempty"`
}

var grafanaInfraDefault = Grafana{
	Enabled:  false,
	Image:    "grafana:7.3.4",
	SvcName:  "grafana",
	Port:     8080,
	NodePort: 30012,
	OauthProxy: OauthProxyServiceConf{
		SkipAuthRegex:        []string{`\/api\/health`},
		TokenValidationRegex: nil,
	},
	CredsRef: "grafana-creds",
}

var grafanaAppDefault = Grafana{
	Enabled:  false,
	Image:    "grafana:7.3.4",
	SvcName:  "grafana",
	Port:     8080,
	NodePort: 30014,
	OauthProxy: OauthProxyServiceConf{
		SkipAuthRegex:        []string{`\/api\/health`},
		TokenValidationRegex: nil,
	},
	CredsRef: "grafana-creds",
}

var prometheusInfraDefault = Prometheus{
	Enabled:             false,
	Image:               "prometheus:v2.22.1",
	Replicas:            1,
	BasicAuthProxyImage: "nginx:1.20",
	Requests: Requests{
		Cpu:    "200m",
		Memory: "500Mi",
	},
	Limits: Limits{
		Cpu:    "2000m",
		Memory: "4Gi",
	},
	SvcName:      "prometheus",
	Port:         9091, // basic auth nginx proxy is enabled by default
	NodePort:     30910,
	Retention:    "8w",
	StorageSize:  "50Gi",
	StorageClass: "",
	CredsRef:     "infra-prom-creds",
	UpstreamRef:  "upstream-prom-static-config",
	NodeSelector: nil,
}

var prometheusAppDefault = Prometheus{
	Enabled:             false,
	Image:               "prometheus:v2.22.1",
	Replicas:            1,
	BasicAuthProxyImage: "nginx:1.20",
	Requests: Requests{
		Cpu:    "200m",
		Memory: "500Mi",
	},
	Limits: Limits{
		Cpu:    "2000m",
		Memory: "4Gi",
	},
	SvcName:      "prometheus",
	Port:         9091, // basic auth nginx proxy is enabled by default
	NodePort:     30909,
	Retention:    "8w",
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
		Enabled: false,
	},
}

var infraMonitoringDefault = CnvrgInfraMonitoring{
	Prometheus: prometheusInfraDefault,
	Grafana:    grafanaInfraDefault,
	PrometheusOperator: PrometheusOperator{
		Enabled:                       false,
		OperatorImage:                 "prometheus-operator:v0.44.1",
		PrometheusConfigReloaderImage: "prometheus-config-reloader:v0.44.1",
		KubeRbacProxyImage:            "kube-rbac-proxy:v0.8.0",
	},
	KubeStateMetrics: KubeStateMetrics{
		Enabled: false,
		Image:   "kube-state-metrics:v1.9.7",
	},
	NodeExporter: NodeExporter{
		Enabled: false,
		Image:   "node-exporter:v1.0.1",
	},
	DcgmExporter: DcgmExporter{
		Enabled: false,
		Image:   "nvcr.io/nvidia/k8s/dcgm-exporter:2.0.13-2.1.2-ubuntu20.04",
	},
	HabanaExporter: HabanaExporter{
		Enabled:   true,
		Image:     "vault.habana.ai/gaudi-metric-exporter/metric-exporter:latest",
		HlmlImage: "vault.habana.ai/gaudi-metric-exporter/hlml-service:latest",
	},
	DefaultServiceMonitors: DefaultServiceMonitors{
		Enabled: false,
	},
	CnvrgIdleMetricsExporter: CnvrgIdleMetricsExporter{
		Enabled: false,
	},
}
