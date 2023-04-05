apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app.kubernetes.io/name: kube-state-metrics
    app.kubernetes.io/version: 1.9.7
    cnvrg-infra-prometheus: {{ .Name }}-{{ ns .}}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
  name: kube-state-metrics
  namespace: {{ ns . }}
spec:
  endpoints:
  - bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
    honorLabels: true
    interval: 30s
    port: https-main
    relabelings:
    - action: labeldrop
      regex: (pod|service|endpoint|namespace)
    scheme: https
    scrapeTimeout: 30s
    tlsConfig:
      insecureSkipVerify: true
  - bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
    interval: 30s
    port: https-self
    scheme: https
    tlsConfig:
      insecureSkipVerify: true
  jobLabel: app.kubernetes.io/name
  namespaceSelector:
    matchNames:
      - {{ ns . }}
  selector:
    matchLabels:
      app.kubernetes.io/name: kube-state-metrics
