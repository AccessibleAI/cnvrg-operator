
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
      securityContext:
        runAsUser: {{ .Spec.Logging.Elastalert.RunAsUser }}
        fsGroup: {{ .Spec.Logging.Elastalert.FsGroup }}
      serviceAccountName: {{ .Spec.ControlPlane.Rbac.ServiceAccountName }}
      {{- if and (ne .Spec.ControlPlane.BaseConfig.HostpathNode "") (not .Spec.ControlPlane.Tenancy.Enabled) }}
      nodeSelector:
        kubernetes.io/hostname: "{{ .Spec.ControlPlane.BaseConfig.HostpathNode }}"
      {{- else if and (eq .Spec.ControlPlane.BaseConfig.HostpathNode "") (.Spec.ControlPlane.Tenancy.Enabled ) }}
      nodeSelector:
        {{ .Spec.ControlPlane.Tenancy.Key }}: "{{ .Spec.ControlPlane.Tenancy.Value }}"
      {{- else if and (ne .Spec.ControlPlane.BaseConfig.HostpathNode "") (.Spec.ControlPlane.Tenancy.Enabled ) }}
      nodeSelector:
        kubernetes.io/hostname: "{{ .Spec.ControlPlane.BaseConfig.HostpathNode }}"
        {{ .Spec.ControlPlane.Tenancy.Key }}: "{{ .Spec.ControlPlane.Tenancy.Value }}"
      {{- end }}
      tolerations:
        - key: {{ .Spec.ControlPlane.Tenancy.Key }}
          operator: Equal
          value: "{{ .Spec.ControlPlane.Tenancy.Value }}"
          effect: "NoSchedule"
      containers:
      - image: {{ .Spec.Logging.Elastalert.Image }}
        name: {{ .Spec.Logging.Elastalert.SvcName }}
        ports:
        - containerPort: {{ .Spec.Logging.Elastalert.ContainerPort }}
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