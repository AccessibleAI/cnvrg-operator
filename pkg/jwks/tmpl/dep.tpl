apiVersion: apps/v1
kind: Deployment
metadata:
  name: cnvrg-jwks
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: cnvrg-jwks
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cnvrg-jwks
  template:
    metadata:
      labels:
        app: cnvrg-jwks
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      serviceAccountName: cnvrg-jwks
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        "{{ .Spec.Tenancy.Key }}": "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - operator: "Exists"
      {{- end }}
      containers:
        - name: cnvrg-jwks
          command:
            - ??
          envFrom:
            - configMapRef:
                name: cnvrg-jwks
          env:
            - name: GIN_MODE
              value: release
          image: {{ image .Spec.ImageHub .Spec.Jwks.Image }}
          imagePullPolicy: Always
          resources:
            requests:
              cpu: 100
              memory: 100Mi
            limits:
              cpu: 1
              memory: 1Gi
          volumeMounts:
            - mountPath: /tmp/jwks-data
              name: jwks-data
      volumes:
        - name: jwks-data
          persistentVolumeClaim:
            claimName: jwks