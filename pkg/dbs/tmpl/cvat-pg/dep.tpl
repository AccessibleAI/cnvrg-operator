apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Spec.Dbs.Cvat.Pg.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{.Spec.Dbs.Cvat.Pg.SvcName }}
    cnvrg-component: pg
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{.Spec.Dbs.Cvat.Pg.SvcName }}
  strategy:
    type: Recreate
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: {{.Spec.Dbs.Cvat.Pg.SvcName }}
        cnvrg-component: pg
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      {{- if (gt (len .Spec.Dbs.Cvat.Pg.NodeSelector) 0) }}
      nodeSelector:
        {{- range $key, $val := .Spec.Dbs.Cvat.Pg.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      {{- end }}
      serviceAccountName: {{ .Spec.Dbs.Cvat.Pg.ServiceAccount }}
      securityContext:
        runAsUser: 26
        fsGroup: 26
      containers:
        - name: postgresql
          envFrom:
            - secretRef:
                name: {{ .Spec.Dbs.Cvat.Pg.CredsRef }}
          image: {{ image .Spec.ImageHub .Spec.Dbs.Cvat.Pg.Image}}
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: {{.Spec.Dbs.Cvat.Pg.Port}}
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
            {{- if isTrue .Spec.Dbs.Cvat.Pg.HugePages.Enabled }}
            - mountPath: "/hugepages"
              name: "hugepage"
            {{- end }}
          resources:
            limits:
              cpu: {{ .Spec.Dbs.Cvat.Pg.Limits.Cpu }}
              memory: {{ .Spec.Dbs.Cvat.Pg.Limits.Memory }}
              {{- if isTrue .Spec.Dbs.Cvat.Pg.HugePages.Enabled }}
              {{- if eq .Spec.Dbs.Cvat.Pg.HugePages.Memory ""}}
              hugepages-{{ .Spec.Dbs.Cvat.Pg.HugePages.Size }}: {{ .Spec.Dbs.Cvat.Pg.Requests.Memory }}
              {{- else }}
              hugepages-{{ .Spec.Dbs.Cvat.Pg.HugePages.Size }}: {{ .Spec.Dbs.Cvat.Pg.HugePages.Memory }}
              {{- end }}
              {{- end }}
            requests:
              cpu: {{ .Spec.Dbs.Cvat.Pg.Requests.Cpu }}
              memory: {{ .Spec.Dbs.Cvat.Pg.Requests.Memory }}
      volumes:
        - name: postgres-data
          persistentVolumeClaim:
            claimName: {{.Spec.Dbs.Cvat.Pg.PvcName}}
        - name: dshm
          emptyDir:
            medium: Memory
            sizeLimit: 2Gi
        {{- if isTrue .Spec.Dbs.Cvat.Pg.HugePages.Enabled }}
        - name: "hugepage"
          emptyDir:
            medium: HugePages
        {{- end}}
