kind: ConfigMap
apiVersion: v1
metadata:
  name: grafana-dashboards
  namespace: {{ ns . }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
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
