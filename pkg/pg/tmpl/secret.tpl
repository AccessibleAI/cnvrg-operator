apiVersion: v1
kind: Secret
metadata:
  name: pg-secret
  namespace: {{ .Spec.CnvrgNs }}
data:
  POSTGRESQL_USER: {{ .Spec.Pg.User | b64enc }}
  POSTGRESQL_PASSWORD: {{ .Spec.Pg.Pass | b64enc }}
  POSTGRESQL_ADMIN_PASSWORD: {{ .Spec.Pg.Pass | b64enc }}
  POSTGRESQL_DATABASE: {{ .Spec.Pg.Dbname | b64enc }}
  POSTGRESQL_MAX_CONNECTIONS: {{ .Spec.Pg.MaxConnections | toString | b64enc }}
  POSTGRESQL_SHARED_BUFFERS: {{ .Spec.Pg.SharedBuffers | b64enc }}
