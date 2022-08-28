apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Spec.Dbs.Cvat.Redis.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{.Spec.Dbs.Cvat.Redis.SvcName }}
    cnvrg-component: redis
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  selector:
    matchLabels:
      app: {{.Spec.Dbs.Cvat.Redis.SvcName }}
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: {{.Spec.Dbs.Cvat.Redis.SvcName }}
        owner: cnvrg-control-plane
        cnvrg-component: redis
    spec:
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      {{- if (gt (len .Spec.Dbs.Cvat.Redis.NodeSelector) 0) }}
      nodeSelector:
        {{- range $key, $val := .Spec.Dbs.Cvat.Redis.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      {{- end }}
      serviceAccountName: {{ .Spec.Dbs.Cvat.Redis.ServiceAccount }}
      containers:
        - image: {{ .Spec.Dbs.Cvat.Redis.Image }}
          name: redis
          ports:
            - containerPort: {{ .Spec.Dbs.Cvat.Redis.Port }}
          resources:
            limits:
              cpu: {{ .Spec.Dbs.Cvat.Redis.Limits.Cpu }}
              memory: {{ .Spec.Dbs.Cvat.Redis.Limits.Memory }}
            requests:
              cpu: {{ .Spec.Dbs.Cvat.Redis.Requests.Cpu }}
              memory: {{ .Spec.Dbs.Cvat.Redis.Requests.Memory }}
