apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Spec.ControlPlane.WebApp.SvcName }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.ControlPlane.WebApp.SvcName }}
    owner: cnvrg-control-plane
    cnvrg-component: webapp
    cnvrg-system-status-check: "true"
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  {{- if eq false (isTrue .Spec.ControlPlane.WebApp.Hpa.Enabled) }}
  replicas: {{ .Spec.ControlPlane.WebApp.Replicas }}
  {{- end }}
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 25%
      maxSurge: 1
  selector:
    matchLabels:
      app: {{.Spec.ControlPlane.WebApp.SvcName}}
  template:
    metadata:
      annotations:
        cnvrg-component: webapp
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        app: {{.Spec.ControlPlane.WebApp.SvcName}}
        owner: cnvrg-control-plane
        cnvrg-component: webapp
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - podAffinityTerm:
              labelSelector:
                matchLabels:
                  app: {{ .Spec.ControlPlane.WebApp.SvcName }}
              namespaces:
              - {{ .Namespace }}
              topologyKey: kubernetes.io/hostname
            weight: 1
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        {{ .Spec.Tenancy.Key }}: "{{ .Spec.Tenancy.Value }}"
      tolerations:
        - key: "{{ .Spec.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.Tenancy.Value }}"
          effect: "NoSchedule"
      {{- end }}
      serviceAccountName: cnvrg-control-plane
      securityContext:
        runAsUser: 1000
        runAsGroup: 1000
      priorityClassName: {{ .Spec.PriorityClass.AppClassRef }}
      enableServiceLinks: false
      containers:
      - image: {{ image .Spec.ImageHub .Spec.ControlPlane.Image }}
        env:
        - name: "CNVRG_RUN_MODE"
          value: "webapp"
        envFrom:
        - configMapRef:
            name: cp-base-config
        - configMapRef:
            name: cp-networking-config
        - secretRef:
            name: cp-base-secret
        {{- if isTrue .Spec.SSO.Enabled }}
        - secretRef:
            name: cp-oauth-proxy-tokens-secret
        {{- end }}
        - secretRef:
            name: cp-ldap
        - secretRef:
            name: cp-object-storage
        - secretRef:
            name: cp-smtp
        - secretRef:
            name: {{ .Spec.Dbs.Es.CredsRef }}
        - secretRef:
            name: {{ .Spec.Dbs.Pg.CredsRef }}
        - secretRef:
            name: {{ .Spec.Dbs.Redis.CredsRef }}
        - secretRef:
            name: {{ .Spec.Dbs.Prom.CredsRef }}
        {{- if isTrue .Spec.Networking.Proxy.Enabled }}
        - configMapRef:
            name: {{ .Spec.Networking.Proxy.ConfigRef }}
        {{- end }}
        name: cnvrg-app
        ports:
          - containerPort: {{ .Spec.ControlPlane.WebApp.Port }}
        readinessProbe:
          httpGet:
            path: "/healthz"
            port: {{ .Spec.ControlPlane.WebApp.Port }}
            scheme: HTTP
          successThreshold: 1
          failureThreshold: {{ .Spec.ControlPlane.WebApp.FailureThreshold }}
          initialDelaySeconds: {{ .Spec.ControlPlane.WebApp.InitialDelaySeconds }}
          periodSeconds: {{ .Spec.ControlPlane.WebApp.ReadinessPeriodSeconds }}
          timeoutSeconds: {{ .Spec.ControlPlane.WebApp.ReadinessTimeoutSeconds }}
        livenessProbe:
          successThreshold: 1
          failureThreshold: {{ .Spec.ControlPlane.WebApp.FailureThreshold }}
          initialDelaySeconds: {{ .Spec.ControlPlane.WebApp.InitialDelaySeconds }}
          periodSeconds: {{ .Spec.ControlPlane.WebApp.ReadinessPeriodSeconds }}
          timeoutSeconds: {{ .Spec.ControlPlane.WebApp.ReadinessTimeoutSeconds }}
          exec:
             command:
             - /bin/bash
             - -c
             - |
               threshold=50
               requests=$(passenger-status | grep '(app)' -A 2 | grep Requests | awk '{print $NF}')
               if (( $requests > $threshold )); then
                   exit 1
               fi
        resources:
          requests:
            cpu: "{{.Spec.ControlPlane.WebApp.Requests.Cpu}}"
            memory: "{{.Spec.ControlPlane.WebApp.Requests.Memory}}"
          limits:
            cpu: "{{.Spec.ControlPlane.WebApp.Limits.Cpu}}"
            memory: "{{.Spec.ControlPlane.WebApp.Limits.Memory}}"
        volumeMounts:
        {{- if eq .Spec.ControlPlane.ObjectStorage.Type "gcp" }}
        - name: {{ .Spec.ControlPlane.ObjectStorage.GcpSecretRef }}
          mountPath: /opt/app-root/conf/gcp-keyfile
          readOnly: true
        {{- end }}
        {{- if and ( isTrue .Spec.Networking.Ingress.OcpSecureRoutes) (eq .Spec.Networking.Ingress.Type "openshift") }}
        - name: tls-secret
          readOnly: true
          mountPath: /opt/app-root/src/tls
        {{- end }}
      volumes:
      {{- if eq .Spec.ControlPlane.ObjectStorage.Type "gcp" }}
      - name: {{ .Spec.ControlPlane.ObjectStorage.GcpSecretRef }}
        secret:
          secretName: {{ .Spec.ControlPlane.ObjectStorage.GcpSecretRef }}
      {{- end }}
      {{- if and ( isTrue .Spec.Networking.Ingress.OcpSecureRoutes) (eq .Spec.Networking.Ingress.Type "openshift") }}
      - name: tls-secret
        secret:
          secretName: {{ .Spec.Networking.HTTPS.CertSecret }}
      {{- end }}
      initContainers:
      - name: ingresscheck
        image: {{ image .Spec.ImageHub .Spec.ControlPlane.Image }}
        envFrom:
        - secretRef:
            name: cp-object-storage
        - secretRef:
            name: {{ .Spec.Dbs.Es.CredsRef }}
        - secretRef:
            name: {{ .Spec.Dbs.Pg.CredsRef }}
        - secretRef:
            name: {{ .Spec.Dbs.Redis.CredsRef }}
        - secretRef:
            name: {{ .Spec.Dbs.Prom.CredsRef }}
        {{- if isTrue .Spec.Networking.Proxy.Enabled }}
        - configMapRef:
            name: {{ .Spec.Networking.Proxy.ConfigRef }}
        {{- end }}
        resources:
          requests:
            cpu: "100m"
            memory: "200Mi"
          limits:
            cpu: "1000m"
            memory: "1Gi"
        command:
        - "/bin/bash"
        - "-lc"
        - |
          #!/bin/bash
          flagFile=/tmp/services_not_ready
          echo true > ${flagFile}
          while $(cat ${flagFile}); do

            timeout 2 bash -c "</dev/tcp/{{.Spec.Dbs.Redis.SvcName}}/{{.Spec.Dbs.Redis.Port}}";
            if [[ $? != 0 ]]; then
              echo "[$(date)] redis not ready"
              sleep 1
              continue
            fi
            echo "[$(date)] redis is ready!"

            timeout 2 bash -c "</dev/tcp/${POSTGRES_HOST}/{{.Spec.Dbs.Pg.Port}}";
            if [[ $? != 0 ]]; then
              echo "[$(date)] postgres [${POSTGRES_HOST}:{{.Spec.Dbs.Pg.Port}}] not ready"
              sleep 1
              continue
            fi
            echo "[$(date)] postgres is ready!"

            if [[ $(curl -s $ELASTICSEARCH_URL/_cluster/health -o /dev/null -w '%{http_code}') != 200 ]]; then
              echo "[$(date)] elasticsearch not ready"
              sleep 1
              continue
            fi
            echo "[$(date)] elasticsearch is ready!"

            {{- if and ( isTrue .Spec.Dbs.Minio.Enabled ) (eq .Spec.ControlPlane.ObjectStorage.Type "minio") }}
            if [[ $(curl -sk $CNVRG_STORAGE_ENDPOINT/minio/health/ready -o /dev/null -w '%{http_code}') != 200 ]]; then
              echo "[$(date)] minio [$CNVRG_STORAGE_ENDPOINT/minio/health/ready] not ready"
              sleep 1
              continue
            fi
            echo "[$(date)] minio is ready!"
            {{- end }}

            echo false > ${flagFile}
            echo "[$(date)] all services are ready!"
          done
      - name: seeder
        image: {{ image .Spec.ImageHub .Spec.ControlPlane.Image }}
        command:
          - /bin/bash
          - -lc
          - |
            if [[ $(kubectl get cm cnvrg-db-init -oname --ignore-not-found | wc -l) == 0 ]]; then
              rails db:migrate \
              && rails db:seed \
              && rails libraries:update \
              && kubectl create cm cnvrg-db-init -n ${KUBE_NAMESPACE}
            fi
        resources:
          requests:
            cpu: "100m"
            memory: "200Mi"
          limits:
            cpu: "1000m"
            memory: "1Gi"
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
        - secretRef:
            name: {{ .Spec.Dbs.Es.CredsRef }}
        - secretRef:
            name: {{ .Spec.Dbs.Pg.CredsRef }}
        - secretRef:
            name: {{ .Spec.Dbs.Redis.CredsRef }}
        - secretRef:
            name: {{ .Spec.Dbs.Prom.CredsRef }}
        {{- if isTrue .Spec.Networking.Proxy.Enabled }}
        - configMapRef:
            name: {{ .Spec.Networking.Proxy.ConfigRef }}
        {{- end }}
        {{- if eq .Spec.ControlPlane.ObjectStorage.Type "gcp" }}
        volumeMounts:
        - name: {{ .Spec.ControlPlane.ObjectStorage.GcpSecretRef }}
          mountPath: /opt/app-root/conf/gcp-keyfile
          readOnly: true
        {{- end }}


