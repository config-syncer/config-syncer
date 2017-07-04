package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/appscode/kubed/pkg/util"
	"github.com/google/certificate-transparency/go/x509"
	"github.com/google/certificate-transparency/go/x509/pkix"
	"github.com/stretchr/testify/assert"
)

var cw = CertWatcher{
	MinRemainingDuration: time.Hour * 24 * 7,
}

func TestInvalidCertificate(t *testing.T) {
	_, _, err := cw.isSoonExpired(bytes.NewBuffer([]byte("an invalid ca certificate :)")), cw.MinRemainingDuration)
	assert.NotNil(t, err)
}

func TestAlreadyExpiredCertificate(t *testing.T) {
	now := time.Now()
	_, _, err := cw.isSoonExpired(bytes.NewBuffer(fakeCert(now.Add(-time.Hour*24*2), now.Add(-time.Hour*24*1))), cw.MinRemainingDuration)
	assert.NotNil(t, err)
	assert.Equal(t, "certificate already expired", err.Error())
}

func TestSoonExpeiredCertificate(t *testing.T) {
	now := time.Now()
	soonExp, days, err := cw.isSoonExpired(bytes.NewBuffer(fakeCert(now, now.Add(time.Hour*24*7))), cw.MinRemainingDuration)
	assert.Nil(t, err)
	assert.True(t, soonExp)
	assert.Equal(t, 6, days)
}

func TestCertificateWithAvailabeExperition(t *testing.T) {
	now := time.Now()
	soonExp, _, err := cw.isSoonExpired(bytes.NewBuffer(fakeCert(now, now.Add(time.Hour*24*30))), cw.MinRemainingDuration)
	assert.Nil(t, err)
	assert.False(t, soonExp)
}

func TestConfiguration(t *testing.T) {
	exp := map[string]string{
		"username":         "username",
		"password":         "password",
		"notify_via":       "plivo",
		"plivo_auth_id":    "auth_id",
		"plivo_auth_token": "auth_token",
		"plivo_to":         "admin,0111",
		"plivo_from":       "server",
	}

	path := os.Getenv("HOME") + "/temp"
	os.MkdirAll(path, 0777)
	util.MapsToFiles(exp, path)
	defer os.RemoveAll(path)

	cw := DefaultCertWatcher(path)
	m, err := cw.configuration()
	assert.Nil(t, err)
	assert.Equal(t, exp, m)
}

func fakeCert(notBefore, notAfter time.Time) []byte {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("failed to generate private key: %s", err)
	}
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
