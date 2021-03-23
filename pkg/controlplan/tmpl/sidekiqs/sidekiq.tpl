apiVersion: apps/v1
kind: Deployment
metadata:
  name: sidekiq
  namespace: {{ ns . }}
  labels:
    app: sidekiq
spec:
  replicas: {{ .Spec.ControlPlan.Sidekiq.Replicas }}
  selector:
    matchLabels:
      app: sidekiq
  template:
    metadata:
      labels:
        app: sidekiq
    spec:
      {{- if eq .Spec.ControlPlan.Tenancy.Enabled "true" }}
      nodeSelector:
        {{ .Spec.ControlPlan.Tenancy.Key }}: "{{ .Spec.ControlPlan.Tenancy.Value }}"
      {{- end }}
      tolerations:
        - key: "{{ .Spec.ControlPlan.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.ControlPlan.Tenancy.Value }}"
          effect: "NoSchedule"
      serviceAccountName: {{ .Spec.ControlPlan.Rbac.ServiceAccountName }}
      terminationGracePeriodSeconds: {{ .Spec.ControlPlan.Sidekiq.KillTimeout }}
      containers:
        - name: sidekiq
          image: {{ .Spec.ControlPlan.WebApp.Image}}
          env:
            - name: "CNVRG_RUN_MODE"
              value: "sidekiq"
          imagePullPolicy: Always
          {{- if eq .Spec.ControlPlan.ObjectStorage.CnvrgStorageType "gcp" }}
          volumeMounts:
            - name: "{{ .Spec.ControlPlan.ObjectStorage.GcpStorageSecret }}"
              mountPath: "{{ .Spec.ControlPlan.ObjectStorage.GcpKeyfileMountPath }}"
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
                name: {{ .Spec.ControlPlan.Pg.SvcName }}
          resources:
            requests:
              cpu: {{ .Spec.ControlPlan.Sidekiq.CPU }}
              memory: {{ .Spec.ControlPlan.Sidekiq.Memory }}
          lifecycle:
            preStop:
              exec:
                command: ["/bin/bash","-lc","sidekiqctl quiet sidekiq-0.pid && sidekiqctl stop sidekiq-0.pid {{ .Spec.ControlPlan.Sidekiq.KillTimeout }}"]
      {{- if eq .Spec.ControlPlan.ObjectStorage.CnvrgStorageType "gcp" }}
      volumes:
        - name: {{ .Spec.ControlPlan.ObjectStorage.GcpStorageSecret }}
          secret:
            secretName: {{ .Spec.ControlPlan.ObjectStorage.GcpStorageSecret }}
      {{- end }}
      initContainers:
        - name: seeder
          image: {{.Spec.ControlPlan.Seeder.Image}}
          command: ["/bin/bash", "-c", "python3 cnvrg-boot.py seeder --mode worker"]
          env:
            - name: "CNVRG_NS"
              value: {{ ns . }}