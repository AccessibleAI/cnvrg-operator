apiVersion: v1
kind: ConfigMap
metadata:
  name: elastalert-auth-proxy
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  nginx.conf: |
    pid        /tmp/nginx.pid;
    http {
      client_body_temp_path /tmp/client_temp;
      proxy_temp_path       /tmp/proxy_temp_path;
      fastcgi_temp_path     /tmp/fastcgi_temp;
      uwsgi_temp_path       /tmp/uwsgi_temp;
      scgi_temp_path        /tmp/scgi_temp;
      server {
        listen 8080;
        location / {
        auth_basic           "Cnvrg's ElastAlert";
        auth_basic_user_file /etc/nginx/htpasswd/htpasswd;
        proxy_pass           http://localhost:3030/;
        }
      }
    }
    events {}