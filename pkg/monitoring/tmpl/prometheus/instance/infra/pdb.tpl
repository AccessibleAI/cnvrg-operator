apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: cnvrg-infra-prometheus
  namespace: {{ ns . }}
  annotations:
  {{- range $k, $v := .Spec.Annotations }}
  {{$k}}: "{{$v}}"
  {{- end }}
  labels:
    app: cnvrg-infra-prometheus
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  minAvailable: 1
  selector:
    matchLabels:
      prometheus: cnvrg-infra-prometheus