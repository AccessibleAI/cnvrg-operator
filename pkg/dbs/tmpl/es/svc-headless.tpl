---
kind: Service
apiVersion: v1
metadata:
  name: {{ .Spec.Dbs.Es.SvcName }}-headless
  namespace: {{ ns . }}
  annotations:
    service.alpha.kubernetes.io/tolerate-unready-endpoints: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.Dbs.Es.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  clusterIP: None # This is needed for statefulset hostnames like elasticsearch-0 to resolve
  # Create endpoints also if the related pod isn't ready
  publishNotReadyAddresses: true
  selector:
    app: {{ .Spec.Dbs.Es.SvcName }}
  ports:
  - name: http
    port: 9200
  - name: transport
    port: 9300