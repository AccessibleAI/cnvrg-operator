apiVersion: v1
kind: Secret
metadata:
  name: {{ .Data.CredsRef }}
  namespace: {{ .Data.Namespace }}
  annotations:
    {{- range $k, $v := .Data.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Data.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
{{- $upstreamConfig := `
    - job_name: 'federate'
      scrape_interval: 10s
      honor_labels: true
      honor_timestamps: false
      metrics_path: '/federate'
      basic_auth:
        username: '%s'
        password: '%s'
      params:
        'match[]':
          - '{namespace="%s"}'
      static_configs:
        - targets:
          - '%s'`
}}
  prometheus-additional.yaml: {{ printf $upstreamConfig .Data.User .Data.Pass .Data.Namespace .Data.Upstream | b64enc }}