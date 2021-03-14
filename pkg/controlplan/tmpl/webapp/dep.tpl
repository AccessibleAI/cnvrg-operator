apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .ControlPlan.WebApp.SvcName }}
  namespace: {{ .CnvrgNs }}
  labels:
    app: {{ .ControlPlan.WebApp.SvcName }}
spec:
  replicas: {{ .ControlPlan.WebApp.Replicas }}
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 25%
      maxSurge: 1
  selector:
    matchLabels:
      app: {{.ControlPlan.WebApp.SvcName}}
  template:
    metadata:
      labels:
        app: {{.ControlPlan.WebApp.SvcName}}
    spec:
      {{- if eq .ControlPlan.Tenancy.Enabled "true" }}
      nodeSelector:
        {{ .ControlPlan.Tenancy.Key }}: "{{ .ControlPlan.Tenancy.Value }}"
      {{- end }}
      tolerations:
      - key: "{{ .ControlPlan.Tenancy.Key }}"
        operator: "Equal"
        value: "{{ .ControlPlan.Tenancy.Value }}"
        effect: "NoSchedule"
      serviceAccountName: {{ .ControlPlan.Rbac.ServiceAccountName }}
      containers:
      {{- if eq .ControlPlan.OauthProxy.Enabled "true" }}
      - name: "cnvrg-oauth-proxy"
        image: {{ .ControlPlan.OauthProxy.Image }}
        command: [ "oauth2-proxy","--config", "/opt/app-root/conf/proxy-config/conf" ]
        volumeMounts:
          - name: "oauth-proxy-config"
            mountPath: "/opt/app-root/conf/proxy-config"
            readOnly: true
      {{- end }}
      - image: {{ .ControlPlan.WebApp.Image }}
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
            name: {{ .Pg.SvcName }}
        name: cnvrg-app
        ports:
          - containerPort: {{ .ControlPlan.WebApp.Port }}
        readinessProbe:
          httpGet:
            path: "/healthz"
            port: {{ .ControlPlan.WebApp.Port }}
            scheme: HTTP
          successThreshold: 1
          failureThreshold: {{ .ControlPlan.WebApp.FailureThreshold }}
          initialDelaySeconds: {{ .ControlPlan.WebApp.InitialDelaySeconds }}
          periodSeconds: {{ .ControlPlan.WebApp.ReadinessPeriodSeconds }}
          timeoutSeconds: {{ .ControlPlan.WebApp.ReadinessTimeoutSeconds }}
        resources:
          requests:
            cpu: "{{.ControlPlan.WebApp.CPU}}"
            memory: "{{.ControlPlan.WebApp.Memory}}"
        {{- if eq .ControlPlan.ObjectStorage.CnvrgStorageType "gcp" }}
        volumeMounts:
        - name: "{{ .ControlPlan.ObjectStorage.GcpStorageSecret }}"
          mountPath: "{{ .ControlPlan.ObjectStorage.GcpKeyfileMountPath }}"
          readOnly: true
        {{- end }}
      {{- if eq .ControlPlan.OauthProxy.Enabled "true" }}
      volumes:
      - name: "oauth-proxy-config"
        secret:
         secretName: "cp-sso"
      {{- end }}
      {{- if eq .ControlPlan.ObjectStorage.CnvrgStorageType "gcp" }}
      - name: {{ .ControlPlan.ObjectStorage.GcpStorageSecret }}
        secret:
          secretName: {{ .ControlPlan.ObjectStorage.GcpStorageSecret }}
      {{- end }}
      initContainers:
      - name: services-check
        image: {{.ControlPlan.Seeder.Image}}
        command: ["/bin/bash", "-c", "python3 cnvrg-boot.py services-check"]
        imagePullPolicy: Always
        env:
        - name: "CNVRG_SERVICE_LIST"
          {{- if and ( eq .Minio.Enabled "true") (eq .ControlPlan.ObjectStorage.CnvrgStorageType "minio") }}
          value: "{{.Pg.SvcName}}:{{.Pg.Port}};{{.ControlPlan.ObjectStorage.CnvrgStorageEndpoint}}/minio/health/ready"
          {{- else }}
          value: "{{.Pg.SvcName}}:{{.Pg.Port}}"
          {{ end }}
      {{- if and ( eq .Minio.Enabled "true") (eq .ControlPlan.ObjectStorage.CnvrgStorageType "minio") }}
      - name: create-cnvrg-bucket
        image: {{ .ControlPlan.Seeder.Image }}
        command: ["/bin/bash","-c", "{{ .ControlPlan.Seeder.CreateBucketCmd }}"]
        imagePullPolicy: Always
        envFrom:
        - secretRef:
            name: cp-object-storage
      {{- end }}
      {{- if eq .Pg.Fixpg "true" }}
      - name: fixpg
        image: {{.ControlPlan.Seeder.Image}}
        command: ["/bin/bash", "-c", "python3 cnvrg-boot.py fixpg"]
        envFrom:
        - secretRef:
            name: {{ .Pg.SvcName }}
        imagePullPolicy: Always
      {{- end }}
      - name: seeder
        image: {{ .ControlPlan.Seeder.Image }}
        command: ["/bin/bash", "-c", "python3 cnvrg-boot.py seeder --mode master"]
        imagePullPolicy: Always
        env:
        - name: "CNVRG_SEEDER_IMAGE"
          value: "{{.ControlPlan.WebApp.Image}}"
        - name: "CNVRG_SEED_CMD"
          value: "{{ .ControlPlan.Seeder.SeedCmd }}"
        - name: "CNVRG_NS"
          value: {{ .CnvrgNs }}
        - name: "CNVRG_SA_NAME"
          value: "cnvrg-control-plan"
        {{- if eq .ControlPlan.ObjectStorage.CnvrgStorageType "gcp" }}
        - name: "CNVRG_GCP_KEYFILE_SECRET"
          value: "{{ .ControlPlan.OjbectStorage.GcpStorageSecret }}"
        - name: "CNVRG_GCP_KEYFILE_MOUNT_PATH"
          value: "{{ .ControlPlan.OjbectStorage.GcpKeyfileMountPath }}"
        {{- end }}


