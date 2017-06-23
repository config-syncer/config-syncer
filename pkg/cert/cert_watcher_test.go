package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"strings"
	"testing"
	"time"

	_ "github.com/appscode/kubed/pkg/notifier/plivo"
	"github.com/google/certificate-transparency/go/x509"
	"github.com/google/certificate-transparency/go/x509/pkix"
	"github.com/stretchr/testify/assert"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func TestIsSoonExpired(t *testing.T) {
	cw := CertWatcher{
		MinRemainingDuration: time.Hour * 24 * 7,
	}
	_, days, err := cw.isSoonExpired(bytes.NewBuffer([]byte("an invalid ca certificate :)")), cw.MinRemainingDuration)
	assert.NotNil(t, err)

	soonExp, days, err := cw.isSoonExpired(bytes.NewBuffer(fakeCert(time.Hour*24*6+time.Second)), cw.MinRemainingDuration)
	assert.Nil(t, err)
	assert.True(t, soonExp)
	assert.Equal(t, 6, days)

	soonExp, days, err = cw.isSoonExpired(bytes.NewBuffer(fakeCert(time.Hour*24*8+time.Second)), cw.MinRemainingDuration)
	assert.Nil(t, err)
	assert.False(t, soonExp)
	assert.Equal(t, 8, days)
}

func TestConfiguration(t *testing.T) {
	s := &apiv1.Secret{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "mysecret",
			Namespace: "kube-system",
		},
		Type: "Opaque",
		Data: map[string][]byte{
			"username":         []byte("username"),
			"password":         []byte("password"),
			"notify_via":       []byte("plivo"),
			"plivo_auth_id":    []byte("auth_id"),
			"plivo_auth_token": []byte("auth_token"),
			"plivo_to":         []byte("admin,0111"),
			"plivo_from":       []byte("server"),
		},
	}
	exp := map[string]string{
		"username":         "username",
		"password":         "password",
		"notify_via":       "plivo",
		"plivo_auth_id":    "auth_id",
		"plivo_auth_token": "auth_token",
		"plivo_to":         "admin,0111",
		"plivo_from":       "server",
	}
	cw := DefaultCertWatcher(fake.NewSimpleClientset(s), s.ObjectMeta.Name, s.ObjectMeta.Namespace)
	m, err := cw.configuration()
	assert.Nil(t, err)
	assert.Equal(t, exp, m)
	cw = DefaultCertWatcher(fake.NewSimpleClientset(s), "dd", s.ObjectMeta.Namespace)
	_, err = cw.configuration()
	assert.NotNil(t, err)

}

func fakeCert(d time.Duration) []byte {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("failed to generate private key: %s", err)
	}
	notBefore := time.Now()
	notAfter := notBefore.Add(d)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("failed to generate serial number: %s", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	hosts := strings.Split("www.example.com", ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}
	certOut1 := bytes.NewBuffer(make([]byte, 0))
	pem.Encode(certOut1, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	return certOut1.Bytes()
}
