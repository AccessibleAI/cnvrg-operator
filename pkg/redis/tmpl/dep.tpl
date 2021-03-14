apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Redis.SvcName }}
  namespace: {{ .CnvrgNs }}
  labels:
    app: {{.Redis.SvcName }}
spec:
  selector:
    matchLabels:
      app: {{.Redis.SvcName }}
  template:
    metadata:
      labels:
        app: {{.Redis.SvcName }}
    spec:
      serviceAccountName: {{ .ControlPlan.Rbac.ServiceAccountName }}
      {{- if and (eq .Storage.Hostpath.Enabled "true") (eq .ControlPlan.Tenancy.Enabled "false") }}
      nodeSelector:
        kubernetes.io/hostname: "{{ .Storage.Hostpath.NodeName }}"
      {{- else if and (eq .Storage.Hostpath.Enabled "false") (eq .ControlPlan.Tenancy.Enabled "true") }}
      nodeSelector:
        {{ .ControlPlan.Tenancy.Key }}: "{{ .ControlPlan.Tenancy.Value }}"
      {{- else if and (eq .Storage.Hostpath.Enabled "true") (eq .ControlPlan.Tenancy.Enabled "true") }}
      nodeSelector:
        kubernetes.io/hostname: "{{ .Storage.Hostpath.NodeName }}"
        {{ .ControlPlan.Tenancy.Key }}: "{{ .ControlPlan.Tenancy.Value }}"
      {{- end }}
      tolerations:
        - key: {{ .ControlPlan.Tenancy.Key }}
          operator: Equal
          value: "{{ .ControlPlan.Tenancy.Value }}"
          effect: "NoSchedule"
      securityContext:
        runAsUser: 1000
        fsGroup: 1000
      containers:
        - image: {{ .Redis.Image }}
          name: redis
          command: [ "/bin/bash", "-lc", "redis-server /config/redis.conf" ]
          ports:
            - containerPort: {{ .Redis.Port }}
          resources:
            limits:
              cpu: {{ .Redis.Limits.CPU }}
              memory: {{ .Redis.Limits.Memory }}
            requests:
              cpu: {{ .Redis.Requests.CPU }}
              memory: {{ .Redis.Requests.Memory }}
          volumeMounts:
            - name: redis-data
              mountPath: /data
            - name: redis-config
              mountPath: /config
      volumes:
        - name: redis-data
          persistentVolumeClaim:
            claimName: {{ .Redis.SvcName }}
        - name: redis-config
          configMap:
            name: redis-conf
