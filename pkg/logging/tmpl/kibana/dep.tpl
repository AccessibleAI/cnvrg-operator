
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
        - name: "kibana-config"
          secret:
            secretName: "kibana-config"
      {{- end }}
      containers:
        {{- if eq .Spec.SSO.Enabled "true" }}
        - name: "cnvrg-oauth-proxy"
          image: {{ .Spec.SSO.Image }}
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
          image: {{ .Spec.Logging.Kibana.Image }}
          volumeMounts:
            - name: "kibana-config"
              mountPath: "/usr/share/kibana/config"
              readOnly: true
          env:
          {{- if eq .Spec.SSO.Enabled "true" }}
          - name: SERVER_PORT
            value: "3000"
          {{- end }}
          {{- if ne .Spec.SSO.Enabled "true" }}
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


