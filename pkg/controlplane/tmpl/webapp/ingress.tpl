apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Spec.ControlPlane.WebApp.SvcName }}
  namespace: {{ ns . }}
spec:
  rules:
    - host: "{{.Spec.ControlPlane.WebApp.SvcName}}.{{ .Spec.ClusterDomain }}"
      http:
        paths:
          - path: /
            backend:
              service:
                name: {{ .Spec.ControlPlane.WebApp.SvcName }}
                port:
                  number:  {{ .Spec.ControlPlane.WebApp.Port }}