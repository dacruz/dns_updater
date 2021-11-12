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
			switch r.Method {
			case "GET":
				rw.Write([]byte(`[{"data":"10.0.0.1","name":"@","ttl":600,"type":"A"}]`))
			case "PUT":
				rw.WriteHeader(http.StatusOK)
			default:
				rw.WriteHeader(http.StatusBadRequest)
			}
			
		} else {
			rw.Write([]byte("sso-key motherf***er, do you speak it?!"))
		}
	},
	"/v1/WRONG/RESPONSE/domains/poiuytre.nl/records/A/@": func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("WRONG_IP"))
	},
	"/v1/INVALID/IP/domains/poiuytre.nl/records/A/@": func(rw http.ResponseWriter, r *http.Request) {
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

	_, err := FetchCurrentRecordValue("http://localhost:7000/v1/INVALID/IP", "poiuytre.nl", "@", "API_KEY")
	
	if err == nil {
		t.Fatal("FetchCurrentRecordValue should not have parsed an invalid ip")
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

func TestUpdateRecordValue(t *testing.T) {
	server := http2xx.StartStubServer(handlers)
	defer http2xx.StopStubServer(server)

	ip, _ := UpdateRecordValue(net.ParseIP("11.0.0.1"),"http://localhost:7000/v1", "poiuytre.nl", "@", "API_KEY")

	if ! net.ParseIP("11.0.0.1").Equal(ip) {
		t.Fatal("record does not have the expected value")
	}
}

func TestFailUpdateRecordValueOnNon2xx(t *testing.T) {
	server := http2xx.StartStubServer(handlers)
	defer http2xx.StopStubServer(server)

	_, err := UpdateRecordValue(net.ParseIP("11.0.0.1"), "http://localhost:7000/NOT_2XX", "poiuytre.nl", "@", "API_KEY")
	
	if err == nil {
		t.Fatal("UpdateRecordValue should fail on non 2xx")
	}

}
