apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Spec.Dbs.Redis.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: {{.Spec.Dbs.Redis.SvcName }}
spec:
  selector:
    matchLabels:
      app: {{.Spec.Dbs.Redis.SvcName }}
  template:
    metadata:
      labels:
        app: {{.Spec.Dbs.Redis.SvcName }}
    spec:
      serviceAccountName: {{ .Spec.Dbs.Redis.ServiceAccount }}
      securityContext:
        runAsUser: 1000
        fsGroup: 1000
      containers:
        - image: {{ .Spec.Dbs.Redis.Image }}
          name: redis
          command: [ "/bin/bash", "-lc", "redis-server /config/redis.conf" ]
          ports:
            - containerPort: {{ .Spec.Dbs.Redis.Port }}
          resources:
            limits:
              cpu: {{ .Spec.Dbs.Redis.Limits.CPU }}
              memory: {{ .Spec.Dbs.Redis.Limits.Memory }}
            requests:
              cpu: {{ .Spec.Dbs.Redis.Requests.CPU }}
              memory: {{ .Spec.Dbs.Redis.Requests.Memory }}
          volumeMounts:
            - name: redis-data
              mountPath: /data
            - name: redis-config
              mountPath: /config
      volumes:
        - name: redis-data
          persistentVolumeClaim:
            claimName: {{ .Spec.Dbs.Redis.SvcName }}
        - name: redis-config
          configMap:
            name: redis-conf
