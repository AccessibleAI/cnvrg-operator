apiVersion: apps/v1
kind: Deployment
metadata:
  name: sidekiq
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: sidekiq
    owner: cnvrg-control-plane
    cnvrg-component: sidekiq
    cnvrg-system-status-check: "true"
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  {{- if eq false (isTrue .Spec.ControlPlane.Sidekiq.Hpa.Enabled) }}
  replicas: {{ .Spec.ControlPlane.Sidekiq.Replicas }}
  {{- end }}
  selector:
    matchLabels:
      app: sidekiq
  template:
    metadata:
      annotations:
        cnvrg-component: sidekiq
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: sidekiq
        owner: cnvrg-control-plane
        cnvrg-component: sidekiq
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
                  app: sidekiq
              namespaces:
              - {{ .Namespace }}
              topologyKey: kubernetes.io/hostname
            weight: 1
      priorityClassName: {{ .Spec.PriorityClass.AppClassRef }}
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
          imagePullPolicy: Always
          volumeMounts:
          {{- if eq .Spec.ControlPlane.ObjectStorage.Type "gcp" }}
            - name: {{ .Spec.ControlPlane.ObjectStorage.GcpSecretRef }}
              mountPath: /opt/app-root/conf/gcp-keyfile
              readOnly: true
          {{- end }}
          {{- if and ( isTrue .Spec.Networking.Ingress.OcpSecureRoutes) (eq .Spec.Networking.Ingress.Type "openshift") }}
            - name: tls-secret
              readOnly: true
              mountPath: /opt/app-root/src/tls
          {{- end }}
          envFrom:
            - configMapRef:
                name: cp-base-config
            - configMapRef:
                name: cp-networking-config
            - secretRef:
                name: cp-base-secret
            {{- if isTrue .Spec.SSO.Enabled }}
            - secretRef:
                name: cp-oauth-proxy-tokens-secret
            {{- end }}
            - secretRef:
                name: cp-ldap
            - secretRef:
                name: cp-object-storage
            - secretRef:
                name: cp-smtp
            - secretRef:
                name: {{ .Spec.Dbs.Es.CredsRef }}
            - secretRef:
                name: {{ .Spec.Dbs.Pg.CredsRef }}
            - secretRef:
                name: {{ .Spec.Dbs.Redis.CredsRef }}
            - secretRef:
                name: {{ .Spec.Dbs.Prom.CredsRef }}
            {{- if isTrue .Spec.Networking.Proxy.Enabled }}
            - configMapRef:
                name: {{ .Spec.Networking.Proxy.ConfigRef }}
            {{- end }}
          resources:
            requests:
              cpu: {{ .Spec.ControlPlane.Sidekiq.Requests.Cpu }}
              memory: {{ .Spec.ControlPlane.Sidekiq.Requests.Memory }}
            limits:
              cpu: {{ .Spec.ControlPlane.Sidekiq.Limits.Cpu }}
              memory: {{ .Spec.ControlPlane.Sidekiq.Limits.Memory }}
          lifecycle:
            preStop:
              exec:
                command: ["/bin/bash","-lc","sidekiqctl quiet sidekiq-0.pid && sidekiqctl stop sidekiq-0.pid 60"]
      volumes:
      {{- if eq .Spec.ControlPlane.ObjectStorage.Type "gcp" }}
        - name: {{ .Spec.ControlPlane.ObjectStorage.GcpSecretRef }}
          secret:
            secretName: {{ .Spec.ControlPlane.ObjectStorage.GcpSecretRef }}
      {{- end }}
      {{- if and ( isTrue .Spec.Networking.Ingress.OcpSecureRoutes) (eq .Spec.Networking.Ingress.Type "openshift") }}
        - name: tls-secret
          secret:
            secretName: {{ .Spec.Networking.HTTPS.CertSecret }}
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
              value: {{ .Namespace }}