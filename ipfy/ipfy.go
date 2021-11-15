package ipfy

import (
	"errors"
	"github.com/dacruz/dns_updater/http2xx"
	"net"
)

func FetchCurrentIp(ipfyUrl string) (net.IP, error) {
	bodyBytes, err := http2xx.Get(ipfyUrl, nil)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(string(bodyBytes))
	if ip == nil {
		return nil, errors.New("invalid IP")
	}

	return ip, nil
}
