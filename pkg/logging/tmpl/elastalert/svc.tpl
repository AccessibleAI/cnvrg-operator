
apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Logging.Elastalert.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: {{ .Spec.Logging.Elastalert.SvcName }}
    owner: cnvrg-control-plane
spec:
  {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
  type: NodePort
  {{- end }}
  ports:
    - port: {{ .Spec.Logging.Elastalert.Port }}
      protocol: TCP
      targetPort: 3030
      {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
      nodePort: {{ .Spec.Logging.Elastalert.NodePort }}
      {{- end }}
  selector:
    app: {{ .Spec.Logging.Elastalert.SvcName }}