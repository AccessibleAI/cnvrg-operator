
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Spec.Dbs.Es.Kibana.SvcName }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.Dbs.Es.Kibana.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ .Spec.Dbs.Es.Kibana.SvcName }}
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: {{ .Spec.Dbs.Es.Kibana.SvcName }}
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      priorityClassName: {{ .Spec.PriorityClass.AppClassRef }}
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
      serviceAccountName: {{ .Spec.Dbs.Es.Kibana.SvcName }}
      volumes:
        - name: "kibana-config"
          secret:
            secretName: "kibana-config"
      enableServiceLinks: false
      containers:
        - name: {{ .Spec.Dbs.Es.Kibana.SvcName }}
          image: {{image .Spec.ImageHub .Spec.Dbs.Es.Kibana.Image }}
          command:
            - /bin/bash
            - -lc
            - |
              #!/bin/bash
              {
                ready=notready
                while [[ "$ready" != "200" ]]; do
                  ready=$(curl -s http://localhost:$SERVER_PORT/api/status -o /dev/null -w '%{http_code}')
                  echo "[$(date)][cnvrg-init] kibana not ready yet.. "
                  sleep 1
                done
                echo "[$(date)][cnvrg-init] kibana is ready!"
                for cnvrgIndexPattern in "cnvrg*" "cnvrg-endpoints*"; do
                  cnvrgIndexPatternExists=$(curl -s http://localhost:$SERVER_PORT/api/saved_objects/index-pattern/$cnvrgIndexPattern -o /dev/null -w '%{http_code}')
                  if [[ "$cnvrgIndexPatternExists" == "200" ]]; then
                    echo "[$(date)][cnvrg-init] cnvrg index pattern found, skip index creation!"
                  fi
                  if [[ "$cnvrgIndexPatternExists" == "404" ]]; then
                    echo "[$(date)][cnvrg-init] cnvrg index pattern not found, going to create one"
                    curl -XPOST "http://localhost:$SERVER_PORT/api/saved_objects/index-pattern/$cnvrgIndexPattern" \
                       -H 'kbn-xsrf: true' \
                       -H 'Content-Type: application/json' \
                       -d '{"attributes":{"title": "'$cnvrgIndexPattern'","timeFieldName": "@timestamp"}}'
                    echo "[$(date)][cnvrg-init] Index $cnvrgIndexPattern created!"
                  fi
                done
              } &
              /usr/local/bin/kibana-docker
          volumeMounts:
            - name: "kibana-config"
              mountPath: "/usr/share/kibana/config"
              readOnly: true
          env:
          - name: SERVER_PORT
            value: "{{ .Spec.Dbs.Es.Kibana.Port }}"
          ports:
          - containerPort: {{ .Spec.Dbs.Es.Kibana.Port }}
          resources:
            limits:
              cpu: {{ .Spec.Dbs.Es.Kibana.Limits.Cpu }}
              memory: {{ .Spec.Dbs.Es.Kibana.Limits.Memory }}
            requests:
              cpu: {{ .Spec.Dbs.Es.Kibana.Requests.Cpu }}
              memory: {{ .Spec.Dbs.Es.Kibana.Requests.Memory }}