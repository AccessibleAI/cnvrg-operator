apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: {{ .Spec.ControlPlan.Minio.SvcName }}
  namespace: {{ ns . }}
spec:
  host: "{{ .Spec.ControlPlan.Minio.SvcName }}.{{ ns . }}.svc.cluster.local"
  trafficPolicy:
    loadBalancer:
      consistentHash:
        {{ .Spec.ControlPlan.Minio.SharedStorage.ConsistentHash.key }}: {{ .Spec.ControlPlan.Minio.SharedStorage.ConsistentHash.Value }}