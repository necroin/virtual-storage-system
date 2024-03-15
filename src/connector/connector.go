package connector

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"vss/src/config"
	"vss/src/logger"
	"vss/src/utils"
)

type Connector struct {
	config *config.Config
}

func NewConnector(config *config.Config) (*Connector, error) {
	return &Connector{
		config: config,
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
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:    connector.config.RootCAs,
				ServerName: "vss",
				VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
					if len(verifiedChains) > 0 {
						logger.Debug("[Connector] Verified certificate chain from peer:")
						for _, certificate := range verifiedChains {
							for i, cert := range certificate {
								logger.Debug("[Connector] Cert %d:\n", i)
								logger.Debug(fmt.Sprintf("[Connector] %s", utils.CertificateInfo(cert)))
							}
						}
					}
					return nil
				},
			},
		},
	}

	request, err := http.NewRequest(
		method,
		"https://"+url,
		bytes.NewReader(data),
	)
	if err != nil {
		return nil, err
	}

	return client.Do(request)
}
