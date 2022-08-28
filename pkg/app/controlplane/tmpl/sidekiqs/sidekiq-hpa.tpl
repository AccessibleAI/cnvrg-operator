apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: sidekiq
  namespace: {{ ns . }}
spec:
  minReplicas: {{ .Spec.ControlPlane.Sidekiq.Replicas }}
  maxReplicas: {{ .Spec.ControlPlane.Sidekiq.Hpa.MaxReplicas }}
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: sidekiq
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
        averageUtilization: {{ .Spec.ControlPlane.Sidekiq.Hpa.Utilization }}
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: {{ .Spec.ControlPlane.Sidekiq.Hpa.Utilization }}