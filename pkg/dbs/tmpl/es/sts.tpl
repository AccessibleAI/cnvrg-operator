apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: {{ .Spec.Dbs.Es.SvcName }}
  namespace: {{ ns . }}
spec:
  serviceName: {{ .Spec.Dbs.Es.SvcName }}
  selector:
    matchLabels:
      app: {{ .Spec.Dbs.Es.SvcName }}
  replicas: 1
  template:
    metadata:
      labels:
        app: {{ .Spec.Dbs.Es.SvcName }}
    spec:
      serviceAccountName: {{ .Spec.Dbs.Es.ServiceAccount }}
      {{- if .Spec.Dbs.Es.PatchEsNodes }}
      initContainers:
      - name: "maxmap"
        image: "docker.io/cnvrg/cnvrg-tools:v0.3"
        imagePullPolicy: "Always"
        command: [ "/bin/bash","-c","sysctl -w vm.max_map_count=262144"]
        securityContext:
          privileged: true
          runAsUser: 0
      {{- end }}
      securityContext:
        runAsUser: {{ .Spec.Dbs.Es.RunAsUser }}
        fsGroup: {{ .Spec.Dbs.Es.FsGroup }}
      containers:
      - name: elastic
        image: {{ .Spec.Dbs.Es.Image }}
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
        ports:
        - containerPort: {{ .Spec.Dbs.Es.Port }}
        resources:
          limits:
            cpu: {{ .Spec.Dbs.Es.CPULimit }}
            memory: {{ .Spec.Dbs.Es.MemoryLimit }}
          requests:
            cpu: {{ .Spec.Dbs.Es.CPURequest }}
            memory: {{ .Spec.Dbs.Es.MemoryRequest }}
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
      name: es-storage
    spec:
      accessModes: [ ReadWriteOnce ]
      resources:
        requests:
          storage: {{ .Spec.Dbs.Es.StorageSize }}
      {{- if ne .Spec.Dbs.Es.StorageClass "" }}
      storageClassName: {{ .Spec.Dbs.Es.StorageClass }}
      {{- end }}