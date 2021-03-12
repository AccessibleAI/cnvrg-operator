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
      {{- if eq .ControlPlan.Conf.Tenancy.Enabled "true" }}
      nodeSelector:
        {{ .ControlPlan.Conf.Tenancy.Key }}: "{{ .ControlPlan.Conf.Tenancy.Value }}"
      {{- end }}
      tolerations:
      - key: "{{ .ControlPlan.Conf.Tenancy.Key }}"
        operator: "Equal"
        value: "{{ .ControlPlan.Conf.Tenancy.Value }}"
        effect: "NoSchedule"
      serviceAccountName: {{ .ControlPlan.Conf.Rbac.ServiceAccountName }}
      containers:
      {{- if eq .ControlPlan.Conf.OauthProxy.Enabled "true" }}
      - name: "cnvrg-oauth-proxy"
        image: {{ .ControlPlan.Conf.OauthProxy.Image }}
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
            name: env-config
        - secretRef:
            name: env-secrets
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
        {{- if eq .ControlPlan.Conf.ResourcesRequestEnabled "true" }}
        resources:
          requests:
            cpu: "{{.ControlPlan.WebApp.CPU}}"
            memory: "{{.ControlPlan.WebApp.Memory}}"
        {{- end }}
        {{- if eq .ControlPlan.Conf.CnvrgStorageType "gcp" }}
        volumeMounts:
        - name: "{{ .ControlPlan.Conf.GcpStorageSecret }}"
          mountPath: "{{ .ControlPlan.Conf.GcpKeyfileMountPath }}"
          readOnly: true
        {{- end }}
      {{- if eq .ControlPlan.Conf.OauthProxy.Enabled "true" }}
      volumes:
      - name: "oauth-proxy-config"
        configMap:
         name: "oauth-proxy-config"
      {{- end }}
      {{- if eq .ControlPlan.Conf.CnvrgStorageType "gcp" }}
      - name: {{ .ControlPlan.Conf.GcpStorageSecret }}
        secret:
          secretName: {{ .ControlPlan.Conf.GcpStorageSecret }}
      {{- end }}
      initContainers:
      - name: services-check
        image: {{.ControlPlan.Seeder.Image}}
        command: ["/bin/bash", "-c", "python3 cnvrg-boot.py services-check"]
        imagePullPolicy: Always
        env:
        - name: "CNVRG_SERVICE_LIST"
          {{- if and ( eq .Minio.Enabled "true") (eq .ControlPlan.Conf.CnvrgStorageType "minio") }}
          value: "{{.Pg.SvcName}}:{{.Pg.Port}};{{.ControlPlan.Conf.CnvrgStorageEndpoint}}/minio/health/ready"
          {{- else }}
          value: "{{.Pg.SvcName}}:{{.Pg.Port}}"
          {{ end }}
      {{- if and ( eq .Minio.Enabled "true") (eq .ControlPlan.Conf.CnvrgStorageType "minio") }}
      - name: create-cnvrg-bucket
        image: {{ .ControlPlan.Seeder.Image }}
        command: ["/bin/bash","-c", "{{ .ControlPlan.Seeder.CreateBucketCmd }}"]
        imagePullPolicy: Always
        envFrom:
        - configMapRef:
            name: "env-config"
        - secretRef:
            name: "env-secrets"
      {{- end }}
      {{- if eq .ControlPlan.Conf.Fixpg "true" }}
      - name: fixpg
        image: {{.ControlPlan.Seeder.Image}}
        command: ["/bin/bash", "-c", "python3 cnvrg-boot.py fixpg"]
        envFrom:
        - configMapRef:
            name: "env-config"
        - secretRef:
            name: "env-secrets"
        imagePullPolicy: Always
      {{- end }}
      - name: seeder
        image: {{ .ControlPlan.Seeder.Image }}
        command: ["/bin/bash", "-c", "python3 cnvrg-boot.py seeder --mode master"]
        imagePullPolicy: Always
        env:
        - name: "CNVRG_SEEDER_IMAGE"
          value: "{{.ControlPlan.Seeder.Image}}"
        - name: "CNVRG_SEED_CMD"
          value: "{{ .ControlPlan.Seeder.SeedCmd }}"
        - name: "CNVRG_NS"
          value: {{ .CnvrgNs }}
        - name: "CNVRG_SA_NAME"
          value: "cnvrg-control-plan"
        {{- if eq .ControlPlan.Conf.CnvrgStorageType "gcp" }}
        - name: "CNVRG_GCP_KEYFILE_SECRET"
          value: "{{ .ControlPlan.Conf.GcpStorageSecret }}"
        - name: "CNVRG_GCP_KEYFILE_MOUNT_PATH"
          value: "{{ .ControlPlan.Conf.GcpKeyfileMountPath }}"
        {{- end }}


