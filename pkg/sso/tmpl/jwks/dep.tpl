apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Spec.SSO.Jwks.SvcName }}
  namespace: {{.Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.SSO.Jwks.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  replicas: {{ .Spec.SSO.Jwks.Replicas }}
  selector:
    matchLabels:
      app: {{ .Spec.SSO.Jwks.SvcName }}
  template:
    metadata:
      labels:
        app: {{ .Spec.SSO.Jwks.SvcName }}
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - podAffinityTerm:
              labelSelector:
                matchLabels:
                  app: {{.Spec.SSO.Jwks.SvcName}}
              namespaces:
              - {{.Namespace}}
              topologyKey: kubernetes.io/hostname
            weight: 1
      priorityClassName: {{ .Spec.PriorityClass.AppClassRef }}
      serviceAccountName: {{ .Spec.SSO.Jwks.SvcName }}
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        "{{ .Spec.Tenancy.Key }}": "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - operator: "Exists"
      {{- else if (gt (len .Spec.SSO.Jwks.NodeSelector) 0) }}
      nodeSelector:
        {{- range $key, $val := .Spec.SSO.Jwks.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      {{- end }}
      enableServiceLinks: false
      containers:
      - name: {{ .Spec.SSO.Jwks.SvcName }}
        command:
          - /usr/bin/cnvrg-jwks
          - start
        image: {{ image .Spec.ImageHub .Spec.SSO.Jwks.Image }}
        imagePullPolicy: Always
        resources:
          requests:
            cpu: 100m
            memory: 500Mi
          limits:
            cpu: 1000m
            memory: 1Gi
        volumeMounts:
          - mountPath: /opt/app-root/config
            name: {{ .Spec.SSO.Jwks.SvcName }}
        ports:
          - containerPort: 8080
        livenessProbe:
          timeoutSeconds: 5
          successThreshold: 1
          failureThreshold: 5
          initialDelaySeconds: 5
          periodSeconds: 10
          httpGet:
            port: 8080
            path: /healthz
        readinessProbe:
          successThreshold: 1
          failureThreshold: 5
          initialDelaySeconds: 5
          periodSeconds: 10
          timeoutSeconds: 5
          httpGet:
            port: 8080
            path: /healthz
      - name: redis-cache
        image: {{ image .Spec.ImageHub .Spec.SSO.Jwks.CacheImage }}
        resources:
          requests:
            cpu: 200m
            memory: 200Mi
          limits:
            cpu: 1000m
            memory: 1Gi
        livenessProbe:
          successThreshold: 1
          failureThreshold: 5
          initialDelaySeconds: 5
          periodSeconds: 10
          exec:
            command: ["redis-cli", "ping"]
        readinessProbe:
          successThreshold: 1
          failureThreshold: 5
          initialDelaySeconds: 5
          periodSeconds: 10
          exec:
            command: [ "redis-cli", "ping" ]
      volumes:
        - name: {{ .Spec.SSO.Jwks.SvcName }}
          configMap:
            name: {{ .Spec.SSO.Jwks.SvcName }}
