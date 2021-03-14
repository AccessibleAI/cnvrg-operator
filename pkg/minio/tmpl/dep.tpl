apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Minio.SvcName }}
  namespace: {{ .CnvrgNs }}
  labels:
    app: {{ .Minio.SvcName }}
spec:
  selector:
    matchLabels:
      app: {{ .Minio.SvcName }}
  replicas: {{ .Minio.Replicas }}
  template:
    metadata:
      labels:
        app: {{ .Minio.SvcName }}
    spec:
      securityContext:
        runAsUser: 1000
        fsGroup: 1000
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
        - name: minio
          image: {{ .Minio.Image }}
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
            - containerPort: {{ .Minio.Port }}
          volumeMounts:
            - name: minio-storage
              mountPath: /data
          readinessProbe:
            httpGet:
              path: /minio/health/ready
              port: {{ .Minio.Port }}
            initialDelaySeconds: 5
            periodSeconds: 5
          livenessProbe:
            httpGet:
              path: /minio/health/live
              port: {{ .Minio.Port }}
            initialDelaySeconds: 60
            periodSeconds: 20
          resources:
            requests:
              cpu: {{ .Minio.CPURequest }}
              memory: {{ .Minio.MemoryRequest }}
      volumes:
        - name: minio-storage
          persistentVolumeClaim:
            claimName: {{ .Minio.SvcName }}
