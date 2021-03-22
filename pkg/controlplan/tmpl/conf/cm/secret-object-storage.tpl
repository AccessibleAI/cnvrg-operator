apiVersion: v1
kind: Secret
metadata:
  name: cp-object-storage
  namespace: {{ ns . }}
data:
  SECRET_KEY_BASE: {{ .Spec.ControlPlan.ObjectStorage.SecretKeyBase | b64enc }}
  STS_IV: {{ .Spec.ControlPlan.ObjectStorage.StsIv | b64enc }}
  STS_KEY: {{ .Spec.ControlPlan.ObjectStorage.StsKey | b64enc }}
  CNVRG_STORAGE_TYPE: {{ .Spec.ControlPlan.ObjectStorage.CnvrgStorageType | b64enc }}
  MINIO_SSE_MASTER_KEY: {{ .Spec.ControlPlan.ObjectStorage.MinioSseMasterKey | b64enc }}
  CNVRG_STORAGE_ENDPOINT: {{ objectStorageUrl . | b64enc }}
  ################## minio/aws storage ObjectStorage  ###########################
  CNVRG_STORAGE_BUCKET: {{ .Spec.ControlPlan.ObjectStorage.CnvrgStorageBucket | b64enc }}
  CNVRG_STORAGE_ACCESS_KEY:  {{ .Spec.ControlPlan.ObjectStorage.CnvrgStorageAccessKey | b64enc }}
  CNVRG_STORAGE_SECRET_KEY: {{ .Spec.ControlPlan.ObjectStorage.CnvrgStorageSecretKey | b64enc }}
  CNVRG_STORAGE_REGION: {{ .Spec.ControlPlan.ObjectStorage.CnvrgStorageRegion | b64enc }}
  ################## azure #########################
  CNVRG_STORAGE_AZURE_ACCESS_KEY: {{ .Spec.ControlPlan.ObjectStorage.CnvrgStorageAzureAccessKey | b64enc }}
  CNVRG_STORAGE_AZURE_ACCOUNT_NAME: {{ .Spec.ControlPlan.ObjectStorage.CnvrgStorageAzureAccountName | b64enc }}
  CNVRG_STORAGE_AZURE_CONTAINER: {{ .Spec.ControlPlan.ObjectStorage.CnvrgStorageAzureContainer | b64enc }}
  ################## gcp ###########################
  CNVRG_STORAGE_KEYFILE: {{ printf "%s/%s" .Spec.ControlPlan.ObjectStorage.GcpKeyfileMountPath .Spec.ControlPlan.ObjectStorage.GcpKeyfileName | b64enc }}
  CNVRG_STORAGE_PROJECT: {{ .Spec.ControlPlan.ObjectStorage.CnvrgStorageProject | b64enc }}