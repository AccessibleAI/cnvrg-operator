### Download and install `cnvrgctl`
* mac: [cnvrgctl-darwin-x86_64](https://cnvrg-public-images.s3-us-west-2.amazonaws.com/cnvrgctl-darwin-x86_64)
  ```shell
  curl -#o /usr/local/bin/cnvrgctl \
    https://cnvrg-public-images.s3-us-west-2.amazonaws.com/cnvrgctl-darwin-x86_64 \
  && chmod +x /usr/local/bin/cnvrgctl \
  && cnvrgctl completion bash > /usr/local/etc/bash_completion.d/cnvrgctl
  ```
* linux: [cnvrgctl-linux-x86_64](https://cnvrg-public-images.s3-us-west-2.amazonaws.com/cnvrgctl-linux-x86_64)
  ```shell
  curl -#o /usr/local/bin/cnvrgctl \
    https://cnvrg-public-images.s3-us-west-2.amazonaws.com/cnvrgctl-linux-x86_64 \
  && chmod +x /usr/local/bin/cnvrgctl \
  && cnvrgctl completion bash > /etc/bash_completion.d/cnvrgctl
  ```

#### Deploy all-in-one single node K8s cluster for cnvrg
Prerequisite
1. VM or Bare metal Ubuntu 20.04 server with 32 CPUs, 64GB memory, 500GB storage
2. In case of VM, [bridged network (preferred) or nat network](https://superuser.com/questions/227505/what-is-the-difference-between-nat-bridged-host-only-networking) between VM and the host
3. either root user or regular user with sudo access
4. SSH access to the server either by ssh key or password

Deploy single node K8s cluster for cnvrg deployment
```shell
# access the server with ssh password  
cnvrgctl cluster up --host=<SERVER-IP> --ssh-user=<SSH-USER> --ssh-pass=<SSH-PASS>
# access the server with ssh key  
cnvrgctl cluster up --host=<SERVER-IP> --ssh-user=<SSH-USER> --ssh-pass=</path/to/private/key>
```

Cleanup
```shell
# access the server with ssh password  
cnvrgctl cluster destroy --host=<SERVER-IP> --ssh-user=<SSH-USER> --ssh-pass=<SSH-PASS>
# access the server with ssh key  
cnvrgctl cluster destroy --host=<SERVER-IP> --ssh-user=<SSH-USER> --ssh-pass=</path/to/private/key>
```