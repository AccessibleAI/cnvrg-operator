apiVersion: v1
kind: ConfigMap
metadata:
  name: routing-config
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
    user www-data;
    worker_processes 2;
    pid /run/nginx.pid;
    worker_rlimit_nofile 65535;
    error_log /var/log/nginx/error.log info;
    events {
      worker_connections 65535;
      accept_mutex off;
    }

    http {
        default_type application/octet-stream;
        access_log /var/log/nginx/access.log combined;
        sendfile on;

        proxy_ignore_client_abort on;
        proxy_connect_timeout       5000;
        proxy_send_timeout          5000;
        proxy_read_timeout          5000;
        fastcgi_buffers 8 16k;
        fastcgi_buffer_size 32k;
        fastcgi_connect_timeout 300;
        fastcgi_send_timeout 300;
        fastcgi_read_timeout 300;
        send_timeout 5000s;
        keepalive_timeout  5000;
        client_body_timeout   5000;


        map $http_upgrade $connection_upgrade {
          default upgrade;
          ''      close;
        }

        server {
          add_header 'Access-Control-Allow-Origin' '*';
          add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS';
          add_header 'Access-Control-Allow-Headers' 'DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Content-Range,Range';
          add_header X-Frame-Options ALLOWALL;
          root /var/www/html;
          index index.html index.htm index.nginx-debian.html;
          error_log /var/log/nginx/error.log debug;
          server_name routing.app.cnvrg.io;
          location ~* "/(.+?)/projects/(.+?)/notebook_sessions/view/(\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})pp(\d{2})[a-z](\d{2})[a-z](?<jup>/?.*)" {
          set $jup_url 'http://$3.$4.$5.$6:$7$8';
          rewrite ^/AccessibleAI/projects/test/notebook_sessions/view/aaaa/$jup break;
          proxy_pass $jup_url;
          client_max_body_size 1G;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header Host $host;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
          proxy_read_timeout 20d;
          proxy_http_version 1.1;
          proxy_set_header Upgrade $http_upgrade;
          proxy_set_header Connection "upgrade";
        }

        location ~* "/(.+?)/projects/(.+?)/notebook_sessions/view/(\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})pp(\d{2})[a-z](\d{2})/(api/kernels/[^/]+/(channels|iopub|shell|stdin)|terminals/websocket)" {
          proxy_pass $jup_url;
          client_max_body_size 1G;
          proxy_http_version    1.1;
          proxy_set_header      Upgrade "websocket";
          proxy_set_header      Connection "Upgrade";
          proxy_read_timeout 20d;
        }

        location ~* "/(.+?)/projects/(.+?)/r_studio_sessions/view/((\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})oo(\d{4})[a-z])/(.*)" {
          set $rstudio_path '$3';
          set $rstudio_url 'http://$4.$5.$6.$7:$8';
          set $suffix '$9';
          rewrite ^/(.+?)/projects/(.+?)/r_studio_sessions/view/([0-9a-zA-Z]+)/(.*)$ /$4 break;
          proxy_pass $rstudio_url;
          proxy_redirect $rstudio_url/ $scheme://$host/$rstudio_path/;
          proxy_hide_header X-Frame-Options;
          proxy_http_version 1.1;
          proxy_set_header Upgrade $http_upgrade;
          proxy_set_header Connection $connection_upgrade;
          proxy_read_timeout 20d;
        }

        location ~* "/(.+?)/projects/(.+?)/terminal_sessioens/view/(\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})pp(\d{2})[a-z](\d{2})[a-z](.*)" {
          set $jup_url 'http://$3.$4.$5.$6:$7$8$9$is_args$query_string';
          add_header X-debug-message $jup_url always;
          proxy_pass $jup_url;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header Host $host;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
          proxy_read_timeout 20d;
          proxy_http_version 1.1;
          proxy_set_header Upgrade $http_upgrade;
          proxy_set_header Connection "upgrade";
        }

        location ~* "/(.+?)/projects/(.+?)/terminal_sessions/view/(\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})pp(\d{2})[a-z](\d{2})[a-z]/wetty/(.*)" {
          set $jup_url 'http://$3.$4.$5.$6:$7$8/wetty/$9$is_args$query_string';
          add_header X-debug-message $jup_url always;
          proxy_pass $jup_url;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header Host $host;
          add_header X-Frame-Options SAMEORIGIN;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
          proxy_read_timeout 20d;
          proxy_http_version 1.1;
          proxy_set_header Upgrade $http_upgrade;
          proxy_set_header Connection "upgrade";
          proxy_set_header X-NginX-Proxy true;
        }

        location ~* "/static/(.*)" {
          proxy_pass $http_referer/static/$1$is_args$query_string;
          proxy_set_header Accept-Encoding "";
          sub_filter "/static/" "/app/spark/master/static/";
          sub_filter_once off;
          add_header X-debug-messagr $http_cookie always;
          proxy_set_header Host $host;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
          proxy_http_version 1.1;
          proxy_set_header Upgrade $http_upgrade;
          proxy_set_header Connection "upgrade";
          proxy_read_timeout 3000;
          proxy_send_timeout 3000;
          proxy_set_header X-NginX-Proxy true;
        }

        location ~* "/wetty/(.*)" {
          add_header X-debug-message $1 always;
          proxy_pass $http_referer/$1$is_args$query_string;
          add_header X-debug-messagr $http_cookie always;
          proxy_set_header Host $host;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
          proxy_http_version 1.1;
          proxy_set_header Upgrade $http_upgrade;
          proxy_set_header Connection "upgrade";
          proxy_read_timeout 3000;
          proxy_send_timeout 3000;
          proxy_set_header X-NginX-Proxy true;
        }

        location ~* "/voila/(.*)" {
          add_header X-debug-message $1 always;
          proxy_pass $http_referer/$1$is_args$query_string;
          add_header X-debug-messagr $http_cookie always;
          proxy_set_header Host $host;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
          proxy_http_version 1.1;
          proxy_set_header Upgrade $http_upgrade;
          proxy_set_header Connection "upgrade";
          proxy_read_timeout 3000;
          proxy_send_timeout 3000;
          proxy_set_header X-NginX-Proxy true;
        }

        location ~* "/((\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})oo(\d{4})[a-z])/(.*)" {
          set $rstudio_path '$1';
          set $rstudio_url 'http://$2.$3.$4.$5:$6';
          set $suffix '$7';
          rewrite ^/([0-9a-zA-Z]+)/(.*)$ /$2 break;
          proxy_pass $rstudio_url;
          proxy_redirect $rstudio_url/ $scheme://$host/$rstudio_path/;
          proxy_hide_header X-Frame-Options;
          proxy_http_version 1.1;
          proxy_set_header Upgrade $http_upgrade;
          proxy_set_header Connection $connection_upgrade;
          proxy_read_timeout 20d;
        }

        location ~* "/(\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})pp(\d{2})[a-z](\d{2})[a-z]/(.*)" {
          proxy_pass http://$1.$2.$3.$4:$5$6/$7$is_args$query_string;
          proxy_set_header Host $host;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
          proxy_http_version 1.1;
          proxy_set_header Upgrade $http_upgrade;
          proxy_set_header Connection "upgrade";
          proxy_read_timeout 86400;
        }

        location ~* "/data/(.*)" {
          proxy_pass $http_referer/data/$1$is_args$query_string;
          proxy_set_header Host $host;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
          proxy_http_version 1.1;
          proxy_set_header Upgrade $http_upgrade;
          proxy_set_header Connection "upgrade";
          proxy_read_timeout 86400;
        }

        location ~* "/(\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})/(.*)" {
          proxy_pass http://$1.$2.$3.$4/$5;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header Host $host;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
          proxy_http_version 1.1;
          proxy_set_header Upgrade $http_upgrade;
          proxy_set_header Connection "upgrade";
          proxy_read_timeout 86400;
          proxy_set_header 'Access-Control-Max-Age' 1728000;
          add_header 'Access-Control-Allow-Origin' "$http_origin" always;
          add_header 'Access-Control-Allow-Credentials' 'true' always;
          add_header Access-Control-Allow-Headers "cnvrg-api-key";
          add_header Access-Control-Allow-Headers "content-type";
        }

        location ~* "/(\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})pp(\d{1})[a-z](\d{1})[a-z]/(.*)" {
          proxy_pass http://$1.$2.$3.$4:$5$6/$7;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header Host $host;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
          proxy_http_version 1.1;
          proxy_set_header Upgrade $http_upgrade;
          proxy_set_header Connection "upgrade";
          proxy_read_timeout 86400;
          proxy_set_header 'Access-Control-Max-Age' 1728000;
          add_header 'Access-Control-Allow-Origin' "$http_origin" always;
          add_header 'Access-Control-Allow-Credentials' 'true' always;
          add_header Access-Control-Allow-Headers "cnvrg-api-key";
          add_header Access-Control-Allow-Headers "content-type";
        }

        location ~* "/(.+?)/projects/(.+?)/terminal_sessions/view/(\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})pp(\d{2})[a-z](\d{2})[a-z](.*)" {
          set $jup_url 'http://$3.$4.$5.$6:$7$8$9$is_args$query_string';
          add_header X-debug-message $jup_url always;
          proxy_pass $jup_url;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header Host $host;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
          proxy_read_timeout 20d;
          proxy_http_version 1.1;
          proxy_set_header Upgrade $http_upgrade;
          proxy_set_header Connection "upgrade";
        }

        location ~* "/(.+?)/projects/(.+?)/r_shiny/view/(\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})[a-z](\d{1,3})rs(\d{2})[a-z](\d{2})[a-z](.*)" {
          set $jup_url 'http://$3.$4.$5.$6:$7$8$9$is_args$query_string';
          add_header X-debug-message $jup_url always;
          proxy_pass $jup_url;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header Host $host;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
          proxy_read_timeout 20d;
          proxy_http_version 1.1;
          proxy_set_header Upgrade $http_upgrade;
          proxy_set_header Connection "upgrade";
        }

        location / {
          try_files $uri $uri/ =404;
        }

        listen 80 default_server;
        listen [::]:80 default_server;
        resolver 8.8.8.8 8.8.4.4 valid=300s;
        resolver_timeout 5s;
        client_max_body_size 1G;
      }
    }
