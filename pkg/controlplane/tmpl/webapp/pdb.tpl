apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: webapp
  namespace: {{ ns . }}
  annotations:
  {{- range $k, $v := .Spec.Annotations }}
  {{$k}}: "{{$v}}"
  {{- end }}
  labels:
    owner: cnvrg-control-plan
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  minAvailable: 1
  selector:
    matchLabels:
      cnvrg-component: webapp