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
  CNVRG_CVAT_POSTGRES_DBNAME: "cvat"
  CNVRG_CVAT_POSTGRES_USER: "root"
  CNVRG_CVAT_REDIS_HOST: {{ .Spec.ControlPlane.Cvat.Redis.SvcName }}
  CNVRG_CVAT_POSTGRES_HOST: {{ .Spec.ControlPlane.Cvat.Pg.SvcName }}
