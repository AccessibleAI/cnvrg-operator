apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  name: cnvrg-istio
  namespace:  {{ ns . }}
spec:
  profile: minimal
  namespace:  {{ ns . }}
  hub: {{ .Spec.Networking.Istio.Hub }}
  tag: {{ .Spec.Networking.Istio.Tag }}
  values:
    global:
      istioNamespace:  {{ ns . }}
      imagePullSecrets:
        - {{ .Spec.Registry.Name }}
    meshConfig:
      rootNamespace:  {{ ns . }}
  components:
    base:
      enabled: true
    pilot:
      enabled: true
    ingressGateways:
    - enabled: true
      name: istio-ingressgateway
      k8s:
        {{- if isTrue .Spec.Tenancy.Enabled }}
        nodeSelector:
          {{ .Spec.Tenancy.Key }}: "{{ .Spec.Tenancy.Value }}"
        tolerations:
          - key: "{{ .Spec.Tenancy.Key }}"
            operator: "Equal"
            value: "{{ .Spec.Tenancy.Value }}"
            effect: "NoSchedule"
        {{- end }}
        serviceAnnotations:
        {{- range $name, $value := .Spec.Networking.Istio.IngressSvcAnnotations }}
          {{ $name }}: {{ $value }}
        {{- end }}
        env:
          - name: ISTIO_META_ROUTER_MODE
            value: sni-dnat
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
            name: istio-ingressgateway
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
        strategy:
          rollingUpdate:
            maxSurge: 100%
            maxUnavailable: 25%