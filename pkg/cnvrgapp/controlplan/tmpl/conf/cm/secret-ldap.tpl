apiVersion: v1
kind: Secret
metadata:
  name: cp-ldap
  namespace: {{ .CnvrgNs }}
data:
  LDAP_ACTIVE: {{ .ControlPlan.Ldap.Enabled | b64enc }}
  LDAP_HOST: {{ .ControlPlan.Ldap.Host | b64enc }}
  LDAP_PORT: {{ .ControlPlan.Ldap.Port | b64enc }}
  LDAP_SSL: {{ .ControlPlan.Ldap.Ssl | b64enc }}
  LDAP_ACCOUNT: {{ .ControlPlan.Ldap.Account | b64enc }}
  LDAP_BASE: {{ .ControlPlan.Ldap.Base | b64enc }}
  LDAP_ADMIN_USER: {{ .ControlPlan.Ldap.AdminUser | b64enc }}
  LDAP_ADMIN_PASSWORD: {{ .ControlPlan.Ldap.AdminPassword | b64enc }}



