
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Spec.Logging.Kibana.SvcName }}
  namespace: {{ ns . }}
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
      securityContext:
        runAsUser: 1000
        fsGroup: 1000
      serviceAccountName: {{ .Spec.Logging.Kibana.ServiceAccount }}
      {{- if eq .Spec.SSO.Enabled "true" }}
      volumes:
        - name: "oauth-proxy-config"
          secret:
            secretName: "oauth-proxy-{{.Spec.Logging.Kibana.SvcName}}"
      {{- end }}
      containers:
        {{- if eq .Spec.SSO.Enabled "true" }}
        - name: "cnvrg-oauth-proxy"
          image: {{ .Spec.SSO.Image }}
          command: [ "oauth2-proxy","--config", "/opt/app-root/conf/proxy-config/conf" ]
          volumeMounts:
            - name: "oauth-proxy-config"
              mountPath: "/opt/app-root/conf/proxy-config"
              readOnly: true
        {{- end }}
        - name: {{ .Spec.Logging.Kibana.SvcName }}
          image: {{ .Spec.Logging.Kibana.Image }}
          env:
          - name: ELASTICSEARCH_URL
            value: {{ esFullInternalUrl .}}
          {{- if eq .Spec.SSO.Enabled "true" }}
          - name: SERVER_HOST
            value: "127.0.0.1"
          - name: SERVER_PORT
            value: "3000"
          {{- end }}
          {{- if ne .Spec.SSO.Enabled "true" }}
          - name: SERVER_HOST
            value: "0.0.0.0"
          - name: SERVER_PORT
            value: "{{ .Spec.Logging.Kibana.Port }}"
          {{- end }}
          ports:
          - containerPort: {{ .Spec.Logging.Kibana.Port }}
          resources:
            limits:
              cpu: {{ .Spec.Logging.Kibana.CPULimit }}
              memory: {{ .Spec.Logging.Kibana.MemoryLimit }}
            requests:
              cpu: {{ .Spec.Logging.Kibana.CPURequest }}
              memory: {{ .Spec.Logging.Kibana.MemoryRequest }}

