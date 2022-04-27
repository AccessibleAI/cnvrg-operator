apiVersion: apps/v1
kind: Deployment
metadata:
  name: cnvrg-exporter
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: cnvrg-exporter
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cnvrg-exporter
  template:
    metadata:
      labels:
        app: cnvrg-exporter
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      serviceAccountName: cnvrg-control-plane
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        "{{ .Spec.Tenancy.Key }}": "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - operator: "Exists"
      {{- end }}
      containers:
        - name: cnvrg-exporter
          image: cnvrg/tenant-metrics-exporter:v1
          imagePullPolicy: Always
          env:
            - name: APP_URL
              value: "{{ httpScheme . }}{{ appDomain . }}"
          resources:
            requests:
              cpu: 100m
              memory: 100Mi
            limits:
              cpu: 200m
              memory: 200Mi