
apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Logging.Kibana.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.Logging.Kibana.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
  type: NodePort
  {{- end }}
  selector:
    app: {{ .Spec.Logging.Kibana.SvcName }}
  ports:
    - port: {{ .Spec.Logging.Kibana.Port }}
      protocol: TCP
      {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
      nodePort: {{ .Spec.Logging.Kibana.NodePort }}
      {{- end }}