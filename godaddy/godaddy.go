package godaddy

import (
	"encoding/json"
	"fmt"
	"github.com/dacruz/dns_updater/http2xx"
	"net"
)

type record struct {
	Data string `json:"data"`
}

func FetchCurrentRecordValue(chnl chan net.IP, godaddyUrl string, domain string, host string, apiKey string) {
	defer close(chnl)

	url := fmt.Sprintf("%s/domains/%s/records/A/%s", godaddyUrl, domain, host)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("sso-key %s", apiKey),
	}

	bodyBytes, err := http2xx.Get(url, headers)
	if err != nil {
		return
	}

	var records []record
	e := json.Unmarshal(bodyBytes, &records)
	if e != nil {
		return
	}

	ip := net.ParseIP(records[0].Data)
	if ip == nil {
		return
	}

	chnl <- ip
}

func UpdateRecordValue(ip net.IP, godaddyUrl string, domain string, host string, apiKey string) (net.IP, error) {
	url := fmt.Sprintf("%s/domains/%s/records/A/%s", godaddyUrl, domain, host)

	headers := map[string]string{
		"Authorization": fmt.Sprintf("sso-key %s", apiKey),
		"Content-Type":  "application/json",
	}

	aRecord := record{Data: ip.String()}
	recordsToUpdate := []record{aRecord}
	jsonBody, _ := json.Marshal(recordsToUpdate)

	_, err := http2xx.Put(url, headers, jsonBody)
	if err != nil {
		return nil, err
	}

	return ip, nil

}
