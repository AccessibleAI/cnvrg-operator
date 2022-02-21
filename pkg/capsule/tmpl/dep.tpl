apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Spec.Capsule.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.Capsule.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ .Spec.Capsule.SvcName }}
  template:
    metadata:
      labels:
        app: {{ .Spec.Capsule.SvcName }}
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      serviceAccountName: cnvrg-capsule
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        "{{ .Spec.Tenancy.Key }}": "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - operator: "Exists"
      {{- end }}
      containers:
        - name: capsule
          command:
            - /opt/app-root/capsule
            - start
          env:
            - name: GIN_MODE
              value: release
          image: {{ image .Spec.ImageHub .Spec.Capsule.Image }}
          imagePullPolicy: Always
          resources:
            requests:
              cpu: "{{.Spec.Capsule.Requests.Cpu}}"
              memory: "{{.Spec.Capsule.Requests.Memory}}"
            limits:
              cpu: "{{.Spec.Capsule.Limits.Cpu}}"
              memory: "{{.Spec.Capsule.Limits.Memory}}"
          volumeMounts:
            - mountPath: /tmp/capsule-data
              name: capsule-data
      volumes:
        - name: capsule-data
          persistentVolumeClaim:
            claimName: capsule