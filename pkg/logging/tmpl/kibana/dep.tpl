
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Logging.Kibana.SvcName }}
  namespace: {{ .CnvrgNs }}
  labels:
    app: {{ .Logging.Kibana.SvcName }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ .Logging.Kibana.SvcName }}
  template:
    metadata:
      labels:
        app: {{ .Logging.Kibana.SvcName }}
    spec:
      serviceAccountName: {{ .ControlPlan.Rbac.ServiceAccountName }}
      {{ if eq .ControlPlan.Tenancy.Enabled "true" }}
      nodeSelector: 
        {{ .ControlPlan.Tenancy.Key }}: "{{ .ControlPlan.Tenancy.Value }}"
      {{- end }}
      tolerations:
        - key: {{ .ControlPlan.Tenancy.Key }}
          operator: Equal
          value: "{{ .ControlPlan.Tenancy.Value }}"
          effect: "NoSchedule"
      containers:
        - name: {{ .Logging.Kibana.SvcName }}
          image: {{ .Logging.Kibana.Image }}
          env:
          - name: ELASTICSEARCH_URL
            value: {{ printf "http://%s.%s.svc.cluster.local:%s" .Logging.Kibana.SvcName .CnvrgNs .Logging.Es.Port }}
          ports:
          - containerPort: {{ .Logging.Kibana.Port }}
          resources:
            limits:
              cpu: {{ .Logging.Kibana.CPULimit }}
              memory: {{ .Logging.Kibana.MemoryLimit }}
            requests:
              cpu: {{ .Logging.Kibana.CPURequest }}
              memory: {{ .Logging.Kibana.MemoryRequest }}

