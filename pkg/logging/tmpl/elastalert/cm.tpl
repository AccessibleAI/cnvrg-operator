apiVersion: v1
kind: ConfigMap
metadata:
  name: elastalert-config
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  config.json: |
    {
      "appName": "elastalert-server",
      "port": 3030,
      "wsport": 3333,
      "elastalertPath": "/opt/elastalert",
      "verbose": true,
      "rulesPath": {
        "relative": true,
        "path": "/rules"
      },
      "templatesPath": {
        "relative": true,
        "path": "/rule_templates"
      },
      "dataPath": {
        "relative": true,
        "path": "/server_data"
      },
      "es_host": "{{ .Spec.Dbs.Es.SvcName }}",
      "es_port": {{ .Spec.Dbs.Es.Port }},
      "es_ssl": false,
      "ea_verify_certs": false,
      "writeback_index": "elastalert_status"
    }
  config.yaml: |
    rules_folder: rules
    run_every:
      minutes: 1
    buffer_time:
      minutes: 15
    es_host: {{ .Spec.Dbs.Es.SvcName }}
    es_port: {{ .Spec.Dbs.Es.Port}}
    use_ssl: False
    verify_certs: False
    writeback_index: elastalert_status
    writeback_alias: elastalert_alerts
    alert_time_limit:
      days: 2
