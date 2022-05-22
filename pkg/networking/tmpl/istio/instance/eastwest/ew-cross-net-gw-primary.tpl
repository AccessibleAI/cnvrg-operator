apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: cnvrg-cross-network-gateway
  namespace:  {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
      {{$k}}: "{{$v}}"
      {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
      {{$k}}: "{{$v}}"
      {{- end }}
spec:
  selector:
    istio: cnvrg-eastwestgateway
  servers:
    - port:
        number: 15443
        name: tls
        protocol: TLS
      tls:
        mode: AUTO_PASSTHROUGH
      hosts:
        - "*.local"