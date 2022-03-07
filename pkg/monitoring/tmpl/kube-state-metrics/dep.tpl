apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app.kubernetes.io/name: kube-state-metrics
    app.kubernetes.io/version: v1.9.7
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
  name: kube-state-metrics
  namespace: {{ ns . }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: kube-state-metrics
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app.kubernetes.io/name: kube-state-metrics
        app.kubernetes.io/version: v1.9.7
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        "{{ .Spec.Tenancy.Key }}": "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - operator: "Exists"
      {{- end }}
      containers:
      - args:
        - --host=127.0.0.1
        - --port=8081
        - --telemetry-host=127.0.0.1
        - --telemetry-port=8082
        image: {{ image .Spec.ImageHub .Spec.Monitoring.KubeStateMetrics.Image }}
        name: kube-state-metrics
        resources:
          requests:
            cpu: 200m
            memory: 200m
          limits:
            cpu: 1000m
            memory: 1Gi
      - args:
        - --logtostderr
        - --secure-listen-address=:8443
        - --tls-cipher-suites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305
        - --upstream=http://127.0.0.1:8081/
        image: {{ image .Spec.ImageHub .Spec.Monitoring.PrometheusOperator.KubeRbacProxyImage }}
        name: kube-rbac-proxy-main
        resources:
          requests:
            cpu: 200m
            memory: 200m
          limits:
            cpu: 1000m
            memory: 1Gi
        ports:
        - containerPort: 8443
          name: https-main
        securityContext:
          runAsGroup: 65532
          runAsNonRoot: true
          runAsUser: 65532
      - args:
        - --logtostderr
        - --secure-listen-address=:9443
        - --tls-cipher-suites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305
        - --upstream=http://127.0.0.1:8082/
        image: {{ image .Spec.ImageHub .Spec.Monitoring.PrometheusOperator.KubeRbacProxyImage }}
        name: kube-rbac-proxy-self
        resources:
          requests:
            cpu: 200m
            memory: 200Mi
          limits:
            cpu: 1000m
            memory: 1Gi
        ports:
        - containerPort: 9443
          name: https-self
        securityContext:
          runAsGroup: 65532
          runAsNonRoot: true
          runAsUser: 65532
      nodeSelector:
        kubernetes.io/os: linux
      serviceAccountName: kube-state-metrics
