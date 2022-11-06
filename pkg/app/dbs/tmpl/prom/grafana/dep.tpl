apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.Dbs.Prom.Grafana.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
  name: {{ .Spec.Dbs.Prom.Grafana.SvcName }}
  namespace: {{ .Namespace }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ .Spec.Dbs.Prom.Grafana.SvcName }}
  template:
    metadata:
      labels:
        app: {{ .Spec.Dbs.Prom.Grafana.SvcName }}
    spec:
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        "{{ .Spec.Tenancy.Key }}": "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - key: "{{ .Spec.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.Tenancy.Value }}"
          effect: "NoSchedule"
      {{- end }}
      priorityClassName: {{ .Spec.PriorityClass.AppClassRef }}
      serviceAccountName: {{ .Spec.Dbs.Prom.Grafana.SvcName }}
      enableServiceLinks: false
      containers:
      - image: {{image .Spec.ImageHub .Spec.Dbs.Prom.Grafana.Image }}
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
            value: "{{ .Spec.Dbs.Prom.Grafana.Port }}"
        ports:
        - containerPort: {{ .Spec.Dbs.Prom.Grafana.Port }}
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
