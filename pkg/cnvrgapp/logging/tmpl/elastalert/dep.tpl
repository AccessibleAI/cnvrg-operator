
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: {{ .Logging.Elastalert.SvcName }}
  name: {{ .Logging.Elastalert.SvcName }}
  namespace: {{ .CnvrgNs }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ .Logging.Elastalert.SvcName }}
  template:
    metadata:
      labels:
        app: {{ .Logging.Elastalert.SvcName }}
    spec:
      securityContext:
        runAsUser: {{ .Logging.Elastalert.RunAsUser }}
        fsGroup: {{ .Logging.Elastalert.FsGroup }}
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
      containers:
      - image: {{ .Logging.Elastalert.Image }}
        name: {{ .Logging.Elastalert.SvcName }}
        ports:
        - containerPort: {{ .Logging.Elastalert.ContainerPort }}
          protocol: TCP
        resources:
          requests:
            cpu: {{.Logging.Elastalert.CPURequest}}
            memory: {{.Logging.Elastalert.MemoryRequest}}
          limits:
            cpu: {{ .Logging.Elastalert.CPULimit }}
            memory: {{ .Logging.Elastalert.MemoryLimit }}
        volumeMounts:
        - mountPath: /opt/elastalert-server/config/config.json
          subPath: config.json
          name: elastalert-config
        - mountPath: /opt/elastalert/config.yaml
          subPath: config.yaml
          name: elastalert-config
        - mountPath: /opt/elastalert/rules
          name: {{ .Logging.Elastalert.SvcName }}
      restartPolicy: Always
      volumes:
      - name: {{ .Logging.Elastalert.SvcName }}
        persistentVolumeClaim:
          claimName: {{ .Logging.Elastalert.SvcName }}
      - configMap:
          name: elastalert-config
        name: elastalert-config