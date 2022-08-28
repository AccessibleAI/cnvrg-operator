
apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Logging.Elastalert.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.Logging.Elastalert.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
  type: NodePort
  {{- end }}
  ports:
    - port: {{ .Spec.Logging.Elastalert.Port }}
      protocol: TCP
      targetPort: 8080
      {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
      nodePort: {{ .Spec.Logging.Elastalert.NodePort }}
      {{- end }}
  selector:
    app: {{ .Spec.Logging.Elastalert.SvcName }}