apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: {{ ns . }}
  name: istio-operator
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  replicas: 1
  selector:
    matchLabels:
      name: istio-operator
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        name: istio-operator
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      imagePullSecrets:
        - name: {{ .Spec.Registry.Name }}
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      serviceAccountName: istio-operator
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        {{ .Spec.Tenancy.Key }}: "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - operator: "Exists"
      {{- end }}
      containers:
        - name: istio-operator
          image: {{ image .Spec.ImageHub .Spec.Networking.Istio.OperatorImage }}
          command:
            - operator
            - server
          securityContext:
            allowPrivilegeEscalation: false
            capabilities:
              drop:
                - ALL
            privileged: false
            readOnlyRootFilesystem: true
            runAsGroup: 1337
            runAsUser: 1337
            runAsNonRoot: true
          imagePullPolicy: IfNotPresent
          resources:
            limits:
              cpu: 1000m
              memory: 2048Mi
            requests:
              cpu: 50m
              memory: 128Mi
          env:
            - name: WATCH_NAMESPACE
              value:  {{ ns . }}
            - name: LEADER_ELECTION_NAMESPACE
              value:  {{ ns . }}
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: OPERATOR_NAME
              value: "istio-operator"
            - name: WAIT_FOR_RESOURCES_TIMEOUT
              value: "300s"
            - name: REVISION
              value: ""