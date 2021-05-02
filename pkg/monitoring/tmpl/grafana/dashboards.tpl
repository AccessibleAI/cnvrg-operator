kind: ConfigMap
apiVersion: v1
metadata:
  name: grafana-dashboards
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
data:
  dashboards.yaml: |-
    {
        "apiVersion": 1,
        "providers": [
            {
                "folder": "Cnvrg - Default",
                "name": "0",
                "options": {
                    "path": "/definitions/0"
                },
                "orgId": 1,
                "type": "file"
            }
        ]
    }
