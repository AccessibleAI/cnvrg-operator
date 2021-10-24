apiVersion: scheduling.k8s.io/v1
kind: PriorityClass
metadata:
  name: {{ .Spec.CnvrgJobPriorityClass.Name }}
value: {{ .Spec.CnvrgJobPriorityClass.Value }}
globalDefault: false
description: {{ .Spec.CnvrgJobPriorityClass.Description }}