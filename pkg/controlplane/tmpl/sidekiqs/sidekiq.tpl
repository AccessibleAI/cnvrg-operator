apiVersion: apps/v1
kind: Deployment
metadata:
  name: sidekiq
  namespace: {{ ns . }}
  labels:
    app: sidekiq
spec:
  replicas: {{ .Spec.ControlPlane.Sidekiq.Replicas }}
  selector:
    matchLabels:
      app: sidekiq
  template:
    metadata:
      labels:
        app: sidekiq
    spec:
      {{- if eq .Spec.ControlPlane.Tenancy.Enabled "true" }}
      nodeSelector:
        {{ .Spec.ControlPlane.Tenancy.Key }}: "{{ .Spec.ControlPlane.Tenancy.Value }}"
      {{- end }}
      tolerations:
        - key: "{{ .Spec.ControlPlane.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.ControlPlane.Tenancy.Value }}"
          effect: "NoSchedule"
      serviceAccountName: {{ .Spec.ControlPlane.Rbac.ServiceAccountName }}
      terminationGracePeriodSeconds: {{ .Spec.ControlPlane.Sidekiq.KillTimeout }}
      containers:
        - name: sidekiq
          image: {{ .Spec.ControlPlane.WebApp.Image}}
          env:
            - name: "CNVRG_RUN_MODE"
              value: "sidekiq"
          imagePullPolicy: Always
          {{- if eq .Spec.ControlPlane.ObjectStorage.CnvrgStorageType "gcp" }}
          volumeMounts:
            - name: "{{ .Spec.ControlPlane.ObjectStorage.GcpStorageSecret }}"
              mountPath: "{{ .Spec.ControlPlane.ObjectStorage.GcpKeyfileMountPath }}"
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
                name: {{ .Spec.Dbs.Pg.SvcName }}
            - secretRef:
                name: cp-smtp
          resources:
            requests:
              cpu: {{ .Spec.ControlPlane.Sidekiq.CPU }}
              memory: {{ .Spec.ControlPlane.Sidekiq.Memory }}
          lifecycle:
            preStop:
              exec:
                command: ["/bin/bash","-lc","sidekiqctl quiet sidekiq-0.pid && sidekiqctl stop sidekiq-0.pid {{ .Spec.ControlPlane.Sidekiq.KillTimeout }}"]
      {{- if eq .Spec.ControlPlane.ObjectStorage.CnvrgStorageType "gcp" }}
      volumes:
        - name: {{ .Spec.ControlPlane.ObjectStorage.GcpStorageSecret }}
          secret:
            secretName: {{ .Spec.ControlPlane.ObjectStorage.GcpStorageSecret }}
      {{- end }}
      initContainers:
        - name: seeder
          image: {{.Spec.ControlPlane.Seeder.Image}}
          command: ["/bin/bash", "-c", "python3 cnvrg-boot.py seeder --mode worker"]
          env:
            - name: "CNVRG_NS"
              value: {{ ns . }}