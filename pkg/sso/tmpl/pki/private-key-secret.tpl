apiVersion: v1
kind: Secret
metadata:
  name: {{ .PrivateKeySecret }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "false"
data:
  CNVRG_PKI_PRIVATE_KEY: {{ .PrivateKey | b64enc }}
