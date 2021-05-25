apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: cnvrg-fluentbit
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: cnvrg-fluentbit
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-fluentbit"
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  selector:
    matchLabels:
      app: cnvrg-fluentbit
  template:
    metadata:
      labels:
        app: cnvrg-fluentbit
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        {{ .Spec.Tenancy.Key }}: {{ .Spec.Tenancy.Value }}
        {{- range $key, $val := .Spec.Logging.Fluentbit.NodeSelector }}
        {{ $key }}: {{ $val }}
      {{- end }}
      {{- else if (gt (len .Spec.Logging.Fluentbit.NodeSelector) 0) }}
      nodeSelector:
        {{- range $key, $val := .Spec.Logging.Fluentbit.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      {{- end }}
      containers:
        - name: fluentbit
          image: {{.Spec.ImageHub }}/{{ .Spec.Logging.Fluentbit.Image }}
          imagePullPolicy: Always
          command:
            - /bin/bash
            - -c
            - /opt/app-root/fluentbit -c /opt/app-root/etc/fluent-bit.conf
          ports:
            - containerPort: 2020
          volumeMounts:
            {{- range $name, $path := .Spec.Logging.Fluentbit.LogsMounts }}
            - name: {{ $name }}
              mountPath: {{ $path }}
              readOnly: true
            {{- end }}
            - name: fluent-bit-config
              mountPath: /opt/app-root/etc/
          securityContext:
            privileged: true
            runAsUser: 0
            runAsGroup: 0
      terminationGracePeriodSeconds: 10
      volumes:
        {{- range $name, $path := .Spec.Logging.Fluentbit.LogsMounts }}
        - name: {{ $name }}
          hostPath:
            path: {{ $path}}
        {{- end }}
        - name: fluent-bit-config
          configMap:
            name: fluent-bit-config
      serviceAccountName: fluent-bit
      tolerations:
        - operator: Exists
