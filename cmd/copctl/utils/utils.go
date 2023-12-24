package utils

import (
	"crypto/rand"
	"crypto/rsa"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"go.uber.org/zap"
	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func Clientset() *kubernetes.Clientset {
	rc, err := config.GetConfig()
	if err != nil {
		zap.S().Fatalf("unable to construct K8s configs, err: %s", err.Error())
	}
	cltset, err := kubernetes.NewForConfig(rc)
	if err != nil {
		zap.S().Fatalf("err: %s, unable to connect to K8s", err.Error())
	}
	return cltset
}

func Kubecrudclient() client.Client {

	rc, err := config.GetConfig()
	if err != nil {
		zap.S().Fatal(err)
	}

	scheme := runtime.NewScheme()
	utilruntime.Must(v1core.AddToScheme(scheme))
	utilruntime.Must(mlopsv1.AddToScheme(scheme))

	cc, err := client.New(rc, client.Options{Scheme: scheme})
	if err != nil {
		zap.S().Fatal(err)
	}
	return cc
}

func PrivateKey() *rsa.PrivateKey {
	zap.S().Info("generating private key")
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		zap.S().Fatal(err)
	}
	return privKey
}
