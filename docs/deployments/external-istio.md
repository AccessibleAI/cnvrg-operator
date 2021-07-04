# External Istio

#### Use Cnvrg Control Plane with external Istio  

Cnvrg Control Plane (CCP) deploys its own Istio instance.
However, cnvrg administrator may change this behavior and instruct cnvrg to not deploy its own Istio instance
but instead use existing one. 

Pass the following flags to use external Istio instance   

```shell
...
 --set networking.istio.enabled=false \
 --set networking.ingress.istioGwEnabled=false \
 --set networking.ingress.istioGwName=<desired-istio-gw-name-to-use-in-vs> \
...
```