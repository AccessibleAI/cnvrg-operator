apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: {{ .Spec.Logging.Es.SvcName }}
  namespace: {{ ns . }}
spec:
  serviceName: {{ .Spec.Logging.Es.SvcName }}
  selector:
    matchLabels:
      app: {{ .Spec.Logging.Es.SvcName }}
  replicas: 1
  template:
    metadata:
      labels:
        app: {{ .Spec.Logging.Es.SvcName }}
    spec:
      {{- if eq .Spec.Logging.Es.PatchEsNodes "true" }}
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
        runAsUser: {{ .Spec.Logging.Es.RunAsUser }}
        fsGroup: {{ .Spec.Logging.Es.FsGroup }}
      serviceAccountName: {{ .Spec.ControlPlan.Rbac.ServiceAccountName }}
      {{- if and (ne .Spec.ControlPlan.BaseConfig.HostpathNode "") (eq .Spec.ControlPlan.Tenancy.Enabled "false") }}
      nodeSelector:
        kubernetes.io/hostname: "{{ .Spec.ControlPlan.BaseConfig.HostpathNode }}"
      {{- else if and (eq .Spec.ControlPlan.BaseConfig.HostpathNode "") (eq .Spec.ControlPlan.Tenancy.Enabled "true") }}
      nodeSelector:
        {{ .Spec.ControlPlan.Tenancy.Key }}: "{{ .Spec.ControlPlan.Tenancy.Value }}"
      {{- else if and (ne .Spec.ControlPlan.BaseConfig.HostpathNode "") (eq .Spec.ControlPlan.Tenancy.Enabled "true") }}
      nodeSelector:
        kubernetes.io/hostname: "{{ .Spec.ControlPlan.BaseConfig.HostpathNode }}"
        {{ .Spec.ControlPlan.Tenancy.Key }}: "{{ .Spec.ControlPlan.Tenancy.Value }}"
      {{- end }}
      tolerations:
        - key: {{ .Spec.ControlPlan.Tenancy.Key }}
          operator: Equal
          value: "{{ .Spec.ControlPlan.Tenancy.Value }}"
          effect: "NoSchedule"
      containers:
      - name: elastic
        image: {{ .Spec.Logging.Es.Image }}
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
          value: "{{ .Spec.Logging.Es.JavaOpts }}"
        ports:
        - containerPort: {{ .Spec.Logging.Es.Port }}
        resources:
          limits:
            cpu: {{ .Spec.Logging.Es.CPULimit }}
            memory: {{ .Spec.Logging.Es.MemoryLimit }}
          requests:
            cpu: {{ .Spec.Logging.Es.CPURequest }}
            memory: {{ .Spec.Logging.Es.MemoryRequest }}
        readinessProbe:
          httpGet:
            path: /_cluster/health
            port: {{ .Spec.Logging.Es.Port }}
          initialDelaySeconds: 30
          periodSeconds: 20
          failureThreshold: 5
        livenessProbe:
          httpGet:
            path: /_cluster/health
            port: {{ .Spec.Logging.Es.Port }}
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
          storage: {{ .Spec.Logging.Es.StorageSize }}
      {{- if ne .Spec.Logging.Es.StorageClass "use-default" }}
      storageClassName: {{ .Spec.Logging.Es.StorageClass }}
      {{- else if ne .Spec.ControlPlan.BaseConfig.CcpStorageClass "" }}
      storageClassName: {{ .Spec.ControlPlan.BaseConfig.CcpStorageClass }}
      {{- end }}