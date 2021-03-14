apiVersion: v1
kind: ConfigMap
metadata:
  name: redis-conf
  namespace: {{ .CnvrgNs }}
data:
  redis.conf: |
    dir /data/
    appendonly "{{ .Redis.Appendonly }}"
    appendfilename "appendonly.aof"
    appendfsync everysec
    auto-aof-rewrite-percentage 100
    auto-aof-rewrite-min-size 128mb