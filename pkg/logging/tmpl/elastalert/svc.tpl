
apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Logging.Elastalert.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: {{ .Spec.Logging.Elastalert.SvcName }}
spec:
  {{- if eq .Spec.Networking.Ingress.IngressType "nodeport" }}
  type: NodePort
  {{- end }}
  ports:
    - port: {{ .Spec.Logging.Elastalert.Port }}
      protocol: TCP
      targetPort: {{ .Spec.Logging.Elastalert.ContainerPort}}
      {{- if eq .Spec.Networking.Ingress.IngressType "nodeport" }}
      nodePort: {{ .Spec.Logging.Elastalert.NodePort }}
      {{- end }}
  selector:
    app: {{ .Spec.Logging.Elastalert.SvcName }}