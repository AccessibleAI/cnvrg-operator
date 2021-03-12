apiVersion: v1
kind: Secret
metadata:
  name: cp-object-storage
  namespace: {{ .CnvrgNs }}
data:
  SECRET_KEY_BASE: {{ .ControlPlan.ObjectStorage.SecretKeyBase | b64enc }}
  STS_IV: {{ .ControlPlan.ObjectStorage.StsIv | b64enc }}
  STS_KEY: {{ .ControlPlan.ObjectStorage.StsKey | b64enc }}
  CNVRG_STORAGE_TYPE: {{ .ControlPlan.ObjectStorage.CnvrgStorageType | b64enc }}
  MINIO_SSE_MASTER_KEY: {{ .ControlPlan.ObjectStorage.MinioSseMasterKey | b64enc }}
  CNVRG_STORAGE_ENDPOINT: {{ .ControlPlan.ObjectStorage.CnvrgStorageEndpoint | b64enc }}
  ################## minio/aws storage ObjectStorageigs  ###########################
  CNVRG_STORAGE_BUCKET: {{ .ControlPlan.ObjectStorage.CnvrgStorageBucket | b64enc }}
  CNVRG_STORAGE_ACCESS_KEY:  {{ .ControlPlan.ObjectStorage.CnvrgStorageAccessKey | b64enc }}
  CNVRG_STORAGE_SECRET_KEY: {{ .ControlPlan.ObjectStorage.CnvrgStorageSecretKey | b64enc }}
  CNVRG_STORAGE_REGION: {{ .ControlPlan.ObjectStorage.CnvrgStorageRegion | b64enc }}
  ################## azure #########################
  CNVRG_STORAGE_AZURE_ACCESS_KEY: {{ .ControlPlan.ObjectStorage.CnvrgStorageAzureAccessKey | b64enc }}
  CNVRG_STORAGE_AZURE_ACCOUNT_NAME: {{ .ControlPlan.ObjectStorage.CnvrgStorageAzureAccountName | b64enc }}
  CNVRG_STORAGE_AZURE_CONTAINER: {{ .ControlPlan.ObjectStorage.CnvrgStorageAzureContainer | b64enc }}
  ################## gcp ###########################
  CNVRG_STORAGE_KEYFILE: {{ printf "%s/%s" .ControlPlan.ObjectStorage.GcpKeyfileMountPath .ControlPlan.ObjectStorage.GcpKeyfileName | b64enc }}
  CNVRG_STORAGE_PROJECT: {{ .ControlPlan.ObjectStorage.CnvrgStorageProject | b64enc }}