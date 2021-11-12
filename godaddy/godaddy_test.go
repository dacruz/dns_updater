package godaddy

import (
	"net"
	"net/http"
	"testing"
	"github.com/dacruz/dns_updater/http2xx"
)

var handlers = map[string]func(http.ResponseWriter, *http.Request) {
	"/v1/domains/poiuytre.nl/records/A/@": func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "sso-key API_KEY" {
			rw.Write([]byte(`[{"data":"10.0.0.1","name":"@","ttl":600,"type":"A"}]`))
		} else {
			rw.Write([]byte("sso-key motherf***er, do you speak it?!"))
		}
	},
	"/v1/WRONG/RESPONSE/domains/poiuytre.nl/records/A/@": func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("WRONG_IP"))
	},
	"/v1/WRONG/IP/domains/poiuytre.nl/records/A/@": func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(`[{"data":"10","name":"@","ttl":600,"type":"A"}]`))
	},
}

func TestFetchCurrentRecordValue(t *testing.T) {
	server := http2xx.StartStubServer(handlers)
	defer http2xx.StopStubServer(server)

	currectIp, _ := FetchCurrentRecordValue("http://localhost:7000/v1", "poiuytre.nl", "@", "API_KEY")
	
	if ! net.ParseIP("10.0.0.1").Equal(currectIp) {
		t.Fatal("record does not have the expected value")
	}
	
}

func TestFailToParseFetchCurrentRecordValueResponse(t *testing.T) {
	server := http2xx.StartStubServer(handlers)
	defer http2xx.StopStubServer(server)

	_, err := FetchCurrentRecordValue("http://localhost:7000/v1/WRONG/RESPONSE", "poiuytre.nl", "@", "API_KEY")
	
	if err == nil {
		t.Fatal("FetchCurrentRecordValue should not have returned a valid json")
	}
	
}

func TestFailToParseFetchCurrentRecordValue(t *testing.T) {
	server := http2xx.StartStubServer(handlers)
	defer http2xx.StopStubServer(server)

	_, err := FetchCurrentRecordValue("http://localhost:7000/v1/WRONG/IP", "poiuytre.nl", "@", "API_KEY")
	
	if err == nil {
		t.Fatal("FetchCurrentRecordValue should not have returned a valid ip")
	}
	
}

func TestFailFetchCurrentRecordValueOnNon2xx(t *testing.T) {
	server := http2xx.StartStubServer(handlers)
	defer http2xx.StopStubServer(server)

	_, err := FetchCurrentRecordValue("http://localhost:7000/NOT_2XX", "poiuytre.nl", "@", "API_KEY")
	
	if err == nil {
		t.Fatal("FetchCurrentRecordValue should fail on non 2xx")
	}
	
}
