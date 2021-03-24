
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
      containers:
      - image: {{ .Spec.Logging.Elastalert.Image }}
        name: {{ .Spec.Logging.Elastalert.SvcName }}
        ports:
        - containerPort: {{ .Spec.Logging.Elastalert.ContainerPort }}
          protocol: TCP
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