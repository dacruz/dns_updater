package ipfy

import (
	"github.com/dacruz/dns_updater/http2xx"
	"net"
)

func FetchCurrentIp(chnl chan net.IP, ipfyUrl string) {
	defer close(chnl)

	bodyBytes, err := http2xx.Get(ipfyUrl, nil)
	if err != nil {
		return
	}

	ip := net.ParseIP(string(bodyBytes))
	if ip == nil {
		return
	}

	chnl <- ip
}
