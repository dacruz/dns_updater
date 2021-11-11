package ipfy

import (
	"fmt"
	"net"
	"net/http"
	"io/ioutil"
	"errors"
)

func FetchCurrentIp(ipfyUrl string) (net.IP, error){
	resp, err := http.Get(ipfyUrl)
	if err != nil {
        return nil, err
    }
	defer resp.Body.Close()

	// if he dies, he dies...
    bodyBytes, _ := ioutil.ReadAll(resp.Body)
    
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		errorMessage := fmt.Sprintf("request failed: %d, %q", resp.StatusCode, string(bodyBytes))
		return nil, errors.New(errorMessage)
	}

	ip := net.ParseIP(string(bodyBytes))
	if ip == nil {
		return nil, errors.New("invalid IP")
	}

	return ip, nil
}