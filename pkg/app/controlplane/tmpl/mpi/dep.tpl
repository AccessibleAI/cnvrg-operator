apiVersion: apps/v1
kind: Deployment
metadata:
  name: mpi-operator
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: mpi-operator
    owner: cnvrg-control-plane
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mpi-operator
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: mpi-operator
        owner: cnvrg-control-plane
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
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
        resources:
          requests:
            cpu: {{.Spec.ControlPlane.Mpi.Requests.Cpu}}
            memory: {{.Spec.ControlPlane.Mpi.Requests.Memory}}
          limits:
            cpu: {{.Spec.ControlPlane.Mpi.Limits.Cpu}}
            memory: {{.Spec.ControlPlane.Mpi.Limits.Memory}}
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

