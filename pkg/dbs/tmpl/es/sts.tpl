apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: {{ .Spec.Dbs.Es.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  serviceName: {{ .Spec.Dbs.Es.SvcName }}
  selector:
    matchLabels:
      app: {{ .Spec.Dbs.Es.SvcName }}
  replicas: 1
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: {{ .Spec.Dbs.Es.SvcName }}
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        {{ .Spec.Tenancy.Key }}: {{ .Spec.Tenancy.Value }}
        {{- range $key, $val := .Spec.Dbs.Es.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      tolerations:
        - operator: "Exists"
      {{- else if (gt (len .Spec.Dbs.Es.NodeSelector) 0) }}
      nodeSelector:
        {{- range $key, $val := .Spec.Dbs.Es.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      {{- end }}
      serviceAccountName: {{ .Spec.Dbs.Es.ServiceAccount }}
      {{- if isTrue .Spec.Dbs.Es.PatchEsNodes }}
      initContainers:
      - name: "maxmap"
        image: {{ image .Spec.ImageHub .Spec.Dbs.Es.Image }}
        imagePullPolicy: "Always"
        command: [ "/bin/bash","-c","sysctl -w vm.max_map_count=262144"]
        securityContext:
          privileged: true
          runAsUser: 0
        resources:
          limits:
            cpu: 200m
            memory: 200Mi
          requests:
            cpu: 100m
            memory: 100Mi
      {{- end }}
      securityContext:
        runAsUser: 1000
        fsGroup: 1000
      containers:
      - name: elastic
        image: {{ image .Spec.ImageHub .Spec.Dbs.Es.Image }}
        env:
        - name: "ES_CLUSTER_NAME"
          value: "cnvrg-es"
        - name: "ES_NODE_NAME"
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: "ES_NETWORK_HOST"
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: "ES_DISCOVERY_TYPE"
          value: "single-node"
        - name: "ES_PATH_DATA"
          value: "/usr/share/elasticsearch/data/data"
        - name: "ES_PATH_LOGS"
          value: "/usr/share/elasticsearch/data/logs"
        - name: "ES_JAVA_OPTS"
          value: "{{ .Spec.Dbs.Es.JavaOpts }}"
        - name: "ES_SECURITY_ENABLED"
          value: "true"
        envFrom:
          - secretRef:
              name: {{ .Spec.Dbs.Es.CredsRef }}
          - configMapRef:
              name: es-ilm-cm
        ports:
        - containerPort: {{ .Spec.Dbs.Es.Port }}
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
              - /bin/bash
              - -c
              - |
                ready=$(curl -s -u$CNVRG_ES_USER:$CNVRG_ES_PASS http://$ES_NETWORK_HOST:9200/_cluster/health -o /dev/null -w '%{http_code}')
                if [ "$ready" == "200" ]; then
                  exit 0
                else
                  exit 1
                fi
          initialDelaySeconds: 30
          periodSeconds: 20
          failureThreshold: 5
        livenessProbe:
          exec:
            command:
              - /bin/bash
              - -c
              - |
                ready=$(curl -s -u$CNVRG_ES_USER:$CNVRG_ES_PASS http://$ES_NETWORK_HOST:9200/_cluster/health -o /dev/null -w '%{http_code}')
                if [ "$ready" == "200" ]; then
                  exit 0;
                else
                  exit 1
                fi
          initialDelaySeconds: 5
          periodSeconds: 20
          failureThreshold: 5
        volumeMounts:
        - name: es-storage
          mountPath: "/usr/share/elasticsearch/data"
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