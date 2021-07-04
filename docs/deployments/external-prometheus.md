# External Prometheus

#### Connect Cnvrg Control Plane to external Prometheus  

Cnvrg Control Plane (CCP) deploys its own Monitoring stack, based on Prometheus Operator and Prometheus instance.
However, cnvrg administrator may change this behavior and instruct cnvrg to not deploy its own Monitoring stack
but instead use existing one. 

CCP depends on the data from the following Exporters and ServiceMonitors, make sure your Prometheus instance holds this data   
1. CCP Prometheus instance have to have metrics from `prometheus-node-exporter` ServiceMonitor
2. CCP Prometheus instance have to have metrics from `kube-state-metrics` ServiceMonitor
3. CCP Prometheus instance have to have metrics from `kubelet` ServiceMonitor
4. CCP Prometheus instance have to have metrics from `dcgm-exporter` ServiceMonitor
5. CCP Prometheus instance have to have metrics from `cnvrg-jobs` ServiceMonitor
Apply the following yamls to create CCP required ServiceMonitors
```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: cnvrg-jobs
  namespace: openshift-monitoring
  labels:
    app: cnvrg-jobs
spec:
  jobLabel: cnvrg-job
  selector:
    matchLabels:
      exporter: cnvrg-job
  endpoints:
    - interval: 30s
      scrapeTimeout: 10s
      port: "http"
---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: dcgm-exporter
  namespace: openshift-monitoring
  labels:
    app: "dcgm-exporter"
spec:
  selector:
    matchLabels:
      app: "dcgm-exporter"
  endpoints:
    - port: "metrics"
      path: "/metrics"
      interval: "15s"
```

To instruct CCP to connect to external Prometheus instance create the following K8s secret 
```shell
oc create secret generic prom-creds -n <cnvrg-namespace> \
 --from-literal=CNVRG_PROMETHEUS_USER=$CNVRG_PROMETHEUS_USER \
 --from-literal=CNVRG_PROMETHEUS_PASS=$CNVRG_PROMETHEUS_PASS \
 --from-literal=CNVRG_PROMETHEUS_URL=$CNVRG_PROMETHEUS_URL
```

Once secret exists, deploy CCP with the following flags (disable Prometheus deployment)
```shell
...
  --set monitoring.prometheusOperator.enabled=false \
  --set monitoring.prometheus.enabled=false \
  --set monitoring.nodeExporter.enabled=false \
  --set monitoring.kubeStateMetrics.enabled=false \
  --set monitoring.defaultServiceMonitors.enabled=false \
  --set monitoring.dcgmExporter.enabled=false 
...
```