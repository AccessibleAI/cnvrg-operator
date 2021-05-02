apiVersion: v1
kind: ConfigMap
metadata:
  name: redis-conf
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
data:
  redis.conf: |
    dir /data/
    appendonly "yes"
    appendfilename "appendonly.aof"
    appendfsync everysec
    auto-aof-rewrite-percentage 100
    auto-aof-rewrite-min-size 128mb
    requirepass foo-bar-password1