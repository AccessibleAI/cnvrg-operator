
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Spec.Logging.Kibana.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.Logging.Kibana.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ .Spec.Logging.Kibana.SvcName }}
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: {{ .Spec.Logging.Kibana.SvcName }}
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
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
      securityContext:
        runAsUser: 1000
        fsGroup: 1000
      serviceAccountName: {{ .Spec.Logging.Kibana.ServiceAccount }}
      volumes:
        - name: "kibana-config"
          secret:
            secretName: "kibana-config"
        {{- if isTrue .Spec.SSO.Enabled }}
        - name: "oauth-proxy-config"
          secret:
            secretName: "oauth-proxy-{{.Spec.Logging.Kibana.SvcName}}"
        {{- end }}
      containers:
        {{- if isTrue .Spec.SSO.Enabled }}
        - name: "cnvrg-oauth-proxy"
          image: {{.Spec.ImageHub }}/{{ .Spec.SSO.Image }}
          command: [ "oauth2-proxy","--config", "/opt/app-root/conf/proxy-config/conf" ]
          envFrom:
            - secretRef:
                name: {{ .Spec.Dbs.Redis.CredsRef }}
          volumeMounts:
            - name: "oauth-proxy-config"
              mountPath: "/opt/app-root/conf/proxy-config"
              readOnly: true
        {{- end }}
        - name: {{ .Spec.Logging.Kibana.SvcName }}
          image: {{.Spec.ImageHub }}/{{ .Spec.Logging.Kibana.Image }}
          volumeMounts:
            - name: "kibana-config"
              mountPath: "/usr/share/kibana/config"
              readOnly: true
          env:
          {{- if isTrue .Spec.SSO.Enabled }}
          - name: SERVER_PORT
            value: "3000"
          {{- else }}
          - name: SERVER_PORT
            value: "{{ .Spec.Logging.Kibana.Port }}"
          {{- end }}
          ports:
          - containerPort: {{ .Spec.Logging.Kibana.Port }}
          resources:
            limits:
              cpu: {{ .Spec.Logging.Kibana.Limits.Cpu }}
              memory: {{ .Spec.Logging.Kibana.Limits.Memory }}
            requests:
              cpu: {{ .Spec.Logging.Kibana.Requests.Cpu }}
              memory: {{ .Spec.Logging.Kibana.Requests.Memory }}
          lifecycle:
            postStart:
              exec:
                command:
                  - /bin/bash
                  - -c
                  - |
                    while [[ "$ready" != "200" ]]; do
                      ready=$(curl -s http://localhost:$SERVER_PORT/api/status -o /dev/null -w '%{http_code}')
                      echo "kibana not ready yet.. "
                      sleep 1
                    done
                    curl -XPOST "http://localhost:$SERVER_PORT/api/saved_objects/index-pattern/cnvrg" -H 'kbn-xsrf: true' -H 'Content-Type: application/json' -d '{"attributes":{"title": "cnvrg","timeFieldName": "@timestamp"}}'


