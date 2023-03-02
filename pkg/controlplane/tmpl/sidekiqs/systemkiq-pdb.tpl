apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: systemkiq
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
  minAvailable: 1
  selector:
    matchLabels:
      cnvrg-component: systemkiq