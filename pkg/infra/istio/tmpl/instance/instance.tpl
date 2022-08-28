apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  name: cnvrg-istio
  namespace:  {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  profile: minimal
  tag: {{ trimPrefix "pilot:" .Spec.Networking.Istio.PilotImage }}
  hub: {{ .Spec.ImageHub }}
  components:
    base:
      enabled: true
    cni:
      enabled: false
    egressGateways:
    - enabled: false
      name: istio-egressgateway
    ingressGateways:
    - enabled: true
      name: cnvrg-ingressgateway
      label:
        istio: cnvrg-ingressgateway
        {{- range $k, $v := .Spec.Labels }}
        {{$k}}: "{{$v}}"
        {{- end }}
      k8s:
        priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
        podAnnotations:
          {{- range $k, $v := .Spec.Annotations }}
          {{$k}}: "{{$v}}"
          {{- end }}
        {{- if isTrue .Spec.Tenancy.Enabled }}
        nodeSelector:
          {{ .Spec.Tenancy.Key }}: "{{ .Spec.Tenancy.Value }}"
        tolerations:
          - operator: "Exists"
        {{- end }}
        serviceAnnotations:
        {{- range $name, $value := .Spec.Networking.Istio.IngressSvcAnnotations }}
          {{ $name }}: "{{ $value }}"
        {{- end }}
        hpaSpec:
          maxReplicas: 5
          metrics:
            - resource:
                name: cpu
                targetAverageUtilization: 80
              type: Resource
          minReplicas: 1
          scaleTargetRef:
            apiVersion: apps/v1
            kind: Deployment
            name: cnvrg-ingressgateway
        resources:
          limits:
            cpu: 2000m
            memory: 1024Mi
          requests:
            cpu: 100m
            memory: 128Mi
        service:
          loadBalancerSourceRanges:
          {{- range $idx, $range := .Spec.Networking.Istio.LBSourceRanges }}
            - {{ $range}}
          {{- end }}
          {{- if gt (len .Spec.Networking.Istio.ExternalIP) 0 }}
          type: ClusterIP
          externalIPs:
          {{- range $idx, $ip := .Spec.Networking.Istio.ExternalIP }}
            - {{$ip}}
          {{- end }}
          {{- end }}
          ports:
          - name: http2
            port: 80
            targetPort: 8080
          - name: https
            port: 443
            targetPort: 8443
          {{- if gt (len .Spec.Networking.Istio.IngressSvcExtraPorts) 0 }}
          {{- range $idx, $port := .Spec.Networking.Istio.IngressSvcExtraPorts }}
          - name: port{{ $port}}
            port: {{ $port}}
            {{- end }}
          {{- end }}
    istiodRemote:
      enabled: false
    pilot:
      enabled: true
      k8s:
        priorityClassName: {{ .Spec.CnvrgAppPriorityClass.Name }}
        podAnnotations:
          {{- range $k, $v := .Spec.Annotations }}
          {{$k}}: "{{$v}}"
          {{- end }}
        {{- if isTrue .Spec.Tenancy.Enabled }}
        nodeSelector:
          {{ .Spec.Tenancy.Key }}: "{{ .Spec.Tenancy.Value }}"
        tolerations:
          - key: "{{ .Spec.Tenancy.Key }}"
            operator: "Equal"
            value: "{{ .Spec.Tenancy.Value }}"
            effect: "NoSchedule"
        {{- end }}
  values:
    global:
      istioNamespace:  {{ ns . }}
      imagePullSecrets:
        - {{ .Spec.Registry.Name }}
    meshConfig:
      rootNamespace: {{ ns . }}
