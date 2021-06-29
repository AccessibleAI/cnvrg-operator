apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Monitoring.Alertmanager.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: cnvrg-alertmanager
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
  type: NodePort
  {{- end }}
  sessionAffinity: ClientIP
  ports:
    - name: web
      port: {{ .Spec.Monitoring.Alertmanager.Port }}
      targetPort: web
      {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
      nodePort: {{ .Spec.Monitoring.Alertmanager.NodePort }}
      {{- end }}
  selector:
    app: alertmanager

