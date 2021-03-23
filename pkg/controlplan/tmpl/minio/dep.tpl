apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Spec.ControlPlan.Minio.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: {{ .Spec.ControlPlan.Minio.SvcName }}
spec:
  selector:
    matchLabels:
      app: {{ .Spec.ControlPlan.Minio.SvcName }}
  replicas: {{ .Spec.ControlPlan.Minio.Replicas }}
  template:
    metadata:
      labels:
        app: {{ .Spec.ControlPlan.Minio.SvcName }}
    spec:
      securityContext:
        runAsUser: 1000
        fsGroup: 1000
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
        - name: minio
          image: {{ .Spec.ControlPlan.Minio.Image }}
          args:
            - gateway
            - nas
            - /data
          env:
            - name: MINIO_SSE_MASTER_KEY
              valueFrom:
                secretKeyRef:
                  name: cp-object-storage
                  key: MINIO_SSE_MASTER_KEY
            - name: MINIO_ACCESS_KEY
              valueFrom:
                secretKeyRef:
                  name: cp-object-storage
                  key: CNVRG_STORAGE_ACCESS_KEY
            - name: MINIO_SECRET_KEY
              valueFrom:
                secretKeyRef:
                  name: cp-object-storage
                  key: CNVRG_STORAGE_SECRET_KEY
          ports:
            - containerPort: {{ .Spec.ControlPlan.Minio.Port }}
          volumeMounts:
            - name: minio-storage
              mountPath: /data
          readinessProbe:
            httpGet:
              path: /minio/health/ready
              port: {{ .Spec.ControlPlan.Minio.Port }}
            initialDelaySeconds: 5
            periodSeconds: 5
          livenessProbe:
            httpGet:
              path: /minio/health/live
              port: {{ .Spec.ControlPlan.Minio.Port }}
            initialDelaySeconds: 60
            periodSeconds: 20
          resources:
            requests:
              cpu: {{ .Spec.ControlPlan.Minio.CPURequest }}
              memory: {{ .Spec.ControlPlan.Minio.MemoryRequest }}
      volumes:
        - name: minio-storage
          persistentVolumeClaim:
            claimName: {{ .Spec.ControlPlan.Minio.SvcName }}
