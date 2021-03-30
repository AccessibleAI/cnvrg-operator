apiVersion: v1
kind: Secret
metadata:
  name: {{ .Spec.Dbs.Pg.SvcName }}
  namespace: {{ ns . }}
data:
  POSTGRESQL_USER: {{ .Spec.Dbs.Pg.User | b64enc }}
  POSTGRESQL_PASSWORD: {{ .Spec.Dbs.Pg.Pass | b64enc }}
  POSTGRESQL_ADMIN_PASSWORD: {{ .Spec.Dbs.Pg.Pass | b64enc }}
  POSTGRESQL_DATABASE: {{ .Spec.Dbs.Pg.Dbname | b64enc }}
  POSTGRESQL_MAX_CONNECTIONS: {{ .Spec.Dbs.Pg.MaxConnections | toString | b64enc }}
  POSTGRESQL_SHARED_BUFFERS: {{ .Spec.Dbs.Pg.SharedBuffers | b64enc }}
  # duplicates for compatibility with webapp/sidekiq
  POSTGRES_DB: {{ .Spec.Dbs.Pg.Dbname | b64enc }}
  POSTGRES_PASSWORD: {{ .Spec.Dbs.Pg.Pass | b64enc }}
  POSTGRES_USER: {{ .Spec.Dbs.Pg.User | b64enc }}
  POSTGRES_HOST: {{ .Spec.Dbs.Pg.SvcName | b64enc }}
