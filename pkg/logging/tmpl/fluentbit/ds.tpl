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
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      nodeSelector:
        {{- range $key, $val := .Spec.Logging.Fluentbit.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      containers:
        - name: fluentbit
          image: {{ image .Spec.ImageHub .Spec.Logging.Fluentbit.Image }}
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
              readOnly: false
            {{- end }}
            - name: fluent-bit-config
              mountPath: /opt/app-root/etc/
          resources:
            requests:
              cpu: {{.Spec.Logging.Fluentbit.Requests.Cpu }}
              memory: {{.Spec.Logging.Fluentbit.Requests.Memory }}
            limits:
              cpu: {{.Spec.Logging.Fluentbit.Limits.Cpu }}
              memory: {{.Spec.Logging.Fluentbit.Limits.Memory }}                              
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
      serviceAccountName: cnvrg-fluentbit
      tolerations:
        - operator: Exists
