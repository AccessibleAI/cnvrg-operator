apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  name: cnvrg-infra-prometheus
  namespace: {{ ns . }}
  labels:
    app: cnvrg-infra-prometheus
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
      app: cnvrg-infra-prometheus
      role: alert-rules
  securityContext:
    fsGroup: 2000
    runAsNonRoot: true
    runAsUser: 1000
  serviceAccountName: cnvrg-infra-prometheus
  podMonitorNamespaceSelector: {}
  podMonitorSelector: {}
  probeNamespaceSelector: {}
  serviceMonitorNamespaceSelector: {}
  serviceMonitorSelector:
    matchLabels:
      cnvrg-infra-prometheus: {{ .Name }}-{{ ns .}}
  version: v2.22.1
  {{- if eq .Spec.SSO.Enabled "true" }}
  listenLocal: true
  containers:
  - name: "cnvrg-oauth-proxy"
    image: {{ .Spec.SSO.Image }}
    command: [ "oauth2-proxy","--config", "/opt/app-root/conf/proxy-config/conf" ]
    ports:
    - containerPort: 9091
      name: web
    volumeMounts:
    - name: "oauth-proxy-config"
      mountPath: "/opt/app-root/conf/proxy-config"
      readOnly: true
    - name: {{ .Spec.Monitoring.Prometheus.BasicAuthRef }}
      mountPath: "/opt/app-root/conf/keys"
      readOnly: true
  volumes:
  - name: "oauth-proxy-config"
    secret:
      secretName: "oauth-proxy-{{.Spec.Monitoring.Prometheus.SvcName}}"
  - name: {{ .Spec.Monitoring.Prometheus.BasicAuthRef }}
    secret:
      secretName: {{ .Spec.Monitoring.Prometheus.BasicAuthRef }}
  {{- end }}
