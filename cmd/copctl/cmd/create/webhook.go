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
	"strings"
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
	Use:     "webhook",
	Aliases: []string{"w"},
	Short:   "Generate certs for Admission Webhook",
	Run: func(cmd *cobra.Command, args []string) {
		NewWebhook(
			viper.GetString("common-name"),
			viper.GetString("certs-dir"),
			viper.GetBool("override"),
		).run()
	},
}

type Webhook struct {
	CommonName string
	CertsDir   string
	Override   bool
}

func NewWebhook(commonName, certsDir string, override bool) *Webhook {
	return &Webhook{
		CommonName: commonName,
		CertsDir:   certsDir,
		Override:   override,
	}
}

func (h *Webhook) run() {
	// re-create the cert folder
	if h.Override {
		h.clean()
	}
	// get key for CA
	cakey := privateKey()
	// create CA certificate
	ca, caPEM := h.createCA(cakey)
	// create certificate and key for server
	crt, key := h.serverCrtAndKey(ca, cakey)
	// create mutation webhook configuration
	ns, svc := h.commonNameToNsAndSvc()
	h.createMutatingWebhookCfg(
		admission.NewAICloudDomainHandler().
			HookCfg(ns, svc, caPEM.Bytes()),
	)
	// dump certificate to disc
	h.dumpToDisk(caPEM, crt, key)
}

func (h *Webhook) commonNameToNsAndSvc() (ns string, svc string) {
	endpoint := strings.Split(h.CommonName, ".")
	if len(endpoint) < 3 {
		zap.S().Fatalf("wrong common name, expected format: <svc-name>.<ns-name>.svc ")
	}
	return endpoint[1], endpoint[0]
}

func (h *Webhook) clean() {
	if err := os.RemoveAll(h.CertsDir); err != nil {
		zap.S().Error(err)
	}
}

func (h *Webhook) createCA(key *rsa.PrivateKey) (ca *x509.Certificate, caPEM *bytes.Buffer) {
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

func (h *Webhook) serverCrtAndKey(ca *x509.Certificate, cakey *rsa.PrivateKey) (serverCrtPEM *bytes.Buffer, serverKeyPEM *bytes.Buffer) {

	serverCrt := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			CommonName: h.CommonName,
		},
		DNSNames:    []string{h.CommonName},
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

func (h *Webhook) dumpToDisk(caCrt, serverCrt, serverKey *bytes.Buffer) {

	// create dir for cert and key
	if err := os.MkdirAll(h.CertsDir, 0755); err != nil {
		zap.S().Fatal(err)
	}
	// dump key to file
	if err := os.WriteFile(h.CertsDir+"/ca.crt", caCrt.Bytes(), 0644); err != nil {
		zap.S().Fatal(err)
	}
	if err := os.WriteFile(h.CertsDir+"/server.crt", serverCrt.Bytes(), 0644); err != nil {
		zap.S().Fatal(err)
	}
	if err := os.WriteFile(h.CertsDir+"/server.key", serverKey.Bytes(), 0644); err != nil {
		zap.S().Fatal(err)
	}

}

func (h *Webhook) createMutatingWebhookCfg(hookCfg *admissionv1.MutatingWebhookConfiguration) {

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
