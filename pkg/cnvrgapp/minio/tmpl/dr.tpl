apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: {{ .Spec.Minio.SvcName }}
  namespace: {{ .Namespace }}
spec:
  host: "{{ .Spec.Minio.SvcName }}.{{ .Namespace }}.svc.cluster.local"
  trafficPolicy:
    loadBalancer:
      consistentHash:
        {{ .Spec.Minio.SharedStorage.ConsistentHash.key }}: {{ .Spec.Minio.SharedStorage.ConsistentHash.Value }}