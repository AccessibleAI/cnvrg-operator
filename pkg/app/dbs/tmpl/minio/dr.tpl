apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: {{ .Spec.Dbs.Minio.SvcName }}
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
  host: "{{ .Spec.Dbs.Minio.SvcName }}.{{ ns . }}.svc.{{ .Spec.ClusterInternalDomain }}"
  trafficPolicy:
    loadBalancer:
      consistentHash:
        {{ .Spec.Dbs.Minio.SharedStorage.ConsistentHash.key }}: {{ .Spec.Dbs.Minio.SharedStorage.ConsistentHash.Value }}