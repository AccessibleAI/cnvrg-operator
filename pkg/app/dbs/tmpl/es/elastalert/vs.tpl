apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Spec.Dbs.Es.Elastalert.SvcName }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  hosts:
    - "{{.Spec.Dbs.Es.Elastalert.SvcName}}{{.Spec.Networking.ClusterDomainPrefix.Prefix}}.{{ .Spec.ClusterDomain }}"
  gateways:
  - {{ .Spec.Networking.Ingress.IstioGwName}}
  http:
  - retries:
      attempts: {{ .Spec.Networking.Ingress.RetriesAttempts }}
      perTryTimeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
    timeout: {{ .Spec.Networking.Ingress.Timeout }}
    route:
    - destination:
        port:
          number: {{ .Spec.Dbs.Es.Elastalert.Port }}
        host: "{{ .Spec.Dbs.Es.Elastalert.SvcName }}.{{ .Namespace }}.svc.{{ .Spec.ClusterInternalDomain }}"