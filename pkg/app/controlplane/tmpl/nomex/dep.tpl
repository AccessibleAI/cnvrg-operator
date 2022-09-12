apiVersion: apps/v1
kind: Deployment
metadata:
  name: nomex
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
  labels:
    app: nomex
    component: nomex
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nomex
  template:
    metadata:
      labels:
        app: nomex
    spec:
      serviceAccountName: cnvrg-nomex
      containers:
      - name: nomex
        image: docker.io/cnvrg/nomex:v1.0.0
        command:
          - /opt/app-root/nomex
        ports:
        - containerPort: 2112
