apiVersion: v1
kind: Secret
metadata:
  name: {{ .Pg.SvcName }}
  namespace: {{ .CnvrgNs }}
data:
  POSTGRESQL_USER: {{ .Pg.User | b64enc }}
  POSTGRESQL_PASSWORD: {{ .Pg.Pass | b64enc }}
  POSTGRESQL_ADMIN_PASSWORD: {{ .Pg.Pass | b64enc }}
  POSTGRESQL_DATABASE: {{ .Pg.Dbname | b64enc }}
  POSTGRESQL_MAX_CONNECTIONS: {{ .Pg.MaxConnections | toString | b64enc }}
  POSTGRESQL_SHARED_BUFFERS: {{ .Pg.SharedBuffers | b64enc }}
  # duplicates for compatibility with webapp/sidekiq
  POSTGRES_DB: {{ .Pg.Dbname | b64enc }}
  POSTGRES_PASSWORD: {{ .Pg.Pass | b64enc }}
  POSTGRES_USER: {{ .Pg.User | b64enc }}
  POSTGRES_HOST: {{ .Pg.SvcName | b64enc }}
