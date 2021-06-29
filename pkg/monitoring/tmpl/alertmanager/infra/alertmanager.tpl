apiVersion: monitoring.coreos.com/v1
kind: Alertmanager
metadata:
  name: cnvrg-infra-alertmanager
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: cnvrg-infra-alertmanager
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  storage:
    disableMountSubPath: true
    volumeClaimTemplate:
      spec:
        resources:
          requests:
            storage: {{ .Spec.Monitoring.Alertmanager.StorageSize }}
        {{- if ne .Spec.Monitoring.Alertmanager.StorageClass "" }}
        storageClassName: {{ .Spec.Monitoring.Alertmanager.StorageClass }}
        {{- end }}
  image: {{ image .Spec.ImageHub .Spec.Monitoring.Alertmanager.Image }}
  replicas: 1
  retention: 240h # 10 days
  resources:
    requests:
      cpu: {{ .Spec.Monitoring.Alertmanager.Requests.Cpu }}
      memory: {{ .Spec.Monitoring.Alertmanager.Requests.Memory }}
    limits:
      cpu: {{ .Spec.Monitoring.Alertmanager.Limits.Cpu }}
      memory: {{ .Spec.Monitoring.Alertmanager.Limits.Memory }}
  podMetadata:
    {{- if .Spec.Annotations }}
    annotations:
      {{- range $k, $v := .Spec.Annotations }}
      {{$k}}: "{{$v}}"
      {{- end }}
    {{- end }}
    {{- if .Spec.Labels}}
    labels:
      {{- range $k, $v := .Spec.Labels }}
      {{$k}}: "{{$v}}"
      {{- end }}
    {{- end }}
  securityContext:
    fsGroup: 2000
    runAsNonRoot: true
    runAsUser: 1000
  serviceAccountName: cnvrg-infra-prometheus
  version: 0.22.2
  listenLocal: true
  containers:
  - name: "auth-proxy"
    image: {{ image .Spec.ImageHub .Spec.Monitoring.Alertmanager.BasicAuthProxyImage }}
    ports:
    - containerPort: 9091
      name: web
    volumeMounts:
    - name: "alertmanager-auth-proxy-config"
      mountPath: "/etc/nginx"
      readOnly: true
    - name: "htpasswd"
      mountPath: "/etc/nginx/htpasswd"
      readOnly: true
  volumes:
  - name: "alertmanager-auth-proxy-config"
    configMap:
      name: "alertmanager-auth-proxy-config"
  - name: "htpasswd"
    secret:
      secretName: {{ .Spec.Monitoring.Alertmanager.CredsRef }}
  {{- if isTrue .Spec.Tenancy.Enabled }}
  nodeSelector:
    {{ .Spec.Tenancy.Key }}: {{ .Spec.Tenancy.Value }}
    {{- range $key, $val := .Spec.Monitoring.Alertmanager.NodeSelector }}
    {{ $key }}: {{ $val }}
    {{- end }}
  tolerations:
    - key: "{{ .Spec.Tenancy.Key }}"
      operator: "Equal"
      value: "{{ .Spec.Tenancy.Value }}"
      effect: "NoSchedule"
  {{- else if (gt (len .Spec.Monitoring.Alertmanager.NodeSelector) 0) }}
  nodeSelector:
    {{- range $key, $val := .Spec.Monitoring.Alertmanager.NodeSelector }}
    {{ $key }}: {{ $val }}
    {{- end }}
  {{- end }}
