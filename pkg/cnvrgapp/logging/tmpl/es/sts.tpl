apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: {{ .Logging.Es.SvcName }}
  namespace: {{ .CnvrgNs }}
spec:
  serviceName: {{ .Logging.Es.SvcName }}
  selector:
    matchLabels:
      app: {{ .Logging.Es.SvcName }}
  replicas: 1
  template:
    metadata:
      labels:
        app: {{ .Logging.Es.SvcName }}
    spec:
      {{- if eq .Logging.Es.PatchEsNodes "true" }}
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
        runAsUser: {{ .Logging.Es.RunAsUser }}
        fsGroup: {{ .Logging.Es.FsGroup }}
      serviceAccountName: {{ .ControlPlan.Rbac.ServiceAccountName }}
      {{- if and (eq .Storage.Hostpath.Enabled "true") (eq .ControlPlan.Tenancy.Enabled "false") }}
      nodeSelector:
        kubernetes.io/hostname: "{{ .Storage.Hostpath.NodeName }}"
      {{- else if and (eq .Storage.Hostpath.Enabled "false") (eq .ControlPlan.Tenancy.Enabled "true") }}
      nodeSelector:
        {{ .ControlPlan.Tenancy.Key }}: "{{ .ControlPlan.Tenancy.Value }}"
      {{- else if and (eq .Storage.Hostpath.Enabled "true") (eq .ControlPlan.Tenancy.Enabled "true") }}
      nodeSelector:
        kubernetes.io/hostname: "{{ .Storage.Hostpath.NodeName }}"
        {{ .ControlPlan.Tenancy.Key }}: "{{ .ControlPlan.Tenancy.Value }}"
      {{- end }}
      tolerations:
        - key: {{ .ControlPlan.Tenancy.Key }}
          operator: Equal
          value: "{{ .ControlPlan.Tenancy.Value }}"
          effect: "NoSchedule"
      containers:
      - name: elastic
        image: {{ .Logging.Es.Image }}
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
          value: "{{ .Logging.Es.JavaOpts }}"
        ports:
        - containerPort: {{ .Logging.Es.Port }}
        resources:
          limits:
            cpu: {{ .Logging.Es.CPULimit }}
            memory: {{ .Logging.Es.MemoryLimit }}
          requests:
            cpu: {{ .Logging.Es.CPURequest }}
            memory: {{ .Logging.Es.MemoryRequest }}
        readinessProbe:
          httpGet:
            path: /_cluster/health
            port: {{ .Logging.Es.Port }}
          initialDelaySeconds: 30
          periodSeconds: 20
          failureThreshold: 5
        livenessProbe:
          httpGet:
            path: /_cluster/health
            port: {{ .Logging.Es.Port }}
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
          storage: {{ .Logging.Es.StorageSize }}
      {{- if ne .Logging.Es.StorageClass "use-default" }}
      storageClassName: {{ .Logging.Es.StorageClass }}
      {{- else if ne .Storage.CcpStorageClass "" }}
      storageClassName: {{ .Storage.CcpStorageClass }}
      {{- end }}