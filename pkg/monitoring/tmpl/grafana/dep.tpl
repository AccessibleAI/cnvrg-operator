apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: grafana
  name: grafana
  namespace: {{ ns . }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: grafana
  template:
    metadata:
      labels:
        app: grafana
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 65534
      serviceAccountName: grafana
      containers:
      - image: {{ .Spec.Monitoring.Grafana.Image }}
        name: grafana
        env:
          - name: GF_AUTH_BASIC_ENABLED
            value: "false"
          - name: GF_AUTH_ANONYMOUS_ENABLED
            value: "true"
          - name: GF_AUTH_ANONYMOUS_ORG_ROLE
            value: Admin
          - name: GF_SECURITY_ALLOW_EMBEDDING
            value: "true"
          {{- if eq .Spec.SSO.Enabled "true" }}
          - name: GF_SERVER_HTTP_ADDR
            value: "127.0.0.1"
          - name: GF_SERVER_HTTP_PORT
            value: "3000"
          {{- end }}
          {{- if ne .Spec.SSO.Enabled "true" }}
          - name: GF_SERVER_HTTP_ADDR
            value: "0.0.0.0"
          - name: GF_SERVER_HTTP_PORT
            value: "8080"
          {{- end }}
        ports:
        - containerPort: 3000
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
        {{- range $_, $dashboard := grafanaDashboards . }}
        - mountPath: /definitions/0/{{ trimSuffix ".json" $dashboard }}
          name: {{ trimSuffix ".json" $dashboard }}
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
      {{- range $_, $dashboard := grafanaDashboards . }}
      - configMap:
          name: {{ trimSuffix ".json" $dashboard }}
        name: {{ trimSuffix ".json" $dashboard }}
      {{- end }}
