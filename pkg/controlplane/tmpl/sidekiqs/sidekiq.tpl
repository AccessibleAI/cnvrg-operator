apiVersion: apps/v1
kind: Deployment
metadata:
  name: sidekiq
  namespace: {{ ns . }}
  labels:
    app: sidekiq
    owner: cnvrg-control-plane
    cnvrg-component: sidekiq
spec:
  replicas: {{ .Spec.ControlPlane.Sidekiq.Replicas }}
  selector:
    matchLabels:
      app: sidekiq
  template:
    metadata:
      labels:
        app: sidekiq
        owner: cnvrg-control-plane
        cnvrg-component: sidekiq
    spec:
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        "{{ .Spec.Tenancy.Key }}": "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - key: "{{ .Spec.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.Tenancy.Value }}"
          effect: "NoSchedule"
      {{- end }}
      serviceAccountName: {{ .Spec.ControlPlane.Rbac.ServiceAccountName }}
      securityContext:
        runAsUser: 1000
        runAsGroup: 1000
      terminationGracePeriodSeconds: 60
      containers:
        - name: sidekiq
          image: {{ .Spec.ControlPlane.Image}}
          env:
            - name: "CNVRG_RUN_MODE"
              value: "sidekiq"
          imagePullPolicy: Always
          {{- if eq (.Spec.ControlPlane.ObjectStorage.Type | toString) "gcp" }}
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
                name: {{ .Spec.Monitoring.Prometheus.CredsRef }}
          resources:
            requests:
              cpu: {{ .Spec.ControlPlane.Sidekiq.Requests.Cpu }}
              memory: {{ .Spec.ControlPlane.Sidekiq.Requests.Memory }}
          lifecycle:
            preStop:
              exec:
                command: ["/bin/bash","-lc","sidekiqctl quiet sidekiq-0.pid && sidekiqctl stop sidekiq-0.pid 60"]
      {{- if eq (.Spec.ControlPlane.ObjectStorage.Type | toString) "gcp" }}
      volumes:
        - name: {{ .Spec.ControlPlane.ObjectStorage.GcpSecretRef }}
          secret:
            secretName: {{ .Spec.ControlPlane.ObjectStorage.GcpSecretRef }}
      {{- end }}
      initContainers:
        - name: seeder
          image: {{.Spec.ControlPlane.Seeder.Image}}
          command: ["/bin/bash", "-c", "python3 cnvrg-boot.py seeder --mode worker"]
          env:
            - name: "CNVRG_NS"
              value: {{ ns . }}