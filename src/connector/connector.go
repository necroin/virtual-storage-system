package connector

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendPostRequest(url string, data any) (*http.Response, error) {
	return SendRequestWithDataEncode(url, data, http.MethodPost)
}

func SendRequestWithDataEncode(url string, data any, method string) (*http.Response, error) {
	encodedMessage, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return SendRequest(url, encodedMessage, method)
}

func SendRequest(url string, data []byte, method string) (*http.Response, error) {
	client := http.Client{}

	request, err := http.NewRequest(
		http.MethodPost,
		"http://"+url,
		bytes.NewReader(data),
	)
	if err != nil {
		return nil, err
	}

	return client.Do(request)
}
