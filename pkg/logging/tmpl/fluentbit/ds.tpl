apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluent-bit
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: fluentbit
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-fluentbit"
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  selector:
    matchLabels:
      app: fluentbit
  template:
    metadata:
      labels:
        app: fluentbit
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      containers:
        - name: fluent-bit
          image: {{.Spec.ImageHub }}/{{ .Spec.Logging.Fluentbit.Image }}
          imagePullPolicy: Always
          command:
            - /bin/bash
            - -c
            - /opt/app-root/fluentbit -c /opt/app-root/etc/fluent-bit.conf
          ports:
            - containerPort: 2020
          volumeMounts:
            - name: varlog
              mountPath: /var/log
            - name: varlibdockercontainers
              mountPath: /var/lib/docker/containers
              readOnly: true
            - name: fluent-bit-config
              mountPath: /opt/app-root/etc/
          securityContext:
            privileged: true
            runAsUser: 0
            runAsGroup: 0
      terminationGracePeriodSeconds: 10
      volumes:
        - name: varlog
          hostPath:
            path: /var/log
        - name: varlibdockercontainers
          hostPath:
            path: /var/lib/docker/containers
        - name: fluent-bit-config
          configMap:
            name: fluent-bit-config
      serviceAccountName: fluent-bit
      tolerations:
        - operator: Exists
