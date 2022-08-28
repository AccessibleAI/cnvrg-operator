apiVersion: v1
kind: Secret
metadata:
  name: cp-ldap
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-ccp"
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  LDAP_ACTIVE: {{ isTrue .Spec.ControlPlane.Ldap.Enabled | toString | b64enc }}
  LDAP_HOST: {{ .Spec.ControlPlane.Ldap.Host | b64enc }}
  LDAP_PORT: {{ .Spec.ControlPlane.Ldap.Port | b64enc }}
  LDAP_SSL: {{ .Spec.ControlPlane.Ldap.Ssl | b64enc }}
  LDAP_ACCOUNT: {{ .Spec.ControlPlane.Ldap.Account | b64enc }}
  LDAP_BASE: {{ .Spec.ControlPlane.Ldap.Base | b64enc }}
  LDAP_ADMIN_USER: {{ .Spec.ControlPlane.Ldap.AdminUser | b64enc }}
  LDAP_ADMIN_PASSWORD: {{ .Spec.ControlPlane.Ldap.AdminPassword | b64enc }}



