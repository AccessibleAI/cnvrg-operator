apiVersion: apps/v1
kind: Deployment
metadata:
  name: systemkiq
  namespace: {{ ns . }}
  labels:
    app: systemkiq
spec:
  replicas: {{ .Spec.ControlPlane.Systemkiq.Replicas }}
  selector:
    matchLabels:
      app: systemkiq
  template:
    metadata:
      labels:
        app: systemkiq
    spec:
      serviceAccountName: {{ .Spec.ControlPlane.Rbac.ServiceAccountName }}
      securityContext:
        runAsUser: 1000
        runAsGroup: 1000
      terminationGracePeriodSeconds: 60
      containers:
        - name: sidekiq
          image: {{ .Spec.ControlPlane.WebApp.Image}}
          env:
            - name: "CNVRG_RUN_MODE"
              value: "sidekiq"
            - name: "SIDEKIQ_SYSTEM"
              value: "true"
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
                name: {{ .Spec.Dbs.Pg.CredsRef }}
            - secretRef:
                name: cp-smtp
            - secretRef:
                name: {{ .Spec.Dbs.Redis.CredsRef }}
          resources:
            requests:
              cpu: {{ .Spec.ControlPlane.Systemkiq.CPU }}
              memory: {{ .Spec.ControlPlane.Systemkiq.Memory }}
          lifecycle:
            preStop:
              exec:
                command: ["/bin/bash","-lc","sidekiqctl quiet sidekiq-0.pid && sidekiqctl stop sidekiq-0.pid 60"]
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