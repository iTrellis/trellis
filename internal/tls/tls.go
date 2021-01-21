/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package tls

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math"
	"math/big"
	"net"
	"time"
)

type Option func(*Options)

type Options struct {
	hosts     []string
	pkixName  *pkix.Name
	expiredIn time.Duration

	algorithm x509.PublicKeyAlgorithm

	rsaBits int
}

func Hosts(hosts ...string) Option {
	return func(o *Options) {
		o.hosts = hosts
	}
}

func Subject(subject *pkix.Name) Option {
	return func(o *Options) {
		o.pkixName = subject
	}
}

func ExpiredIn(expiredIn time.Duration) Option {
	return func(o *Options) {
		o.expiredIn = expiredIn
	}
}

func PublicKeyAlgorithm(algorithm x509.PublicKeyAlgorithm) Option {
	return func(o *Options) {
		o.algorithm = algorithm
	}
}

func RSABits(bits int) Option {
	return func(o *Options) {
		o.rsaBits = bits
	}
}

func (p *Options) check() {
	if p.pkixName == nil {
		p.pkixName = &pkix.Name{
			Organization: []string{"Go Trellis"},
		}
	}

	if p.expiredIn <= 0 {
		p.expiredIn = time.Hour * 24 * 365
	}

	if p.algorithm != x509.ECDSA && p.algorithm != x509.RSA {
		p.algorithm = x509.ECDSA
	}

	if p.rsaBits == 0 {
		p.rsaBits = 11
	}
	p.rsaBits = int(math.Pow(float64(2), float64(p.rsaBits)))

}

// Certificate gen tls Certificate
func Certificate(ofs ...Option) (tls.Certificate, error) {
	opts := &Options{}

	for _, o := range ofs {
		o(opts)
	}

	opts.check()

	notBefore := time.Now()
	notAfter := notBefore.Add(opts.expiredIn)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      *opts.pkixName,
		NotBefore:    notBefore,
		NotAfter:     notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, h := range opts.hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	switch opts.algorithm {
	case x509.ECDSA:
		return ecdsaCertificate(&template)
	case x509.RSA:
		return rsaCertificate(&template, opts.rsaBits)
	default:
		return tls.Certificate{}, errors.New("unsupported public key algorithm")
	}
}

func ecdsaCertificate(template *x509.Certificate) (tls.Certificate, error) {

	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &pk.PublicKey, pk)
	if err != nil {
		return tls.Certificate{}, err
	}

	// create public key
	certOut := bytes.NewBuffer(nil)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	// create private key
	keyOut := bytes.NewBuffer(nil)
	b, err := x509.MarshalECPrivateKey(pk)
	if err != nil {
		return tls.Certificate{}, err
	}

	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b})

	return tls.X509KeyPair(certOut.Bytes(), keyOut.Bytes())
}

func rsaCertificate(template *x509.Certificate, rsaBits int) (tls.Certificate, error) {

	pk, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		return tls.Certificate{}, err
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &pk.PublicKey, pk)
	if err != nil {
		return tls.Certificate{}, err
	}

	// create public key
	certOut := bytes.NewBuffer(nil)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	// create private key
	keyOut := bytes.NewBuffer(nil)

	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})

	return tls.X509KeyPair(certOut.Bytes(), keyOut.Bytes())
}
