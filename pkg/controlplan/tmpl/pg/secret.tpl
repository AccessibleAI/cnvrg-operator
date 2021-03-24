apiVersion: v1
kind: Secret
metadata:
  name: {{ .Spec.ControlPlan.Pg.SvcName }}
  namespace: {{ ns . }}
data:
  POSTGRESQL_USER: {{ .Spec.ControlPlan.Pg.User | b64enc }}
  POSTGRESQL_PASSWORD: {{ .Spec.ControlPlan.Pg.Pass | b64enc }}
  POSTGRESQL_ADMIN_PASSWORD: {{ .Spec.ControlPlan.Pg.Pass | b64enc }}
  POSTGRESQL_DATABASE: {{ .Spec.ControlPlan.Pg.Dbname | b64enc }}
  POSTGRESQL_MAX_CONNECTIONS: {{ .Spec.ControlPlan.Pg.MaxConnections | toString | b64enc }}
  POSTGRESQL_SHARED_BUFFERS: {{ .Spec.ControlPlan.Pg.SharedBuffers | b64enc }}
  # duplicates for compatibility with webapp/sidekiq
  POSTGRES_DB: {{ .Spec.ControlPlan.Pg.Dbname | b64enc }}
  POSTGRES_PASSWORD: {{ .Spec.ControlPlan.Pg.Pass | b64enc }}
  POSTGRES_USER: {{ .Spec.ControlPlan.Pg.User | b64enc }}
  POSTGRES_HOST: {{ .Spec.ControlPlan.Pg.SvcName | b64enc }}
