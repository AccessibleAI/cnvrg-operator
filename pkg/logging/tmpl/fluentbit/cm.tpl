apiVersion: v1
kind: ConfigMap
metadata:
  name: fluent-bit-config
  namespace: {{ .Namespace }}
  annotations:
    {{- range $k, $v := .Data.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-fluentbit"
    k8s-app: fluent-bit
    {{- range $k, $v := .Data.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  fluent-bit.conf: |
    [SERVICE]
        Flush                     1
        Log_Level                 info
        Daemon                    off
        Parsers_File              parsers.conf
        HTTP_Server               On
        HTTP_Listen               0.0.0.0
        HTTP_Port                 2020
    {{- range $_, $app := .Data.AppInstance }}
    @INCLUDE {{ $app.SpecName }}-{{ $app.SpecNs }}-input.conf
    @INCLUDE {{ $app.SpecName }}-{{ $app.SpecNs }}-filter.conf
    @INCLUDE {{ $app.SpecName }}-{{ $app.SpecNs }}-output.conf
    {{- end }}

  {{- range $_, $app := .Data.AppInstance }}

  {{ $app.SpecName }}-{{ $app.SpecNs }}-input.conf: |
    [INPUT]
        Name              tail
        Tag               kube.{{ $app.SpecNs }}.*
        Path              /var/log/containers/*_{{ $app.SpecNs }}_*.log
        Parser            {{ $.Data.CriType }}
        DB                /var/log/cnvrg_flb_kube.db
        Skip_Long_Lines   On
        Refresh_Interval  10
  {{ $app.SpecName }}-{{ $app.SpecNs }}-filter.conf: |
    [FILTER]
        Name                  kubernetes
        Match                 kube.{{ $app.SpecNs }}.*
        Kube_URL              https://kubernetes.default.svc:443
        Kube_CA_File          /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        Kube_Token_File       /var/run/secrets/kubernetes.io/serviceaccount/token
        Kube_Tag_Prefix        kube.{{ $app.SpecNs }}.var.log.containers.
        Merge_Log             On
        K8S-Logging.Parser    On
        K8S-Logging.Exclude   Off

  {{ $app.SpecName }}-{{ $app.SpecNs }}-output.conf: |
    [OUTPUT]
        Name                      es
        Match                     kube.{{ $app.SpecNs }}.*
        Host                      elasticsearch.{{ $app.SpecNs }}.svc.{{ $.Data.ClusterInternalDomain }}
        Port                      9200
        Logstash_Format           On
        Logstash_DateFormat       %m.%Y
        Replace_Dots              On
        Retry_Limit               False
        Trace_Error               On
        Index                     cnvrg-all
        Logstash_Prefix            cnvrg-all
        HTTP_User                 {{ $app.EsUser }}
        HTTP_Passwd               {{ $app.EsPass }}

    [OUTPUT]
        Name                      es
        Match                     kube.{{ $app.SpecNs }}.*app*
        Host                      elasticsearch.{{ $app.SpecNs }}.svc.{{ $.Data.ClusterInternalDomain }}
        Port                      9200
        Logstash_Format           On
        Logstash_DateFormat       %m.%Y
        Replace_Dots              On
        Retry_Limit               False
        Trace_Error               On
        Index                     cnvrg-app
        Logstash_Prefix            cnvrg-app
        HTTP_User                 {{ $app.EsUser }}
        HTTP_Passwd               {{ $app.EsPass }}

    [OUTPUT]
        Name                      es
        Match                     kube.{{ $app.SpecNs }}.*kiq*
        Host                      elasticsearch.{{ $app.SpecNs }}.svc.{{ $.Data.ClusterInternalDomain }}
        Port                      9200
        Logstash_Format           On
        Logstash_DateFormat       %m.%Y
        Replace_Dots              On
        Retry_Limit               False
        Trace_Error               On
        Index                     cnvrg-app
        Logstash_Prefix            cnvrg-app
        HTTP_User                 {{ $app.EsUser }}
        HTTP_Passwd               {{ $app.EsPass }}

    [OUTPUT]
        Name                      es
        Match                     kube.{{ $app.SpecNs }}.*hyper*
        Host                      elasticsearch.{{ $app.SpecNs }}.svc.{{ $.Data.ClusterInternalDomain }}
        Port                      9200
        Logstash_Format           On
        Logstash_DateFormat       %m.%Y
        Replace_Dots              On
        Retry_Limit               False
        Trace_Error               On
        Index                     cnvrg-app
        Logstash_Prefix            cnvrg-app
        HTTP_User                 {{ $app.EsUser }}
        HTTP_Passwd               {{ $app.EsPass }}

    [OUTPUT]
        Name                      es
        Match                     kube.{{ $app.SpecNs }}.*scheduler*
        Host                      elasticsearch.{{ $app.SpecNs }}.svc.{{ $.Data.ClusterInternalDomain }}
        Port                      9200
        Logstash_Format           On
        Logstash_DateFormat       %m.%Y
        Replace_Dots              On
        Retry_Limit               False
        Trace_Error               On
        Index                     cnvrg-app
        Logstash_Prefix            cnvrg-app
        HTTP_User                 {{ $app.EsUser }}
        HTTP_Passwd               {{ $app.EsPass }}

    [OUTPUT]
        Name                      es
        Match                     kube.{{ $app.SpecNs }}.*job*
        Host                      elasticsearch.{{ $app.SpecNs }}.svc.{{ $.Data.ClusterInternalDomain }}
        Port                      9200
        Logstash_Format           On
        Logstash_DateFormat       %m.%Y
        Replace_Dots              On
        Retry_Limit               False
        Trace_Error               On
        Index                     cnvrg-jobs
        Logstash_Prefix            cnvrg-jobs
        HTTP_User                 {{ $app.EsUser }}
        HTTP_Passwd               {{ $app.EsPass }}

    [OUTPUT]
        Name                      es
        Match                     kube.{{ $app.SpecNs }}.*cnvrg-je*
        Host                      elasticsearch.{{ $app.SpecNs }}.svc.{{ $.Data.ClusterInternalDomain }}
        Port                      9200
        Logstash_Format           On
        Logstash_DateFormat       %m.%Y
        Replace_Dots              On
        Retry_Limit               False
        Trace_Error               On
        Index                     cnvrg-endpoints
        Logstash_Prefix            cnvrg-endpoints
        HTTP_User                 {{ $app.EsUser }}
        HTTP_Passwd               {{ $app.EsPass }}
        
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
        Name containerd
        Format regex
        Regex ^(?<time>[^ ]+) (?<stream>stdout|stderr) (?<logtag>[^ ]*) (?<log>.*)$
        Time_Key    time
        Time_Format %Y-%m-%dT%H:%M:%S.%L%z

    [PARSER]
        # http://rubular.com/r/tjUt3Awgg4
        Name cri-o
        Format regex
        Regex ^(?<time>[^ ]+) (?<stream>stdout|stderr) (?<logtag>[^ ]*) (?<log>.*)$
        Time_Key    time
        Time_Format %Y-%m-%dT%H:%M:%S.%L%z

    [PARSER]
        Name        syslog
        Format      regex
        Regex       ^\<(?<pri>[0-9]+)\>(?<time>[^ ]* {1,2}[^ ]* [^ ]*) (?<host>[^ ]*) (?<ident>[a-zA-Z0-9_\/\.\-]*)(?:\[(?<pid>[0-9]+)\])?(?:[^\:]*\:)? *(?<message>.*)$
        Time_Key    time
        Time_Format %b %d %H:%M:%S