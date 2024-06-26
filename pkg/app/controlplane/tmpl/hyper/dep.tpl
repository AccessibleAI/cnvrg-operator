apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Spec.ControlPlane.Hyper.SvcName }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    owner: cnvrg-control-plane
    app: {{ .Spec.ControlPlane.Hyper.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  replicas: {{ .Spec.ControlPlane.Hyper.Replicas }}
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 25%
      maxSurge: 1
  selector:
    matchLabels:
      app: {{ .Spec.ControlPlane.Hyper.SvcName }}
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: {{ .Spec.ControlPlane.Hyper.SvcName }}
        owner: cnvrg-control-plane
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      priorityClassName: {{ .Spec.PriorityClass.AppClassRef }}
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        "{{ .Spec.Tenancy.Key }}": "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - key: "{{ .Spec.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.Tenancy.Value }}"
          effect: "NoSchedule"
      {{- end }}
      serviceAccountName: cnvrg-control-plane
      enableServiceLinks: false
      containers:
        - image: {{ image .Spec.ImageHub .Spec.ControlPlane.Hyper.Image }}
          name: {{ .Spec.ControlPlane.Hyper.SvcName }}
          envFrom:
            - configMapRef:
                name: cp-base-config
            - configMapRef:
                name: cp-networking-config
            - secretRef:
                name: cp-base-secret
            - secretRef:
                name: cp-ldap
            - secretRef:
                name: cp-object-storage
            {{- if isTrue .Spec.Networking.Proxy.Enabled }}
            - configMapRef:
                name: {{ .Spec.Networking.Proxy.ConfigRef }}
            {{- end }}
          ports:
            - containerPort: {{ .Spec.ControlPlane.Hyper.Port }}
          readinessProbe:
            failureThreshold: 3
            httpGet:
              path: "/?key={{.Spec.ControlPlane.Hyper.Token}}"
              port: {{.Spec.ControlPlane.Hyper.Port}}
              scheme: HTTP
            initialDelaySeconds: 20
            successThreshold: 1
            periodSeconds: {{ .Spec.ControlPlane.Hyper.ReadinessPeriodSeconds }}
            timeoutSeconds: {{ .Spec.ControlPlane.Hyper.ReadinessTimeoutSeconds }}
          resources:
            requests:
              cpu: {{.Spec.ControlPlane.Hyper.Requests.Cpu}}
              memory: {{.Spec.ControlPlane.Hyper.Requests.Memory}}
            limits:
              cpu: {{ .Spec.ControlPlane.Hyper.Limits.Cpu }}
              memory: {{ .Spec.ControlPlane.Hyper.Limits.Memory }}