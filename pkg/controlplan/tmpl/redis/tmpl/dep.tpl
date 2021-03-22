apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Spec.Redis.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: {{.Spec.Redis.SvcName }}
spec:
  selector:
    matchLabels:
      app: {{.Spec.Redis.SvcName }}
  template:
    metadata:
      labels:
        app: {{.Spec.Redis.SvcName }}
    spec:
      serviceAccountName: {{ .Spec.ControlPlan.Rbac.ServiceAccountName }}
      {{- if and (ne .Spec.ControlPlan.BaseConfig.HostpathNode "") (eq .Spec.ControlPlan.Tenancy.Enabled "false") }}
      nodeSelector:
        kubernetes.io/hostname: "{{ .Spec.ControlPlan.BaseConfig.HostpathNode }}"
      {{- else if and (eq .Spec.ControlPlan.BaseConfig.HostpathNode "") (eq .Spec.ControlPlan.Tenancy.Enabled "true") }}
      nodeSelector:
        {{ .Spec.ControlPlan.Tenancy.Key }}: "{{ .Spec.ControlPlan.Tenancy.Value }}"
      {{- else if and (ne .Spec.ControlPlan.BaseConfig.HostpathNode "") (eq .Spec.ControlPlan.Tenancy.Enabled "true") }}
      nodeSelector:
        kubernetes.io/hostname: "{{ .Spec.ControlPlan.BaseConfig.HostpathNode }}"
        {{ .Spec.ControlPlan.Tenancy.Key }}: "{{ .Spec.ControlPlan.Tenancy.Value }}"
      {{- end }}
      tolerations:
        - key: {{ .Spec.ControlPlan.Tenancy.Key }}
          operator: Equal
          value: "{{ .Spec.ControlPlan.Tenancy.Value }}"
          effect: "NoSchedule"
      securityContext:
        runAsUser: 1000
        fsGroup: 1000
      containers:
        - image: {{ .Spec.Redis.Image }}
          name: redis
          command: [ "/bin/bash", "-lc", "redis-server /config/redis.conf" ]
          ports:
            - containerPort: {{ .Spec.Redis.Port }}
          resources:
            limits:
              cpu: {{ .Spec.Redis.Limits.CPU }}
              memory: {{ .Spec.Redis.Limits.Memory }}
            requests:
              cpu: {{ .Spec.Redis.Requests.CPU }}
              memory: {{ .Spec.Redis.Requests.Memory }}
          volumeMounts:
            - name: redis-data
              mountPath: /data
            - name: redis-config
              mountPath: /config
      volumes:
        - name: redis-data
          persistentVolumeClaim:
            claimName: {{ .Spec.Redis.SvcName }}
        - name: redis-config
          configMap:
            name: redis-conf
