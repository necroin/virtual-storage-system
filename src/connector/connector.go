package connector

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
)

type Connector struct {
	certificate *tls.Certificate
}

func NewConnector(certificate *tls.Certificate) (*Connector, error) {
	return &Connector{
		certificate: certificate,
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
				ServerName:         "localhost",
				InsecureSkipVerify: true,
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
