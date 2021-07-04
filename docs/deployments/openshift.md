# Cnvrg on OpenShift


#### Deploy Cnvrg Control Plane on OpenShift
By default, Cnvrg Control Plane (CCP) deploys Prometheus Operator and Prometheus Instance to collect, 
store and analyze cluster metrics. OpenShift came out-of-the-box with his own Monitoring Stack, 
which is based also on Prometheus Operator and Prometheus Instance. As a result, CCP can't deploy it's
own Monitoring Stack, instead it's have to rely on OCP Monitoring Stack.
There are different ways OCP administrator can fulfill the CCP monitoring requirements
* Create dedicated Prometheus instances for user-defined projects 
* Gave access to Infrastructure Prometheus
* Patch OCP Monitoring Stack, and watch for cnvrg namespace for Prometheus Instances

Finally: CCP is agnostic to how OCP admin diced to expose Prometheus Metrics, CCP only care about the following 
1. CCP Prometheus instance have to have metrics from `prometheus-node-exporter` ServiceMonitor
2. CCP Prometheus instance have to have metrics from `kube-state-metrics` ServiceMonitor
3. CCP Prometheus instance have to have metrics from `kubelet` ServiceMonitor
4. CCP Prometheus instance have to have metrics from `dcgm-exporter` ServiceMonitor
5. CCP Prometheus instance have to have metrics from `cnvrg-jobs` ServiceMonitor

While `prometheus-node-exporter`, `kube-state-metrics` and `kubelet` ServiceMonitors 
are deployed by default with OCP Monitoring stack, 
the `dcgm-exporter` and `cnvrg-jobs` ServiceMonitors have to deployed manually on the OCP.
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

At the following OCP based deployment, we'll keep as much as possible simple.
We'll instruct CCP to use OCP's infrastructure prometheus instance to query the monitoring metrics. 
1. Get access to OCP Infra Prometheus instance (please note, the bellow commands assume you've [jq](https://stedolan.github.io/jq/) installed on your local machine)
```bash
CNVRG_PROMETHEUS_USER=$(oc get secret -nopenshift-monitoring grafana-datasources -ojson | jq -r '.data."prometheus.yaml"' | base64 -D | jq -r '.datasources[].basicAuthUser')
CNVRG_PROMETHEUS_PASS=$(oc get secret -nopenshift-monitoring grafana-datasources -ojson | jq -r '.data."prometheus.yaml"' | base64 -D | jq -r '.datasources[].basicAuthPassword')
CNVRG_PROMETHEUS_URL=$(oc get secret -nopenshift-monitoring grafana-datasources -ojson | jq -r '.data."prometheus.yaml"' | base64 -D | jq -r '.datasources[].url')
```
2. Create `prom-creds` secret to instruct CCP how to connect to Prometheus instance 
```bash
oc create secret generic prom-creds -ncnvrg \
 --from-literal=CNVRG_PROMETHEUS_USER=$CNVRG_PROMETHEUS_USER \
 --from-literal=CNVRG_PROMETHEUS_PASS=$CNVRG_PROMETHEUS_PASS \
 --from-literal=CNVRG_PROMETHEUS_URL=$CNVRG_PROMETHEUS_URL
```
3. Install CCP

Install with helm3
```bash

helm install cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
  --create-namespace -n cnvrg  \
  --set clusterDomain=<apps-openshift-route-domain> \
  --set controlPlane.baseConfig.featureFlags.OCP_ENABLED="true" \
  --set networking.ingress.type="openshift" \
  --set monitoring.prometheusOperator.enabled=false \
  --set monitoring.prometheus.enabled=false \
  --set monitoring.nodeExporter.enabled=false \
  --set monitoring.kubeStateMetrics.enabled=false \
  --set monitoring.defaultServiceMonitors.enabled=false \
  --set monitoring.dcgmExporter.enabled=false 
```

Install with raw ocp yamls
```bash
oc new-project cnvrg
helm template cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
  --create-namespace -n cnvrg  \
  --set clusterDomain=apps.<openshift-route-domain> \
  --set controlPlane.baseConfig.featureFlags.OCP_ENABLED="true" \
  --set networking.ingress.type="openshift" \
  --set monitoring.prometheus.enabled=false \
  --no-hooks | oc apply -f -
```

#### Use these containers for testing the CCP deployment on OCP

CPU runtime: docker.io/cnvrg/cnvrg-cpu-runtime:1.5

GPU runtime: docker.io/cnvrg/cnvrg-gpu-runtime:tf-2.3.2.c1
