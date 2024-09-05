apiVersion: batch/v1
kind: Job
metadata:
  name: run-elasticsearch-script
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.Dbs.Es.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  template:
    spec:
      containers:
        - name: script-runner
          image: byrnedo/alpine-curl:3.19
          command: ["sh", "/tmp/elastic/elastic_cleanup.sh"]
          envFrom:
            - secretRef:
                name: {{ .Spec.Dbs.Es.CredsRef }}
            - configMapRef:
                name: es-ilm-cm
          volumeMounts:
            - name: es-ilm
              mountPath: "/tmp/elastic/"
      restartPolicy: OnFailure
      volumes:
        - name: es-ilm
          configMap:
            name: "es-ilm"
            defaultMode: 0755

