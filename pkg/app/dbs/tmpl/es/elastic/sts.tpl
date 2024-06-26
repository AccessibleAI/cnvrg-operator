{{- define "endpoints" -}}
{{- $replicas := int (toString ( .Spec.Dbs.Es.Replicas)) }}
  {{- range $i, $e := untilStep 0 $replicas 1 -}}
{{ $.Spec.Dbs.Es.SvcName }}-{{ $i }},
  {{- end -}}
{{- end -}}
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: {{ .Spec.Dbs.Es.SvcName }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.Dbs.Es.SvcName }}
    cnvrg-system-status-check: "true"
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  podManagementPolicy: Parallel
  replicas: {{ .Spec.Dbs.Es.Replicas }}
  updateStrategy:
    type: RollingUpdate
  selector:
    matchLabels:
      app: {{ .Spec.Dbs.Es.SvcName }}
  serviceName: {{ .Spec.Dbs.Es.SvcName }}-headless
  volumeClaimTemplates:
  - metadata:
      name: {{ .Spec.Dbs.Es.PvcName  }}
    spec:
      accessModes: [ ReadWriteOnce ]
      resources:
        requests:
          storage: {{ .Spec.Dbs.Es.StorageSize }}
      {{- if ne .Spec.Dbs.Es.StorageClass "" }}
      storageClassName: {{ .Spec.Dbs.Es.StorageClass }}
      {{- end }}
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: {{ .Spec.Dbs.Es.SvcName }}
        {{- range $k, $v := .Spec.Annotations }}
        {{- if eq $k "eastwest_custom_name" }}
        sidecar.istio.io/inject: "true"
        {{- end }}
        {{- end }}
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
      name: {{ .Spec.Dbs.Es.SvcName }}
    spec:
      securityContext:
        runAsUser: 1000
        fsGroup: 1000
      automountServiceAccountToken: true
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - {{ .Spec.Dbs.Es.SvcName }}
              topologyKey: kubernetes.io/hostname
            weight: 1
      volumes:
        - name: es-ilm
          configMap:
            name: "es-ilm"
            defaultMode: 0755
        - name: {{ .Spec.Dbs.Es.SvcName }}-certs
          secret:
           secretName: {{ .Spec.Dbs.Es.SvcName }}-certs
      priorityClassName: {{ .Spec.PriorityClass.AppClassRef }}
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        "{{ .Spec.Tenancy.Key }}": "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - key: "{{ .Spec.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.Tenancy.Value }}"
          effect: "NoSchedule"
      {{- else if (gt (len .Spec.Dbs.Es.NodeSelector) 0) }}
      nodeSelector:
        {{- range $key, $val := .Spec.Dbs.Es.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      {{- end }}
      serviceAccountName: {{ .Spec.Dbs.Es.ServiceAccount }}
      enableServiceLinks: false
      containers:
      - name: elastic
        image: {{ image .Spec.ImageHub .Spec.Dbs.Es.Image }}
        env:
        - name: cluster.name
          value: "cnvrg-es"
        - name: node.name
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: ELASTIC_PASSWORD
          valueFrom:
            secretKeyRef:
              name: {{ .Spec.Dbs.Es.CredsRef }}
              key: CNVRG_ES_PASS
        - name: ingest.geoip.downloader.enabled
          value: "false"
        - name: network.host
          value: "0.0.0.0"
        - name: cluster.initial_master_nodes
          value: "{{ template "endpoints" . }}"
        - name: discovery.seed_hosts
          value: "{{ .Spec.Dbs.Es.SvcName }}-headless"
        - name: cluster.deprecation_indexing.enabled
          value: "false"
        - name: "path.data"
          value: "/usr/share/elasticsearch/data/data"
        - name: "path.logs"
          value: "/usr/share/elasticsearch/data/logs"
        - name: "ES_JAVA_OPTS"
          value: "{{ .Spec.Dbs.Es.JavaOpts }}"
        - name: xpack.security.enabled
          value: "true"
        - name: xpack.security.authc.realms.native.native1.order
          value: "0"
        - name: xpack.security.authc.realms.file.file1.order
          value: "1"
        - name: xpack.security.transport.ssl.enabled
          value: "true"
        - name: xpack.security.transport.ssl.verification_mode
          value: certificate
        - name: xpack.security.transport.ssl.key
          value: /usr/share/elasticsearch/config/certs/tls.key
        - name: xpack.security.transport.ssl.certificate
          value: /usr/share/elasticsearch/config/certs/tls.crt
        - name: xpack.security.transport.ssl.certificate_authorities
          value: /usr/share/elasticsearch/config/certs/ca.crt
        - name: node.store.allow_mmap
          value: "false"
        envFrom:
          - secretRef:
              name: {{ .Spec.Dbs.Es.CredsRef }}
          - configMapRef:
              name: es-ilm-cm
          {{- if isTrue .Spec.Networking.Proxy.Enabled }}
          - configMapRef:
              name: {{ .Spec.Networking.Proxy.ConfigRef }}
          {{- end }}
        ports:
        - name: http
          containerPort: {{ .Spec.Dbs.Es.Port }}
        - name: transport
          containerPort: 9300
        resources:
          limits:
            cpu: {{ .Spec.Dbs.Es.Limits.Cpu }}
            memory: {{ .Spec.Dbs.Es.Limits.Memory }}
          requests:
            cpu: {{ .Spec.Dbs.Es.Requests.Cpu }}
            memory: {{ .Spec.Dbs.Es.Requests.Memory }}
        readinessProbe:
          exec:
            command:
              - bash
              - -c
              - |
                set -e
                # If the node is starting up wait for the cluster to be ready (request params: "wait_for_status=yellow&timeout=1s" )
                # Once it has started only check that the node itself is responding
                START_FILE=/tmp/.es_start_file

                # Disable nss cache to avoid filling dentry cache when calling curl
                # This is required with Elasticsearch Docker using nss < 3.52
                export NSS_SDB_USE_CACHE=no

                http () {
                  local path="${1}"
                  local args="${2}"
                  set -- -XGET -s

                  if [ "$args" != "" ]; then
                    set -- "$@" $args
                  fi

                  if [ -n "${CNVRG_ES_PASS}" ]; then
                    set -- "$@" -u "${CNVRG_ES_USER}:${CNVRG_ES_PASS}"
                  fi

                  curl --output /dev/null -k "$@" "http://127.0.0.1:9200${path}"
                }

                if [ -f "${START_FILE}" ]; then
                  echo 'Elasticsearch is already running, lets check the node is healthy'
                  HTTP_CODE=$(http "/" "-w %{http_code}")
                  RC=$?
                  if [[ ${RC} -ne 0 ]]; then
                    echo "curl --output /dev/null -k -XGET -s -w '%{http_code}' \${BASIC_AUTH} http://127.0.0.1:9200/ failed with RC ${RC}"
                    exit ${RC}
                  fi
                  # ready if HTTP code 200, 503 is tolerable if ES version is 6.x
                  if [[ ${HTTP_CODE} == "200" ]]; then
                    exit 0
                  elif [[ ${HTTP_CODE} == "503" && "7" == "6" ]]; then
                    exit 0
                  else
                    echo "curl --output /dev/null -k -XGET -s -w '%{http_code}' \${BASIC_AUTH} http://127.0.0.1:9200/ failed with HTTP code ${HTTP_CODE}"
                    exit 1
                  fi

                else
                  echo 'Waiting for elasticsearch cluster to become ready (request params: "wait_for_status=yellow&timeout=1s" )'
                  if http "/_cluster/health?wait_for_status=yellow&timeout=1s" "--fail" ; then
                    touch ${START_FILE}
                    exit 0
                  else
                    echo 'Cluster is not yet ready (request params: "wait_for_status=yellow&timeout=1s" )'
                    exit 1
                  fi
                fi
          failureThreshold: 3
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 3
          timeoutSeconds: 5
        lifecycle:
          postStart:
            exec:
              command:
              - /bin/bash
              - -c
              - /tmp/elastic/elastic_cleanup.sh&
        volumeMounts:
        - name: {{ .Spec.Dbs.Es.PvcName  }}
          mountPath: "/usr/share/elasticsearch/data"
        - name: {{ .Spec.Dbs.Es.SvcName }}-certs
          mountPath: /usr/share/elasticsearch/config/certs
          readOnly: true
        - name: es-ilm
          mountPath: "/tmp/elastic/"