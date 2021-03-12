apiVersion: v1
kind: Secret
metadata:
  name: cp-ldap
  namespace: {{ .CnvrgNs }}
data:
  LDAP_ACTIVE: {{ .ControlPlan.Conf.Ldap.Enabled | b64enc }}
  LDAP_HOST: {{ .ControlPlan.Conf.Ldap.Host | b64enc }}
  LDAP_PORT: {{ .ControlPlan.Conf.Ldap.Port | b64enc }}
  LDAP_SSL: {{ .ControlPlan.Conf.Ldap.Ssl | b64enc }}
  LDAP_ACCOUNT: {{ .ControlPlan.Conf.Ldap.Account | b64enc }}
  LDAP_BASE: {{ .ControlPlan.Conf.Ldap.Base | b64enc }}
  LDAP_ADMIN_USER: {{ .ControlPlan.Conf.Ldap.AdminUser | b64enc }}
  LDAP_ADMIN_PASSWORD: {{ .ControlPlan.Conf.Ldap.AdminPassword | b64enc }}



