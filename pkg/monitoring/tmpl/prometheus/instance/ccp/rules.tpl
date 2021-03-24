apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  labels:
    app: cnvrg-ccp-prometheus
    role: alert-rules
  name: prometheus-k8s-rules
  namespace: {{ ns . }}
spec:
  groups:
    - name: k8s.rules
      rules:
        - expr: |
            sum(rate(container_cpu_usage_seconds_total{job="kubelet", metrics_path="/metrics/cadvisor", image!="", container!="POD"}[5m])) by (namespace)
          record: namespace:container_cpu_usage_seconds_total:sum_rate
        - expr: |
            sum by (cluster, namespace, pod, container) (
              rate(container_cpu_usage_seconds_total{job="kubelet", metrics_path="/metrics/cadvisor", image!="", container!="POD"}[5m])
            ) * on (cluster, namespace, pod) group_left(node) topk by (cluster, namespace, pod) (
              1, max by(cluster, namespace, pod, node) (kube_pod_info{node!=""})
            )
          record: node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate
        - expr: |
            container_memory_working_set_bytes{job="kubelet", metrics_path="/metrics/cadvisor", image!=""}
            * on (namespace, pod) group_left(node) topk by(namespace, pod) (1,
              max by(namespace, pod, node) (kube_pod_info{node!=""})
            )
          record: node_namespace_pod_container:container_memory_working_set_bytes
        - expr: |
            container_memory_rss{job="kubelet", metrics_path="/metrics/cadvisor", image!=""}
            * on (namespace, pod) group_left(node) topk by(namespace, pod) (1,
              max by(namespace, pod, node) (kube_pod_info{node!=""})
            )
          record: node_namespace_pod_container:container_memory_rss
        - expr: |
            container_memory_cache{job="kubelet", metrics_path="/metrics/cadvisor", image!=""}
            * on (namespace, pod) group_left(node) topk by(namespace, pod) (1,
              max by(namespace, pod, node) (kube_pod_info{node!=""})
            )
          record: node_namespace_pod_container:container_memory_cache
        - expr: |
            container_memory_swap{job="kubelet", metrics_path="/metrics/cadvisor", image!=""}
            * on (namespace, pod) group_left(node) topk by(namespace, pod) (1,
              max by(namespace, pod, node) (kube_pod_info{node!=""})
            )
          record: node_namespace_pod_container:container_memory_swap
        - expr: |
            sum(container_memory_usage_bytes{job="kubelet", metrics_path="/metrics/cadvisor", image!="", container!="POD"}) by (namespace)
          record: namespace:container_memory_usage_bytes:sum
        - expr: |
            sum by (namespace) (
                sum by (namespace, pod) (
                    max by (namespace, pod, container) (
                        kube_pod_container_resource_requests_memory_bytes{job="kube-state-metrics"}
                    ) * on(namespace, pod) group_left() max by (namespace, pod) (
                        kube_pod_status_phase{phase=~"Pending|Running"} == 1
                    )
                )
            )
          record: namespace:kube_pod_container_resource_requests_memory_bytes:sum
        - expr: |
            sum by (namespace) (
                sum by (namespace, pod) (
                    max by (namespace, pod, container) (
                        kube_pod_container_resource_requests_cpu_cores{job="kube-state-metrics"}
                    ) * on(namespace, pod) group_left() max by (namespace, pod) (
                      kube_pod_status_phase{phase=~"Pending|Running"} == 1
                    )
                )
            )
          record: namespace:kube_pod_container_resource_requests_cpu_cores:sum
        - expr: |
            max by (cluster, namespace, workload, pod) (
              label_replace(
                label_replace(
                  kube_pod_owner{job="kube-state-metrics", owner_kind="ReplicaSet"},
                  "replicaset", "$1", "owner_name", "(.*)"
                ) * on(replicaset, namespace) group_left(owner_name) topk by(replicaset, namespace) (
                  1, max by (replicaset, namespace, owner_name) (
                    kube_replicaset_owner{job="kube-state-metrics"}
                  )
                ),
                "workload", "$1", "owner_name", "(.*)"
              )
            )
          labels:
            workload_type: deployment
          record: namespace_workload_pod:kube_pod_owner:relabel
        - expr: |
            max by (cluster, namespace, workload, pod) (
              label_replace(
                kube_pod_owner{job="kube-state-metrics", owner_kind="DaemonSet"},
                "workload", "$1", "owner_name", "(.*)"
              )
            )
          labels:
            workload_type: daemonset
          record: namespace_workload_pod:kube_pod_owner:relabel
        - expr: |
            max by (cluster, namespace, workload, pod) (
              label_replace(
                kube_pod_owner{job="kube-state-metrics", owner_kind="StatefulSet"},
                "workload", "$1", "owner_name", "(.*)"
              )
            )
          labels:
            workload_type: statefulset
          record: namespace_workload_pod:kube_pod_owner:relabel
