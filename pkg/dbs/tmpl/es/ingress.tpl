apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Spec.Dbs.Es.SvcName }}
  namespace: {{ ns . }}
spec:
  rules:
    - host: "{{.Spec.Dbs.Es.SvcName}}.{{ .Spec.ClusterDomain }}"
      http:
        paths:
          - path: /
            backend:
              service:
                name: {{ .Spec.Dbs.Es.SvcName }}
                port:
                  number:  {{ .Spec.Dbs.Es.Port }}