package v1

type Images struct {
	OperatorImage                 string `json:"operatorImage,omitempty"`
	ConfigReloaderImage           string `json:"configReloaderImage,omitempty"`
	PrometheusConfigReloaderImage string `json:"prometheusConfigReloaderImage,omitempty"`
	KubeRbacProxyImage            string `json:"kubeRbacProxyImage,omitempty"`
}
type PrometheusOperator struct {
	Enabled string `json:"enabled,omitempty"`
	Images  Images `json:"images,omitempty"`
}
type Prometheus struct {
	Enabled       string `json:"enabled,omitempty"`
	Image         string `json:"image,omitempty"`
	CPURequest    int    `json:"cpuRequest,omitempty"`
	MemoryRequest string `json:"memoryRequest,omitempty"`
	SvcName       string `json:"svcName,omitempty"`
	Port          int    `json:"port,omitempty"`
	NodePort      int    `json:"nodePort,omitempty"`
	StorageSize   string `json:"storageSize,omitempty"`
	StorageClass  string `json:"storageClass,omitempty"`
}

type NodeExporter struct {
	Enabled string `json:"enabled,omitempty"`
	Port    int    `json:"port,omitempty"`
	Image   string `json:"image,omitempty"`
}
type KubeStateMetrics struct {
	Enabled string `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}
type Grafana struct {
	Enabled  string `json:"enabled,omitempty"`
	Image    string `json:"image,omitempty"`
	SvcName  string `json:"svcName,omitempty"`
	Port     int    `json:"port,omitempty"`
	NodePort int    `json:"nodePort,omitempty"`
}
type DefaultServiceMonitors struct {
	Enabled string `json:"enabled,omitempty"`
}
type SidekiqExporter struct {
	Enabled string `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}
type MinioExporter struct {
	Enabled string `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}
type DcgmExporter struct {
	Enabled string `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
	Port    int    `json:"port,omitempty"`
}
type IdleMetricsExporter struct {
	Enabled string `json:"enabled,omitempty"`
}
type MetricsServer struct {
	Enabled string `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type Monitoring struct {
	Enabled               string             `json:"enabled,omitempty"`
	PrometheusOperator    PrometheusOperator `json:"prometheusOperator,omitempty"`
	Prometheus            Prometheus         `json:"prometheus,omitempty"`
	KubeletServiceMonitor string             `json:"kubeletServiceMonitor,omitempty"`
	//NodeExporter           NodeExporter           `json:"nodeExporter,omitempty"`
	KubeStateMetrics KubeStateMetrics `json:"kubeStateMetrics,omitempty"`
	//Grafana                Grafana                `json:"grafana,omitempty"`
	//DefaultServiceMonitors DefaultServiceMonitors `json:"defaultServiceMonitors,omitempty"`
	//SidekiqExporter        SidekiqExporter        `json:"sidekiqExporter,omitempty"`
	//MinioExporter          MinioExporter          `json:"minioExporter,omitempty"`
	//DcgmExporter           DcgmExporter           `json:"dcgmExporter,omitempty"`
	//IdleMetricsExporter    IdleMetricsExporter    `json:"idleMetricsExporter,omitempty"`
	//MetricsServer          MetricsServer          `json:"metricsServer,omitempty"`
}

var grafanaDefault = Grafana{
	Enabled:  "true",
	Image:    "grafana/grafana:7.3.4",
	SvcName:  "grafana",
	Port:     3000,
	NodePort: 30012,
}

var prometheusDefault = Prometheus{
	Enabled:       "true",
	Image:         "quay.io/prometheus/prometheus:v2.22.1",
	CPURequest:    1,
	MemoryRequest: "1Gi",
	SvcName:       "prometheus",
	Port:          9090,
	NodePort:      30909,
	StorageSize:   "100Gi",
	StorageClass:  "use-default",
}

var monitoringDefault = Monitoring{
	Enabled: "true",
	PrometheusOperator: PrometheusOperator{
		Enabled: "true",
		Images: Images{
			OperatorImage:                 "quay.io/prometheus-operator/prometheus-operator:v0.44.1",
			PrometheusConfigReloaderImage: "quay.io/prometheus-operator/prometheus-config-reloader:v0.44.1",
			KubeRbacProxyImage:            "quay.io/brancz/kube-rbac-proxy:v0.8.0",
		},
	},
	Prometheus:            prometheusDefault,
	KubeletServiceMonitor: "true",
	KubeStateMetrics: KubeStateMetrics{
		Enabled: "true",
		Image:   "quay.io/coreos/kube-state-metrics:v1.9.7",
	},
	//NodeExporter: NodeExporter{
	//	Enabled: "true",
	//	Port:    9100,
	//	Image:   "quay.io/prometheus/node-exporter:v0.18.1",
	//},
	//KubeStateMetrics: KubeStateMetrics{
	//	Enabled: "true",
	//	Image:   "quay.io/coreos/kube-state-metrics:v1.9.5",
	//},
	//Grafana: Grafana{
	//	Enabled:  "true",
	//	Image:    "grafana/grafana:7.2.0",
	//	SvcName:  "grafana",
	//	Port:     3000,
	//	NodePort: 30012,
	//},
	//DefaultServiceMonitors: DefaultServiceMonitors{Enabled: "true"},
	//SidekiqExporter: SidekiqExporter{
	//	Enabled: "true",
	//	Image:   "docker.io/strech/sidekiq-prometheus-exporter:0.1.13",
	//},
	//MinioExporter: MinioExporter{
	//	Enabled: "true",
	//	Image:   "docker.io/cnvrg/cnvrg-boot:v0.24",
	//},
	//DcgmExporter: DcgmExporter{
	//	Enabled: "true",
	//	Image:   "nvidia/dcgm-exporter:1.7.2",
	//	Port:    9400,
	//},
	//IdleMetricsExporter: IdleMetricsExporter{Enabled: "true"},
	//MetricsServer: MetricsServer{
	//	Enabled: "true",
	//	Image:   "k8s.gcr.io/metrics-server/metrics-server:v0.3.7",
	//},
}
