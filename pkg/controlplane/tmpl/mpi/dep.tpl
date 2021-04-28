apiVersion: apps/v1
kind: Deployment
metadata:
  name: mpi-operator
  namespace: {{ ns . }}
  labels:
    app: mpi-operator
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mpi-operator
  template:
    metadata:
      labels:
        app: mpi-operator
    spec:
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        "{{ .Spec.Tenancy.Key }}": "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - key: "{{ .Spec.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.Tenancy.Value }}"
          effect: "NoSchedule"
      {{- end }}
      serviceAccountName: mpi-operator
      imagePullSecrets:
        - name: {{ .Spec.ControlPlane.Mpi.Registry.Name }}
      containers:
      - name: mpi-operator
        imagePullPolicy: Always
        image: {{ .Spec.ControlPlane.Mpi.Image }}
        args:
        - -alsologtostderr
        - --kubectl-delivery-image
        - {{ .Spec.ControlPlane.Mpi.KubectlDeliveryImage }}
        - --lock-namespace
        - {{ ns . }}
        - --namespace
        - {{ ns . }}
        {{- range $extraArgName, $extraArgValue := .Spec.ControlPlane.Mpi.ExtraArgs }}
        - "{{ $extraArgName }}"
        - "{{ $extraArgValue }}"
        {{- end }}

