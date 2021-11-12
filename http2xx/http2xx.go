package http2xx

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(url string, headers map[string]string) ([]byte, error) {

	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)

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

	return 	bodyBytes, nil

}