package shared

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

func GeneratePrivateKey() (*rsa.PrivateKey, error) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}

	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func EncodeKeysToPEM(privateKey *rsa.PrivateKey) ([]byte, []byte, error) {

	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	privBlock := pem.Block{Type: "RSA PRIVATE KEY", Bytes: privDER}

	privatePEM := pem.EncodeToMemory(&privBlock)

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		logrus.Error(err, "error when dumping public key")
		return nil, nil, err
	}

	publicBlock := pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyBytes}

	publicPEM := pem.EncodeToMemory(&publicBlock)

	return privatePEM, publicPEM, nil
}

func CreateSSOToken(jwtAudience string, privateKeyBytes []byte) (string, error) {
	claims := &jwt.StandardClaims{Audience: jwtAudience, Issuer: "http://auth.cnvrg"}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
	key, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	return token.SignedString(key)

}
