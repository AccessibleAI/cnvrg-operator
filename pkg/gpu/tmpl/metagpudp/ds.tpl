apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: metagpu-device-plugin
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
      name: metagpu-device-plugin
  template:
    metadata:
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ""
        {{- range $k, $v := .Data.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        name: metagpu-device-plugin
        {{- range $k, $v := .Data.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      tolerations:
        - key: CriticalAddonsOnly
          operator: Exists
        - key: nvidia.com/gpu
          operator: Exists
          effect: NoSchedule
      priorityClassName: "system-node-critical"
      imagePullSecrets:
        - name: regcred
      hostPID: true
      hostNetwork: true
      serviceAccountName: metagpu-device-plugin
      nodeSelector:
        accelerator: nvidia
      containers:
        - name: metagpu-device-plugin
          image: "docker.io/cnvrg/metagpu-device-plugin:main"
          imagePullPolicy: Always
          command:
            - /usr/bin/mgdp
            - start
            - -c
            - /etc/metagpu-device-plugin
          ports:
            - containerPort: 50052
          securityContext:
            privileged: true
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: MG_CTL_TOKEN
              value: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1ldGFncHVAaW5zdGFuY2UiLCJ2aXNpYmlsaXR5TGV2ZWwiOiJsMCJ9.2rHykHFcHoIr-OCoPA5Am4ubf31-RJcayZnOTK6db94
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
          image: "docker.io/cnvrg/metagpu-device-plugin:main"
          imagePullPolicy: Always
          command:
            - /usr/bin/mgex
            - start
            - -t
            - eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1ldGFncHVAaW5zdGFuY2UiLCJ2aXNpYmlsaXR5TGV2ZWwiOiJsMCJ9.2rHykHFcHoIr-OCoPA5Am4ubf31-RJcayZnOTK6db94
          ports:
            - containerPort: 2112
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