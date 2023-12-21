package create

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	certsv1 "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

func init() {
	kubeCertsCmd.PersistentFlags().StringP("common-name", "c", "",
		"certificate common name (will be included to alternative DNS names as well)")
	kubeCertsCmd.PersistentFlags().BoolP("persist-on-disk", "d", true,
		"save certs as files locally")
	kubeCertsCmd.PersistentFlags().BoolP("persist-on-k8s", "k", true,
		"create K8s tls secret")
	kubeCertsCmd.PersistentFlags().BoolP("override", "w", true,
		"override csr,files and secrets if they already exists")
	kubeCertsCmd.PersistentFlags().StringP("certs-dir", "f", "certs",
		"path to dump key & crt")

	viper.BindPFlag("common-name", kubeCertsCmd.PersistentFlags().Lookup("common-name"))
	viper.BindPFlag("persist-on-disk", kubeCertsCmd.PersistentFlags().Lookup("persist-on-disk"))
	viper.BindPFlag("persist-on-k8s", kubeCertsCmd.PersistentFlags().Lookup("persist-on-k8s"))
	viper.BindPFlag("override", kubeCertsCmd.PersistentFlags().Lookup("override"))
	viper.BindPFlag("certs-dir", kubeCertsCmd.PersistentFlags().Lookup("certs-dir"))

	Cmd.AddCommand(kubeCertsCmd)
}

var kubeCertsCmd = &cobra.Command{
	Use:     "kube-certs",
	Aliases: []string{"c"},
	Short:   "create K8s signed certificate",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("override") {
			clean(viper.GetString("certs-dir"), viper.GetString("common-name"))
		}
		pkey := privateKey()
		approveCsr(
			createCsr(
				csrPem(
					pkey,
					viper.GetString("common-name"),
				),
				viper.GetString("common-name"),
			),
		)
		dumpKeyAndCert(pkey, viper.GetString("certs-dir"), viper.GetString("common-name"))

	},
}

func clean(certsDir string, commonName string) {
	zap.S().Info("cleaning up exiting certs")
	err := clientset().
		CertificatesV1().
		CertificateSigningRequests().
		Delete(context.Background(), commonName, metav1.DeleteOptions{})

	if err != nil && !errors.IsNotFound(err) {
		zap.S().Fatalf("err: %s, failed to delete csr: %s", err.Error(), commonName)
	}

	if err = os.RemoveAll(certsDir); err != nil {
		zap.S().Error(err)
	}

}

func dumpKeyAndCert(pkey *rsa.PrivateKey, certsDir, commonName string) {
	pkeyPem := new(bytes.Buffer)
	_ = pem.Encode(pkeyPem,
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(pkey),
		},
	)
	// create dir for cert and key
	if err := os.MkdirAll(certsDir, 0755); err != nil {
		zap.S().Fatal(err)
	}
	// dump key to file
	if err := os.WriteFile(certsDir+"/server.key", pkeyPem.Bytes(), 0644); err != nil {
		zap.S().Fatal(err)
	}
	// dump certificate to file
	if err := os.WriteFile(certsDir+"/server.crt", fetchCertificateFromCsr(commonName), 0644); err != nil {
		zap.S().Fatal(err)
	}

}

func csrPem(pKey *rsa.PrivateKey, commonName string) *bytes.Buffer {
	zap.S().Info("creating csr pem")
	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"system:nodes"},
		},
		DNSNames: []string{"cnvrg-operator-admission.cnvrg-shim.svc"},
	}
	csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, pKey)
	clientCsrPem := new(bytes.Buffer)
	if err := pem.Encode(clientCsrPem, &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	}); err != nil {
		zap.S().Fatal(err)
	}
	return clientCsrPem
}

func createCsr(csrPem *bytes.Buffer, commonName string) *certsv1.CertificateSigningRequest {
	zap.S().Infof("creating CertificateSigningRequest %s", commonName)
	//signerName := "kubernetes.io/kube-apiserver-client-kubelet"
	//signerName := "kubernetes.io/kube-apiserver-client"
	signerName := "kubernetes.io/kubelet-serving"
	csr, err := clientset().CertificatesV1().
		CertificateSigningRequests().
		Create(context.Background(), &certsv1.CertificateSigningRequest{
			ObjectMeta: metav1.ObjectMeta{Name: commonName},
			Spec: certsv1.CertificateSigningRequestSpec{
				SignerName: signerName,
				Request:    csrPem.Bytes(),
				Usages: []certsv1.KeyUsage{
					certsv1.UsageDigitalSignature,
					certsv1.UsageKeyEncipherment,
					certsv1.UsageServerAuth,
				},
				Groups: []string{"system:authenticated"},
				//Groups: []string{"system:nodes"},
			},
		}, metav1.CreateOptions{})
	if err != nil {
		zap.S().Fatal(err)
	}
	return csr
}

func approveCsr(csr *certsv1.CertificateSigningRequest) {
	zap.S().Infof("appriving  csr %s", csr.Name)
	csr.Status.Conditions = append(csr.Status.Conditions, certsv1.CertificateSigningRequestCondition{
		Status:             corev1.ConditionTrue,
		Type:               certsv1.CertificateApproved,
		Reason:             "approved by copctl kube-cert command",
		Message:            "This CSR was approved by kubectl kube-cert tool",
		LastTransitionTime: metav1.Now(),
	})

	if _, err := clientset().CertificatesV1().
		CertificateSigningRequests().
		UpdateApproval(context.Background(), csr.ObjectMeta.Name, csr, metav1.UpdateOptions{}); err != nil {
		zap.S().Fatal(err)
	}
}

func fetchCertificateFromCsr(commonName string) []byte {
	csr, err := clientset().
		CertificatesV1().
		CertificateSigningRequests().
		Get(context.Background(), commonName, metav1.GetOptions{})
	if err != nil {
		zap.S().Fatal(err)
	}
	return csr.Status.Certificate
}
