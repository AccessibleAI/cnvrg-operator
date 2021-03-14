kind: Deployment
apiVersion: apps/v1
metadata:
  name: nfs-client-provisioner
  namespace: {{ .CnvrgNs }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nfs-client-provisioner
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: nfs-client-provisioner
    spec:
      {{- if eq .ControlPlan.Tenancy.Enabled "true" }}
      nodeSelector:
        {{ .ControlPlan.Tenancy.Key }}: "{{ .ControlPlan.Tenancy.Value }}"
      {{- end }}
      tolerations:
        - key: "{{ .ControlPlan.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .ControlPlan.Tenancy.Value }}"
          effect: "NoSchedule"
      serviceAccountName: nfs-client-provisioner
      containers:
        - name: nfs-client-provisioner
          image: {{ .Storage.Nfs.Image }}
          volumeMounts:
            - name: nfs-client-root
              mountPath: /persistentvolumes
          env:
            - name: PROVISIONER_NAME
              value: {{ .Storage.Nfs.Provisioner }}
            - name: NFS_SERVER
              value: {{ .Storage.Nfs.Server }}
            - name: NFS_PATH
              value: "{{ .Storage.Nfs.Path }}"
          resources:
            limits:
              cpu: {{ .Storage.Nfs.CPULimit }}
              memory: {{ .Storage.Nfs.MemoryLimit }}
            requests:
              cpu: {{ .Storage.Nfs.CPURequest }}
              memory: {{ .Storage.Nfs.MemoryRequest }}
      volumes:
        - name: nfs-client-root
          nfs:
            server: {{ .Storage.Nfs.Server }}
            path: {{ .Storage.Nfs.Path }}