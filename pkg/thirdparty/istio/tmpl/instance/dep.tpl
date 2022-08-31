apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: {{ .Namespace }}
  name: istio-operator
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "false"
    mlops.cnvrg.io/updatable: "false"
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
        - name: {{ .Spec.Registry.Name }}
      serviceAccountName: istio-operator
      containers:
        - name: istio-operator
          image: {{ image .Spec.ImageHub .Spec.Istio.OperatorImage }}
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
              value:  {{ .Namespace }}
            - name: LEADER_ELECTION_NAMESPACE
              value:  {{ .Namespace }}
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