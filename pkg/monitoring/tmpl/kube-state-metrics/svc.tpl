apiVersion: v1
kind: Service
metadata:
  labels:
    app.kubernetes.io/name: kube-state-metrics
    app.kubernetes.io/version: v1.9.7
    owner: cnvrg-control-plane
  name: kube-state-metrics
  namespace: {{ ns . }}
spec:
  clusterIP: None
  ports:
  - name: https-main
    port: 8443
    targetPort: https-main
  - name: https-self
    port: 9443
    targetPort: https-self
  selector:
    app.kubernetes.io/name: kube-state-metrics
