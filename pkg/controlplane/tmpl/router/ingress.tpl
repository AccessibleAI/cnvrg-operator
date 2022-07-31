apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/proxy-send-timeout: 18000s
    nginx.ingress.kubernetes.io/proxy-read-timeout: 18000s
    nginx.ingress.kubernetes.io/proxy-body-size: 5G
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
  name: {{ .Spec.ControlPlane.CnvrgRouter.SvcName }}
  namespace: {{ ns . }}
spec:
  rules:
  - host: "{{.Spec.ControlPlane.CnvrgRouter.SvcName}}.{{ .Spec.ClusterDomain }}"
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: {{ .Spec.ControlPlane.CnvrgRouter.SvcName }}
            port:
              number: {{ .Spec.ControlPlane.CnvrgRouter.Port }}
