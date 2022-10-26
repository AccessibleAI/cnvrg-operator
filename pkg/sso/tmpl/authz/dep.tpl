apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Spec.SSO.Authz.SvcName }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
  labels:
    app: {{ .Spec.SSO.Authz.SvcName }}
spec:
  selector:
    matchLabels:
      app: {{ .Spec.SSO.Authz.SvcName }}
  template:
    metadata:
      labels:
        app: {{ .Spec.SSO.Authz.SvcName }}
    spec:
      serviceAccountName: {{ .Spec.SSO.Authz.SvcName }}
      containers:
      - name: authz
        imagePullPolicy: Always
        image: {{  image .Spec.ImageHub .Spec.SSO.Authz.Image }}
        command:
          - /opt/app-root/authz
          - --ingress-type={{.Spec.Networking.Ingress.Type}}