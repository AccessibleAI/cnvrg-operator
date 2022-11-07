
apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Dbs.Es.Kibana.SvcName }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.Dbs.Es.Kibana.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
  type: NodePort
  {{- end }}
  selector:
    app: {{ .Spec.Dbs.Es.Kibana.SvcName }}
  ports:
    - name: http
      port: {{ .Spec.Dbs.Es.Kibana.Port }}
      targetPort: {{ .Spec.Dbs.Es.Kibana.Port }}
      protocol: TCP
      {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
      nodePort: {{ .Spec.Dbs.Es.Kibana.NodePort }}
      {{- end }}