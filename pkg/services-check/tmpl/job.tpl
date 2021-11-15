apiVersion: batch/v1
kind: Job
metadata:
  name: services-check
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
      {{$k}}: "{{$v}}"
      {{- end }}
  labels:
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-ccp"
    app: services-check
    owner: cnvrg-control-plane
    cnvrg-component: services-check
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  template:
    metadata:
      name: services-check
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
          {{$k}}: "{{$v}}"
          {{- end }}
      labels:
        cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-ccp"
        app: services-check
        owner: cnvrg-control-plane
        cnvrg-component: services-check
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      containers:
        - name: services-check
          image: ubuntu:20.10
          envFrom:
            - configMapRef:
                name: cp-base-config
            - configMapRef:
                name: cp-networking-config
            - secretRef:
                name: cp-base-secret
            - secretRef:
                name: cp-ldap
            - secretRef:
                name: cp-object-storage
            - secretRef:
                name: cp-smtp
            {{- if isTrue .Spec.Dbs.Es.Enabled }}
            - secretRef:
                name: {{ .Spec.Dbs.Es.CredsRef }}
            {{- end }}
            {{- if isTrue .Spec.Dbs.Redis.Enabled }}
            - secretRef:
                name: {{ .Spec.Dbs.Redis.CredsRef }}
            {{- end }}
            {{- if isTrue .Spec.Monitoring.Prometheus.Enabled }}
            - secretRef:
                name: {{ .Spec.Monitoring.Prometheus.CredsRef }}
            {{- end }}
            {{- if isTrue .Spec.Networking.Proxy.Enabled }}
            - configMapRef:
                name: {{ .Spec.Networking.Proxy.ConfigRef }}
            {{- end }}
          command:
            - "/bin/bash"
            - "-lc"
            - |
              #!/bin/bash
              apt-get update && apt-get install -y curl
              flagFile=/tmp/services_not_ready
              echo true > ${flagFile}
              while $(cat ${flagFile}); do

                {{- if ( isTrue .Spec.Dbs.Es.Enabled ) }}
                if [[ $(curl -sk {{ httpScheme . }}$CNVRG_ES_USER:$CNVRG_ES_PASS@{{ .Spec.Dbs.Es.SvcName }}.{{ .Spec.ClusterDomain }} -o /dev/null -w '%{http_code}') != 200 ]]; then
                  echo "[$(date)] elasticsearch [{{ httpScheme . }}$CNVRG_ES_USER:$CNVRG_ES_PASS@{{ .Spec.Dbs.Es.SvcName }}.{{ .Spec.ClusterDomain }}] not ready"
                  sleep 1
                  continue
                fi
                echo "[$(date)] elasticsearch is ready!"
                {{- end }}

                {{- if ( isTrue .Spec.Monitoring.Grafana.Enabled ) }}
                if [[ $(curl -sk {{ httpScheme . }}{{ .Spec.Monitoring.Grafana.SvcName }}.{{ .Spec.ClusterDomain }} -o /dev/null -w '%{http_code}') != 200 ]]; then
                  echo "[$(date)] grafana [{{ httpScheme . }}{{ .Spec.Monitoring.Grafana.SvcName }}.{{ .Spec.ClusterDomain }}] not ready"
                  sleep 1
                  continue
                fi
                echo "[$(date)] grafana is ready!"
                {{- end }}

                {{- if ( isTrue .Spec.Logging.Kibana.Enabled ) }}
                if [[ $(curl -sk {{ httpScheme . }}{{ .Spec.Logging.Kibana.SvcName }}.{{ .Spec.ClusterDomain }} -o /dev/null -w '%{http_code}') != 200 ]]; then
                  echo "[$(date)] kibana [{{ httpScheme . }}{{ .Spec.Logging.Kibana.SvcName }}.{{ .Spec.ClusterDomain }}] not ready"
                  sleep 1
                  continue
                fi
                echo "[$(date)] kibana is ready!"
                {{- end }}

                {{- if ( isTrue .Spec.Monitoring.Prometheus.Enabled ) }}
                if [[ $(curl -sk {{ httpScheme . }}$CNVRG_PROMETHEUS_USER:$CNVRG_PROMETHEUS_PASS@{{ .Spec.Monitoring.Prometheus.SvcName }}.{{ .Spec.ClusterDomain }} -o /dev/null -w '%{http_code}') != 200 ]]; then
                  echo "[$(date)] prometheus [{{ httpScheme . }}$CNVRG_PROMETHEUS_USER:$CNVRG_PROMETHEUS_PASS@{{ .Spec.Monitoring.Prometheus.SvcName }}.{{ .Spec.ClusterDomain }}] not ready"
                  sleep 1
                  continue
                fi
                echo "[$(date)] prometheus is ready!"
                {{- end }}

                {{- if and ( isTrue .Spec.Dbs.Minio.Enabled ) (eq .Spec.ControlPlane.ObjectStorage.Type "minio") }}
                if [[ $(curl -sk $CNVRG_STORAGE_ENDPOINT/minio/health/ready -o /dev/null -w '%{http_code}') != 200 ]]; then
                  echo "[$(date)] minio [$CNVRG_STORAGE_ENDPOINT/minio/health/ready] not ready"
                  sleep 1
                  continue
                fi
                echo "[$(date)] minio is ready!"
                {{- end }}

                echo false > ${flagFile}
                echo "[$(date)] all services are ready!"
              done
      restartPolicy: Never
  backoffLimit: 4