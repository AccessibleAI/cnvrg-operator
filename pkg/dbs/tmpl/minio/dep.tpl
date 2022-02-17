apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Spec.Dbs.Minio.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.Dbs.Minio.SvcName }}
    cnvrg-component: minio
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  selector:
    matchLabels:
      app: {{ .Spec.Dbs.Minio.SvcName }}
  replicas: {{ .Spec.Dbs.Minio.Replicas }}
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: {{ .Spec.Dbs.Minio.SvcName }}
        owner: cnvrg-control-plane
        cnvrg-component: minio
    spec:
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        {{ .Spec.Tenancy.Key }}: {{ .Spec.Tenancy.Value }}
        {{- range $key, $val := .Spec.Dbs.Minio.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      tolerations:
        - operator: "Exists"
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
          image: {{ image .Spec.ImageHub .Spec.Dbs.Minio.Image }}
          envFrom:
          - secretRef:
              name: cp-object-storage
          env:
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
              mountPath: /data # the Minio mount path have to be /data -> it's hardcoded into Minio server startup command
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
              cpu: {{ .Spec.Dbs.Minio.Requests.Cpu }}
              memory: {{ .Spec.Dbs.Minio.Requests.Memory }}
            limits:
              cpu: {{ .Spec.Dbs.Minio.Limits.Cpu }}
              memory: {{ .Spec.Dbs.Minio.Limits.Memory }}
      volumes:
        - name: minio-storage
          persistentVolumeClaim:
            claimName: {{ .Spec.Dbs.Minio.PvcName }}
