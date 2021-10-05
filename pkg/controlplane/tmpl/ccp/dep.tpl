apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    control-plane: controller-manager
  name: cnvrg-ccp-operator-controller-manager
  namespace: {{ ns . }}
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
      containers:
        - args:
            - --secure-listen-address=0.0.0.0:8443
            - --upstream=http://127.0.0.1:8080/
            - --logtostderr=true
            - --v=10
          image: gcr.io/kubebuilder/kube-rbac-proxy:v0.8.0
          name: kube-rbac-proxy
          ports:
            - containerPort: 8443
              name: https
        - args:
            - --health-probe-bind-address=:8081
            - --metrics-bind-address=127.0.0.1:8080
            - --leader-elect
          command:
            - /manager
          env:
            - name: USER_ACCESS_KEY
              valueFrom:
                secretKeyRef:
                  key: accesskey
                  name: cnvrg-aws-sec
            - name: USER_SECRET_KEY
              valueFrom:
                secretKeyRef:
                  key: secretkey
                  name: cnvrg-aws-sec
          image: {{ .Spec.ControlPlane.CnvrgClusterProvisionerOperator.Image }}
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