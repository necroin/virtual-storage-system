package utils

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"path"
	"strings"
	"time"
)

func GetMapKeys[K comparable, V any](value map[K]V) []K {
	result := []K{}
	for key := range value {
		result = append(result, key)
	}
	return result
}

func GenerateSecureToken(length int) string {
	token := make([]byte, length)
	if _, err := rand.Read(token); err != nil {
		return ""
	}
	return hex.EncodeToString(token)
}

func FormatTokemizedEndpoint(endpoint string, token string) string {
	return strings.NewReplacer("{token}", token).Replace(endpoint)
}

func HandleFilesystemPath(value string) string {
	value = strings.Trim(value, "\"")
	value = path.Clean(value)
	return value
}

func IfNotNil[T any](object *T, handler func(object *T) error) error {
	if object != nil {
		return handler(object)
	}
	return nil
}

func CertificateInfo(cert *x509.Certificate) []string {
	result := []string{}
	if cert.Subject.CommonName == cert.Issuer.CommonName {
		result = append(result, fmt.Sprintf("(Self-signed certificate) %v", cert.Issuer.CommonName))
		return result
	}

	result = append(result, fmt.Sprintf("(Subject) %v", cert.DNSNames))
	result = append(result, fmt.Sprintf("(Usage) %v", cert.ExtKeyUsage))
	result = append(result, fmt.Sprintf("(Issued by) %s", cert.Issuer.CommonName))
	result = append(result, fmt.Sprintf("(Issued by) %s", cert.Issuer.SerialNumber))

	return result
}

func Try(handler func() error, count int, sleep time.Duration) error {
	var err error = nil
	for i := 0; i < count; i++ {
		err = handler()
		if err == nil {
			return nil
		}
		time.Sleep(sleep)
	}
	return err
}
