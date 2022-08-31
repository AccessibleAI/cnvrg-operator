apiVersion: v1
kind: ConfigMap
metadata:
  name: prom-config
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "false"
data:
  prometheus.yml: |
    global:
      scrape_interval: 5s
      evaluation_interval: 10s
    scrape_configs:
      - job_name: cnvrg-metrics
        honor_labels: false
        relabel_configs:
          - source_labels: [__meta_kubernetes_pod_name]
            action: replace
            target_label: pod_name
          - source_labels: [__meta_kubernetes_pod_container_name]
            action: replace
            target_label: container_name
          - source_labels: [__meta_kubernetes_namespace]
            action: replace
            target_label: namespace
        kubernetes_sd_configs:
          - role: pod
            selectors:
              - role: pod
                label: "component=cnvrg-workload"
            namespaces:
              own_namespace: true
  web-config.yml: |
    basic_auth_users:
      cnvrg: {{ .PassHash }}