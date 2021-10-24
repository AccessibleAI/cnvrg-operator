apiVersion: scheduling.k8s.io/v1
kind: PriorityClass
metadata:
  name: {{ .Spec.CnvrgAppPriorityClass.Name }}
value: {{ .Spec.CnvrgAppPriorityClass.Value }}
globalDefault: false
description: {{ .Spec.CnvrgAppPriorityClass.Description }}