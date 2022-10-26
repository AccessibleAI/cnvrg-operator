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
    {{- if isTrue .Spec.SSO.Enabled }}
    sso.cnvrg.io/enabled: "true"
    sso.cnvrg.io/skipAuthRoutes: ""
    sso.cnvrg.io/central: "{{ .Spec.SSO.Central.PublicUrl }}"
    sso.cnvrg.io/upstream: "{{ .Spec.Dbs.Es.Kibana.SvcName }}.{{ ns . }}.svc:{{.Spec.Dbs.Es.Kibana.Port}}"
    {{- end }}
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  name: {{ .Spec.Dbs.Es.Kibana.SvcName }}
  namespace: {{ .Namespace }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  ingressClassName: nginx
  rules:
    - host: "{{.Spec.Dbs.Es.Kibana.SvcName}}.{{ .Spec.ClusterDomain }}"
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: {{ .Spec.Dbs.Es.Kibana.SvcName }}
                port:
                 number: {{ .Spec.Dbs.Es.Kibana.Port }}