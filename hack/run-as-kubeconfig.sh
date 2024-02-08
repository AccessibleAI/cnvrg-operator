export SECRET_NAME_SA=cnvrg-operator
export TOKEN_SA=`kubectl get secret cnvrg-operator-secret-debug -n cnvrg -ojsonpath='{.data.token}' | base64 -d`
kubectl config view --raw --minify > kubeconfig.txt
kubectl config unset users --kubeconfig=kubeconfig.txt
kubectl config set-credentials ${SECRET_NAME_SA} --kubeconfig=kubeconfig.txt --token=${TOKEN_SA}
kubectl config set-context --current --kubeconfig=kubeconfig.txt --user=${SECRET_NAME_SA}