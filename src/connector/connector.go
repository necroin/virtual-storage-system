package connector

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"net/http"
)

type Connector struct {
	rootCAs *x509.CertPool
}

func NewConnector(rootCAs *x509.CertPool) (*Connector, error) {
	return &Connector{
		rootCAs: rootCAs,
	}, nil
}

func (connector *Connector) SendPostRequest(url string, data any) (*http.Response, error) {
	return connector.SendRequestWithDataEncode(url, data, http.MethodPost)
}

func (connector *Connector) SendGetRequest(url string, data []byte, result any) error {
	response, err := connector.SendRequest(url, data, http.MethodGet)
	if err != nil {
		return err
	}

	return json.NewDecoder(response.Body).Decode(result)
}

func (connector *Connector) SendRequestWithDataEncode(url string, data any, method string) (*http.Response, error) {
	encodedMessage, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return connector.SendRequest(url, encodedMessage, method)
}

func (connector *Connector) SendRequest(url string, data []byte, method string) (*http.Response, error) {
	proto := "http://"
	var tlsConfig *tls.Config = nil
	if connector.rootCAs != nil {
		tlsConfig = &tls.Config{
			RootCAs:    connector.rootCAs,
			ServerName: "vss",
			// VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			// 	if len(verifiedChains) > 0 {
			// 		logger.Debug("[Connector] Verified certificate chain from peer:")
			// 		for _, certificate := range verifiedChains {
			// 			for i, cert := range certificate {
			// 				logger.Debug(fmt.Sprintf("[Connector] [Cert %d] %s", i, utils.CertificateInfo(cert)))
			// 			}
			// 		}
			// 	}
			// 	return nil
			// },
		}
		proto = "https://"
	}

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	request, err := http.NewRequest(
		method,
		proto+url,
		bytes.NewReader(data),
	)
	if err != nil {
		return nil, err
	}

	return client.Do(request)
}
