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
      priorityClassName: {{ .Spec.PriorityClass.AppClassRef }}
      serviceAccountName: cnvrg-prom
      securityContext:
        runAsUser: 1000
        fsGroup: 1000
      containers:
      - name: prometheus
        image: {{ image .Spec.ImageHub .Spec.Dbs.Prom.Image }}
        volumeMounts:
          - mountPath: /prometheus/config/scrape
            name: prom-scrape-configs
          - mountPath: /prometheus/config/web
            name: prom-web-configs
          - mountPath: /data
            name: prom-data
        command:
        - prometheus
        - --storage.tsdb.path=/data
        - --config.file=/prometheus/config/scrape/prometheus.yml
        - --web.config.file=/prometheus/config/web/web-config.yml
        ports:
          - containerPort: 9090
      volumes:
        - name: prom-scrape-configs
          configMap:
            name: prom-scrape-configs
        - name: prom-web-configs
          configMap:
            name: prom-web-configs
        - name: prom-data
          persistentVolumeClaim:
            claimName: prom