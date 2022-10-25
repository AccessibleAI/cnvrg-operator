apiVersion: v1
kind: ConfigMap
metadata:
  name: es-ilm
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
  elastic_cleanup.sh: |+
    #!/bin/bash
    ES_URL=http://localhost:9200
    while [[ "$(curl -u "elastic:${CNVRG_ES_PASS}" -s -o /dev/null -w '%{http_code}\n' $ES_URL)" != "200" ]]; do sleep 2; done
    while [[ "$(curl -X PUT -u "elastic:${CNVRG_ES_PASS}" "${ES_URL}/_security/user/${CNVRG_ES_USER}?pretty"  -H 'Content-Type: application/json'  -d "{ \"password\" : \"${CNVRG_ES_PASS}\", \"roles\" : [ \"superuser\" ] }"  -s -o /dev/null -w '%{http_code}\n')" != "200" ]]; do sleep 2; done
    curl -X PUT -u "${CNVRG_ES_USER}:${CNVRG_ES_PASS}" "${ES_URL}/_ilm/policy/cleanup_policy_app?pretty"  -H 'Content-Type: application/json'  -d '{ "policy": { "phases": { "hot": { "actions": {} }, "delete": { "min_age": "30d", "actions": { "delete": {} } } } } }'
    curl -X PUT -u "${CNVRG_ES_USER}:${CNVRG_ES_PASS}" "${ES_URL}/_ilm/policy/cleanup_policy_jobs?pretty"  -H 'Content-Type: application/json'  -d '{ "policy": { "phases": { "hot": { "actions": {} }, "delete": { "min_age": "14d", "actions": { "delete": {} } } } } }'
    curl -X PUT -u "${CNVRG_ES_USER}:${CNVRG_ES_PASS}" "${ES_URL}/_ilm/policy/cleanup_policy_all?pretty"  -H 'Content-Type: application/json'  -d '{ "policy": { "phases": { "hot": { "actions": {} }, "delete": { "min_age": "3d", "actions": { "delete": {} } } } } }'
    curl -X PUT -u "${CNVRG_ES_USER}:${CNVRG_ES_PASS}" "${ES_URL}/_ilm/policy/cleanup_policy_endpoints?pretty"  -H 'Content-Type: application/json'  -d '{ "policy": { "phases": { "hot": { "actions": {} }, "delete": { "min_age": "1825d", "actions": { "delete": {} } } } } }'
    curl -X PUT -u "${CNVRG_ES_USER}:${CNVRG_ES_PASS}" "${ES_URL}/cnvrg-app*/_settings?pretty"  -H 'Content-Type: application/json'  -d '{ "lifecycle.name": "cleanup_policy_app" }'
    curl -X PUT -u "${CNVRG_ES_USER}:${CNVRG_ES_PASS}" "${ES_URL}/cnvrg-jobs*/_settings?pretty"  -H 'Content-Type: application/json'  -d '{ "lifecycle.name": "cleanup_policy_jobs" }'
    curl -X PUT -u "${CNVRG_ES_USER}:${CNVRG_ES_PASS}" "${ES_URL}/cnvrg-all*/_settings?pretty"  -H 'Content-Type: application/json'  -d '{ "lifecycle.name": "cleanup_policy_all" }'
    curl -X PUT -u "${CNVRG_ES_USER}:${CNVRG_ES_PASS}" "${ES_URL}/cnvrg-endpoints*/_settings?pretty"  -H 'Content-Type: application/json'  -d '{ "lifecycle.name": "cleanup_policy_endpoints" }'
    curl -X PUT -u "${CNVRG_ES_USER}:${CNVRG_ES_PASS}" "${ES_URL}/_template/logging_policy_template_app?pretty"  -H 'Content-Type: application/json'  -d '{ "index_patterns": ["cnvrg-app*"], "settings": { "index.lifecycle.name": "cleanup_policy_app" } }'
    curl -X PUT -u "${CNVRG_ES_USER}:${CNVRG_ES_PASS}" "${ES_URL}/_template/logging_policy_template_jobs?pretty"  -H 'Content-Type: application/json'  -d '{ "index_patterns": ["cnvrg-jobs*"], "settings": { "index.lifecycle.name": "cleanup_policy_jobs" } }'
    curl -X PUT -u "${CNVRG_ES_USER}:${CNVRG_ES_PASS}" "${ES_URL}/_template/logging_policy_template_all?pretty"  -H 'Content-Type: application/json'  -d '{ "index_patterns": ["cnvrg-all*"], "settings": { "index.lifecycle.name": "cleanup_policy_all" } }'
    curl -X PUT -u "${CNVRG_ES_USER}:${CNVRG_ES_PASS}" "${ES_URL}/_template/logging_policy_template_endpoints?pretty"  -H 'Content-Type: application/json'  -d '{ "index_patterns": ["cnvrg-endpoints*"], "settings": { "index.lifecycle.name": "cleanup_policy_endpoints" } }'
