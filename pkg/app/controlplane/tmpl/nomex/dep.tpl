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
      priorityClassName: {{ .Spec.PriorityClass.AppClassRef }}
      serviceAccountName: cnvrg-nomex
      enableServiceLinks: false
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        "{{ .Spec.Tenancy.Key }}": "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - operator: "Exists"
      {{- end }}
      containers:
      - name: nomex
        imagePullPolicy: Always
        image: {{  image .Spec.ImageHub .Spec.ControlPlane.Nomex.Image }}
        command:
          - /opt/app-root/nomex
        ports:
        - containerPort: 2112
        envFrom:
          {{- if isTrue .Spec.Networking.Proxy.Enabled }}
          - configMapRef:
              name: {{ .Spec.Networking.Proxy.ConfigRef }}
          {{- end }}
