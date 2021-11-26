package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/dacruz/dns_updater/http2xx"
)

func TestRunFailsOnConfigNotCorrect(t *testing.T) {
	os.Unsetenv(DNS_UPDATER_GO_DADDY_API_URL)
	os.Unsetenv(DNS_UPDATER_GO_DADDY_API_KEY)
	os.Unsetenv(DNS_UPDATER_HOST)
	os.Unsetenv(DNS_UPDATER_DOMAIN)
	os.Unsetenv(DNS_UPDATER_IPIFY_URL)

	err := run()
	if err == nil {
		t.Fatal("it shoudl have failed on missing config")
	}
}

func TestRunFailsOnIpifyFailure(t *testing.T) {
	setEnvVars()

	var handlers = map[string]func(http.ResponseWriter, *http.Request){
		"/ipify": func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusBadRequest)
		},
	}

	server := http2xx.StartStubServer(handlers)
	defer http2xx.StopStubServer(server)

	err := run()
	if err == nil {
		t.Fatal("it shoudl have failed on ipify failure")
	}
}

func TestRunFailsOnGetGodaddyFailure(t *testing.T) {
	setEnvVars()

	var handlers = map[string]func(http.ResponseWriter, *http.Request){
		"/ipify": func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte("11.0.0.1"))
		},
	}

	server := http2xx.StartStubServer(handlers)
	defer http2xx.StopStubServer(server)

	err := run()
	if err == nil {
		t.Fatal("it shoudl have failed on get godaddy failure")
	}
}

func TestRunFailsOnPutGodaddyFailure(t *testing.T) {
	setEnvVars()

	var handlers = map[string]func(http.ResponseWriter, *http.Request){
		"/godaddy/domains/poiuytre.nl/records/A/@": func(rw http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				rw.Write([]byte(`[{"data":"10.0.0.1","name":"@","ttl":600,"type":"A"}]`))
			}

			if r.Method == "PUT" {
				rw.WriteHeader(http.StatusBadRequest)
			}
		},
		"/ipify": func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte("11.0.0.1"))
		},
	}

	server := http2xx.StartStubServer(handlers)
	defer http2xx.StopStubServer(server)

	err := run()
	if err == nil {
		t.Fatal("it shoudl have failed on put godaddy failure")
	}
}

func TestUpdateRecord(t *testing.T) {
	setEnvVars()

	var updated bool
	var handlers = map[string]func(http.ResponseWriter, *http.Request){
		"/godaddy/domains/poiuytre.nl/records/A/@": func(rw http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				rw.Write([]byte(`[{"data":"10.0.0.1","name":"@","ttl":600,"type":"A"}]`))
			}

			if r.Method == "PUT" {
				bodyBytes, _ := ioutil.ReadAll(r.Body)
				body := string(bodyBytes)
				updated = strings.Contains(body, "11.0.0.1")
			}
		},
		"/ipify": func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte("11.0.0.1"))
		},
	}

	server := http2xx.StartStubServer(handlers)
	defer http2xx.StopStubServer(server)

	main()

	if !updated {
		t.Fatal("record was not updated")
	}

}

func TestDoNotUpdateRecordIfTheSame(t *testing.T) {
	setEnvVars()

	var updated bool
	var handlers = map[string]func(http.ResponseWriter, *http.Request){
		"/godaddy/domains/poiuytre.nl/records/A/@": func(rw http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				rw.Write([]byte(`[{"data":"10.0.0.1","name":"@","ttl":600,"type":"A"}]`))
			}

			if r.Method == "PUT" {
				updated = true
			}
		},
		"/ipify": func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte("10.0.0.1"))
		},
	}

	server := http2xx.StartStubServer(handlers)
	defer http2xx.StopStubServer(server)

	main()

	if updated {
		t.Fatal("record was updated when it should not")
	}

}

func TestLoadConfReadsGoDaddyAPIUrl(t *testing.T) {
	setEnvVars()
	conf, _ := loadConf()

	expectedValue, _ := os.LookupEnv(DNS_UPDATER_GO_DADDY_API_URL)
	if conf.GoDaddyAPIUrl != expectedValue {
		t.Fatal("conf.GoDaddyAPIUrl does not have the expected value")
	}

	if conf.GoDaddyAPIUrl == "" {
		t.Fatal("conf.GoDaddyAPIUrl was not loaded")
	}
}

func TestLoadConfFailsForAbsentGoDaddyAPIUrl(t *testing.T) {
	unsetEnvVars(DNS_UPDATER_GO_DADDY_API_URL)
	_, err := loadConf()

	if err == nil {
		t.Fatal("loadConf must fail if conf.GoDaddyAPIUrl is absent")
	}
}

func TestLoadConfReadsGoDaddyAPIKey(t *testing.T) {
	setEnvVars()
	conf, _ := loadConf()

	expectedValue, _ := os.LookupEnv(DNS_UPDATER_GO_DADDY_API_KEY)
	if conf.GoDaddyAPIKey != expectedValue {
		t.Fatal("conf.GoDaddyAPIKey does not have the expected value")
	}

	if conf.GoDaddyAPIKey == "" {
		t.Fatal("conf.GoDaddyAPIKey was not loaded")
	}
}

func TestLoadConfFailsForAbsentGoDaddyAPIKey(t *testing.T) {
	unsetEnvVars(DNS_UPDATER_GO_DADDY_API_KEY)
	_, err := loadConf()

	if err == nil {
		t.Fatal("loadConf must fail if conf.GoDaddyAPIKey is absent")
	}
}

func TestLoadConfReadsHost(t *testing.T) {
	setEnvVars()
	conf, _ := loadConf()

	expectedValue, _ := os.LookupEnv(DNS_UPDATER_HOST)
	if conf.Host != expectedValue {
		t.Fatal("conf.Host does not have the expected value")
	}

	if conf.GoDaddyAPIKey == "" {
		t.Fatal("conf.Host was not loaded")
	}
}

func TestLoadConfFailsForAbsentHost(t *testing.T) {
	unsetEnvVars(DNS_UPDATER_HOST)
	_, err := loadConf()

	if err == nil {
		t.Fatal("loadConf must fail if conf.Host is absent")
	}
}

func TestLoadConfReadsDomain(t *testing.T) {
	setEnvVars()
	conf, _ := loadConf()

	expectedValue, _ := os.LookupEnv(DNS_UPDATER_DOMAIN)
	if conf.Domain != expectedValue {
		t.Fatal("conf.Domain does not have the expected value")
	}

	if conf.Domain == "" {
		t.Fatal("conf.Domain was not loaded")
	}
}

func TestLoadConfFailsForAbsentDomain(t *testing.T) {
	unsetEnvVars(DNS_UPDATER_DOMAIN)
	_, err := loadConf()

	if err == nil {
		t.Fatal("loadConf must fail if conf.Domain is absent")
	}
}

func TestLoadConfReadsIpifyUrl(t *testing.T) {
	setEnvVars()
	conf, _ := loadConf()

	expectedValue, _ := os.LookupEnv(DNS_UPDATER_IPIFY_URL)
	if conf.IpifyAPIUrl != expectedValue {
		t.Fatal("conf.IpifyAPIUrl does not have the expected value")
	}

	if conf.IpifyAPIUrl == "" {
		t.Fatal("conf.IpifyAPIUrl was not loaded")
	}
}

func TestLoadConfFailsForAbsentIpifyUrl(t *testing.T) {
	unsetEnvVars(DNS_UPDATER_IPIFY_URL)
	_, err := loadConf()

	if err == nil {
		t.Fatal("loadConf must fail if conf.IpifyAPIUrl is absent")
	}
}

func setEnvVars() {
	os.Setenv(DNS_UPDATER_GO_DADDY_API_URL, "http://localhost:7000/godaddy")
	os.Setenv(DNS_UPDATER_GO_DADDY_API_KEY, "SOME_API_KEY")
	os.Setenv(DNS_UPDATER_HOST, "@")
	os.Setenv(DNS_UPDATER_DOMAIN, "poiuytre.nl")
	os.Setenv(DNS_UPDATER_IPIFY_URL, "http://localhost:7000/ipify")
}

func unsetEnvVars(varName string) {
	os.Unsetenv(varName)
}
