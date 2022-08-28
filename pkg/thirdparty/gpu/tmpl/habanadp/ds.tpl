apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: habanalabs-device-plugin-daemonset-hpu
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
      name: habanalabs-device-plugin-ds
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      # This annotation is deprecated. Kept here for backward compatibility
      # See https://kubernetes.io/docs/tasks/administer-cluster/guaranteed-scheduling-critical-addon-pods/
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ""
        {{- range $k, $v := .Data.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        name: habanalabs-device-plugin-ds
        {{- range $k, $v := .Data.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      priorityClassName: "system-node-critical"
      serviceAccountName: habana-device-plugin
      tolerations:
        - operator: Exists
      nodeSelector:
        node.kubernetes.io/instance-type: dl1.24xlarge
      containers:
        - image: {{ image .Data.ImageHub .Data.HabanaDp.Image }}
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
          env:
          - name: LD_LIBRARY_PATH
            value: "/usr/lib/habanalabs/"
          securityContext:
            privileged: true
          volumeMounts:
            - name: device-plugin
              mountPath: /var/lib/kubelet/device-plugins
            - name: habana-lib
              mountPath: /usr/lib/habanalabs
      volumes:
        - name: device-plugin
          hostPath:
            path: /var/lib/kubelet/device-plugins
        - name: habana-lib
          hostPath:
            path: /usr/lib/habanalabs