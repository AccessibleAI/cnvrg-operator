istio:
	# deploy istio CRDs & istiod
	helmfile apply -lstage=istio-infra --wait
	# deploy istio ingress controller
	helmfile apply -lstage=istio-gateway --wait
nginx:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.10.1/deploy/static/provider/cloud/deploy.yaml