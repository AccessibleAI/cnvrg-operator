apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: habanalabs-device-plugin-daemonset-gaudi
  namespace: {{ .Namespace }}
  annotations:
    {{ - range $k, $v := .Data.Annotations }}
    {{$k}}: "{{$v}}"
    {{ - end }}
  labels:
    {{ - range $k, $v := .Data.Labels }}
    {{$k}}: "{{$v}}"
    {{ - end }}
spec:
  selector:
    matchLabels:
      name: habanalabs-device-plugin-ds
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      # This annotation is deprecated. Kept here for backward compatibility
      # See https://kubernetes.io/docs/tasks/administer-cluster/guaranteed-scheduling-critical-addon-pods/
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ""
        {{ - range $k, $v := .Data.Annotations }}
        {{$k}}: "{{$v}}"
        {{ - end }}
      labels:
        name: habanalabs-device-plugin-ds
        {{ - range $k, $v := .Data.Labels }}
        {{$k}}: "{{$v}}"
        {{ - end }}
    spec:
      priorityClassName: "system-node-critical"
      serviceAccountName: cnvrg-control-plane
      nodeSelector:
        node.kubernetes.io/instance-type: dl1.24xlarge
      containers:
        - image: {{ image .Data.ImageHub .Data.HpuDp.Image }}
          name: habanalabs-device-plugin-ctr
          resources:
            requests:
              cpu: 100m
              memory: 100Mi
            limits:
              cpu: 500m
              memory: 500Mi
          command: ["habanalabs-device-plugin"]
          args: ["--dev_type", " gaudi"]
          securityContext:
            privileged: true
          volumeMounts:
            - name: device-plugin
              mountPath: /var/lib/kubelet/device-plugins
      volumes:
        - name: device-plugin
          hostPath:
            path: /var/lib/kubelet/device-plugins