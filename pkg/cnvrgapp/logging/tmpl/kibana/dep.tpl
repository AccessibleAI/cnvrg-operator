
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Spec.Logging.Kibana.SvcName }}
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Spec.Logging.Kibana.SvcName }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ .Spec.Logging.Kibana.SvcName }}
  template:
    metadata:
      labels:
        app: {{ .Spec.Logging.Kibana.SvcName }}
    spec:
      serviceAccountName: {{ .Spec.ControlPlan.Rbac.ServiceAccountName }}
      {{ if eq .Spec.ControlPlan.Tenancy.Enabled "true" }}
      nodeSelector: 
        {{ .Spec.ControlPlan.Tenancy.Key }}: "{{ .Spec.ControlPlan.Tenancy.Value }}"
      {{- end }}
      tolerations:
        - key: {{ .Spec.ControlPlan.Tenancy.Key }}
          operator: Equal
          value: "{{ .Spec.ControlPlan.Tenancy.Value }}"
          effect: "NoSchedule"
      containers:
        - name: {{ .Spec.Logging.Kibana.SvcName }}
          image: {{ .Spec.Logging.Kibana.Image }}
          env:
          - name: ELASTICSEARCH_URL
            value: {{ printf "http://%s.%s.svc.cluster.local:%s" .Spec.Logging.Kibana.SvcName .Namespace .Spec.Logging.Es.Port }}
          ports:
          - containerPort: {{ .Spec.Logging.Kibana.Port }}
          resources:
            limits:
              cpu: {{ .Spec.Logging.Kibana.CPULimit }}
              memory: {{ .Spec.Logging.Kibana.MemoryLimit }}
            requests:
              cpu: {{ .Spec.Logging.Kibana.CPURequest }}
              memory: {{ .Spec.Logging.Kibana.MemoryRequest }}

