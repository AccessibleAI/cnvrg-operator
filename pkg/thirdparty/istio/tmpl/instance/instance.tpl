apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  name: cnvrg-istio
  namespace:  {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "false"
    mlops.cnvrg.io/updatable: "false"
spec:
  profile: minimal
  tag: {{ trimPrefix "pilot:" .Spec.Istio.PilotImage }}
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
      k8s:
        serviceAnnotations:
        {{- range $name, $value := .Spec.Istio.IngressSvcAnnotations }}
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
          {{- range $idx, $range := .Spec.Istio.LBSourceRanges }}
            - {{ $range}}
          {{- end }}
          {{- if gt (len .Spec.Istio.ExternalIP) 0 }}
          type: ClusterIP
          externalIPs:
          {{- range $idx, $ip := .Spec.Istio.ExternalIP }}
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
          {{- if gt (len .Spec.Istio.IngressSvcExtraPorts) 0 }}
          {{- range $idx, $port := .Spec.Istio.IngressSvcExtraPorts }}
          - name: port{{ $port}}
            port: {{ $port}}
            {{- end }}
          {{- end }}
    istiodRemote:
      enabled: false
    pilot:
      enabled: true
  values:
    global:
      istioNamespace:  {{ .Namespace }}
      imagePullSecrets:
        - {{ .Spec.Registry.Name }}
    meshConfig:
      rootNamespace: {{ .Namespace }}
