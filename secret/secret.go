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
	"path/filepath"
	"strings"
	"time"

	"github.com/msw-x/moon/uerr"
	"github.com/msw-x/moon/ufs"
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

func (o *Gen) init() {
	if o.Organization == "" {
		o.Organization = "Goodlife"
	}
	if o.ValidFrom.IsZero() {
		o.ValidFrom = time.Now()
	}
	if o.ValidFor == 0 {
		o.ValidFor = 365 * 24 * time.Hour
	}
	if o.IsCA == nil {
		o.isCA = true
	} else {
		o.isCA = o.IsCA.(bool)
	}
	if o.RsaBits == 0 {
		o.RsaBits = 2048
	}
}

func (o Gen) Generate(certPath string, keyPath string) error {
	o.init()
	var priv any
	var err error
	switch o.EcdsaCurve {
	case "":
		priv, err = rsa.GenerateKey(rand.Reader, o.RsaBits)
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		return fmt.Errorf("unrecognized elliptic curve: %q", o.EcdsaCurve)
	}
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	notBefore := o.ValidFrom
	notAfter := notBefore.Add(o.ValidFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{o.Organization},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(o.Host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if o.isCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %v", err)
	}

	ufs.MakeDir(filepath.Dir(certPath))
	certOut, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %v", certPath, err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	ufs.MakeDir(filepath.Dir(keyPath))
	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %v", keyPath, err)
	}
	pem.Encode(keyOut, pemBlockForKey(priv))
	keyOut.Close()
	return nil
}

func (o Gen) Ensure(certPath string, keyPath string) {
	if !ufs.Exist(certPath) || !ufs.Exist(keyPath) {
		err := o.Generate(certPath, keyPath)
		uerr.Strict(err, "certificate generation")
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

func pemBlockForKey(priv any) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}
