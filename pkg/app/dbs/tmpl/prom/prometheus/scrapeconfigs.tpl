apiVersion: v1
kind: ConfigMap
metadata:
  name: prom-scrape-configs
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
data:
  prometheus.yml: |
    global:
      scrape_interval: 5s
      evaluation_interval: 10s
    scrape_configs:
      - job_name: cnvrg-metrics
        relabel_configs:
          - source_labels: [__meta_kubernetes_pod_name]
            action: replace
            target_label: pod_name
          - source_labels: [__meta_kubernetes_pod_container_name]
            action: replace
            target_label: container_name
          - source_labels: [__meta_kubernetes_namespace]
            action: replace
            target_label: scraper_namespace
        kubernetes_sd_configs:
          - role: pod
            selectors:
              - role: pod
                label: "component=cnvrg-workload"
            namespaces:
              own_namespace: true
          - role: pod
            selectors:
              - role: pod
                label: "component=nomex"
            namespaces:
              own_namespace: true
          - role: endpoints
            selectors:
              - role: endpoints
                label: "exporter=cnvrg-job"
            namespaces:
              own_namespace: true
          {{- range $_, $cfg := .Spec.Dbs.Prom.ExtraScrapeConfigs }}
          - role: "{{$cfg.Role}}"
            selectors:
              - role: "{{$cfg.Role}}"
                label: "{{$cfg.LabelSelector}}"
            namespaces:
              names:
                - {{ $cfg.Namespace }}
          {{- end }}