apiVersion: v1
kind: ConfigMap
metadata:
  name: cvat-pg-config
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-ccp"
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  CVAT_ENABLED: "{{ .Spec.ControlPlane.BaseConfig.AgentCustomTag }}"
  CNVRG_CVAT_POSTGRES_DB: "cvat"
  CNVRG_CVAT_POSTGRES_USER: "root"