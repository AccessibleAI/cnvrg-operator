apiVersion: v1
kind: Secret
metadata:
  name: {{ .Spec.Pg.SvcName }}
  namespace: {{ .Namespace }}
data:
  POSTGRESQL_USER: {{ .Spec.Pg.User | b64enc }}
  POSTGRESQL_PASSWORD: {{ .Spec.Pg.Pass | b64enc }}
  POSTGRESQL_ADMIN_PASSWORD: {{ .Spec.Pg.Pass | b64enc }}
  POSTGRESQL_DATABASE: {{ .Spec.Pg.Dbname | b64enc }}
  POSTGRESQL_MAX_CONNECTIONS: {{ .Spec.Pg.MaxConnections | toString | b64enc }}
  POSTGRESQL_SHARED_BUFFERS: {{ .Spec.Pg.SharedBuffers | b64enc }}
  # duplicates for compatibility with webapp/sidekiq
  POSTGRES_DB: {{ .Spec.Pg.Dbname | b64enc }}
  POSTGRES_PASSWORD: {{ .Spec.Pg.Pass | b64enc }}
  POSTGRES_USER: {{ .Spec.Pg.User | b64enc }}
  POSTGRES_HOST: {{ .Spec.Pg.SvcName | b64enc }}
