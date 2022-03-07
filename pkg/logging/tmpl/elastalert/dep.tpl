apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: {{ .Spec.Logging.Elastalert.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  name: {{ .Spec.Logging.Elastalert.SvcName }}
  namespace: {{ ns . }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ .Spec.Logging.Elastalert.SvcName }}
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: {{ .Spec.Logging.Elastalert.SvcName }}
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        {{ .Spec.Tenancy.Key }}: {{ .Spec.Tenancy.Value }}
        {{- range $key, $val := .Spec.Logging.Elastalert.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      tolerations:
        - operator: "Exists"
      {{- else if (gt (len .Spec.Logging.Elastalert.NodeSelector) 0) }}
      nodeSelector:
        {{- range $key, $val := .Spec.Logging.Elastalert.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      {{- end }}
      serviceAccountName: {{ .Spec.Logging.Elastalert.SvcName }}
      securityContext:
        runAsUser: 1000
        fsGroup: 1000
      containers:
      - image: {{ image .Spec.ImageHub .Spec.Logging.Elastalert.Image }}
        name: {{ .Spec.Logging.Elastalert.SvcName }}
        ports:
        - containerPort: 3030
          protocol: TCP
        envFrom:
          - secretRef:
              name: {{ .Spec.Dbs.Es.CredsRef }}
        resources:
          requests:
            cpu: {{.Spec.Logging.Elastalert.Requests.Cpu}}
            memory: {{.Spec.Logging.Elastalert.Requests.Memory}}
          limits:
            cpu: {{ .Spec.Logging.Elastalert.Limits.Cpu }}
            memory: {{ .Spec.Logging.Elastalert.Limits.Memory }}
        volumeMounts:
        - mountPath: /opt/elastalert-server/config/config.json
          subPath: config.json
          name: elastalert-config
        - mountPath: /opt/elastalert/config.yaml
          subPath: config.yaml
          name: elastalert-config
        - mountPath: /opt/elastalert/rules
          name: {{ .Spec.Logging.Elastalert.SvcName }}
      - name: "elastalert-auth-proxy"
        image: {{ image .Spec.ImageHub .Spec.Logging.Elastalert.AuthProxyImage }}
        ports:
        - containerPort: 8080
        volumeMounts:
        - name: elastalert-auth-proxy
          mountPath: "/etc/nginx"
          readOnly: true
        - name: htpasswd
          mountPath: "/etc/nginx/htpasswd"
          readOnly: true
        resources:
          requests:
            cpu: 100m
            memory: 100Mi
          limits:
            cpu: 1000m
            memory: 1Gi
      restartPolicy: Always
      volumes:
      - name: {{ .Spec.Logging.Elastalert.SvcName }}
        persistentVolumeClaim:
          claimName: {{ .Spec.Logging.Elastalert.PvcName }}
      - name: elastalert-config
        configMap:
          name: elastalert-config
      - name: elastalert-auth-proxy
        configMap:
          name: elastalert-auth-proxy
      - name: htpasswd
        secret:
          secretName: {{ .Spec.Logging.Elastalert.CredsRef }}
