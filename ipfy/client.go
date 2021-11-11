package ipfy

import (
	"errors"
	"net"
	"github.com/dacruz/dns_updater/http_2xx_only"
)

func FetchCurrentIp(ipfyUrl string) (net.IP, error){
    bodyBytes, err := http_2xx_only.Get(ipfyUrl, nil)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(string(bodyBytes))
	if ip == nil {
		return nil, errors.New("invalid IP")
	}

	return ip, nil
}