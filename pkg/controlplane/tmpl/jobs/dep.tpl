apiVersion: apps/v1
kind: Deployment
metadata:
  name: workflowmap-controller-manager
  namespace: {{ ns . }}
  annotations:
      mlops.cnvrg.io/default-loader: "true"
      mlops.cnvrg.io/own: "true"
      mlops.cnvrg.io/updatable: "true"
      {{- range $k, $v := .Spec.Annotations }}
      {{$k}}: "{{$v}}"
      {{- end }}
    labels:
      control-plane: controller-manager
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
      annotations:
        kubectl.kubernetes.io/default-container: manager
      labels:
        app: workflowmap-operator
        control-plane: controller-manager
    spec:
      containers:
      - args:
        - --secure-listen-address=0.0.0.0:8443
        - --upstream=http://127.0.0.1:8080/
        - --logtostderr=true
        - --v=0
        image: gcr.io/kubebuilder/kube-rbac-proxy:v0.11.0
        name: kube-rbac-proxy
        ports:
        - containerPort: 8443
          name: https
          protocol: TCP
        resources:
          limits:
            cpu: 500m
            memory: 128Mi
          requests:
            cpu: 5m
            memory: 64Mi
      - args:
        - --health-probe-bind-address=:8081
        - --metrics-bind-address=127.0.0.1:8080
        - --leader-elect
        command:
        - /manager
        env:
        - name: OPERATOR_ENV
          value: dev
        - name: AGENT_DEFAULT_PORT
          value: "4000"
        - name: MAX_CO_RECONCILES
          value: "19"
        - name: COMMANDS_REST_TIME
          value: 5s
        - name: FORMATION_REST_TIME
          value: 10s
        - name: GRPC_TIMEOUT
          value: 5s
        image: {{  image .Spec.ImageHub .Spec.ControlPlane.CnvrgJobsOperator.Image }}
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
            cpu: {{ .Spec.ControlPlane.CnvrgJobsOperator.Limits.Cpu }}
            memory: {{ .Spec.ControlPlane.CnvrgJobsOperator.Limits.Memory }}
          requests:
            cpu: {{ .Spec.ControlPlane.CnvrgJobsOperator.Requests.Cpu }}
            memory: {{ .Spec.ControlPlane.CnvrgJobsOperator.Requests.Memory }}
        securityContext:
          allowPrivilegeEscalation: false
      securityContext:
        runAsNonRoot: true
      serviceAccountName: workflowmap-controller-manager
      terminationGracePeriodSeconds: 10
