apiVersion: v1
kind: ServiceAccount
metadata:
  name: mpi-operator
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane