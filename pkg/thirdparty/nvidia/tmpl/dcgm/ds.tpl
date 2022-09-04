apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: dcgm-exporter
  namespace: {{ ns . }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "false"
    mlops.cnvrg.io/updatable: "false"
  labels:
    app: "dcgm-exporter"
spec:
  selector:
    matchLabels:
      app: "dcgm-exporter"
  template:
    metadata:
      labels:
        app: "dcgm-exporter"
    spec:
      serviceAccountName: nvidia
      volumes:
        - name: "pod-gpu-resources"
          hostPath:
            path: "/var/lib/kubelet/pod-resources"
      tolerations:
        - operator: Exists
      nodeSelector:
        {{.Spec.Nvidia.NodeSelector.Key}}: {{.Spec.Nvidia.NodeSelector.Value}}
      containers:
        - name: exporter
          securityContext:
            capabilities:
              add:
                - SYS_ADMIN
            runAsNonRoot: false
            runAsUser: 0
          image: {{ image .Spec.ImageHub .Spec.Nvidia.DevicePlugin.Image }}
          imagePullPolicy: "IfNotPresent"
          args:
            - -f
            - /etc/dcgm-exporter/dcp-metrics-included.csv
          env:
            - name: "DCGM_EXPORTER_KUBERNETES"
              value: "true"
            - name: "DCGM_EXPORTER_LISTEN"
              value: ":9400"
            - name: "DCGM_EXPORTER_INTERVAL"
              value: "5000"
          ports:
            - name: "metrics"
              containerPort: 9400
          volumeMounts:
            - name: "pod-gpu-resources"
              readOnly: true
              mountPath: "/var/lib/kubelet/pod-resources"
          resources:
            requests:
              cpu: 100m
              memory: 100Mi
            limits:
              cpu: 500m
              memory: 1Gi
          livenessProbe:
            httpGet:
              path: /health
              port: 9400
            initialDelaySeconds: 15
            periodSeconds: 5
          readinessProbe:
            httpGet:
              path: /health
              port: 9400
            initialDelaySeconds: 15