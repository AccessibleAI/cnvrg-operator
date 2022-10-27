apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Data.Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Data.Spec.Dbs.Prom.Grafana.SvcName }}
    {{- range $k, $v := .Data.Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
  name: {{ .Data.Spec.Dbs.Prom.Grafana.SvcName }}
  namespace: {{ .Data.Namespace }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ .Data.Spec.Dbs.Prom.Grafana.SvcName }}
  template:
    metadata:
      labels:
        app: {{ .Data.Spec.Dbs.Prom.Grafana.SvcName }}
    spec:
      {{- if isTrue .Data.Spec.Tenancy.Enabled }}
      nodeSelector:
        "{{ .Data.Spec.Tenancy.Key }}": "{{ .Data.Spec.Tenancy.Value }}"
      tolerations:
        - key: "{{ .Data.Spec.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Data.Spec.Tenancy.Value }}"
          effect: "NoSchedule"
      {{- end }}
      priorityClassName: {{ .Spec.PriorityClass.AppClassRef }}
      serviceAccountName: {{ .Data.Spec.Dbs.Prom.Grafana.SvcName }}
      containers:
      - image: {{image .Data.Spec.ImageHub .Data.Spec.Dbs.Prom.Grafana.Image }}
        name: grafana
        securityContext:
          runAsNonRoot: true
          runAsUser: 65534
        env:
          - name: GF_AUTH_BASIC_ENABLED
            value: "false"
          - name: GF_AUTH_ANONYMOUS_ENABLED
            value: "true"
          - name: GF_AUTH_ANONYMOUS_ORG_ROLE
            value: Admin
          - name: GF_SECURITY_ALLOW_EMBEDDING
            value: "true"
          - name: GF_SERVER_HTTP_ADDR
            value: "0.0.0.0"
          - name: GF_SERVER_HTTP_PORT
            value: "8080"
        ports:
        - containerPort: 8080
          name: http
        readinessProbe:
          httpGet:
            path: /api/health
            port: http
        resources:
          limits:
            cpu: 200m
            memory: 200Mi
          requests:
            cpu: 100m
            memory: 100Mi
        volumeMounts:
        - mountPath: /var/lib/grafana
          name: grafana-storage
          readOnly: false
        - mountPath: /etc/grafana/provisioning/datasources
          name: grafana-datasources
          readOnly: false
        - mountPath: /etc/grafana/provisioning/dashboards
          name: grafana-dashboards
          readOnly: false
        {{- range $_, $dashboard := .Dashboards }}
        - mountPath: /definitions/0/{{ $dashboard }}
          name: {{ $dashboard }}
          readOnly: false
        {{- end }}
      volumes:
      - emptyDir: {}
        name: grafana-storage
      - name: grafana-datasources
        secret:
          secretName: grafana-datasources
      - configMap:
          name: grafana-dashboards
        name: grafana-dashboards
      {{- range $_, $dashboard := .Dashboards }}
      - configMap:
          name: {{ $dashboard }}
        name: {{ $dashboard }}
      {{- end }}
