apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Spec.Dbs.Pg.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{.Spec.Dbs.Pg.SvcName }}
    cnvrg-component: pg
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{.Spec.Dbs.Pg.SvcName }}
  strategy:
    type: Recreate
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: {{.Spec.Dbs.Pg.SvcName }}
        cnvrg-component: pg
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        {{ .Spec.Tenancy.Key }}: {{ .Spec.Tenancy.Value }}
        {{- range $key, $val := .Spec.Dbs.Pg.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      tolerations:
        - operator: "Exists"
      {{- else if (gt (len .Spec.Dbs.Pg.NodeSelector) 0) }}
      nodeSelector:
        {{- range $key, $val := .Spec.Dbs.Pg.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      {{- end }}
      serviceAccountName: {{ .Spec.Dbs.Pg.ServiceAccount }}
      securityContext:
        runAsUser: 26
        fsGroup: 26
      containers:
        - name: postgresql
          envFrom:
            - secretRef:
                name: {{ .Spec.Dbs.Pg.CredsRef }}
          image: {{ image .Spec.ImageHub .Spec.Dbs.Pg.Image}}
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: {{.Spec.Dbs.Pg.Port}}
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
            {{- if isTrue .Spec.Dbs.Pg.HugePages.Enabled }}
            - mountPath: "/hugepages"
              name: "hugepage"
            {{- end }}
          resources:
            limits:
              cpu: {{ .Spec.Dbs.Pg.Limits.Cpu }}
              memory: {{ .Spec.Dbs.Pg.Limits.Memory }}
              {{- if isTrue .Spec.Dbs.Pg.HugePages.Enabled }}
              {{- if eq .Spec.Dbs.Pg.HugePages.Memory ""}}
              hugepages-{{ .Spec.Dbs.Pg.HugePages.Size }}: {{ .Spec.Dbs.Pg.Requests.Memory }}
              {{- else }}
              hugepages-{{ .Spec.Dbs.Pg.HugePages.Size }}: {{ .Spec.Dbs.Pg.HugePages.Memory }}
              {{- end }}
              {{- end }}
            requests:
              cpu: {{ .Spec.Dbs.Pg.Requests.Cpu }}
              memory: {{ .Spec.Dbs.Pg.Requests.Memory }}
      volumes:
        - name: postgres-data
          persistentVolumeClaim:
            claimName: {{.Spec.Dbs.Pg.PvcName}}
        - name: dshm
          emptyDir:
            medium: Memory
            sizeLimit: 2Gi
        {{- if isTrue .Spec.Dbs.Pg.HugePages.Enabled }}
        - name: "hugepage"
          emptyDir:
            medium: HugePages
        {{- end}}
