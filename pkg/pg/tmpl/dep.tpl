apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Spec.Pg.SvcName}}
  namespace: {{ .ObjectMeta.Namespace }}
  labels:
    app: {{.Spec.Pg.SvcName}}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{.Spec.Pg.SvcName}}
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: {{.Spec.Pg.SvcName}}
    spec:
      serviceAccountName: {{ .Spec.ControlPlan.Conf.Rbac.ServiceAccountName }}
      {{- if and (eq .Spec.Tenancy.Enabled "true") (eq .Spec.Tenancy.DedicatedNodes "true") }}
      tolerations:
        - key: {{ .Spec.Tenancy.Cnvrg.Key }}
          operator: Equal
          value: {{ .Spec.Tenancy.Cnvrg.Value }}
          effect: "NoSchedule"
      {{- end }}
      securityContext:
        runAsUser: {{ .Spec.Pg.RunAsUser }}
        fsGroup: {{ .Spec.Pg.FsGroup }}
      {{- if and (eq .Spec.Storage.Hostpath.Enabled "true") (eq .Spec.Tenancy.Enabled "false") }}
      nodeSelector:
        kubernetes.io/hostname: "{{ .Spec.Storage.Hostpath.NodeName }}"
      {{- else if and (eq .Spec.Storage.Hostpath.Enabled "false") (eq .Spec.Tenancy.Enabled "true") }}
      nodeSelector:
      {{ .Spec.Tenancy.Cnvrg.Key }}: "{{ .Spec.Tenancy.Cnvrg.Value }}"
      {{- else if and (eq .Spec.Storage.Hostpath.Enabled "true") (eq .Spec.Tenancy.Enabled "true") }}
      nodeSelector:
        kubernetes.io/hostname: "{{ .Spec.Storage.Hostpath.NodeName }}"
        {{ .Spec.Tenancy.Cnvrg.Key }}: "{{ .Spec.Tenancy.Cnvrg.Value }}"
      {{- end }}
      containers:
        - name: postgresql
          envFrom:
            - secretRef:
                name: "pg-secret"
          image: {{.Spec.Pg.Image}}
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: {{.Spec.Pg.Port}}
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
            {{- if eq .Spec.Pg.HugePages.Enabled "true" -}}
            - mountPath: "/hugepages"
              name: "hugepage"
            {{- end}}
          resources:
            {{- if eq .Spec.Pg.HugePages.Enabled "true" }}
            limits:
              {{- if eq .Spec.HugePages.memory ""}}
              hugepages-{{ .Spec.Pg.HugePages.Size }}: {{ .Spec.Pg.MemoryRequest }}
              {{- else }}
              hugepages-{{ .Spec.Pg.HugePages.Size }}: {{ .Spec.Pg.HugePages.Memory }}
              {{- end }}
            {{- end}}
            requests:
              cpu: {{ .Spec.Pg.CPURequest }}
              memory: {{ .Spec.Pg.MemoryRequest }}
      volumes:
        - name: postgres-data
          persistentVolumeClaim:
            claimName: {{.Spec.Pg.SvcName}}
        - name: dshm
          emptyDir:
            medium: Memory
            sizeLimit: 2Gi
        {{- if eq .Spec.Pg.HugePages.Enabled "true" }}
        - name: "hugepage"
          emptyDir:
            medium: HugePages
        {{- end}}
