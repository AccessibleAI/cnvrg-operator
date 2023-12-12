apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: cnvrg-ccp-prometheus
  namespace: {{ ns . }}
  annotations:
  {{- range $k, $v := .Spec.Annotations }}
  {{$k}}: "{{$v}}"
  {{- end }}
  labels:
    app: cnvrg-ccp-prometheus
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  minAvailable: 1
  selector:
    matchLabels:
      prometheus: cnvrg-ccp-prometheus