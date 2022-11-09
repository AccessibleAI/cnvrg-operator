apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Spec.SSO.Central.SvcName}}
  namespace: {{.Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
  labels:
    app: {{.Spec.SSO.Central.SvcName}}
spec:
  replicas: {{.Spec.SSO.Central.Replicas}}
  selector:
    matchLabels:
      app: {{.Spec.SSO.Central.SvcName}}
  template:
    metadata:
      labels:
        app: {{.Spec.SSO.Central.SvcName}}
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - podAffinityTerm:
              labelSelector:
                matchLabels:
                  app: {{.Spec.SSO.Central.SvcName}}
              namespaces:
              - {{.Namespace}}
              topologyKey: kubernetes.io/hostname
            weight: 1
      priorityClassName: {{ .Spec.PriorityClass.AppClassRef }}
      serviceAccountName: {{.Spec.SSO.Central.SvcName}}-central
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        "{{ .Spec.Tenancy.Key }}": "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - operator: "Exists"
      {{- else if (gt (len .Spec.SSO.Central.NodeSelector) 0) }}
      nodeSelector:
        {{- range $key, $val := .Spec.SSO.Central.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      {{- end }}
      enableServiceLinks: false
      containers:
      - name: {{.Spec.SSO.Central.SvcName}}
        imagePullPolicy: Always
        image: {{  image .Spec.ImageHub .Spec.SSO.Central.CentralUiImage }}
        env:
          - name: CNVRG_CENTRAL_SSO_BIND_ADDR
            value: "0.0.0.0:8000"
          - name: CNVRG_CENTRAL_SSO_DOMAIN_ID
            value: {{ .SsoDomainId }}
          - name: CNVRG_CENTRAL_SSO_SIGN_KEY
            value: "config/CNVRG_PKI_PRIVATE_KEY"
          - name: CNVRG_CENTRAL_SSO_JWT_IIS
            value: "{{ .Spec.SSO.Central.JwksURL }}"
        volumeMounts:
          - name: "private-key"
            mountPath: "/opt/app-root/config"
            readOnly: true
        {{- if isTrue .Spec.SSO.Central.Readiness }}
        readinessProbe:
          httpGet:
            path: /ready
            port: 8000
          initialDelaySeconds: 5
          periodSeconds: 20
        {{- end }}
        resources:
          limits:
            cpu: 100m
            memory: 128Mi
          requests:
            cpu: 100m
            memory: 128Mi
      - name: oauth2-proxy
        image: {{  image .Spec.ImageHub .Spec.SSO.Central.OauthProxyImage }}
        command: [ "oauth2-proxy", "--config", "/opt/app-root/conf/proxy-config/conf" ]
        envFrom:
          - secretRef:
              name: {{ .Spec.Dbs.Redis.CredsRef }}
        volumeMounts:
          - name: "proxy-config"
            mountPath: "/opt/app-root/conf/proxy-config"
            readOnly: true
        resources:
          limits:
            cpu: {{ .Spec.SSO.Central.Limits.Cpu }}
            memory: {{ .Spec.SSO.Central.Limits.Memory }}
          requests:
            cpu: {{ .Spec.SSO.Central.Requests.Cpu }}
            memory: {{ .Spec.SSO.Central.Requests.Memory }}
      volumes:
      - name: "proxy-config"
        configMap:
         name: "proxy-config"
      - name: "private-key"
        secret:
          secretName: {{ .Spec.SSO.Pki.PrivateKeySecret }}