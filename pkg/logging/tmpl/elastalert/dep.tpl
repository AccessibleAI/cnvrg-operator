
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: {{ .Spec.Logging.Elastalert.SvcName }}
  name: {{ .Spec.Logging.Elastalert.SvcName }}
  namespace: {{ ns . }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ .Spec.Logging.Elastalert.SvcName }}
  template:
    metadata:
      labels:
        app: {{ .Spec.Logging.Elastalert.SvcName }}
    spec:
      serviceAccountName: {{ .Spec.Logging.Elastalert.SvcName }}
      securityContext:
        runAsUser: 1000
        fsGroup: 1000
      containers:
      - image: {{ .Spec.Logging.Elastalert.Image }}
        name: {{ .Spec.Logging.Elastalert.SvcName }}
        ports:
        - containerPort: 3030
          protocol: TCP
        envFrom:
          - secretRef:
              name: {{ .Spec.Dbs.Es.CredsRef }}
        resources:
          requests:
            cpu: {{.Spec.Logging.Elastalert.CPURequest}}
            memory: {{.Spec.Logging.Elastalert.MemoryRequest}}
          limits:
            cpu: {{ .Spec.Logging.Elastalert.CPULimit }}
            memory: {{ .Spec.Logging.Elastalert.MemoryLimit }}
        volumeMounts:
        - mountPath: /opt/elastalert-server/config/config.json
          subPath: config.json
          name: elastalert-config
        - mountPath: /opt/elastalert/config.yaml
          subPath: config.yaml
          name: elastalert-config
        - mountPath: /opt/elastalert/rules
          name: {{ .Spec.Logging.Elastalert.SvcName }}
      restartPolicy: Always
      volumes:
      - name: {{ .Spec.Logging.Elastalert.SvcName }}
        persistentVolumeClaim:
          claimName: {{ .Spec.Logging.Elastalert.SvcName }}
      - configMap:
          name: elastalert-config
        name: elastalert-config