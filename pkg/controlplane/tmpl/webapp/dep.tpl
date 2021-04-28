apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Spec.ControlPlane.WebApp.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: {{ .Spec.ControlPlane.WebApp.SvcName }}
spec:
  replicas: {{ .Spec.ControlPlane.WebApp.Replicas }}
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 25%
      maxSurge: 1
  selector:
    matchLabels:
      app: {{.Spec.ControlPlane.WebApp.SvcName}}
  template:
    metadata:
      labels:
        app: {{.Spec.ControlPlane.WebApp.SvcName}}
    spec:
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        "{{ .Spec.Tenancy.Key }}": "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - key: "{{ .Spec.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.Tenancy.Value }}"
          effect: "NoSchedule"
      {{- end }}
      serviceAccountName: {{ .Spec.ControlPlane.Rbac.ServiceAccountName }}
      securityContext:
        runAsUser: 1000
        runAsGroup: 1000
      containers:
      {{- if isTrue .Spec.SSO.Enabled }}
      - name: "cnvrg-oauth-proxy"
        image: {{ .Spec.SSO.Image }}
        command: [ "oauth2-proxy","--config", "/opt/app-root/conf/proxy-config/conf" ]
        envFrom:
          - secretRef:
              name: {{ .Spec.Dbs.Redis.CredsRef }}
        volumeMounts:
          - name: "oauth-proxy-webapp"
            mountPath: "/opt/app-root/conf/proxy-config"
            readOnly: true
      {{- end }}
      - image: {{ .Spec.ControlPlane.WebApp.Image }}
        env:
        - name: "CNVRG_RUN_MODE"
          value: "webapp"
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
        - secretRef:
            name: {{ .Spec.Dbs.Es.CredsRef }}
        - secretRef:
            name: {{ .Spec.Dbs.Pg.CredsRef }}
        - secretRef:
            name: {{ .Spec.Dbs.Redis.CredsRef }}
        - secretRef:
            name: {{ .Spec.Monitoring.Prometheus.CredsRef }}
        name: cnvrg-app
        ports:
          - containerPort: {{ .Spec.ControlPlane.WebApp.Port }}
        readinessProbe:
          httpGet:
            path: "/healthz"
            port: {{ .Spec.ControlPlane.WebApp.Port }}
            scheme: HTTP
          successThreshold: 1
          failureThreshold: {{ .Spec.ControlPlane.WebApp.FailureThreshold }}
          initialDelaySeconds: {{ .Spec.ControlPlane.WebApp.InitialDelaySeconds }}
          periodSeconds: {{ .Spec.ControlPlane.WebApp.ReadinessPeriodSeconds }}
          timeoutSeconds: {{ .Spec.ControlPlane.WebApp.ReadinessTimeoutSeconds }}
        resources:
          requests:
            cpu: "{{.Spec.ControlPlane.WebApp.CPU}}"
            memory: "{{.Spec.ControlPlane.WebApp.Memory}}"
        {{- if eq .Spec.ControlPlane.ObjectStorage.CnvrgStorageType "gcp" }}
        volumeMounts:
        - name: "{{ .Spec.ControlPlane.ObjectStorage.GcpStorageSecret }}"
          mountPath: "{{ .Spec.ControlPlane.ObjectStorage.GcpKeyfileMountPath }}"
          readOnly: true
        {{- end }}
      {{- if isTrue .Spec.SSO.Enabled }}
      volumes:
      - name: "oauth-proxy-webapp"
        secret:
         secretName: "oauth-proxy-webapp"
      {{- end }}
      {{- if eq .Spec.ControlPlane.ObjectStorage.CnvrgStorageType "gcp" }}
      - name: {{ .Spec.ControlPlane.ObjectStorage.GcpStorageSecret }}
        secret:
          secretName: {{ .Spec.ControlPlane.ObjectStorage.GcpStorageSecret }}
      {{- end }}
      initContainers:
      - name: services-check
        image: {{.Spec.ControlPlane.Seeder.Image}}
        command: ["/bin/bash", "-c", "python3 cnvrg-boot.py services-check"]
        imagePullPolicy: Always
        env:
        - name: "CNVRG_SERVICE_LIST"
          {{- if and ( isTrue .Spec.Dbs.Minio.Enabled ) (eq .Spec.ControlPlane.ObjectStorage.CnvrgStorageType "minio") }}
          value: "{{ .Spec.Dbs.Pg.SvcName }}:{{ .Spec.Dbs.Pg.Port }};{{ objectStorageUrl . }}/minio/health/ready"
          {{- else }}
          value: "{{ .Spec.Dbs.Pg.SvcName }}:{{ .Spec.Dbs.Pg.Port }}"
          {{ end }}
      {{- if and ( isTrue .Spec.Dbs.Minio.Enabled ) (eq .Spec.ControlPlane.ObjectStorage.CnvrgStorageType "minio") }}
      - name: create-cnvrg-bucket
        image: {{ .Spec.ControlPlane.Seeder.Image }}
        command: ["/bin/bash","-c", "{{ .Spec.ControlPlane.Seeder.CreateBucketCmd }}"]
        imagePullPolicy: Always
        envFrom:
        - secretRef:
            name: cp-object-storage
      {{- end }}
      {{- if isTrue .Spec.Dbs.Pg.Fixpg }}
      - name: fixpg
        image: {{.Spec.ControlPlane.Seeder.Image}}
        command: ["/bin/bash", "-c", "python3 cnvrg-boot.py fixpg"]
        envFrom:
        - secretRef:
            name: {{ .Spec.Dbs.Pg.CredsRef }}
        imagePullPolicy: Always
      {{- end }}
      - name: seeder
        image: {{ .Spec.ControlPlane.Seeder.Image }}
        command: ["/bin/bash", "-c", "python3 cnvrg-boot.py seeder --mode master"]
        imagePullPolicy: Always
        env:
        - name: "CNVRG_SEEDER_IMAGE"
          value: "{{.Spec.ControlPlane.WebApp.Image}}"
        - name: "CNVRG_SEED_CMD"
          value: "{{ .Spec.ControlPlane.Seeder.SeedCmd }}"
        - name: "CNVRG_NS"
          value: {{ ns . }}
        - name: "CNVRG_SA_NAME"
          value: {{ .Spec.ControlPlane.Rbac.ServiceAccountName }}
        - name: "CNVRG_PG_CREDS"
          value: {{ .Spec.Dbs.Pg.CredsRef }}
        - name: "CNVRG_REDIS_CREDS"
          value: {{ .Spec.Dbs.Redis.CredsRef }}
        - name: CNVRG_TENANCY
          value: "{{ isTrue .Spec.Tenancy.Enabled }}"
        - name: CNVRG_TENANCY_KEY
          value: "{{ .Spec.Tenancy.Key }}"
        - name: CNVRG_TENANCY_VALUE
          value: "{{ .Spec.Tenancy.Value }}"
        {{- if eq .Spec.ControlPlane.ObjectStorage.CnvrgStorageType "gcp" }}
        - name: "CNVRG_GCP_KEYFILE_SECRET"
          value: "{{ .Spec.ControlPlane.OjbectStorage.GcpStorageSecret }}"
        - name: "CNVRG_GCP_KEYFILE_MOUNT_PATH"
          value: "{{ .Spec.ControlPlane.OjbectStorage.GcpKeyfileMountPath }}"
        {{- end }}


