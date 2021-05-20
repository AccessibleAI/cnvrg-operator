apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Spec.ControlPlane.WebApp.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.ControlPlane.WebApp.SvcName }}
    owner: cnvrg-control-plane
    cnvrg-component: webapp
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
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
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: {{.Spec.ControlPlane.WebApp.SvcName}}
        owner: cnvrg-control-plane
        cnvrg-component: webapp
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        {{ .Spec.Tenancy.Key }}: "{{ .Spec.Tenancy.Value }}"
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
        image: {{.Spec.ImageHub }}/{{ .Spec.SSO.Image }}
        command: [ "oauth2-proxy","--config", "/opt/app-root/conf/proxy-config/conf" ]
        envFrom:
          - secretRef:
              name: {{ .Spec.Dbs.Redis.CredsRef }}
        volumeMounts:
          - name: "oauth-proxy-webapp"
            mountPath: "/opt/app-root/conf/proxy-config"
            readOnly: true
      {{- end }}
      - image: {{.Spec.ImageHub }}/{{ .Spec.ControlPlane.Image }}
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
            cpu: "{{.Spec.ControlPlane.WebApp.Requests.Cpu}}"
            memory: "{{.Spec.ControlPlane.WebApp.Requests.Memory}}"
        {{- if eq .Spec.ControlPlane.ObjectStorage.Type "gcp" }}
        volumeMounts:
        - name: {{ .Spec.ControlPlane.ObjectStorage.GcpSecretRef }}
          mountPath: /opt/app-root/conf/gcp-keyfile
          readOnly: true
        {{- end }}
      volumes:
      {{- if isTrue .Spec.SSO.Enabled }}
      - name: "oauth-proxy-webapp"
        secret:
         secretName: "oauth-proxy-webapp"
      {{- end }}
      {{- if eq .Spec.ControlPlane.ObjectStorage.Type "gcp" }}
      - name: {{ .Spec.ControlPlane.ObjectStorage.GcpSecretRef }}
        secret:
          secretName: {{ .Spec.ControlPlane.ObjectStorage.GcpSecretRef }}
      {{- end }}
      initContainers:
      - name: services-check
        image: {{.Spec.ImageHub }}/{{.Spec.ControlPlane.Seeder.Image}}
        command: ["/bin/bash", "-c", "python3 cnvrg-boot.py services-check"]
        imagePullPolicy: Always
        env:
        - name: "CNVRG_SERVICE_LIST"
          {{- if and ( isTrue .Spec.Dbs.Minio.Enabled ) (eq .Spec.ControlPlane.ObjectStorage.Type "minio") }}
          value: "{{ .Spec.Dbs.Pg.SvcName }}:{{ .Spec.Dbs.Pg.Port }};{{ objectStorageUrl . }}/minio/health/ready"
          {{- else }}
          value: "{{ .Spec.Dbs.Pg.SvcName }}:{{ .Spec.Dbs.Pg.Port }}"
          {{ end }}
      {{- if and ( isTrue .Spec.Dbs.Minio.Enabled ) (eq .Spec.ControlPlane.ObjectStorage.Type "minio") }}
      - name: create-cnvrg-bucket
        image: {{.Spec.ImageHub }}/{{ .Spec.ControlPlane.Seeder.Image }}
        command: ["/bin/bash","-c", "{{ .Spec.ControlPlane.Seeder.CreateBucketCmd }}"]
        imagePullPolicy: Always
        envFrom:
        - secretRef:
            name: cp-object-storage
      {{- end }}
      - name: seeder
        image: {{.Spec.ImageHub }}/{{ .Spec.ControlPlane.Image }}
        command: ["/bin/bash", "-lc", "{{ .Spec.ControlPlane.Seeder.SeedCmd }}"]
        imagePullPolicy: Always
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


