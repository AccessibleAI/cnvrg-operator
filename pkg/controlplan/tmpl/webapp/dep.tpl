apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Spec.ControlPlan.WebApp.SvcName }}
  namespace: {{ .Spec.CnvrgNs }}
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
      {{- if eq .Spec.ControlPlan.Conf.Tenancy.Enabled "true" }}
      nodeSelector:
        {{ .Spec.ControlPlan.Conf.Tenancy.Key }}: "{{ .Spec.ControlPlan.Conf.Tenancy.Value }}"
      {{- end }}
      tolerations:
      - key: "{{ .Spec.ControlPlan.Conf.Tenancy.Key }}"
        operator: "Equal"
        value: "{{ .Spec.ControlPlan.Conf.Tenancy.Value }}"
        effect: "NoSchedule"
      serviceAccountName: {{ .Spec.ControlPlan.Conf.Rbac.ServiceAccountName }}
      containers:
      {{- if eq .Spec.ControlPlan.Conf.OauthProxy.Enabled "true" }}
      - name: "cnvrg-oauth-proxy"
        image: {{ .Spec.ControlPlan.Conf.OauthProxy.Image }}
        command: [ "oauth2-proxy","--config", "/opt/app-root/conf/proxy-config/conf" ]
        volumeMounts:
          - name: "oauth-proxy-config"
            mountPath: "/opt/app-root/conf/proxy-config"
            readOnly: true
      {{- end }}
      - image: {{ .Spec.ControlPlan.WebApp.Image }}
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
          - containerPort: {{ .Spec.ControlPlan.WebApp.Port}}
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
        {{- if eq .Spec.ControlPlan.Conf.ResourcesRequestEnabled "true" }}
        resources:
          requests:
            cpu: "{{.Spec.ControlPlan.WebApp.CPU}}"
            memory: "{{.Spec.ControlPlan.WebApp.Memory}}"
        {{- end }}
        {{- if eq .Spec.ControlPlan.Conf.CnvrgStorageType "gcp" }}
        volumeMounts:
        - name: "{{ .Spec.ControlPlan.Conf.GcpStorageSecret }}"
          mountPath: "{{ .Spec.ControlPlan.Conf.GcpKeyfileMountPath }}"
          readOnly: true
        {{- end }}
      {{- if eq .Spec.ControlPlan.Conf.OauthProxy.Enabled "true" }}
      volumes:
      - name: "oauth-proxy-config"
        configMap:
         name: "oauth-proxy-config"
      {{- end }}
      {{- if eq .Spec.ControlPlan.Conf.CnvrgStorageType "gcp" }}
      - name: {{ .Spec.ControlPlan.Conf.GcpStorageSecret }}
        secret:
          secretName: {{ .Spec.ControlPlan.Conf.GcpStorageSecret }}
      {{- end }}
      initContainers:
      - name: services-check
        image: {{.Spec.ControlPlan.Seeder.Image}}
        command: ["/bin/bash", "-c", "python3 cnvrg-boot.py services-check"]
        imagePullPolicy: Always
        env:
        - name: "CNVRG_SERVICE_LIST"
          {{- if and ( eq .Spec.Minio.Enabled "true") (eq .Spec.ControlPlan.Conf.CnvrgStorageType "minio") }}
          value: "{{.Spec.Pg.SvcName}}:{{.Spec.Pg.Port}};{{.Spec.ControlPlan.Conf.CnvrgStorageEndpoint}}/minio/health/ready"
          {{- else }}
          value: "{{.Spec.Pg.SvcName}}:{{.Spec.Pg.Port}}"
          {{ end }}
      {{- if and ( eq .Spec.Minio.Enabled "true") (eq .Spec.ControlPlan.Conf.CnvrgStorageType "minio") }}
      - name: create-cnvrg-bucket
        image: {{ .Spec.ControlPlan.Seeder.Image }}
        command: ["/bin/bash","-c", "{{ .Spec.ControlPlan.Seeder.CreateBucketCmd }}"]
        imagePullPolicy: Always
        envFrom:
        - configMapRef:
            name: "env-config"
        - secretRef:
            name: "env-secrets"
      {{- end }}
      {{- if eq .Spec.ControlPlan.Conf.Fixpg "true" }}
      - name: fixpg
        image: {{.Spec.ControlPlan.Seeder.Image}}
        command: ["/bin/bash", "-c", "python3 cnvrg-boot.py fixpg"]
        envFrom:
        - configMapRef:
            name: "env-config"
        - secretRef:
            name: "env-secrets"
        imagePullPolicy: Always
      {{- end }}
      - name: seeder
        image: {{.Spec.ControlPlan.Seeder.Image}}
        command: ["/bin/bash", "-c", "python3 cnvrg-boot.py seeder --mode master"]
        imagePullPolicy: Always
        env:
        - name: "CNVRG_SEEDER_IMAGE"
          value: "{{.Spec.ControlPlan.Seeder.Image}}"
        - name: "CNVRG_SEED_CMD"
          value: "{{ .Spec.ControlPlan.Seeder.SeedCmd }}"
        - name: "CNVRG_NS"
          value: {{ .Spec.CnvrgNs }}
        - name: "CNVRG_SA_NAME"
          value: "{{.Spec.ControlPlan.Conf.Rbac.ServiceAccountName}}"
        {{- if eq .Spec.ControlPlan.Conf.CnvrgStorageType "gcp" }}
        - name: "CNVRG_GCP_KEYFILE_SECRET"
          value: "{{ .Spec.ControlPlan.Conf.GcpStorageSecret }}"
        - name: "CNVRG_GCP_KEYFILE_MOUNT_PATH"
          value: "{{ .Spec.ControlPlan.Conf.GcpKeyfileMountPath }}"
        {{- end }}


