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
  selector:
    matchLabels:
      app: nomex
      component: nomex
  template:
    metadata:
      labels:
        app: nomex
        component: nomex
    spec:
      serviceAccountName: cnvrg-nomex
      containers:
      - name: nomex
        image: {{ .Spec.ControlPlane.Nomex.Image }}
        command:
          - /opt/app-root/nomex
        ports:
        - containerPort: 2112