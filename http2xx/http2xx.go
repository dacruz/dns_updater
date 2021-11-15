package http2xx

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(url string, headers map[string]string) ([]byte, error) {
	request, _ := http.NewRequest("GET", url, nil)
	return exec(request, headers)
}

func Put(url string, headers map[string]string, body []byte) ([]byte, error) {
	request, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	return exec(request, headers)
}

func exec(request *http.Request, headers map[string]string) ([]byte, error) {
	client := &http.Client{}

	for name, value := range headers {
		request.Header.Set(name, value)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// if he dies, he dies...
	bodyBytes, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode < 200 || response.StatusCode > 299 {
		errorMessage := fmt.Sprintf("request failed: %d, %q", response.StatusCode, string(bodyBytes))
		return nil, errors.New(errorMessage)
	}

	return bodyBytes, nil
}
