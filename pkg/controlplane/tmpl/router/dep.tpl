apiVersion: apps/v1
kind: Deployment
metadata:
  name: cnvrg-router
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.ControlPlane.CnvrgRouter.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  selector:
    matchLabels:
      app: {{ .Spec.ControlPlane.CnvrgRouter.SvcName }}
  replicas: 1
  template:
    metadata:
      labels:
        app: {{ .Spec.ControlPlane.CnvrgRouter.SvcName }}
    spec:
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        {{ .Spec.Tenancy.Key }}: "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - key: "{{ .Spec.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.Tenancy.Value }}"
          effect: "NoSchedule"
      {{- end }}
      serviceAccountName: cnvrg-control-plane
      containers:
        - name: nginx
          image: {{ .Spec.ControlPlane.CnvrgRouter.Image }}
          ports:
            - containerPort: 80
          volumeMounts:
            - mountPath: /etc/nginx
              readOnly: true
              name: routing-config
            - mountPath: /var/log/nginx
              name: log
      volumes:
        - name: routing-config
          configMap:
            name: routing-config
            items:
              - key: nginx.conf
                path: nginx.conf
        - name: log
          emptyDir: {}
