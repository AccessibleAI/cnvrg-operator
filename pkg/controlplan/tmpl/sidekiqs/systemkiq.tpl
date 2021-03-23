apiVersion: apps/v1
kind: Deployment
metadata:
  name: systemkiq
  namespace: {{ ns . }}
  labels:
    app: systemkiq
spec:
  replicas: {{ .Spec.ControlPlan.Systemkiq.Replicas }}
  selector:
    matchLabels:
      app: systemkiq
  template:
    metadata:
      labels:
        app: systemkiq
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
      terminationGracePeriodSeconds: {{ .Spec.ControlPlan.Systemkiq.KillTimeout }}
      containers:
        - name: sidekiq
          image: {{ .Spec.ControlPlan.WebApp.Image}}
          env:
            - name: "CNVRG_RUN_MODE"
              value: "sidekiq"
            - name: "SIDEKIQ_SYSTEM"
              value: "true"
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
              cpu: {{ .Spec.ControlPlan.Systemkiq.CPU }}
              memory: {{ .Spec.ControlPlan.Systemkiq.Memory }}
          lifecycle:
            preStop:
              exec:
                command: ["/bin/bash","-lc","sidekiqctl quiet sidekiq-0.pid && sidekiqctl stop sidekiq-0.pid {{ .Spec.ControlPlan.Systemkiq.KillTimeout }}"]
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