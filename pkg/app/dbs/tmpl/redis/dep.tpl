apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Spec.Dbs.Redis.SvcName }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{.Spec.Dbs.Redis.SvcName }}
    cnvrg-component: redis
    cnvrg-system-status-check: "true"
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  selector:
    matchLabels:
      app: {{.Spec.Dbs.Redis.SvcName }}
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: {{.Spec.Dbs.Redis.SvcName }}
        owner: cnvrg-control-plane
        cnvrg-component: redis
    spec:
      priorityClassName: {{ .Spec.PriorityClass.AppClassRef }}
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        {{ .Spec.Tenancy.Key }}: {{ .Spec.Tenancy.Value }}
        {{- range $key, $val := .Spec.Dbs.Redis.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      tolerations:
        - operator: "Exists"
      {{- else if (gt (len .Spec.Dbs.Redis.NodeSelector) 0) }}
      nodeSelector:
        {{- range $key, $val := .Spec.Dbs.Redis.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      {{- end }}
      serviceAccountName: {{ .Spec.Dbs.Redis.ServiceAccount }}
      securityContext:
        runAsUser: 1000
        fsGroup: 1000
      enableServiceLinks: false
      containers:
        - image: {{ image .Spec.ImageHub .Spec.Dbs.Redis.Image }}
          name: redis
          command: [ "/bin/bash", "-lc", "redis-server /config/redis.conf" ]
          envFrom:
            {{- if isTrue .Spec.Networking.Proxy.Enabled }}
            - configMapRef:
                name: {{ .Spec.Networking.Proxy.ConfigRef }}
            {{- end }}
          ports:
            - containerPort: {{ .Spec.Dbs.Redis.Port }}
          resources:
            limits:
              cpu: {{ .Spec.Dbs.Redis.Limits.Cpu }}
              memory: {{ .Spec.Dbs.Redis.Limits.Memory }}
            requests:
              cpu: {{ .Spec.Dbs.Redis.Requests.Cpu }}
              memory: {{ .Spec.Dbs.Redis.Requests.Memory }}
          volumeMounts:
            - name: redis-data
              mountPath: /data
            - name: redis-config
              mountPath: /config
      volumes:
        - name: redis-data
          persistentVolumeClaim:
            claimName: {{ .Spec.Dbs.Redis.PvcName }}
        - name: redis-config
          secret:
            secretName: {{ .Spec.Dbs.Redis.CredsRef }}
