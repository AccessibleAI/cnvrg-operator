apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: hostpath-provisioner
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    k8s-app: hostpath-provisioner
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  selector:
    matchLabels:
      k8s-app: hostpath-provisioner
  template:
    metadata:
      annotations:
        {{- range $k, $v := .Spec.Annotations }}
        {{$k}}: "{{$v}}"
        {{- end }}
      labels:
        k8s-app: hostpath-provisioner
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
    spec:
      priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
      {{- if isTrue .Spec.Tenancy.Enabled }}
      nodeSelector:
        {{ .Spec.Tenancy.Key }}: {{ .Spec.Tenancy.Value }}
        {{- range $key, $val := .Spec.Storage.Hostpath.NodeSelector }}
        {{ $key }}: {{ $val }}
      {{- end }}
      tolerations:
        - key: "{{ .Spec.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.Tenancy.Value }}"
          effect: "NoSchedule"
      {{- else if (gt (len .Spec.Storage.Hostpath.NodeSelector) 0) }}
      nodeSelector:
        {{- range $key, $val := .Spec.Storage.Hostpath.NodeSelector }}
        {{ $key }}: {{ $val }}
        {{- end }}
      {{- end }}
      serviceAccountName: hostpath-provisioner-admin
      containers:
        - name: hostpath-provisioner
          image: {{image .Spec.ImageHub .Spec.Storage.Hostpath.Image }}
          imagePullPolicy: Always
          env:
            - name: USE_NAMING_PREFIX
              value: "true" # change to true, to have the name of the pvc be part of the directory
            - name: NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: PV_DIR
              value: {{ .Spec.Storage.Hostpath.Path }}
          volumeMounts:
            - name: pv-volume # root dir where your bind mounts will be on the node
              mountPath: {{ .Spec.Storage.Hostpath.Path }}
          resources:
            limits:
              cpu: {{ .Spec.Storage.Hostpath.Limits.Cpu }}
              memory: {{ .Spec.Storage.Hostpath.Limits.Memory }}
            requests:
              cpu: {{ .Spec.Storage.Hostpath.Requests.Cpu }}
              memory: {{ .Spec.Storage.Hostpath.Requests.Memory }}
      volumes:
        - name: pv-volume
          hostPath:
            path: {{ .Spec.Storage.Hostpath.Path }}
