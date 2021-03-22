apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Spec.ControlPlan.Hyper.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: {{ .Spec.ControlPlan.Hyper.SvcName }}
spec:
  replicas: {{ .Spec.ControlPlan.Hyper.Replicas }}
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 25%
      maxSurge: 1
  selector:
    matchLabels:
      app: {{ .Spec.ControlPlan.Hyper.SvcName }}
  template:
    metadata:
      labels:
        app: {{ .Spec.ControlPlan.Hyper.SvcName }}
    spec:
      serviceAccountName: {{ .Spec.ControlPlan.Rbac.ServiceAccountName }}
      {{- if eq .Spec.ControlPlan.Tenancy.Enabled "true" }}
      nodeSelector:
        {{ .Spec.ControlPlan.Tenancy.Key }}: "{{ .Spec.ControlPlan.Tenancy.Value }}"
      {{- end }}
      tolerations:
        - key: "{{ .Spec.ControlPlan.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.ControlPlan.Tenancy.Value }}"
          effect: "NoSchedule"
      containers:
        - image: {{ .Spec.ControlPlan.Hyper.Image }}
          name: {{ .Spec.ControlPlan.Hyper.SvcName }}
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
            - containerPort: {{ .Spec.ControlPlan.Hyper.Port }}
          readinessProbe:
            failureThreshold: 3
            httpGet:
              path: "/?key={{.Spec.ControlPlan.Hyper.Token}}"
              port: {{.Spec.ControlPlan.Hyper.Port}}
              scheme: HTTP
            initialDelaySeconds: 20
            successThreshold: 1
            periodSeconds: {{ .Spec.ControlPlan.Hyper.ReadinessPeriodSeconds }}
            timeoutSeconds: {{ .Spec.ControlPlan.Hyper.ReadinessTimeoutSeconds }}
          resources:
            requests:
              cpu: {{.Spec.ControlPlan.Hyper.CPURequest}}
              memory: {{.Spec.ControlPlan.Hyper.MemoryRequest}}
            limits:
              cpu: {{ .Spec.ControlPlan.Hyper.CPULimit }}
              memory: {{ .Spec.ControlPlan.Hyper.MemoryLimit }}