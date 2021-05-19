apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  name: cnvrg-ccp-prometheus
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
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
  podMetadata:
    annotations:
      {{- range $k, $v := .Spec.Annotations }}
      {{$k}}: "{{$v}}"
      {{- end }}
    labels:
      {{- range $k, $v := .Spec.Labels }}
      {{$k}}: "{{$v}}"
      {{- end }}
  resources:
    requests:
      cpu: {{ .Spec.Monitoring.Prometheus.Requests.Cpu }}
      memory: {{ .Spec.Monitoring.Prometheus.Requests.Memory }}
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
  {{- if isTrue .Spec.Tenancy.Enabled }}
  nodeSelector:
    {{ .Spec.Tenancy.Key }}: {{ .Spec.Tenancy.Value }}
    {{- range $key, $val := .Spec.Dbs.Es.NodeSelector }}
    {{ $key }}: {{ $val }}
    {{- end }}
  tolerations:
    - key: "{{ .Spec.Tenancy.Key }}"
      operator: "Equal"
      value: "{{ .Spec.Tenancy.Value }}"
      effect: "NoSchedule"
  {{- else if (gt (len .Spec.Dbs.Es.NodeSelector) 0) }}
  nodeSelector:
    {{- range $key, $val := .Spec.Dbs.Es.NodeSelector }}
    {{ $key }}: {{ $val }}
    {{- end }}
  {{- end }}
