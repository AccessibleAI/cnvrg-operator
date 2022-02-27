apiVersion: apps/v1
kind: Deployment
metadata:
  name: systemkiq
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-ccp"
    app: systemkiq
    owner: cnvrg-control-plane
    cnvrg-component: systemkiq
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  {{- if eq false (isTrue .Spec.ControlPlane.Systemkiq.Hpa.Enabled) }}
  replicas: {{ .Spec.ControlPlane.Systemkiq.Replicas }}
  {{- end }}
  selector:
    matchLabels:
      app: systemkiq
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: systemkiq
        cnvrg-component: systemkiq
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
                  app: systemkiq
              namespaces:
              - {{ ns . }}
              topologyKey: kubernetes.io/hostname
            weight: 1
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        "{{ .Spec.Tenancy.Key }}": "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - key: "{{ .Spec.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.Tenancy.Value }}"
          effect: "NoSchedule"
      {{- end }}
      serviceAccountName: cnvrg-control-plane
      securityContext:
        runAsUser: 1000
        runAsGroup: 1000
      terminationGracePeriodSeconds: 60
      containers:
        - name: sidekiq
          image: {{ image .Spec.ImageHub .Spec.ControlPlane.Image}}
          env:
            - name: "CNVRG_RUN_MODE"
              value: "sidekiq"
            - name: "SIDEKIQ_SYSTEM"
              value: "true"
          imagePullPolicy: Always
          {{- if eq .Spec.ControlPlane.ObjectStorage.Type "gcp" }}
          volumeMounts:
            - name: {{ .Spec.ControlPlane.ObjectStorage.GcpSecretRef }}
              mountPath: /opt/app-root/conf/gcp-keyfile
              readOnly: true
          {{- end }}
          envFrom:
            - configMapRef:
                name: cp-base-config
            - configMapRef:
                name: cp-networking-config
            - secretRef:
                name: cp-base-secret
            - secretRef:
                name: cp-oauth-proxy-tokens-secret
            - secretRef:
                name: cp-ldap
            - secretRef:
                name: cp-object-storage
            - secretRef:
                name: cp-smtp
            {{- if isTrue .Spec.Dbs.Cvat.Enabled }}
            - secretRef:
                name: {{ .Spec.Dbs.Cvat.Pg.CredsRef }}
            - secretRef:
                name: {{ .Spec.Dbs.Cvat.Redis.CredsRef }}
            {{- end }}
            - secretRef:
                name: {{ .Spec.Dbs.Es.CredsRef }}
            - secretRef:
                name: {{ .Spec.Dbs.Pg.CredsRef }}
            - secretRef:
                name: {{ .Spec.Dbs.Redis.CredsRef }}
            - secretRef:
                name: {{ .Spec.Monitoring.Prometheus.CredsRef }}
            {{- if isTrue .Spec.Networking.Proxy.Enabled }}
            - configMapRef:
                name: {{ .Spec.Networking.Proxy.ConfigRef }}
            {{- end }}
          resources:
            requests:
              cpu: {{ .Spec.ControlPlane.Systemkiq.Requests.Cpu }}
              memory: {{ .Spec.ControlPlane.Systemkiq.Requests.Memory }}
            limits:
              cpu: {{ .Spec.ControlPlane.Systemkiq.Limits.Cpu }}
              memory: {{ .Spec.ControlPlane.Systemkiq.Limits.Memory }}
          lifecycle:
            preStop:
              exec:
                command: ["/bin/bash","-lc","sidekiqctl quiet sidekiq-0.pid && sidekiqctl stop sidekiq-0.pid 60"]
      {{- if eq .Spec.ControlPlane.ObjectStorage.Type "gcp" }}
      volumes:
        - name: {{ .Spec.ControlPlane.ObjectStorage.GcpSecretRef }}
          secret:
            secretName: {{ .Spec.ControlPlane.ObjectStorage.GcpSecretRef }}
      {{- end }}
      initContainers:
        - name: seeder
          image:  {{ image .Spec.ImageHub .Spec.ControlPlane.Image }}
          command: ["/bin/bash", "-lc", "while true; do if [[ $(kubectl get cm cnvrg-db-init -oname --ignore-not-found | wc -l) == 0 ]]; then echo 'cnvrg seed not ready'; sleep 1; else echo 'cnvrg seed is done'; exit 0; fi; done"]
          resources:
            requests:
              cpu: "100m"
              memory: "100Mi"
            limits:
              cpu: "1000m"
              memory: "1Gi"
          env:
            - name: "CNVRG_NS"
              value: {{ ns . }}