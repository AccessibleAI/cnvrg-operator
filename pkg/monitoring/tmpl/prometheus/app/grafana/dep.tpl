apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: grafana
  name: grafana
  namespace: {{ ns . }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: grafana
  template:
    metadata:
      labels:
        app: grafana
    spec:
      containers:
      - image: {{ .Spec.Grafana.Image }}
        name: grafana
        env:
          - name: GF_AUTH_BASIC_ENABLED
            value: "false"
          - name: GF_AUTH_ANONYMOUS_ENABLED
            value: "true"
          - name: GF_AUTH_ANONYMOUS_ORG_ROLE
            value: Admin
          - name: GF_SECURITY_ALLOW_EMBEDDING
            value: "true"
        ports:
        - containerPort: 3000
          name: http
        readinessProbe:
          httpGet:
            path: /api/health
            port: http
        resources:
          limits:
            cpu: 200m
            memory: 200Mi
          requests:
            cpu: 100m
            memory: 100Mi
        volumeMounts:
        - mountPath: /var/lib/grafana
          name: grafana-storage
          readOnly: false
        - mountPath: /etc/grafana/provisioning/datasources
          name: grafana-datasources
          readOnly: false
        - mountPath: /etc/grafana/provisioning/dashboards
          name: grafana-dashboards
          readOnly: false
        - mountPath: /definitions/0/k8s-resources-namespace
          name: k8s-resources-namespace
          readOnly: false
        - mountPath: /definitions/0/k8s-resources-pod
          name: k8s-resources-pod
          readOnly: false
        - mountPath: /definitions/0/k8s-resources-workload
          name: k8s-resources-workload
          readOnly: false
        - mountPath: /definitions/0/k8s-resources-workloads-namespace
          name: k8s-resources-workloads-namespace
          readOnly: false
        - mountPath: /definitions/0/namespace-by-pod
          name: namespace-by-pod
          readOnly: false
        - mountPath: /definitions/0/namespace-by-workload
          name: namespace-by-workload
          readOnly: false
        - mountPath: /definitions/0/persistentvolumesusage
          name: persistentvolumesusage
          readOnly: false
        - mountPath: /definitions/0/pod-total
          name: pod-total
          readOnly: false
        - mountPath: /definitions/0/statefulset
          name: statefulset
          readOnly: false
        - mountPath: /definitions/0/workload-total
          name: workload-total
          readOnly: false
      nodeSelector:
        beta.kubernetes.io/os: linux
      securityContext:
        runAsNonRoot: true
        runAsUser: 65534
      serviceAccountName: grafana
      volumes:
      - emptyDir: {}
        name: grafana-storage
      - name: grafana-datasources
        secret:
          secretName: grafana-datasources
      - configMap:
          name: grafana-dashboards
        name: grafana-dashboards
      - configMap:
          name: k8s-resources-namespace
        name: k8s-resources-namespace
      - configMap:
          name: k8s-resources-pod
        name: k8s-resources-pod
      - configMap:
          name: k8s-resources-workload
        name: k8s-resources-workload
      - configMap:
          name: k8s-resources-workloads-namespace
        name: k8s-resources-workloads-namespace
      - configMap:
          name: namespace-by-pod
        name: namespace-by-pod
      - configMap:
          name: namespace-by-workload
        name: namespace-by-workload
      - configMap:
          name: persistentvolumesusage
        name: persistentvolumesusage
      - configMap:
          name: pod-total
        name: pod-total
      - configMap:
          name: statefulset
        name: statefulset
      - configMap:
          name: workload-total
        name: workload-total
