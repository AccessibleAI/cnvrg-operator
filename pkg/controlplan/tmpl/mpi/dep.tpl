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
      serviceAccountName: mpi-operator
      imagePullSecrets:
        - name: {{ .Spec.ControlPlan.Mpi.Registry.Name }}
      containers:
      - name: mpi-operator
        imagePullPolicy: Always
        image: {{ .Spec.ControlPlan.Mpi.Image }}
        args:
        - -alsologtostderr
        - --kubectl-delivery-image
        - {{ .Spec.ControlPlan.Mpi.KubectlDeliveryImage }}
        - --lock-namespace
        - {{ ns . }}
        - --namespace
        - {{ ns . }}
        {{- range $extraArgName, $extraArgValue := .Spec.ControlPlan.Mpi.ExtraArgs }}
        - "{{ $extraArgName }}"
        - "{{ $extraArgValue }}"
        {{- end }}

