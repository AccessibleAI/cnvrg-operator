apiVersion: apps/v1
kind: Deployment
metadata:
  name:  {{ .Spec.Dbs.Prom.SvcName }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
spec:
  strategy:
    type: Recreate
  selector:
    matchLabels:
      app:  {{ .Spec.Dbs.Prom.SvcName }}
  template:
    metadata:
      labels:
        app: {{ .Spec.Dbs.Prom.SvcName }}
        {{- range $k, $v := .ObjectMeta.Annotations }}
        {{- if eq $k "eastwest_custom_name" }}
        sidecar.istio.io/inject: "true"
        {{- end }}
        {{- end }}
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - podAffinityTerm:
              labelSelector:
                matchLabels:
                  app: {{.Spec.Dbs.Prom.SvcName}}
              namespaces:
              - {{.Namespace}}
              topologyKey: kubernetes.io/hostname
            weight: 1
      priorityClassName: {{ .Spec.PriorityClass.AppClassRef }}
      serviceAccountName: cnvrg-prom
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        "{{ .Spec.Tenancy.Key }}": "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - operator: "Exists"
      {{- end }}
      securityContext:
        runAsUser: 1000
        fsGroup: 1000
      enableServiceLinks: false
      containers:
      - name: prometheus
        image: {{ image .Spec.ImageHub .Spec.Dbs.Prom.Image }}
        envFrom:
          {{- if isTrue .Spec.Networking.Proxy.Enabled }}
          - configMapRef:
              name: {{ .Spec.Networking.Proxy.ConfigRef }}
          {{- end }}
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
          - containerPort: {{ .Spec.Dbs.Prom.Port }}
      volumes:
        - name: prom-scrape-configs
          configMap:
            name: prom-scrape-configs
        - name: prom-web-configs
          configMap:
            name: prom-web-configs
        - name: prom-data
          persistentVolumeClaim:
            claimName: {{ .Spec.Dbs.Prom.SvcName }}