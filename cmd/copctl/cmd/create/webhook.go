package create

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/AccessibleAI/cnvrg-operator/pkg/admission"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/big"
	"os"
	"time"
)

func init() {
	webhookCmd.PersistentFlags().StringP("common-name", "c", "",
		"certificate common name (will be included to alternative DNS names as well)")
	webhookCmd.PersistentFlags().BoolP("override", "w", true,
		"override csr, files and secrets if they already exists")
	webhookCmd.PersistentFlags().StringP("certs-dir", "f", "certs",
		"path to dump key & crt")

	viper.BindPFlag("common-name", webhookCmd.PersistentFlags().Lookup("common-name"))
	viper.BindPFlag("override", webhookCmd.PersistentFlags().Lookup("override"))
	viper.BindPFlag("certs-dir", webhookCmd.PersistentFlags().Lookup("certs-dir"))

	Cmd.AddCommand(webhookCmd)
}

var webhookCmd = &cobra.Command{
	Use:     "webhook-certs",
	Aliases: []string{"w"},
	Short:   "Generate certs for Admission Webhook",
	Run: func(cmd *cobra.Command, args []string) {
		commonName := viper.GetString("common-name")
		certsDir := viper.GetString("certs-dir")
		// re-create the cert folder
		if viper.GetBool("override") {
			clean(certsDir, commonName)
		}
		// get key for CA
		cakey := privateKey()
		// create CA certificate
		ca, caPEM := createCA(cakey)
		// create certificate and key for server
		crt, key := serverCrtAndKey(commonName, ca, cakey)
		// create mutation webhook configuration
		ns, svc := commonNameToNsAndSvc(commonName)
		createMutatingWebhookCfg(
			admission.NewAICloudDomainHandler().
				HookCfg(ns, svc, caPEM.Bytes()),
		)
		// dump certificate to disc
		dumpToDisk(caPEM, crt, key, certsDir)
	},
}

func createCA(key *rsa.PrivateKey) (ca *x509.Certificate, caPEM *bytes.Buffer) {
	ca = &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"cnvrg-ai-cloud"},
			Country:      []string{"IL"},
			CommonName:   "root.localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &key.PublicKey, key)
	if err != nil {
		zap.S().Fatalf("error creating certificate, err: %s", err.Error())
	}
	caPEM = new(bytes.Buffer)
	if err := pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}); err != nil {
		zap.S().Fatalf("error encoding certificate, err: %s", err.Error())
	}
	return ca, caPEM
}

func serverCrtAndKey(commonName string, ca *x509.Certificate, cakey *rsa.PrivateKey) (serverCrtPEM *bytes.Buffer, serverKeyPEM *bytes.Buffer) {

	serverCrt := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			CommonName: commonName,
		},
		DNSNames:    []string{commonName},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(10, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	serverKey := privateKey()
	certBytes, err := x509.CreateCertificate(rand.Reader, serverCrt, ca, &serverKey.PublicKey, cakey)
	if err != nil {
		zap.S().Fatalf("error creating server certificate, err: %s ", err.Error())
	}

	serverCrtPEM = new(bytes.Buffer)
	if err := pem.Encode(serverCrtPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		zap.S().Fatalf("error encoding server certificate, err: %s ", err.Error())
	}

	serverKeyPEM = new(bytes.Buffer)
	if err := pem.Encode(serverKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(serverKey),
	}); err != nil {
		zap.S().Fatalf("error encoding server key, err: %s ", err.Error())
	}

	return
}

func dumpToDisk(caCrt, serverCrt, serverKey *bytes.Buffer, certsDir string) {

	// create dir for cert and key
	if err := os.MkdirAll(certsDir, 0755); err != nil {
		zap.S().Fatal(err)
	}
	// dump key to file
	if err := os.WriteFile(certsDir+"/ca.crt", caCrt.Bytes(), 0644); err != nil {
		zap.S().Fatal(err)
	}
	if err := os.WriteFile(certsDir+"/server.crt", serverCrt.Bytes(), 0644); err != nil {
		zap.S().Fatal(err)
	}
	if err := os.WriteFile(certsDir+"/server.key", serverKey.Bytes(), 0644); err != nil {
		zap.S().Fatal(err)
	}

}

func createAICloudDomainWebhookCfg(ns, svc string, caBundle []byte) {
	zap.S().Infof("creating ai cloud domain mutation webhook")
	createMutatingWebhookCfg(
		admission.NewAICloudDomainHandler().
			HookCfg(ns, svc, caBundle),
	)
}

func createMutatingWebhookCfg(hookCfg *admissionv1.MutatingWebhookConfiguration) {
	zap.S().Infof("creating webhook: %s", hookCfg.Name)

	err := clientset().
		AdmissionregistrationV1().
		MutatingWebhookConfigurations().
		Delete(context.Background(), hookCfg.Name, metav1.DeleteOptions{})

	if err != nil && !errors.IsNotFound(err) {
		zap.S().Fatalf("error deleting webhook: %s, err: %s ", hookCfg.Name, err.Error())
	}

	if _, err := clientset().
		AdmissionregistrationV1().
		MutatingWebhookConfigurations().
		Create(context.Background(), hookCfg, metav1.CreateOptions{}); err != nil {
		zap.S().Fatalf("error creating mutating webhook: %s cofigurations, err: %s", hookCfg.Name, err.Error())
	}
}
