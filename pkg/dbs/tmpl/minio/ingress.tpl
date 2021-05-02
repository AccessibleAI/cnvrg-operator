apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Spec.Dbs.Minio.SvcName }}
  namespace: {{ ns . }}
spec:
  rules:
    - host: "{{.Spec.Dbs.Minio.SvcName}}.{{ .Spec.ClusterDomain }}"
      http:
        paths:
          - path: /
            backend:
              service:
                name: {{ .Spec.Dbs.Minio.SvcName }}
                port:
                  number: {{ .Spec.Dbs.Minio.Port }}