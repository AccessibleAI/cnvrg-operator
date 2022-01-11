apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.ControlPlane.Cvat.Redis.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{.Spec.ControlPlane.Cvat.Redis.SvcName }}
    owner: cnvrg-control-plane
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  type: ClusterIP
  ports:
    - port: {{ .Spec.ControlPlane.Cvat.Redis.Port }}
      targetPort: {{ .Spec.ControlPlane.Cvat.Redis.Port }}
      protocol: TCP
      name: http
  selector:
    app: {{ .Spec.ControlPlane.Cvat.Redis.SvcName }}