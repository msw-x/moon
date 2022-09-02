package secret

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/msw-x/moon"
	"github.com/msw-x/moon/fs"
)

func Ensure(certPath string, keyPath string) {
	Gen{}.Ensure(certPath, keyPath)
}

func EnsureDir(dir string) {
	Gen{}.Ensure(FileNames(dir))
}

func FileNames(dir string) (string, string) {
	return path.Join(dir, "cert.pem"), path.Join(dir, "key.pem")
}

type Gen struct {
	Host         string
	Organization string
	ValidFrom    time.Time
	ValidFor     time.Duration
	IsCA         any
	RsaBits      int
	EcdsaCurve   string

	isCA bool
}

func (this *Gen) init() {
	if this.Organization == "" {
		this.Organization = "Goodlife"
	}
	if this.ValidFrom == 0 {
		this.ValidFrom = time.Now()
	}
	if this.ValidFor == 0 {
		this.ValidFor = 365 * 24 * time.Hour
	}
	if IsCA == nil {
		this.isCA = true
	} else {
		this.isCA = IsCA.(bool)
	}
	if RsaBits == 0 {
		this.RsaBits = 2048
	}
}

func (this Gen) Generate(certPath string, keyPath string) error {
	var priv any
	var err error
	switch this.EcdsaCurve {
	case "":
		priv, err = rsa.GenerateKey(rand.Reader, rsaBits)
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		return fmt.Errorf("unrecognized elliptic curve: %q", this.EcdsaCurve)
	}
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	notBefore := this.ValidFrom
	notAfter := notBefore.Add(this.ValidFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{this.Organization},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if this.isCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %v", err)
	}

	certOut, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %v", certPath, err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %v", keyPath, err)
	}
	pem.Encode(keyOut, pemBlockForKey(priv))
	keyOut.Close()
	return nil
}

func (this Gen) Ensure(certPath string, keyPath string) {
	if !fs.Exist(certPath) || !fs.Exist(keyPath) {
		err := this.Generate(certPath, keyPath)
		moon.Check(err, "certificate generation")
	}
}

func publicKey(priv any) any {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}
