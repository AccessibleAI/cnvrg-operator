apiVersion: apps/v1
kind: Deployment
metadata:
  name: config-reloader
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: config-reloader
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: config-reloader
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: config-reloader
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      serviceAccountName: cnvrg-operator
      securityContext:
        runAsUser: 1000
        runAsGroup: 1000
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        {{ .Spec.Tenancy.Key }}: "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - key: "{{ .Spec.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.Tenancy.Value }}"
          effect: "NoSchedule"
      {{- end }}
      containers:
        - name: config-reloader
          image: {{.Spec.ImageHub}}/{{.Spec.ConfigReloader.Image}}
          imagePullPolicy: Always
          command:
            - /opt/app-root/config-reloader
            - --match-label
            - cnvrg-config-reloader.mlops.cnvrg.io
          resources:
            limits:
              cpu: 1000m
              memory: 1000Mi
            requests:
              cpu: 100m
              memory: 200Mi
