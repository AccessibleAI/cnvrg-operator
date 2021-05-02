apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: systemkiq
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plan
spec:
  minAvailable: 1
  selector:
    matchLabels:
      cnvrg-component: systemkiq