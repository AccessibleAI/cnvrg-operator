apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: {{ .Spec.ControlPlane.WebApp.SvcName }}
  namespace: {{ ns . }}
spec:
  minReplicas: {{ .Spec.ControlPlane.WebApp.Replicas }}
  maxReplicas: {{ .Spec.ControlPlane.WebApp.Hpa.MaxReplicas }}
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: {{ .Spec.ControlPlane.WebApp.SvcName }}
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
        averageUtilization: {{ .Spec.ControlPlane.WebApp.Hpa.Utilization }}
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: {{ .Spec.ControlPlane.WebApp.Hpa.Utilization }}