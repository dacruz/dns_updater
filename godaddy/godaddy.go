package godaddy

import (
	"fmt"
	"encoding/json"
	"errors"
	"net"
	"github.com/dacruz/dns_updater/http2xx"
)

type record struct {
	Data string
}

func FetchCurrentRecordValue(godaddyUrl string, domain string, host string, apiKey string) (net.IP, error) {

	url := fmt.Sprintf("%s/domains/%s/records/A/%s", godaddyUrl, domain, host)
	headers := map[string]string {
		"Authorization": fmt.Sprintf("sso-key %s", apiKey),
	}

	bodyBytes, err := http2xx.Get(url, headers)
	if err != nil {
		return nil, err
	}

	var records []record
	e := json.Unmarshal(bodyBytes, &records)
	if e != nil {
		return nil, errors.New("invalid json response")
	}

	ip := net.ParseIP(records[0].Data)
	if ip == nil {
		return nil, errors.New("invalid IP")
	}

	return 	ip, nil

}

func UpdateRecordValue(ip net.IP, godaddyUrl string, domain string, host string, apiKey string) (net.IP, error) {
	url := fmt.Sprintf("%s/domains/%s/records/A/%s", godaddyUrl, domain, host)

	headers := map[string]string {
		"Authorization": fmt.Sprintf("sso-key %s", apiKey),
		"Content-Type": "application/json",
	}

	aRecord := record{Data: ip.String()} 
	jsonBody, _ := json.Marshal(aRecord)

	bodyBytes, err := http2xx.Put(url, headers, jsonBody)
	if err != nil {
		return nil, err
	}

	var records []record
	e := json.Unmarshal(bodyBytes, &records)
	if e != nil {
		return nil, errors.New("invalid json response")
	}

	updatedIp := net.ParseIP(records[0].Data)
	if updatedIp == nil {
		return nil, errors.New("invalid IP")
	}

	return 	updatedIp, nil

}
