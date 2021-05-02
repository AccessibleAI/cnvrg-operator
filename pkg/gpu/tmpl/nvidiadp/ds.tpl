apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: nvidia-device-plugin-daemonset
  namespace: {{ .Namespace }}
  labels:
    owner: cnvrg-control-plane
spec:
  selector:
    matchLabels:
      name: nvidia-device-plugin-ds
  template:
    metadata:
      labels:
        name: nvidia-device-plugin-ds
        owner: cnvrg-control-plane
    spec:
      serviceAccountName: nvidia-device-plugin
      tolerations:
        - key: nvidia.com/gpu
          operator: Exists
          effect: NoSchedule
      priorityClassName: "system-node-critical"
      nodeSelector:
        accelerator: nvidia
      containers:
        - image: {{ .Data.NvidiaDp.Image }}
          name: nvidia-device-plugin-ctr
          args: ["--fail-on-init-error=true"]
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