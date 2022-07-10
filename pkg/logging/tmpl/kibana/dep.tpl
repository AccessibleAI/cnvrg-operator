
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
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-kibana-oauth"
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
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
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
      serviceAccountName: {{ .Spec.Logging.Kibana.SvcName }}
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
          image: {{image .Spec.ImageHub .Spec.SSO.Image }}
          command: [ "oauth2-proxy","--config", "/opt/app-root/conf/proxy-config/conf" ]
          resources:
            requests:
              cpu: 100m
              memory: 100m
            limits:
              cpu: 500m
              memory: 1Gi
          envFrom:
            - secretRef:
                name: {{ .Spec.Dbs.Redis.CredsRef }}
          volumeMounts:
            - name: "oauth-proxy-config"
              mountPath: "/opt/app-root/conf/proxy-config"
              readOnly: true
        {{- end }}
        - name: {{ .Spec.Logging.Kibana.SvcName }}
          image: {{image .Spec.ImageHub .Spec.Logging.Kibana.Image }}
          command:
            - /bin/bash
            - -lc
            - |
              #!/bin/bash
              {
                cnvrgIndexPattern=cnvrg*
                ready=notready
                while [[ "$ready" != "200" ]]; do
                  ready=$(curl -s http://localhost:$SERVER_PORT/api/status -o /dev/null -w '%{http_code}')
                  echo "[$(date)][cnvrg-init] kibana not ready yet.. "
                  sleep 1
                done
                echo "[$(date)][cnvrg-init] kibana is ready!"
                cnvrgIndexPatternExists=$(curl -s http://localhost:$SERVER_PORT/api/saved_objects/index-pattern/$cnvrgIndexPattern -o /dev/null -w '%{http_code}')
                if [[ "$cnvrgIndexPatternExists" == "200" ]]; then
                  echo "[$(date)][cnvrg-init] cnvrg index pattern found, skip index creation!"
                fi
                if [[ "$cnvrgIndexPatternExists" == "404" ]]; then
                  echo "[$(date)][cnvrg-init] cnvrg index pattern not found, going to create one"
                  curl -XPOST "http://localhost:$SERVER_PORT/api/saved_objects/index-pattern/$cnvrgIndexPattern" \
                     -H 'kbn-xsrf: true' \
                     -H 'Content-Type: application/json' \
                     -d '{"attributes":{"title": "cnvrg","timeFieldName": "@timestamp"}}'
                  echo "[$(date)][cnvrg-init] Index created!"
                fi
              } &
              /usr/local/bin/kibana-docker
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