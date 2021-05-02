apiVersion: v1
kind: Secret
metadata:
  name: cp-object-storage
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
data:
  SECRET_KEY_BASE: {{ .Spec.ControlPlane.ObjectStorage.SecretKeyBase | b64enc }}
  STS_IV: {{ .Spec.ControlPlane.ObjectStorage.StsIv | b64enc }}
  STS_KEY: {{ .Spec.ControlPlane.ObjectStorage.StsKey | b64enc }}
  CNVRG_STORAGE_TYPE: {{ .Spec.ControlPlane.ObjectStorage.CnvrgStorageType | b64enc }}
  MINIO_SSE_MASTER_KEY: {{ .Spec.ControlPlane.ObjectStorage.MinioSseMasterKey | b64enc }}
  CNVRG_STORAGE_ENDPOINT: {{ objectStorageUrl . | b64enc }}
  ################## minio/aws storage ObjectStorage  ###########################
  CNVRG_STORAGE_BUCKET: {{ .Spec.ControlPlane.ObjectStorage.CnvrgStorageBucket | b64enc }}
  CNVRG_STORAGE_ACCESS_KEY:  {{ .Spec.ControlPlane.ObjectStorage.CnvrgStorageAccessKey | b64enc }}
  CNVRG_STORAGE_SECRET_KEY: {{ .Spec.ControlPlane.ObjectStorage.CnvrgStorageSecretKey | b64enc }}
  CNVRG_STORAGE_REGION: {{ .Spec.ControlPlane.ObjectStorage.CnvrgStorageRegion | b64enc }}
  ################## azure #########################
  CNVRG_STORAGE_AZURE_ACCESS_KEY: {{ .Spec.ControlPlane.ObjectStorage.CnvrgStorageAzureAccessKey | b64enc }}
  CNVRG_STORAGE_AZURE_ACCOUNT_NAME: {{ .Spec.ControlPlane.ObjectStorage.CnvrgStorageAzureAccountName | b64enc }}
  CNVRG_STORAGE_AZURE_CONTAINER: {{ .Spec.ControlPlane.ObjectStorage.CnvrgStorageAzureContainer | b64enc }}
  ################## gcp ###########################
  CNVRG_STORAGE_KEYFILE: {{ printf "%s/%s" .Spec.ControlPlane.ObjectStorage.GcpKeyfileMountPath .Spec.ControlPlane.ObjectStorage.GcpKeyfileName | b64enc }}
  CNVRG_STORAGE_PROJECT: {{ .Spec.ControlPlane.ObjectStorage.CnvrgStorageProject | b64enc }}