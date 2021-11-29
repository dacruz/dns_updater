package ipify

import (
	"github.com/dacruz/dns_updater/http2xx"
	"net"
)

func FetchCurrentIp(chnl chan net.IP, ipifyUrl string) {
	defer close(chnl)

	bodyBytes, err := http2xx.Get(ipifyUrl, nil)
	if err != nil {
		return
	}

	ip := net.ParseIP(string(bodyBytes))
	if ip == nil {
		return
	}

	chnl <- ip
}
