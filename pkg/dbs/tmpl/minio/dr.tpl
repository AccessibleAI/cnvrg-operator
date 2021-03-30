apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: {{ .Spec.Dbs.Minio.SvcName }}
  namespace: {{ ns . }}
spec:
  host: "{{ .Spec.Dbs.Minio.SvcName }}.{{ ns . }}.svc.cluster.local"
  trafficPolicy:
    loadBalancer:
      consistentHash:
        {{ .Spec.Dbs.Minio.SharedStorage.ConsistentHash.key }}: {{ .Spec.Dbs.Minio.SharedStorage.ConsistentHash.Value }}