apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: nvidia-device-plugin-daemonset
  namespace: {{ .Namespace }}
  annotations:
    {{- range $k, $v := .Data.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Data.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  selector:
    matchLabels:
      name: nvidia-device-plugin-ds
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Data.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        name: nvidia-device-plugin-ds
        {{- range $k, $v := .Data.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      serviceAccountName: nvidia-device-plugin
      tolerations:
        - operator: Exists
        - key: nvidia.com/gpu
          operator: Exists
          effect: NoSchedule
      priorityClassName: "system-node-critical"
      nodeSelector:
        accelerator: nvidia
      containers:
        - image: {{ image .Data.ImageHub .Data.NvidiaDp.Image }}
          name: nvidia-device-plugin-ctr
          args: ["--fail-on-init-error=true"]
          resources:
            requests:
              cpu: 100m
              memory: 100Mi
            limits:
              cpu: 500m
              memory: 500Mi
          securityContext:
            allowPrivilegeEscalation: false
            capabilities:
              drop: ["ALL"]
          volumeMounts:
            - name: device-plugin
              mountPath: /var/lib/kubelet/device-plugins
      volumes:
        - name: device-plugin
          hostPath:
            path: /var/lib/kubelet/device-plugins