package godaddy

import (
	"fmt"
	"net"
	"net/http"
	"io/ioutil"
	"errors"
)

func FetchCurrentRecordValue(godaddyUrl string, domain string, host string, apiKey string) (net.IP, error) {

	url := fmt.Sprintf("%s/domains/%s/records/A/%s", godaddyUrl, domain, host)

	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)

	auth := fmt.Sprintf("sso-key %s", apiKey)
	request.Header.Set("Authorization", auth)
	
	response, err := client.Do(request)
	if err != nil {
        return nil, err
    }
	
	// if he dies, he dies... again!
    bodyBytes, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode < 200 || response.StatusCode > 299 {
		errorMessage := fmt.Sprintf("request failed: %d, %q", response.StatusCode, string(bodyBytes))
		return nil, errors.New(errorMessage)
	}

	ip := net.ParseIP(string(bodyBytes))
	if ip == nil {
		return nil, errors.New("invalid IP")
	}

	return 	ip, nil

}