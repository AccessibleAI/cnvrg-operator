apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Spec.Dbs.Minio.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: {{ .Spec.Dbs.Minio.SvcName }}
spec:
  selector:
    matchLabels:
      app: {{ .Spec.Dbs.Minio.SvcName }}
  replicas: {{ .Spec.Dbs.Minio.Replicas }}
  template:
    metadata:
      labels:
        app: {{ .Spec.Dbs.Minio.SvcName }}
    spec:
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        {{ .Spec.Tenancy.Key }}: {{ .Spec.Tenancy.Value }}
        {{- range $key, $val := .Spec.Dbs.Minio.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      tolerations:
        - key: "{{ .Spec.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.Tenancy.Value }}"
          effect: "NoSchedule"
      {{- else if (gt (len .Spec.Dbs.Minio.NodeSelector) 0) }}
      nodeSelector:
        {{- range $key, $val := .Spec.Dbs.Minio.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      {{- end }}
      securityContext:
        runAsUser: 1000
        fsGroup: 1000
      serviceAccountName: {{ .Spec.Dbs.Minio.ServiceAccount }}
      containers:
        - name: minio
          image: {{ .Spec.Dbs.Minio.Image }}
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
            - containerPort: {{ .Spec.Dbs.Minio.Port }}
          volumeMounts:
            - name: minio-storage
              mountPath: /data
          readinessProbe:
            httpGet:
              path: /minio/health/ready
              port: {{ .Spec.Dbs.Minio.Port }}
            initialDelaySeconds: 5
            periodSeconds: 5
          livenessProbe:
            httpGet:
              path: /minio/health/live
              port: {{ .Spec.Dbs.Minio.Port }}
            initialDelaySeconds: 60
            periodSeconds: 20
          resources:
            requests:
              cpu: {{ .Spec.Dbs.Minio.CPURequest }}
              memory: {{ .Spec.Dbs.Minio.MemoryRequest }}
      volumes:
        - name: minio-storage
          persistentVolumeClaim:
            claimName: {{ .Spec.Dbs.Minio.SvcName }}
