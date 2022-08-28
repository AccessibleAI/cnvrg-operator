apiVersion: v1
kind: Secret
metadata:
  name: cp-object-storage
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
  CNVRG_STORAGE_TYPE: {{ .Spec.ControlPlane.ObjectStorage.Type | toString | b64enc }}
  CNVRG_STORAGE_ENDPOINT: {{ objectStorageUrl . | b64enc }}
  ################## minio/aws storage ObjectStorage  ###########################
  CNVRG_STORAGE_BUCKET: {{ .Spec.ControlPlane.ObjectStorage.Bucket | b64enc }}
  {{- if and ( isTrue .Spec.Dbs.Minio.Enabled ) (eq .Spec.ControlPlane.ObjectStorage.Type "minio") (eq .Spec.ControlPlane.ObjectStorage.AccessKey "") (eq .Spec.ControlPlane.ObjectStorage.SecretKey "") }}
  CNVRG_STORAGE_ACCESS_KEY:  {{ randAlpha 20 | b64enc }}
  CNVRG_STORAGE_SECRET_KEY: {{ randAlpha 40 | b64enc }}
  {{- else }}
  CNVRG_STORAGE_ACCESS_KEY:  {{ .Spec.ControlPlane.ObjectStorage.AccessKey | b64enc }}
  CNVRG_STORAGE_SECRET_KEY: {{ .Spec.ControlPlane.ObjectStorage.SecretKey | b64enc }}
  {{- end }}

  CNVRG_STORAGE_REGION: {{ .Spec.ControlPlane.ObjectStorage.Region | b64enc }}
  ################## azure #########################
  CNVRG_STORAGE_AZURE_ACCESS_KEY: {{ .Spec.ControlPlane.ObjectStorage.AccessKey | b64enc }}
  CNVRG_STORAGE_AZURE_ACCOUNT_NAME: {{ .Spec.ControlPlane.ObjectStorage.AzureAccountName | b64enc }}
  CNVRG_STORAGE_AZURE_CONTAINER: {{ .Spec.ControlPlane.ObjectStorage.AzureContainer | b64enc }}
  ################## gcp ###########################
  CNVRG_STORAGE_KEYFILE: {{ "/opt/app-root/conf/gcp-keyfile/key.json"  | b64enc }}
  CNVRG_STORAGE_PROJECT: {{ .Spec.ControlPlane.ObjectStorage.GcpProject | b64enc }}
  CNVRG_STORAGE_GCP_SECRET_REF: {{ .Spec.ControlPlane.ObjectStorage.GcpSecretRef | b64enc }}