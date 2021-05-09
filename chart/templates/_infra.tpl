{{- define "infra" }}
---
apiVersion: mlops.cnvrg.io/v1
kind: CnvrgInfra
metadata:
  name: cnvrg-infra
spec:
  clusterDomain: {{ .Values.clusterDomain }}
  infraNamespace: {{ template "spec.cnvrgNs" . }}
  {{- include "spec.registry" . | indent 2 }}
  {{- include "spec.storage" . | indent 2 }}
  {{- include "spec.infra_dbs" . | indent 2 }}
  {{- include "spec.sso" . | indent 2 }}
  {{- include "spec.logging_infra" . | indent 2 }}
  {{- include "spec.monitoring_infra" . | indent 2 }}
  {{- include "spec.networking_infra" . | indent 2 }}
  {{- include "spec.tenancy" . | indent 2 }}
---
{{- end }}