apiVersion: apps/v1
kind: Deployment
metadata:
  name: cnvrg-authz
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
  labels:
    app: cnvrg-authz
spec:
  selector:
    matchLabels:
      app: cnvrg-authz
  template:
    metadata:
      labels:
        app: cnvrg-authz
    spec:
      containers:
      - name: authz
        imagePullPolicy: Always
        image: {{ .Image }}
        command:
          - /opt/app-root/authz
          - --ingress-type={{ .IngressType }}