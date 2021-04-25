apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  name: cnvrg-ccp-prometheus
  namespace: {{ ns . }}
  labels:
    app: cnvrg-ccp-prometheus
spec:
  storage:
    disableMountSubPath: true
    volumeClaimTemplate:
      spec:
        resources:
          requests:
            storage: {{ .Spec.Monitoring.Prometheus.StorageSize }}
        {{- if ne .Spec.Monitoring.Prometheus.StorageClass "" }}
        storageClassName: {{ .Spec.Monitoring.Prometheus.StorageClass }}
        {{- end }}
  image: {{ .Spec.Monitoring.Prometheus.Image }}
  replicas: 1
  resources:
    requests:
      cpu: {{ .Spec.Monitoring.Prometheus.CPURequest }}
      memory: {{ .Spec.Monitoring.Prometheus.MemoryRequest }}
  ruleSelector:
    matchLabels:
      app: cnvrg-ccp-prometheus
      role: alert-rules
  securityContext:
    fsGroup: 2000
    runAsNonRoot: true
    runAsUser: 1000
  serviceAccountName: cnvrg-ccp-prometheus
  podMonitorNamespaceSelector: {}
  podMonitorSelector: {}
  probeNamespaceSelector: {}
  serviceMonitorNamespaceSelector: {}
  serviceMonitorSelector:
    matchLabels:
      cnvrg-ccp-prometheus: {{ .Name }}-{{ ns .}}
  version: v2.22.1
  additionalScrapeConfigs:
    name: {{ .Spec.Monitoring.Prometheus.UpstreamRef }}
    key: prometheus-additional.yaml
  listenLocal: true
  containers:
    - name: "prom-auth-proxy"
      image: {{ .Spec.Monitoring.Prometheus.BasicAuthProxyImage }}
      ports:
        - containerPort: 9091
          name: web
      volumeMounts:
        - name: "prom-auth-proxy"
          mountPath: "/etc/nginx"
          readOnly: true
        - name: "htpasswd"
          mountPath: "/etc/nginx/htpasswd"
          readOnly: true
  volumes:
    - name: "prom-auth-proxy"
      configMap:
        name: "prom-auth-proxy"
    - name: "htpasswd"
      secret:
        secretName: {{ .Spec.Monitoring.Prometheus.CredsRef }}
