apiVersion: v1
kind: ConfigMap
metadata:
  name: redis-conf
  namespace: {{ .Namespace }}
data:
  redis.conf: |
    dir /data/
    appendonly "{{ .Spec.Redis.Appendonly }}"
    appendfilename "appendonly.aof"
    appendfsync everysec
    auto-aof-rewrite-percentage 100
    auto-aof-rewrite-min-size 128mb