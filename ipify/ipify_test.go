package ipify

import (
	"github.com/dacruz/dns_updater/http2xx"
	"net"
	"net/http"
	"testing"
)

var handlers = map[string]func(http.ResponseWriter, *http.Request){
	"/current/ip": func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("10.0.0.1"))
	},
	"/WRONG/IP": func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("WRONG_IP"))
	},
}

func TestFetchCurrentIp(t *testing.T) {
	server := http2xx.StartStubServer(handlers)
	defer http2xx.StopStubServer(server)

	currentIpChannel := make(chan net.IP)
	go FetchCurrentIp(currentIpChannel, "http://localhost:7000/current/ip")

	currectIp := <-currentIpChannel
	if !net.ParseIP("10.0.0.1").Equal(currectIp) {
		t.Fatal("currectIp does not have the expected value")
	}

}

func TestFailToParseFetchCurrentIpResponse(t *testing.T) {
	server := http2xx.StartStubServer(handlers)
	defer http2xx.StopStubServer(server)

	currentIpChannel := make(chan net.IP)
	go FetchCurrentIp(currentIpChannel, "http://localhost:7000/WRONG/IP")

	_, ok := <-currentIpChannel
	if ok {
		t.Fatal("FetchCurrentIp should not have returned a valid ip")
	}

}

func TestFailFetchCurrentIpOnNon2xx(t *testing.T) {
	server := http2xx.StartStubServer(handlers)
	defer http2xx.StopStubServer(server)
	currentIpChannel := make(chan net.IP)
	go FetchCurrentIp(currentIpChannel, "http://localhost:7000/NOT_2XX")

	_, ok := <-currentIpChannel
	if ok {
		t.Fatal("FetchCurrentIp should fail on non 2xx")
	}

}
