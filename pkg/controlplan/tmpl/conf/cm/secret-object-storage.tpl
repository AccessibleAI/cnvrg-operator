apiVersion: v1
kind: Secret
metadata:
  name: cp-object-storage
  namespace: {{ .CnvrgNs }}
data:
  SECRET_KEY_BASE: {{ .ControlPlan.Conf.SecretKeyBase | b64enc }}
  STS_IV: {{ .ControlPlan.Conf.StsIv | b64enc }}
  STS_KEY: {{ .ControlPlan.Conf.StsKey | b64enc }}
  CNVRG_STORAGE_TYPE: {{ .ControlPlan.Conf.CnvrgStorageType | b64enc }}
  MINIO_SSE_MASTER_KEY: {{ .ControlPlan.Conf.MinioSseMasterKey | b64enc }}
  CNVRG_STORAGE_ENDPOINT: {{ .ControlPlan.Conf.CnvrgStorageEndpoint | b64enc }}
  ################## minio/aws storage configs  ###########################
  CNVRG_STORAGE_BUCKET: {{ .ControlPlan.Conf.CnvrgStorageBucket | b64enc }}
  CNVRG_STORAGE_ACCESS_KEY:  {{ .ControlPlan.Conf.CnvrgStorageAccessKey | b64enc }}
  CNVRG_STORAGE_SECRET_KEY: {{ .ControlPlan.Conf.CnvrgStorageSecretKey | b64enc }}
  CNVRG_STORAGE_REGION: {{ .ControlPlan.Conf.CnvrgStorageRegion | b64enc }}
  ################## azure #########################
  CNVRG_STORAGE_AZURE_ACCESS_KEY: {{ .ControlPlan.Conf.CnvrgStorageAzureAccessKey | b64enc }}
  CNVRG_STORAGE_AZURE_ACCOUNT_NAME: {{ .ControlPlan.Conf.CnvrgStorageAzureAccountName | b64enc }}
  CNVRG_STORAGE_AZURE_CONTAINER: {{ .ControlPlan.Conf.CnvrgStorageAzureContainer | b64enc }}
  ################## gcp ###########################
  CNVRG_STORAGE_KEYFILE: {{ printf "%s/%s" .ControlPlan.Conf.GcpKeyfileMountPath .ControlPlan.Conf.GcpKeyfileName | b64enc }}
  CNVRG_STORAGE_PROJECT: {{ .ControlPlan.Conf.CnvrgStorageProject | b64enc }}