apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Dbs.Minio.SvcName }}
  namespace: {{.Namespace  }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.Dbs.Minio.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
  type: NodePort
  {{- end }}
  ports:
  - name: http
    port: {{ .Spec.Dbs.Minio.Port }}
    targetPort: {{ .Spec.Dbs.Minio.Port }}
    {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
    nodePort: {{ .Spec.Dbs.Minio.NodePort }}
    {{- end }}
  selector:
    app: {{ .Spec.Dbs.Minio.SvcName }}
