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
      serviceAccountName: {{ .Spec.SSO.Proxy.SvcName}}
      enableServiceLinks: false
      containers:
      - name: proxy-central
        imagePullPolicy: Always
        image: {{image .Spec.ImageHub .Spec.SSO.Proxy.Image}}
        command:
          - /opt/app-root/proxy
          - --authz-addr=127.0.0.1:50052
          - --ingress-type={{.Spec.Networking.Ingress.Type}}
        ports:
          - containerPort: 8888
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
          - --ingress-type={{.Spec.Networking.Ingress.Type}}
        resources:
          limits:
            cpu: {{ .Spec.SSO.Proxy.Limits.Cpu }}
            memory: {{ .Spec.SSO.Proxy.Limits.Memory }}
          requests:
            cpu: {{ .Spec.SSO.Proxy.Requests.Cpu }}
            memory: {{ .Spec.SSO.Proxy.Requests.Memory }}