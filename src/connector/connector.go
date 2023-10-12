package connector

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
)

func SendPostRequest(url string, data any) (*http.Response, error) {
	return SendRequestWithDataEncode(url, data, http.MethodPost)
}

func SendGetRequest[T any](url string, data []byte) (*T, error) {
	result := new(T)

	response, err := SendRequest(url, data, http.MethodGet)
	if err != nil {
		return nil, err
	}

	json.NewDecoder(response.Body).Decode(result)

	return result, nil
}

func SendRequestWithDataEncode(url string, data any, method string) (*http.Response, error) {
	encodedMessage, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return SendRequest(url, encodedMessage, method)
}

func SendRequest(url string, data []byte, method string) (*http.Response, error) {
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
