apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Spec.SSO.Proxy.SvcName}}
  namespace: {{.Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
  labels:
    app: {{.Spec.SSO.Proxy.SvcName}}
spec:
  replicas: {{ .Spec.SSO.Proxy.Replicas }}
  selector:
    matchLabels:
      app: {{.Spec.SSO.Proxy.SvcName}}
  template:
    metadata:
      labels:
        app: {{.Spec.SSO.Proxy.SvcName}}
        {{- range $k, $v := .ObjectMeta.Annotations }}
        {{- if eq $k "eastwest_custom_name" }}
        sidecar.istio.io/inject: "true"
        {{- end }}
        {{- end }}
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - podAffinityTerm:
              labelSelector:
                matchLabels:
                  app: {{.Spec.SSO.Proxy.SvcName}}
              namespaces:
              - {{.Namespace}}
              topologyKey: kubernetes.io/hostname
            weight: 1
      priorityClassName: {{ .Spec.PriorityClass.AppClassRef }}
      securityContext:
        runAsUser: 1000
        runAsGroup: 1000
      serviceAccountName: {{ .Spec.SSO.Proxy.SvcName}}
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        "{{ .Spec.Tenancy.Key }}": "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - operator: "Exists"
      {{- else if (gt (len .Spec.SSO.Proxy.NodeSelector) 0) }}
      nodeSelector:
        {{- range $key, $val := .Spec.SSO.Proxy.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      {{- end }}
      enableServiceLinks: false
      containers:
      - name: proxy-central
        imagePullPolicy: Always
        image: {{image .Spec.ImageHub .Spec.SSO.Proxy.Image}}
        envFrom:
          {{- if isTrue .Spec.Networking.Proxy.Enabled }}
          - configMapRef:
              name: {{ .Spec.Networking.Proxy.ConfigRef }}
          {{- end }}
        command:
          - /opt/app-root/proxy
          - start
          - --authz-addr=127.0.0.1:50052
          - --ingress-type={{.Spec.Networking.Ingress.Type}}
        ports:
          - containerPort: 8888
        {{- if isTrue .Spec.SSO.Proxy.Readiness }}
        readinessProbe:
          httpGet:
            path: /ready
            port: 2112
          initialDelaySeconds: 10
          periodSeconds: 20
        {{- end }}
        resources:
          limits:
            cpu: {{ .Spec.SSO.Proxy.Limits.Cpu }}
            memory: {{ .Spec.SSO.Proxy.Limits.Memory }}
          requests:
            cpu: {{ .Spec.SSO.Proxy.Requests.Cpu }}
            memory: {{ .Spec.SSO.Proxy.Requests.Memory }}
      - name: authz
        imagePullPolicy: Always
        image: {{  image .Spec.ImageHub .Spec.SSO.Proxy.Image }}
        command:
          - /opt/app-root/authz
          - start
          - --ingress-type={{.Spec.Networking.Ingress.Type}}
        resources:
          limits:
            cpu: {{ .Spec.SSO.Proxy.Limits.Cpu }}
            memory: {{ .Spec.SSO.Proxy.Limits.Memory }}
          requests:
            cpu: {{ .Spec.SSO.Proxy.Requests.Cpu }}
            memory: {{ .Spec.SSO.Proxy.Requests.Memory }}