apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: hostpath-provisioner
  namespace: {{ ns . }}
  labels:
    k8s-app: hostpath-provisioner
spec:
  selector:
    matchLabels:
      k8s-app: hostpath-provisioner
  template:
    metadata:
      labels:
        k8s-app: hostpath-provisioner
    spec:
      nodeSelector:
      {{- if .Spec.Tenancy.Enabled }}
        kubernetes.io/hostname: {{ .Spec.Storage.Hostpath.NodeName }}
        {{.Spec.Tenancy.Key }}: {{ .Spec.Tenancy.Value }}
      {{- else }}
        kubernetes.io/hostname: {{ .Spec.Storage.Hostpath.NodeName }}
      {{- end }}
      tolerations:
        - key: {{.Spec.Tenancy.Key }}
          operator: "Equal"
          value: {{ .Spec.Tenancy.Value }}
          effect: "NoSchedule"
      serviceAccountName: hostpath-provisioner-admin
      containers:
        - name: hostpath-provisioner
          image: {{ .Spec.Storage.Hostpath.Image }}
          imagePullPolicy: Always
          env:
            - name: USE_NAMING_PREFIX
              value: "true"
            - name: NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: PV_DIR
              value: {{ .Spec.Storage.Hostpath.HostPath }}
          volumeMounts:
            - name: pv-volume
              mountPath: {{ .Spec.Storage.Hostpath.HostPath }}
          resources:
            limits:
              cpu: {{ .Spec.Storage.Hostpath.CPULimit }}
              memory: {{ .Spec.Storage.Hostpath.MemoryLimit }}
            requests:
              cpu: {{ .Spec.Storage.Hostpath.CPURequest }}
              memory: {{ .Spec.Storage.Hostpath.MemoryRequest }}
      volumes:
        - name: pv-volume
          hostPath:
            path: {{ .Spec.Storage.Hostpath.HostPath }}
