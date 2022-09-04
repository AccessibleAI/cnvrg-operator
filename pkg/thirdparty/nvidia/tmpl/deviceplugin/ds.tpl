apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: nvidia-device-plugin-daemonset
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "false"
    mlops.cnvrg.io/updatable: "false"
spec:
  selector:
    matchLabels:
      name: nvidia-device-plugin-ds
  template:
    metadata:
      labels:
        name: nvidia-device-plugin-ds
    spec:
      serviceAccountName: nvidia
      tolerations:
        - operator: Exists
      priorityClassName: "system-node-critical"
      nodeSelector:
        {{.Spec.Nvidia.NodeSelector.Key}}: {{.Spec.Nvidia.NodeSelector.Value}}
      containers:
        - image: {{ image .Spec.ImageHub .Spec.Nvidia.MetricsExporter.Image }}
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