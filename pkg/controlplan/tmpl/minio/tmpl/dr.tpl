apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: {{ .Spec.Minio.SvcName }}
  namespace: {{ ns . }}
spec:
  host: "{{ .Spec.Minio.SvcName }}.{{ ns . }}.svc.cluster.local"
  trafficPolicy:
    loadBalancer:
      consistentHash:
        {{ .Spec.Minio.SharedStorage.ConsistentHash.key }}: {{ .Spec.Minio.SharedStorage.ConsistentHash.Value }}