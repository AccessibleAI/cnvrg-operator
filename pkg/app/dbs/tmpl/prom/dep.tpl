apiVersion: apps/v1
kind: Deployment
metadata:
  name: prom
  namespace: {{ .Data.Namespace }}
spec:
  selector:
    matchLabels:
      app: prom
  template:
    metadata:
      labels:
        app: prom
    spec:
      serviceAccountName: prom
      containers:
      - name: prometheus
        image: prom/prometheus
        volumeMounts:
          - mountPath: /prometheus/config
            name: config
          - mountPath: /data
            name: prom-data
        command:
        - prometheus
        - --storage.tsdb.path=/data
        - --config.file=/prometheus/config/prometheus.yml
        - --web.config.file=/prometheus/config/web-config.yml
        ports:
          - containerPort: 9090
      volumes:
        - name: config
          configMap:
            name: prom-config
        - name: prom-data
          persistentVolumeClaim:
            claimName: prom