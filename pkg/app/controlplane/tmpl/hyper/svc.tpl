apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.ControlPlane.Hyper.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.ControlPlane.Hyper.SvcName }}
    owner: cnvrg-control-plane
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  ports:
    - port: {{ .Spec.ControlPlane.Hyper.Port }}
  selector:
    app: {{ .Spec.ControlPlane.Hyper.SvcName }}