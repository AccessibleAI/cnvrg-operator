apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  name: cnvrg-istio
  namespace:  {{ .Spec.CnvrgNs }}
spec:
  profile: minimal
  namespace:  {{ .Spec.CnvrgNs }}
  hub: {{ .Spec.Networking.Istio.Hub }}
  tag: {{ .Spec.Networking.Istio.Tag }}
  values:
    global:
      istioNamespace:  {{ .Spec.CnvrgNs }}
      imagePullSecrets:
        - {{ .Spec.ControlPlan.Conf.Registry.Name }}
    meshConfig:
      rootNamespace:  {{ .Spec.CnvrgNs }}
  components:
    base:
      enabled: true
    pilot:
      enabled: true
      k8s:
        {{- if eq .Spec.ControlPlan.Conf.Tenancy.Enabled "true" }}
        nodeSelector:
          {{ .Spec.ControlPlan.Conf.Tenancy.Key }}: "{{ .Spec.ControlPlan.Conf.Tenancy.Value }}"
        {{- end }}
        tolerations:
        - key: "{{ .Spec.ControlPlan.Conf.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.ControlPlan.Conf.Tenancy.Value }}"
          effect: "NoSchedule"
    ingressGateways:
    - enabled: true
      name: istio-ingressgateway
      k8s:
        {{- if eq .Spec.ControlPlan.Conf.Tenancy.Enabled "true" }}
        nodeSelector:
          {{ .Spec.ControlPlan.Conf.Tenancy.Key }}: "{{ .Spec.ControlPlan.Conf.Tenancy.Value }}"
        {{- end }}
        tolerations:
        - key: "{{ .Spec.ControlPlan.Conf.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.ControlPlan.Conf.Tenancy.Value }}"
          effect: "NoSchedule"
        {{- if ne .Spec.Networking.Istio.IngressSvcAnnotations "" }}
        serviceAnnotations:
          {{- $annotations := split ";" .Spec.Networking.Istio.IngressSvcAnnotations }}
            {{- range $idx, $annotation := $annotations }}
              {{- $annotationItem := split ":" $annotation}}
              {{- if eq (len $annotationItem) 2 }}
          {{ trim $annotationItem._0 }}: {{ trim $annotationItem._1 }}
              {{- end }}
            {{- end }}
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
          {{- if ne .Spec.Networking.Istio.LoadBalancerSourceRanges "" }}
          loadBalancerSourceRanges:
            {{- $srouceRanges := split ";" .Spec.Networking.Istio.LoadBalancerSourceRanges }}
            {{- range $idx, $range := $srouceRanges }}
              {{- if ne (trim $range) "" }}
            - {{trim $range}}
              {{- end }}
            {{- end }}
          {{- end }}
          {{- if ne .Spec.Networking.Istio.ExternalIP "" }}
          type: ClusterIP
          externalIPs:
          {{- $ips := split ";" .Spec.Networking.Istio.ExternalIP }}
          {{- range $idx, $ip := $ips }}
            {{- if ne (trim $ip) "" }}
            - {{trim $ip}}
            {{- end }}
          {{- end }}
          {{- end }}
          ports:
          - name: http2
            port: 80
            targetPort: 8080
          - name: https
            port: 443
            targetPort: 8443
          {{- if ne .Spec.Networking.Istio.IngressSvcExtraPorts "" }}
          {{- $ports := split ";" .Spec.Networking.Istio.IngressSvcExtraPorts }}
          {{- range $idx, $port := $ports }}
            {{- if ne (trim $port) "" }}
          - name: port{{trim $port}}
            port: {{trim $port}}
            {{- end }}
          {{- end }}
          {{- end}}
        strategy:
          rollingUpdate:
            maxSurge: 100%
            maxUnavailable: 25%