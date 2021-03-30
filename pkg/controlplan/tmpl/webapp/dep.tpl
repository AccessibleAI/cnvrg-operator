apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Spec.ControlPlan.WebApp.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: {{ .Spec.ControlPlan.WebApp.SvcName }}
spec:
  replicas: {{ .Spec.ControlPlan.WebApp.Replicas }}
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 25%
      maxSurge: 1
  selector:
    matchLabels:
      app: {{.Spec.ControlPlan.WebApp.SvcName}}
  template:
    metadata:
      labels:
        app: {{.Spec.ControlPlan.WebApp.SvcName}}
    spec:
      {{- if eq .Spec.ControlPlan.Tenancy.Enabled "true" }}
      nodeSelector:
        {{ .Spec.ControlPlan.Tenancy.Key }}: "{{ .Spec.ControlPlan.Tenancy.Value }}"
      {{- end }}
      tolerations:
      - key: "{{ .Spec.ControlPlan.Tenancy.Key }}"
        operator: "Equal"
        value: "{{ .Spec.ControlPlan.Tenancy.Value }}"
        effect: "NoSchedule"
      serviceAccountName: {{ .Spec.ControlPlan.Rbac.ServiceAccountName }}
      containers:
      {{- if eq .Spec.SSO.Enabled "true" }}
      - name: "cnvrg-oauth-proxy"
        image: {{ .Spec.SSO.Image }}
        command: [ "oauth2-proxy","--config", "/opt/app-root/conf/proxy-config/conf" ]
        volumeMounts:
          - name: "oauth-proxy-webapp"
            mountPath: "/opt/app-root/conf/proxy-config"
            readOnly: true
      {{- end }}
      - image: {{ .Spec.ControlPlan.WebApp.Image }}
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
            name: {{ .Spec.Dbs.Pg.SvcName }}
        name: cnvrg-app
        ports:
          - containerPort: {{ .Spec.ControlPlan.WebApp.Port }}
        readinessProbe:
          httpGet:
            path: "/healthz"
            port: {{ .Spec.ControlPlan.WebApp.Port }}
            scheme: HTTP
          successThreshold: 1
          failureThreshold: {{ .Spec.ControlPlan.WebApp.FailureThreshold }}
          initialDelaySeconds: {{ .Spec.ControlPlan.WebApp.InitialDelaySeconds }}
          periodSeconds: {{ .Spec.ControlPlan.WebApp.ReadinessPeriodSeconds }}
          timeoutSeconds: {{ .Spec.ControlPlan.WebApp.ReadinessTimeoutSeconds }}
        resources:
          requests:
            cpu: "{{.Spec.ControlPlan.WebApp.CPU}}"
            memory: "{{.Spec.ControlPlan.WebApp.Memory}}"
        {{- if eq .Spec.ControlPlan.ObjectStorage.CnvrgStorageType "gcp" }}
        volumeMounts:
        - name: "{{ .Spec.ControlPlan.ObjectStorage.GcpStorageSecret }}"
          mountPath: "{{ .Spec.ControlPlan.ObjectStorage.GcpKeyfileMountPath }}"
          readOnly: true
        {{- end }}
      {{- if eq .Spec.SSO.Enabled "true" }}
      volumes:
      - name: "oauth-proxy-webapp"
        secret:
         secretName: "oauth-proxy-webapp"
      {{- end }}
      {{- if eq .Spec.ControlPlan.ObjectStorage.CnvrgStorageType "gcp" }}
      - name: {{ .Spec.ControlPlan.ObjectStorage.GcpStorageSecret }}
        secret:
          secretName: {{ .Spec.ControlPlan.ObjectStorage.GcpStorageSecret }}
      {{- end }}
      initContainers:
      - name: services-check
        image: {{.Spec.ControlPlan.Seeder.Image}}
        command: ["/bin/bash", "-c", "python3 cnvrg-boot.py services-check"]
        imagePullPolicy: Always
        env:
        - name: "CNVRG_SERVICE_LIST"
          {{- if and ( eq .Spec.Dbs.Minio.Enabled "true") (eq .Spec.ControlPlan.ObjectStorage.CnvrgStorageType "minio") }}
          value: "{{ .Spec.Dbs.Pg.SvcName }}:{{ .Spec.Dbs.Pg.Port }};{{ objectStorageUrl . }}/minio/health/ready"
          {{- else }}
          value: "{{ .Spec.Dbs.Pg.SvcName }}:{{ .Spec.Dbs.Pg.Port }}"
          {{ end }}
      {{- if and ( eq .Spec.Dbs.Minio.Enabled "true") (eq .Spec.ControlPlan.ObjectStorage.CnvrgStorageType "minio") }}
      - name: create-cnvrg-bucket
        image: {{ .Spec.ControlPlan.Seeder.Image }}
        command: ["/bin/bash","-c", "{{ .Spec.ControlPlan.Seeder.CreateBucketCmd }}"]
        imagePullPolicy: Always
        envFrom:
        - secretRef:
            name: cp-object-storage
      {{- end }}
      {{- if eq .Spec.Dbs.Pg.Fixpg "true" }}
      - name: fixpg
        image: {{.Spec.ControlPlan.Seeder.Image}}
        command: ["/bin/bash", "-c", "python3 cnvrg-boot.py fixpg"]
        envFrom:
        - secretRef:
            name: {{ .Spec.Dbs.Pg.SvcName }}
        imagePullPolicy: Always
      {{- end }}
      - name: seeder
        image: {{ .Spec.ControlPlan.Seeder.Image }}
        command: ["/bin/bash", "-c", "python3 cnvrg-boot.py seeder --mode master"]
        imagePullPolicy: Always
        env:
        - name: "CNVRG_SEEDER_IMAGE"
          value: "{{.Spec.ControlPlan.WebApp.Image}}"
        - name: "CNVRG_SEED_CMD"
          value: "{{ .Spec.ControlPlan.Seeder.SeedCmd }}"
        - name: "CNVRG_NS"
          value: {{ ns . }}
        - name: "CNVRG_SA_NAME"
          value: {{ .Spec.ControlPlan.Rbac.ServiceAccountName }}
        {{- if eq .Spec.ControlPlan.ObjectStorage.CnvrgStorageType "gcp" }}
        - name: "CNVRG_GCP_KEYFILE_SECRET"
          value: "{{ .Spec.ControlPlan.OjbectStorage.GcpStorageSecret }}"
        - name: "CNVRG_GCP_KEYFILE_MOUNT_PATH"
          value: "{{ .Spec.ControlPlan.OjbectStorage.GcpKeyfileMountPath }}"
        {{- end }}


