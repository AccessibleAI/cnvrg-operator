apiVersion: v1
kind: Service
metadata:
  labels:
    control-plane: controller-manager
  name: cnvrg-ccp-operator-controller-manager-metrics-service
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
  ports:
    - name: https
      port: 8443
      targetPort: https
  selector:
    control-plane: controller-manager