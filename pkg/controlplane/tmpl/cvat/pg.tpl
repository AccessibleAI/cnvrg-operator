apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Spec.ControlPlane.Cvat.Pg.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{.Spec.ControlPlane.Cvat.Pg.SvcName }}
    cnvrg-component: cvat
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{.Spec.ControlPlane.Cvat.Pg.SvcName }}
  strategy:
    type: Recreate
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: {{.Spec.ControlPlane.Cvat.Pg.SvcName }}
        cnvrg-component: cvat
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        {{ .Spec.Tenancy.Key }}: {{ .Spec.Tenancy.Value }}
        {{- range $key, $val := .Spec.ControlPlane.Cvat.Pg.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      tolerations:
        - key: "{{ .Spec.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.Tenancy.Value }}"
          effect: "NoSchedule"
      {{- else if (gt (len .Spec.ControlPlane.Cvat.Pg.NodeSelector) 0) }}
      nodeSelector:
        {{- range $key, $val := .Spec.ControlPlane.Cvat.Pg.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      {{- end }}
      serviceAccountName: {{ .Spec.ControlPlane.Cvat.Pg.ServiceAccount }}
      securityContext:
        runAsUser: 26
        fsGroup: 26
      containers:
        - name: cvat-postgres
          env:
            - name: POSTGRES_DB
              valueFrom:
                configMapKeyRef:
                  name: cvat-pg-config
                  key: CNVRG_CVAT_POSTGRES_DBNAME
            - name: POSTGRES_USER
              valueFrom:
                configMapKeyRef:
                  name: cvat-pg-config
                  key: CNVRG_CVAT_POSTGRES_USER
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: {{ .Spec.ControlPlane.Cvat.Pg.CredsRef }}
                  key: CNVRG_CVAT_POSTGRES_PASSWORD
          image: {{ .Spec.ControlPlane.Cvat.Pg.Image }}
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: {{.Spec.ControlPlane.Cvat.Pg.Port}}
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
                - sh
                - -c
                - su - postgres -c "pg_isready --host=$POD_IP"
            initialDelaySeconds: 5
            timeoutSeconds: 1
          securityContext:
            capabilities: {}
            privileged: false
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          volumeMounts:
            - mountPath: /var/lib/pgsql/data
              name: postgres-data
            - mountPath: /dev/shm
              name: dshm
            {{- if isTrue .Spec.ControlPlane.Cvat.Pg.HugePages.Enabled }}
            - mountPath: "/hugepages"
              name: "hugepage"
            {{- end }}
          resources:
            limits:
              cpu: {{ .Spec.ControlPlane.Cvat.Pg.Limits.Cpu }}
              memory: {{ .Spec.ControlPlane.Cvat.Pg.Limits.Memory }}
              {{- if isTrue .Spec.ControlPlane.Cvat.Pg.HugePages.Enabled }}
              {{- if eq .Spec.ControlPlane.Cvat.Pg.HugePages.Memory ""}}
              hugepages-{{ .Spec.ControlPlane.Cvat.Pg.HugePages.Size }}: {{ .Spec.ControlPlane.Cvat.Pg.Requests.Memory }}
              {{- else }}
              hugepages-{{ .Spec.ControlPlane.Cvat.Pg.HugePages.Size }}: {{ .Spec.ControlPlane.Cvat.Pg.HugePages.Memory }}
              {{- end }}
              {{- end }}
            requests:
              cpu: {{ .Spec.ControlPlane.Cvat.Pg.Requests.Cpu }}
              memory: {{ .Spec.ControlPlane.Cvat.Pg.Requests.Memory }}
      volumes:
        - name: postgres-data
          persistentVolumeClaim:
            claimName: {{.Spec.ControlPlane.Cvat.Pg.PvcName}}
        - name: dshm
          emptyDir:
            medium: Memory
            sizeLimit: 2Gi
        {{- if isTrue .Spec.ControlPlane.Cvat.Pg.HugePages.Enabled }}
        - name: "hugepage"
          emptyDir:
            medium: HugePages
        {{- end}}
