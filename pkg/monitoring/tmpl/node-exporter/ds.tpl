apiVersion: apps/v1
kind: DaemonSet
metadata:
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app.kubernetes.io/name: node-exporter
    app.kubernetes.io/version: v1.0.1
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
  name: node-exporter
  namespace: {{ ns .}}
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: node-exporter
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app.kubernetes.io/name: node-exporter
        app.kubernetes.io/version: v1.0.1
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      containers:
      - args:
        - --web.listen-address=127.0.0.1:9100
        - --path.procfs=/host/proc
        - --path.sysfs=/host/sys
        - --path.rootfs=/host/root
        - --no-collector.wifi
        - --no-collector.hwmon
        - --collector.filesystem.ignored-mount-points=^/(dev|proc|sys|var/lib/docker/.+|var/lib/kubelet/pods/.+)($|/)
        image: {{ image .Spec.ImageHub .Spec.Monitoring.NodeExporter.Image }}
        name: node-exporter
        resources:
          limits:
            cpu: 500m
            memory: 1Gi
          requests:
            cpu: 102m
            memory: 180Mi
        volumeMounts:
        - mountPath: /host/proc
          mountPropagation: HostToContainer
          name: proc
          readOnly: true
        - mountPath: /host/sys
          mountPropagation: HostToContainer
          name: sys
          readOnly: true
        - mountPath: /host/root
          mountPropagation: HostToContainer
          name: root
          readOnly: true
      - args:
        - --logtostderr
        - --secure-listen-address=[$(IP)]:9100
        - --tls-cipher-suites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305
        - --upstream=http://127.0.0.1:9100/
        env:
        - name: IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        image: {{ image .Spec.ImageHub .Spec.Monitoring.PrometheusOperator.KubeRbacProxyImage }}
        name: kube-rbac-proxy
        ports:
        - containerPort: 9100
          hostPort: 9100
          name: https
        resources:
          limits:
            cpu: 20m
            memory: 40Mi
          requests:
            cpu: 10m
            memory: 20Mi
        securityContext:
          runAsGroup: 65532
          runAsNonRoot: true
          runAsUser: 65532
      hostNetwork: true
      hostPID: true
      nodeSelector:
        kubernetes.io/os: linux
      securityContext:
        runAsNonRoot: true
        runAsUser: 65534
      serviceAccountName: node-exporter
      tolerations:
      - operator: Exists
      volumes:
      - hostPath:
          path: /proc
        name: proc
      - hostPath:
          path: /sys
        name: sys
      - hostPath:
          path: /
        name: root
  updateStrategy:
    rollingUpdate:
      maxUnavailable: 10%
    type: RollingUpdate
