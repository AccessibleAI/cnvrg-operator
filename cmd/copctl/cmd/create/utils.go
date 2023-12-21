package create

import (
	"crypto/rand"
	"crypto/rsa"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func clientset() *kubernetes.Clientset {
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

func privateKey() *rsa.PrivateKey {
	zap.S().Info("generating private key")
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		zap.S().Fatal(err)
	}
	return privKey
}
