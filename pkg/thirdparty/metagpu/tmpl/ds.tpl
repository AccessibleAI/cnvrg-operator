apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: metagpu-device-plugin
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
spec:
  selector:
    matchLabels:
      name: metagpu-device-plugin
  template:
    metadata:
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ""
      labels:
        name: metagpu-device-plugin
    spec:
      tolerations:
        - operator: Exists
      priorityClassName: "system-node-critical"
      imagePullSecrets:
        - name: regcred
      hostPID: true
      hostNetwork: true
      serviceAccountName: metagpu-device-plugin
      nodeSelector:
        {{- range $k, $v := .Spec.Metagpu.NodeSelector }}
        {{$k}}: "{{$v}}"
        {{- end }}
      containers:
        - name: metagpu-device-plugin
          image: {{ image .Spec.ImageHub .Spec.Metagpu.Image }}
          imagePullPolicy: Always
          command:
            - /usr/bin/mgdp
            - start
            - -c
            - /etc/metagpu-device-plugin
          securityContext:
            privileged: true
          envFrom:
            - configMapRef:
                name: metagpu-device-plugin-config
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          volumeMounts:
            - name: device-plugin
              mountPath: /var/lib/kubelet/device-plugins
            - name: config
              mountPath: /etc/metagpu-device-plugin
            - mountPath: /host/proc
              mountPropagation: HostToContainer
              name: proc
              readOnly: true
        - name: metagpu-exporter
          image: {{ image .Spec.ImageHub .Spec.Metagpu.Image }}
          imagePullPolicy: Always
          command:
            - /usr/bin/mgex
            - start
          envFrom:
            - configMapRef:
                name: metagpu-device-plugin-config
      volumes:
        - name: device-plugin
          hostPath:
            path: /var/lib/kubelet/device-plugins
        - name: config
          configMap:
            name: metagpu-device-plugin-config
        - hostPath:
            path: /proc
          name: proc