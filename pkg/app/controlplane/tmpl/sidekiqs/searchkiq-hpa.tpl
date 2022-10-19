apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: searchkiq
  namespace: {{ ns . }}
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
spec:
  minReplicas: {{ .Spec.ControlPlane.Searchkiq.Replicas }}
  maxReplicas: {{ .Spec.ControlPlane.Searchkiq.Hpa.MaxReplicas }}
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: searchkiq
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Pods
        value: 1
        periodSeconds: 15
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
      - type: Pods
        value: 1
        periodSeconds: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: {{ .Spec.ControlPlane.Searchkiq.Hpa.Utilization }}
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: {{ .Spec.ControlPlane.Searchkiq.Hpa.Utilization }}