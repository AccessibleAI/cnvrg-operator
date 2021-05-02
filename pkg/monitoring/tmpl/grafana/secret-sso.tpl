apiVersion: v1
kind: Secret
metadata:
  name: cp-sso-grafana
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
data:
  conf: {{ oauthProxyConfig . .Spec.Monitoring.Grafana.SvcName | b64enc }}

