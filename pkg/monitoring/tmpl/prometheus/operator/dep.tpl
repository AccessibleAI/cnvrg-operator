apiVersion: apps/v1
kind: Deployment
metadata:
  name: cnvrg-prometheus-operator
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: cnvrg-prometheus-operator
    version: v0.44.1
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cnvrg-prometheus-operator
      version: v0.44.1
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: cnvrg-prometheus-operator
        version: v0.44.1
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
      serviceAccountName: cnvrg-prometheus-operator
      securityContext:
        runAsNonRoot: true
        runAsUser: 65534
      containers:
      - args:
        - --kubelet-service=kube-system/kubelet
        - --prometheus-config-reloader={{ image .Spec.ImageHub .Spec.Monitoring.PrometheusOperator.PrometheusConfigReloaderImage }}
        image: {{image .Spec.ImageHub .Spec.Monitoring.PrometheusOperator.OperatorImage }}
        name: prometheus-operator
        ports:
        - containerPort: 8080
          name: http
        resources:
          limits:
            cpu: 1000m
            memory: 2048Mi
          requests:
            cpu: 100m
            memory: 100Mi
        securityContext:
          allowPrivilegeEscalation: false
      - args:
        - --logtostderr
        - --secure-listen-address=:8443
        - --tls-cipher-suites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305
        - --upstream=http://127.0.0.1:8080/
        image: {{image .Spec.ImageHub .Spec.Monitoring.PrometheusOperator.KubeRbacProxyImage }}
        name: kube-rbac-proxy
        ports:
        - containerPort: 8443
          name: https
        securityContext:
          runAsGroup: 65532
          runAsNonRoot: true
          runAsUser: 65532


