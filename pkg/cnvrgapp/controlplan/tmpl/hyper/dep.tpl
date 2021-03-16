apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .ControlPlan.Hyper.SvcName }}
  namespace: {{ .CnvrgNs }}
  labels:
    app: {{ .ControlPlan.Hyper.SvcName }}
spec:
  replicas: {{ .ControlPlan.Hyper.Replicas }}
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 25%
      maxSurge: 1
  selector:
    matchLabels:
      app: {{ .ControlPlan.Hyper.SvcName }}
  template:
    metadata:
      labels:
        app: {{ .ControlPlan.Hyper.SvcName }}
    spec:
      serviceAccountName: {{ .ControlPlan.Rbac.ServiceAccountName }}
      {{- if eq .ControlPlan.Tenancy.Enabled "true" }}
      nodeSelector:
        {{ .ControlPlan.Tenancy.Key }}: "{{ .ControlPlan.Tenancy.Value }}"
      {{- end }}
      tolerations:
        - key: "{{ .ControlPlan.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .ControlPlan.Tenancy.Value }}"
          effect: "NoSchedule"
      containers:
        - image: {{ .ControlPlan.Hyper.Image }}
          name: {{ .ControlPlan.Hyper.SvcName }}
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
          ports:
            - containerPort: {{ .ControlPlan.Hyper.Port }}
          readinessProbe:
            failureThreshold: 3
            httpGet:
              path: "/?key={{.ControlPlan.Hyper.Token}}"
              port: {{.ControlPlan.Hyper.Port}}
              scheme: HTTP
            initialDelaySeconds: 20
            successThreshold: 1
            periodSeconds: {{ .ControlPlan.Hyper.ReadinessPeriodSeconds }}
            timeoutSeconds: {{ .ControlPlan.Hyper.ReadinessTimeoutSeconds }}
          resources:
            requests:
              cpu: {{.ControlPlan.Hyper.CPURequest}}
              memory: {{.ControlPlan.Hyper.MemoryRequest}}
            limits:
              cpu: {{ .ControlPlan.Hyper.CPULimit }}
              memory: {{ .ControlPlan.Hyper.MemoryLimit }}