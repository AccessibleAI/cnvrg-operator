# cnvrg.io operator (v3) - KCP Demo
---
0. Install prerequisite
```shell
sudo apt install dnsutils ruby-full
```

1. Download KCP linux binary from cnvrg public bucket
```shell
curl -# -o /usr/local/bin/kcp https://cnvrg-public-data.s3-us-west-2.amazonaws.com/kcp
chmod +x /usr/local/bin/kcp
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
mkdir -p ~/.kube
cp ~/kcp/.kcp/data/admin.kubeconfig ~/.kube/config
sed -i "s/$(dig +short $(hostname -f))/apiserver-loopback-client/g" ~/.kube/config
sudo echo "$(dig +short $(hostname -f)) apiserver-loopback-client" >> /etc/hosts
kubectl apply -f https://cnvrg-public-data.s3-us-west-2.amazonaws.com/executer-crd.yaml
kubectl create ns cnvrg
```


4. Download and start cnvrg OnPremExecutor
```shell
sudo curl -# -o /usr/local/bin/cnvrg-operator https://cnvrg-public-data.s3-us-west-2.amazonaws.com/cnvrg-operator
sudo chmod +x /usr/local/bin/cnvrg-operator
```