apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: {{ .Minio.SvcName }}
  namespace: {{ .CnvrgNs }}
spec:
  host: "{{ .Minio.SvcName }}.{{ .CnvrgNs }}.svc.cluster.local"
  trafficPolicy:
    loadBalancer:
      consistentHash:
        {{ .Minio.SharedStorage.ConsistentHash.key }}: {{ .Minio.SharedStorage.ConsistentHash.Value }}