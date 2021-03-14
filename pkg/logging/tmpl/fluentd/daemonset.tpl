apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluentd
  namespace: {{ .CnvrgNs }}
  labels:
    app: fluentd-logging
spec:
  selector:
    matchLabels:
      app: fluentd-logging
  template:
    metadata:
      labels:
        app: fluentd-logging
    spec:
      serviceAccountName: fluentd
      tolerations:
        - key: node-role.kubernetes.io/master
          effect: NoSchedule
        - key: nvidia.com/gpu
          operator: Exists
          effect: NoSchedule
        - key: "kubernetes.azure.com/scalesetpriority"
          operator: "Equal"
          value: "spot"
          effect: "NoSchedule"
        - key: {{ .ControlPlan.Tenancy.Key }}
          operator: "Equal"
          value: "{{ .ControlPlan.Tenancy.Value }}"
          effect: "NoSchedule"
      containers:
        - name: fluentd
          image: {{ .Logging.Fluentd.Image }}
          securityContext:
            privileged: true
          env:
            - name:  FLUENT_ELASTICSEARCH_HOST
              value: {{ printf "%s.%s.svc.cluster.local" .Logging.Es.SvcName .CnvrgNs }}
            - name:  FLUENT_ELASTICSEARCH_PORT
              value: "9200"
            - name: FLUENT_ELASTICSEARCH_SCHEME
              value: "http"
            - name: FLUENT_UID
              value: "0"
            - name:  FLUENT_ELASTICSEARCH_LOGSTASH_PREFIX
              value: "cnvrg"
            - name: FLUENT_ELASTICSEARCH_LOGSTASH_INDEX_NAME
              value: "cnvrg"
          resources:
            limits:
              memory: {{ .Logging.Fluentd.MemoryLimit }}
            requests:
              cpu: {{ .Logging.Fluentd.CPURequest }}
              memory: {{ .Logging.Fluentd.MemoryRequest }}
          volumeMounts:
            - name: config-volume
              mountPath: /fluentd/etc/fluent.conf
              subPath: fluent.conf
            - name: journal
              mountPath: /var/log/journal
              readOnly: true
            - name: varlog
              mountPath: /var/log
            - name: varlibdockercontainers
              mountPath: /var/lib/docker/containers
              readOnly: true
      terminationGracePeriodSeconds: 30
      volumes:
        - name: config-volume
          configMap:
            name: fluentd-conf
        - name: varlog
          hostPath:
            path: /var/log
        - name: varlibdockercontainers
          hostPath:
            path: {{ .Logging.Fluentd.ContainersPath }}
        - name: journal
          hostPath:
            path: {{ .Logging.Fluentd.JournalPath }}
