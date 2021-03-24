apiVersion: v1
kind: Secret
metadata:
  name: cp-ldap
  namespace: {{ ns . }}
data:
  LDAP_ACTIVE: {{ .Spec.ControlPlan.Ldap.Enabled | b64enc }}
  LDAP_HOST: {{ .Spec.ControlPlan.Ldap.Host | b64enc }}
  LDAP_PORT: {{ .Spec.ControlPlan.Ldap.Port | b64enc }}
  LDAP_SSL: {{ .Spec.ControlPlan.Ldap.Ssl | b64enc }}
  LDAP_ACCOUNT: {{ .Spec.ControlPlan.Ldap.Account | b64enc }}
  LDAP_BASE: {{ .Spec.ControlPlan.Ldap.Base | b64enc }}
  LDAP_ADMIN_USER: {{ .Spec.ControlPlan.Ldap.AdminUser | b64enc }}
  LDAP_ADMIN_PASSWORD: {{ .Spec.ControlPlan.Ldap.AdminPassword | b64enc }}



