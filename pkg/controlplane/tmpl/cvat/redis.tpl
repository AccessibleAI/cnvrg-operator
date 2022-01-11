apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Spec.ControlPlane.Cvat.Redis.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{.Spec.ControlPlane.Cvat.Redis.SvcName }}
    cnvrg-component: redis
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  selector:
    matchLabels:
      app: {{.Spec.ControlPlane.Cvat.Redis.SvcName }}
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
          {{$k}}: "{{$v}}"
          {{- end }}
      labels:
        app: {{.Spec.ControlPlane.Cvat.Redis.SvcName }}
        owner: cnvrg-control-plane
        cnvrg-component: redis
    spec:
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        {{ .Spec.Tenancy.Key }}: {{ .Spec.Tenancy.Value }}
        {{- range $key, $val := .Spec.ControlPlane.Cvat.Redis.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      tolerations:
        - key: "{{ .Spec.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.Tenancy.Value }}"
          effect: "NoSchedule"
      {{- else if (gt (len .Spec.ControlPlane.Cvat.Redis.NodeSelector) 0) }}
      nodeSelector:
        {{- range $key, $val := .Spec.ControlPlane.Cvat.Redis.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      {{- end }}
      serviceAccountName: {{ .Spec.ControlPlane.Cvat.Redis.ServiceAccount }}
      containers:
        - image: {{ .Spec.ControlPlane.Cvat.Redis.Image }}
          name: cvat-redis
          ports:
            - containerPort: {{ .Spec.ControlPlane.Cvat.Redis.Port }}
          resources:
            limits:
              cpu: {{ .Spec.ControlPlane.Cvat.Redis.Limits.Cpu }}
              memory: {{ .Spec.ControlPlane.Cvat.Redis.Limits.Memory }}
            requests:
              cpu: {{ .Spec.ControlPlane.Cvat.Redis.Requests.Cpu }}
              memory: {{ .Spec.ControlPlane.Cvat.Redis.Requests.Memory }}
