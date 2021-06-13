apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.ControlPlane.CnvrgRouter.SvcName }}
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
  ports:
  - port: 80
  selector:
    app: {{ .Spec.ControlPlane.CnvrgRouter.SvcName }}