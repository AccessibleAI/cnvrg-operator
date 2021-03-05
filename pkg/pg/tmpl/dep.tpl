apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Pg.SvcName}}
  namespace: {{ .CnvrgNs }}
  labels:
    app: {{.Pg.SvcName}}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{.Pg.SvcName}}
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: {{.Pg.SvcName}}
    spec:
      serviceAccountName: {{ .ControlPlan.Conf.Rbac.ServiceAccountName }}
      {{- if and (eq .ControlPlan.Conf.Tenancy.Enabled "true") (eq .ControlPlan.Conf.Tenancy.DedicatedNodes "true") }}
      tolerations:
        - key: {{ .ControlPlan.Conf.Tenancy.Key }}
          operator: Equal
          value: {{ .ControlPlan.Conf.Tenancy.Value }}
          effect: "NoSchedule"
      {{- end }}
      securityContext:
        runAsUser: {{ .Pg.RunAsUser }}
        fsGroup: {{ .Pg.FsGroup }}
      {{- if and (eq .Storage.Hostpath.Enabled "true") (eq .ControlPlan.Conf.Tenancy.Enabled "false") }}
      nodeSelector:
        kubernetes.io/hostname: "{{ .Storage.Hostpath.NodeName }}"
      {{- else if and (eq .Storage.Hostpath.Enabled "false") (eq .ControlPlan.Conf.Tenancy.Enabled "true") }}
      nodeSelector:
      {{ .ControlPlan.Conf.Tenancy.Key }}: "{{ .ControlPlan.Conf.Tenancy.Value }}"
      {{- else if and (eq .Storage.Hostpath.Enabled "true") (eq .ControlPlan.Conf.Tenancy.Enabled "true") }}
      nodeSelector:
        kubernetes.io/hostname: "{{ .Storage.Hostpath.NodeName }}"
        {{ .ControlPlan.Conf.Tenancy.Key }}: "{{ .ControlPlan.Conf.Tenancy.Value }}"
      {{- end }}
      containers:
        - name: postgresql
          envFrom:
            - secretRef:
                name: "pg-secret"
          image: {{.Pg.Image}}
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: {{.Pg.Port}}
              protocol: TCP
          livenessProbe:
            exec:
              command:
                - /usr/libexec/check-container
                - --live
            initialDelaySeconds: 120
            timeoutSeconds: 10
          readinessProbe:
            exec:
              command:
                - /usr/libexec/check-container
            initialDelaySeconds: 5
            timeoutSeconds: 1
          securityContext:
            capabilities: {}
            privileged: false
          terminationMessagePath: /dev/termination-log
          volumeMounts:
            - mountPath: /var/lib/pgsql/data
              name: postgres-data
            - mountPath: /dev/shm
              name: dshm
            {{- if eq .Pg.HugePages.Enabled "true" -}}
            - mountPath: "/hugepages"
              name: "hugepage"
            {{- end}}
          resources:
            {{- if eq .Pg.HugePages.Enabled "true" }}
            limits:
              {{- if eq .HugePages.memory ""}}
              hugepages-{{ .Pg.HugePages.Size }}: {{ .Pg.MemoryRequest }}
              {{- else }}
              hugepages-{{ .Pg.HugePages.Size }}: {{ .Pg.HugePages.Memory }}
              {{- end }}
            {{- end}}
            requests:
              cpu: {{ .Pg.CPURequest }}
              memory: {{ .Pg.MemoryRequest }}
      volumes:
        - name: postgres-data
          persistentVolumeClaim:
            claimName: {{.Pg.SvcName}}
        - name: dshm
          emptyDir:
            medium: Memory
            sizeLimit: 2Gi
        {{- if eq .Pg.HugePages.Enabled "true" }}
        - name: "hugepage"
          emptyDir:
            medium: HugePages
        {{- end}}
