apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: nvidia-device-plugin-daemonset
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
spec:
  selector:
    matchLabels:
      name: nvidia-device-plugin-ds
  template:
    metadata:
      labels:
        name: nvidia-device-plugin-ds
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
        - image: {{ image .Spec.ImageHub .Spec.Nvidia.DevicePlugin.Image }}
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