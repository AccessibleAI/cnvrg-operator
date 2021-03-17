
apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Logging.Elastalert.SvcName }}
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Spec.Logging.Elastalert.SvcName }}
spec:
  {{- if eq .Spec.Ingress.IngressType "nodeport" }}
  type: NodePort
  {{- end }}
  ports:
    - port: {{ .Spec.Logging.Elastalert.Port }}
      protocol: TCP
      targetPort: {{ .Spec.Logging.Elastalert.ContainerPort}}
      {{- if eq .Spec.Ingress.IngressType "nodeport" }}
      nodePort: {{ .Spec.Logging.Elastalert.NodePort }}
      {{- end }}
  selector:
    app: {{ .Spec.Logging.Elastalert.SvcName }}