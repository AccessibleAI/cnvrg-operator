apiVersion: batch/v1
kind: Job
metadata:
  name: ingresscheck
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
      {{$k}}: "{{$v}}"
      {{- end }}
  labels:
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-ccp"
    app: ingresscheck
    owner: cnvrg-control-plane
    cnvrg-component: ingresscheck
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  template:
    metadata:
      name: ingresscheck
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
          {{$k}}: "{{$v}}"
          {{- end }}
      labels:
        cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-ccp"
        app: ingresscheck
        owner: cnvrg-control-plane
        cnvrg-component: ingresscheck
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      containers:
        - name: python-server
          image: {{ image .Spec.ImageHub .Spec.ControlPlane.Image }}
          imagePullPolicy: Always
          command:
            - "/bin/bash"
            - "-lc"
            - "python -m http.server"
        - name: ingresscheck
          image: {{ image .Spec.ImageHub .Spec.ControlPlane.Image }}
          imagePullPolicy: Always
          envFrom:
            - configMapRef:
                name: cp-base-config
            - configMapRef:
                name: cp-networking-config
            - secretRef:
                name: cp-base-secret
            - secretRef:
                name: cp-ldap
            - secretRef:
                name: cp-object-storage
            - secretRef:
                name: cp-smtp
            {{- if isTrue .Spec.Dbs.Es.Enabled }}
            - secretRef:
                name: {{ .Spec.Dbs.Es.CredsRef }}
            {{- end }}
            {{- if isTrue .Spec.Dbs.Redis.Enabled }}
            - secretRef:
                name: {{ .Spec.Dbs.Redis.CredsRef }}
            {{- end }}
            {{- if isTrue .Spec.Monitoring.Prometheus.Enabled }}
            - secretRef:
                name: {{ .Spec.Monitoring.Prometheus.CredsRef }}
            {{- end }}
            {{- if isTrue .Spec.Networking.Proxy.Enabled }}
            - configMapRef:
                name: {{ .Spec.Networking.Proxy.ConfigRef }}
            {{- end }}
          command:
            - "/bin/bash"
            - "-lc"
            - |
              #!/bin/bash
              apt-get update && apt-get install -y curl
              flagFile=/tmp/services_not_ready
              echo true > ${flagFile}
              while $(cat ${flagFile}); do

                if [[ $(curl -sk {{ httpScheme . }}cnvrg-ingress-test.{{ .Spec.ClusterDomain }} -o /dev/null -w '%{http_code}') != 200 ]]; then
                  echo "[$(date)] grafana [{{ httpScheme . }}cnvrg-ingress-test.{{ .Spec.ClusterDomain }}] not ready"
                  sleep 1
                  continue
                fi

                echo false > ${flagFile}
                echo "[$(date)] test service is ready!"
              done
      restartPolicy: Never
  backoffLimit: 4