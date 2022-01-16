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
      securityContext:
        runAsUser: 1000
        fsGroup: 1000
      containers:
        - image: {{ image .Spec.ImageHub .Spec.Dbs.Cvat.Redis.Image }}
          name: redis
          command: [ "/bin/bash", "-lc", "redis-server /config/redis.conf" ]
          ports:
            - containerPort: {{ .Spec.Dbs.Cvat.Redis.Port }}
          resources:
            limits:
              cpu: {{ .Spec.Dbs.Cvat.Redis.Limits.Cpu }}
              memory: {{ .Spec.Dbs.Cvat.Redis.Limits.Memory }}
            requests:
              cpu: {{ .Spec.Dbs.Cvat.Redis.Requests.Cpu }}
              memory: {{ .Spec.Dbs.Cvat.Redis.Requests.Memory }}
          volumeMounts:
            - name: redis-data
              mountPath: /data
            - name: redis-config
              mountPath: /config
      volumes:
        - name: redis-data
          persistentVolumeClaim:
            claimName: {{ .Spec.Dbs.Cvat.Redis.PvcName }}
        - name: redis-config
          secret:
            secretName: {{ .Spec.Dbs.Cvat.Redis.CredsRef }}