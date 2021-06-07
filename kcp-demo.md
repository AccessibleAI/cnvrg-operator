# cnvrg.io operator (v3) - KCP Demo
---
0. Install prerequisite
```shell
sudo apt update -y \
 && sudo apt install dnsutils ubuntu-dev-tools ruby-full python3 python3-pip python3-testresources -y \
 && gem install cnvrg \
 && curl -# -L0 -o /usr/local/bin/kubectl https://dl.k8s.io/release/v1.21.0/bin/linux/amd64/kubectl \
 && chmod +x /usr/local/bin/kubectl \
 && wget https://cnvrg-public-data.s3-us-west-2.amazonaws.com/cnvrg-cli-kcp.gem \
 && gem install cnvrg-cli-kcp.gem \
 && curl -# -o /usr/local/bin/tiny https://cnvrg-public-data.s3-us-west-2.amazonaws.com/tiny \
 && chmod +x /usr/local/bin/tiny \
 && pip3 install watchdog prometheus_client \
 && pip3 install --upgrade tensorflow 
```

1. Download KCP linux binary from cnvrg public bucket
```shell
curl -# -o /usr/local/bin/kcp https://cnvrg-public-data.s3-us-west-2.amazonaws.com/kcp \
 && chmod +x /usr/local/bin/kcp
```

2. start KCP server 
```shell
mkdir -p ~/kcp \
 && cd ~/kcp \
 && rm -fr .kcp \
 && kcp start --listen=$(dig +short $(hostname -f)):6443
```

3. add new K8s cluster to cnvrg control plane (copy/paste kubeconfig) 
```shell
cat ~/kcp/.kcp/data/admin.kubeconfig
```

4. Prepare KCP for OnPremExecutor startup 
```shell
mkdir -p ~/.kube \
 && cp ~/kcp/.kcp/data/admin.kubeconfig ~/.kube/config \
 && sed -i "s/$(dig +short $(hostname -f))/apiserver-loopback-client/g" ~/.kube/config \
 && sudo echo "$(dig +short $(hostname -f)) apiserver-loopback-client" >> /etc/hosts \
 && kubectl apply -f https://cnvrg-public-data.s3-us-west-2.amazonaws.com/executer-crd.yaml \
 && kubectl create ns cnvrg
```


4. Download and start cnvrg OnPremExecutor
```shell
sudo curl -# -o /usr/local/bin/cnvrg-operator https://cnvrg-public-data.s3-us-west-2.amazonaws.com/cnvrg-operator \
 && sudo chmod +x /usr/local/bin/cnvrg-operator \
 && cnvrg-operator start
``` 