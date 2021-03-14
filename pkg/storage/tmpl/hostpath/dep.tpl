apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: hostpath-provisioner
  namespace: {{ .CnvrgNs }}
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
      tolerations:
        - key: {{ .ControlPlan.Tenancy.Key }}
          operator: Equal
          value: "{{ .ControlPlan.Tenancy.Value }}"
          effect: "NoSchedule"
      {{- if and (eq .Storage.Hostpath.Enabled "true") (eq .ControlPlan.Tenancy.Enabled "false") }}
      nodeSelector:
        kubernetes.io/hostname: "{{ .Storage.Hostpath.NodeName }}"
      {{- else if and (eq .Storage.Hostpath.Enabled "false") (eq .ControlPlan.Tenancy.Enabled "true") }}
      nodeSelector:
        {{ .ControlPlan.Tenancy.Key }}: "{{ .ControlPlan.Tenancy.Value }}"
      {{- else if and (eq .Storage.Hostpath.Enabled "true") (eq .ControlPlan.Tenancy.Enabled "true") }}
      nodeSelector:
        kubernetes.io/hostname: "{{ .Storage.Hostpath.NodeName }}"
        {{ .ControlPlan.Tenancy.Key }}: "{{ .ControlPlan.Tenancy.Value }}"
      {{- end }}
      serviceAccountName: hostpath-provisioner-admin
      containers:
        - name: hostpath-provisioner
          image: {{ .Storage.Hostpath.Image }}
          imagePullPolicy: Always
          env:
            - name: USE_NAMING_PREFIX
              value: "true" # change to true, to have the name of the pvc be part of the directory
            - name: NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: PV_DIR
              value: {{ .Storage.Hostpath.HostPath }}
          volumeMounts:
            - name: pv-volume # root dir where your bind mounts will be on the node
              mountPath: {{ .Storage.Hostpath.HostPath }}
              nodeSelector:
              - name: {{ .Storage.Hostpath.NodeName }}
          resources:
            limits:
              cpu: {{ .Storage.Hostpath.CPULimit }}
              memory: {{ .Storage.Hostpath.MemoryLimit }}
            requests:
              cpu: {{ .Storage.Hostpath.CPURequest }}
              memory: {{ .Storage.Hostpath.MemoryRequest }}
      volumes:
        - name: pv-volume
          hostPath:
            path: {{ .Storage.Hostpath.HostPath }}
