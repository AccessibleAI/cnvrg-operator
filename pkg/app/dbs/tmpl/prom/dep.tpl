apiVersion: apps/v1
kind: Deployment
metadata:
  name: prom
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
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
      securityContext:
        runAsUser: 1000
        fsGroup: 1000
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