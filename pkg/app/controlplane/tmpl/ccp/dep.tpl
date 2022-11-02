apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    control-plane: controller-manager
  name: cnvrg-ccp-operator-controller-manager
  namespace: {{ ns . }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
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
      control-plane: controller-manager
  template:
    metadata:
      labels:
        control-plane: controller-manager
    spec:
      priorityClassName: {{ .Spec.PriorityClass.AppClassRef }}
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        {{ .Spec.Tenancy.Key }}: "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - key: "{{ .Spec.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.Tenancy.Value }}"
          effect: "NoSchedule"
      {{- end }}
      containers:
        - args:
            - --health-probe-bind-address=:8081
            - --leader-elect
          command:
            - /manager
          env:
            - name: TENANT_NAMESPACE
              value: {{ ns . }}
            - name: USER_ACCESS_KEY
              valueFrom:
                secretKeyRef:
                  key: accesskey
                  name: {{ .Spec.ControlPlane.CnvrgClusterProvisionerOperator.AwsCredsRef }}
            - name: USER_SECRET_KEY
              valueFrom:
                secretKeyRef:
                  key: secretkey
                  name: {{ .Spec.ControlPlane.CnvrgClusterProvisionerOperator.AwsCredsRef }}
          image: {{  image .Spec.ImageHub .Spec.ControlPlane.CnvrgClusterProvisionerOperator.Image }}
          imagePullPolicy: Always
          livenessProbe:
            httpGet:
              path: /healthz
              port: 8081
            initialDelaySeconds: 15
            periodSeconds: 20
          name: manager
          readinessProbe:
            httpGet:
              path: /readyz
              port: 8081
            initialDelaySeconds: 5
            periodSeconds: 10
          resources:
            limits:
              cpu: {{ .Spec.ControlPlane.CnvrgClusterProvisionerOperator.Limits.Cpu }}
              memory: {{ .Spec.ControlPlane.CnvrgClusterProvisionerOperator.Limits.Memory }}
            requests:
              cpu: {{ .Spec.ControlPlane.CnvrgClusterProvisionerOperator.Requests.Cpu }}
              memory: {{ .Spec.ControlPlane.CnvrgClusterProvisionerOperator.Requests.Memory }}
          securityContext:
            allowPrivilegeEscalation: false
      serviceAccountName: cnvrg-ccp-operator-controller-manager
      terminationGracePeriodSeconds: 10