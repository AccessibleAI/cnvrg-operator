apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: habana-exporter
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app.kubernetes.io/name: habana-exporter
    app.kubernetes.io/version: v0.0.1
    app: "habana-exporter"
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: habana-exporter
      app: "habana-exporter"
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app.kubernetes.io/name: habana-exporter
        app.kubernetes.io/version: v0.0.1
        app: "habana-exporter"
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      serviceAccountName: habana-exporter
      volumes:
        - name: "pod-resources"
          hostPath:
            path: "/var/lib/kubelet/pod-resources"
        - name: "habana-lib"
          hostPath:
            path: "/usr/lib/habanalabs"
      tolerations:
        - operator: Exists
      nodeSelector:
        node.kubernetes.io/instance-type: dl1.24xlarge
      containers:
        - name: exporter
          securityContext:
            privileged: true
          image: {{image .Spec.ImageHub .Spec.Monitoring.HabanaExporter.Image }}
          imagePullPolicy: "IfNotPresent"
          env:
          - name: LD_LIBRARY_PATH
            value: "/usr/lib/habanalabs"
          volumeMounts:
            - name: pod-resources
              mountPath: /var/lib/kubelet/pod-resources
            - name: habana-lib
              mountPath: /usr/lib/habanalabs
          ports:
          - name: habana-metrics
            containerPort: 41611
            protocol: TCP
          resources:
            limits:
              cpu: 1
              memory: 200Mi
            requests:
              cpu: 100m
              memory: 200Mi
        - name: hlml-service
          image: {{ .Spec.Monitoring.HabanaExporter.HlmlImage }}
          env:
          - name: LD_LIBRARY_PATH
            value: "/usr/lib/habanalabs"
          securityContext:
            privileged: true
          resources:
            limits:
              cpu: 1
              memory: 200Mi
            requests:
              cpu: 100m
              memory: 200Mi
          livenessProbe:
            exec:
              command:
                - ls
                - /tmp
            initialDelaySeconds: 5
            periodSeconds: 5
          readinessProbe:
            exec:
              command:
                - ls
                - /tmp
            initialDelaySeconds: 5
            periodSeconds: 5
          volumeMounts:
            - name: habana-lib
              mountPath: /usr/lib/habanalabs