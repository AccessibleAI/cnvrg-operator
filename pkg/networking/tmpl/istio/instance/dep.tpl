apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: {{ .CnvrgNs }}
  name: istio-operator
spec:
  replicas: 1
  selector:
    matchLabels:
      name: istio-operator
  template:
    metadata:
      labels:
        name: istio-operator
    spec:
      imagePullSecrets:
        - name: {{ .ControlPlan.Registry.Name }}
      serviceAccountName: istio-operator
      {{- if eq .ControlPlan.Tenancy.Enabled "true" }}
      nodeSelector:
        {{ .ControlPlan.Tenancy.Key }}: "{{ .ControlPlan.Tenancy.Value }}"
      {{- end }}
      tolerations:
        - key: "{{ .ControlPlan.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .ControlPlan.Tenancy.Value }}"
          effect: "NoSchedule"
      containers:
        - name: istio-operator
          image: {{ .Networking.Istio.OperatorImage }}
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
              cpu: 200m
              memory: 256Mi
            requests:
              cpu: 50m
              memory: 128Mi
          env:
            - name: WATCH_NAMESPACE
              value:  {{ .CnvrgNs }}
            - name: LEADER_ELECTION_NAMESPACE
              value:  {{ .CnvrgNs }}
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