apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.SvcName}}
  namespace: {{.Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
  labels:
    app: {{.SvcName}}
spec:
  replicas: {{ .Replicas }}
  selector:
    matchLabels:
      app: {{.SvcName}}
  template:
    metadata:
      labels:
        app: {{.SvcName}}
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - podAffinityTerm:
              labelSelector:
                matchLabels:
                  app: {{.SvcName}}
              namespaces:
              - {{.Namespace}}
              topologyKey: kubernetes.io/hostname
            weight: 1
      priorityClassName: {{ .AppClassRef }}
      serviceAccountName: cnvrg-{{.SvcName}}
      enableServiceLinks: false
      containers:
      - name: {{.SvcName}}
        imagePullPolicy: Always
        image: {{ image .ImageHub  .CentralUIImage }}
        env:
          - name: CNVRG_CENTRAL_SSO_BIND_ADDR
            value: "127.0.0.1:8000"
          - name: CNVRG_CENTRAL_SSO_DOMAIN_ID
            value: {{ .SsoDomainId }}
          - name: CNVRG_CENTRAL_SSO_SIGN_KEY
            value: "config/CNVRG_PKI_PRIVATE_KEY"
          - name: CNVRG_CENTRAL_SSO_JWT_IIS
            value: "{{ .JwksUrl }}"
        volumeMounts:
          - name: "private-key"
            mountPath: "/opt/app-root/config"
            readOnly: true
        resources:
          limits:
            cpu: 100m
            memory: 128Mi
          requests:
            cpu: 100m
            memory: 128Mi
      - name: oauth2-proxy
        image: {{  image .ImageHub  .OauthProxyImage }}
        command: [ "oauth2-proxy", "--config", "/opt/app-root/conf/proxy-config/conf" ]
        envFrom:
          - secretRef:
              name: {{ .RedisCredsRef }}
        volumeMounts:
          - name: "proxy-config"
            mountPath: "/opt/app-root/conf/proxy-config"
            readOnly: true
        resources:
          limits:
            cpu: {{ .Limits.Cpu }}
            memory: {{ .Limits.Memory }}
          requests:
            cpu: {{ .Requests.Cpu }}
            memory: {{ .Requests.Memory }}
      volumes:
      - name: "proxy-config"
        configMap:
         name: "proxy-config"
      - name: "private-key"
        secret:
          secretName: {{ .PrivateKeySecret }}