#!/bin/bash

dump_job_env_vars(){
  rm -fr /root/envs
  touch /root/envs
  {{- range $idx, $v := .Data.JobEnvVars }}
  echo "{{$v}}" >> /root/envs
  {{- end }}
}

cleanup(){
  rm -fr /cnvrg && mkdir /cnvrg
  rm -fr /data && mkdir /data
  rm -fr /script && mkdir /script
  rm -fr /conf && mkdir /conf
  if [[ $(pgrep -f tiny | wc -l) > 0 ]]; then
    pgrep -f tiny | xargs kill -9
  fi
  if [[ $(pgrep -f metrics | wc -l) > 0 ]]; then
    pgrep -f metrics | xargs kill -9
  fi
}

start_cnvrg_tiny_server(){
  cd /cnvrg && /conf/tiny &
}

start_cnvrg_job(){
  cd /cnvrg && cnvrg job start
}

. /root/envs

dump_job_env_vars
cleanup
start_cnvrg_tiny_server
start_cnvrg_job
