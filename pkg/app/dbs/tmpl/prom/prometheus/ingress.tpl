apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    nginx.ingress.kubernetes.io/proxy-send-timeout: 18000s
    nginx.ingress.kubernetes.io/proxy-read-timeout: 18000s
    nginx.ingress.kubernetes.io/proxy-body-size: 5G
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  name: prom
  namespace: {{ .Namespace }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  ingressClassName: nginx
  rules:
    - host: "prom.{{ .Spec.ClusterDomain }}"
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: prom
                port:
                 number: 9090