apiVersion: v1
kind: ConfigMap
metadata:
  name: fluent-bit-config
  namespace: {{ ns . }}
  labels:
    k8s-app: fluent-bit
data:
  fluent-bit.conf: |
    [SERVICE]
        Flush         1
        Log_Level     info
        Daemon        off
        Parsers_File  parsers.conf
        HTTP_Server   On
        HTTP_Listen   0.0.0.0
        HTTP_Port     2020
          {{- range $_, $cnvrgAppInstance := .Spec.CnvrgAppInstances }}
    @INCLUDE {{ $cnvrgAppInstance.Name }}-{{ $cnvrgAppInstance.Namespace }}-input.conf
    @INCLUDE {{ $cnvrgAppInstance.Name }}-{{ $cnvrgAppInstance.Namespace }}-filter.conf
    @INCLUDE {{ $cnvrgAppInstance.Name }}-{{ $cnvrgAppInstance.Namespace }}-output.conf
          {{- end }}

          {{- range $_, $cnvrgAppInstance := .Spec.CnvrgAppInstances }}

  {{ $cnvrgAppInstance.Name }}-{{ $cnvrgAppInstance.Namespace }}-input.conf: |
    [INPUT]
        Name              tail
        Tag               kube.{{ $cnvrgAppInstance.Namespace }}.*
        Path              /var/log/containers/*_{{ $cnvrgAppInstance.Namespace }}_*.log
        Parser            docker
        DB                /var/log/flb_kube.db
        Mem_Buf_Limit     5MB
        Skip_Long_Lines   On
        Refresh_Interval  10

  {{ $cnvrgAppInstance.Name }}-{{ $cnvrgAppInstance.Namespace }}-filter.conf: |
    [FILTER]
        Name                kubernetes
        Match               kube.{{ $cnvrgAppInstance.Namespace }}.*
        Kube_URL            https://kubernetes.default.svc:443
        Kube_CA_File        /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        Kube_Token_File     /var/run/secrets/kubernetes.io/serviceaccount/token
        Kube_Tag_Prefix     kube.{{ $cnvrgAppInstance.Namespace }}.var.log.containers.
        Merge_Log           On
        Merge_Log_Key       log_processed
        K8S-Logging.Parser  On
        K8S-Logging.Exclude Off

  {{ $cnvrgAppInstance.Name }}-{{ $cnvrgAppInstance.Namespace }}-output.conf: |
    [OUTPUT]
        Name            es
        Match           kube.{{ $cnvrgAppInstance.Namespace }}.*
        Host            elasticsearch.{{ $cnvrgAppInstance.Namespace }}.svc.cluster.local
        Port            9200
        Logstash_Format Off
        Replace_Dots    On
        Retry_Limit     False
        Index           cnvrg
        Logstash_Prefix cnvrg

          {{- end }}

  parsers.conf: |
    [PARSER]
        Name   apache
        Format regex
        Regex  ^(?<host>[^ ]*) [^ ]* (?<user>[^ ]*) \[(?<time>[^\]]*)\] "(?<method>\S+)(?: +(?<path>[^\"]*?)(?: +\S*)?)?" (?<code>[^ ]*) (?<size>[^ ]*)(?: "(?<referer>[^\"]*)" "(?<agent>[^\"]*)")?$
        Time_Key time
        Time_Format %d/%b/%Y:%H:%M:%S %z

    [PARSER]
        Name   apache2
        Format regex
        Regex  ^(?<host>[^ ]*) [^ ]* (?<user>[^ ]*) \[(?<time>[^\]]*)\] "(?<method>\S+)(?: +(?<path>[^ ]*) +\S*)?" (?<code>[^ ]*) (?<size>[^ ]*)(?: "(?<referer>[^\"]*)" "(?<agent>[^\"]*)")?$
        Time_Key time
        Time_Format %d/%b/%Y:%H:%M:%S %z

    [PARSER]
        Name   apache_error
        Format regex
        Regex  ^\[[^ ]* (?<time>[^\]]*)\] \[(?<level>[^\]]*)\](?: \[pid (?<pid>[^\]]*)\])?( \[client (?<client>[^\]]*)\])? (?<message>.*)$

    [PARSER]
        Name   nginx
        Format regex
        Regex ^(?<remote>[^ ]*) (?<host>[^ ]*) (?<user>[^ ]*) \[(?<time>[^\]]*)\] "(?<method>\S+)(?: +(?<path>[^\"]*?)(?: +\S*)?)?" (?<code>[^ ]*) (?<size>[^ ]*)(?: "(?<referer>[^\"]*)" "(?<agent>[^\"]*)")?$
        Time_Key time
        Time_Format %d/%b/%Y:%H:%M:%S %z

    [PARSER]
        Name   json
        Format json
        Time_Key time
        Time_Format %d/%b/%Y:%H:%M:%S %z

    [PARSER]
        Name        docker
        Format      json
        Time_Key    time
        Time_Format %Y-%m-%dT%H:%M:%S.%L
        Time_Keep   On

    [PARSER]
        # http://rubular.com/r/tjUt3Awgg4
        Name cri
        Format regex
        Regex ^(?<time>[^ ]+) (?<stream>stdout|stderr) (?<logtag>[^ ]*) (?<message>.*)$
        Time_Key    time
        Time_Format %Y-%m-%dT%H:%M:%S.%L%z

    [PARSER]
        Name        syslog
        Format      regex
        Regex       ^\<(?<pri>[0-9]+)\>(?<time>[^ ]* {1,2}[^ ]* [^ ]*) (?<host>[^ ]*) (?<ident>[a-zA-Z0-9_\/\.\-]*)(?:\[(?<pid>[0-9]+)\])?(?:[^\:]*\:)? *(?<message>.*)$
        Time_Key    time
        Time_Format %b %d %H:%M:%S