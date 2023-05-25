apiVersion: v1
kind: Secret
metadata:
  name: {{ .PublicKeySecret }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "false"
  labels:
    domainId: {{ .DomainID }}
data:
  CNVRG_PKI_PUBLIC_KEY: {{ .PublicKey | b64enc }}
