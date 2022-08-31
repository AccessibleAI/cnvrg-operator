
apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Dbs.Es.Elastalert.SvcName }}
  namespace: {{ ns . }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.Dbs.Es.Elastalert.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
  type: NodePort
  {{- end }}
  ports:
    - port: {{ .Spec.Dbs.Es.Elastalert.Port }}
      protocol: TCP
      targetPort: 8080
      {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
      nodePort: {{ .Spec.Dbs.Es.Elastalert.NodePort }}
      {{- end }}
  selector:
    app: {{ .Spec.Dbs.Es.Elastalert.SvcName }}