package main

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

)

func TestLoadConfReadsGoDaddyAPIUrl(t *testing.T) {
	setEnvVars()
	conf, _ := loadConf()
	
	expectedValue, _ := os.LookupEnv(DNS_UPDATER_GO_DADDY_API_URL)
	if conf.GoDaddyAPIUrl != expectedValue  {
		t.Fatal("conf.GoDaddyAPIUrl does not have the expected value")
	}

	if conf.GoDaddyAPIUrl == ""  {
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
	if conf.GoDaddyAPIKey != expectedValue  {
		t.Fatal("conf.GoDaddyAPIKey does not have the expected value")
	}

	if conf.GoDaddyAPIKey == ""  {
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
	if conf.Host != expectedValue  {
		t.Fatal("conf.Host does not have the expected value")
	}

	if conf.GoDaddyAPIKey == ""  {
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
	if conf.Domain != expectedValue  {
		t.Fatal("conf.Domain does not have the expected value")
	}

	if conf.Domain == ""  {
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

func TestLoadConfReadsIpfyUrl(t *testing.T) {
	setEnvVars()
	conf, _ := loadConf()
	
	expectedValue, _ := os.LookupEnv(DNS_UPDATER_IPFY_URL)
	if conf.IpfyAPIUrl != expectedValue  {
		t.Fatal("conf.IpfyAPIUrl does not have the expected value")
	}

	if conf.IpfyAPIUrl == ""  {
		t.Fatal("conf.IpfyAPIUrl was not loaded")
	}
}

func TestLoadConfFailsForAbsentIpfyUrl(t *testing.T) {
	unsetEnvVars(DNS_UPDATER_IPFY_URL)
	_, err := loadConf()
	
	if err == nil {
		t.Fatal("loadConf must fail if conf.IpfyAPIUrl is absent")
	}
}

func TestFetchCurrentIp(t *testing.T) {
	server := startServer()
	
	currectIp, _ := fetchCurrentIp("http://localhost:7000")
	
	if currectIp != "10.0.0.1" {
		t.Fatal("currectIp does not have the expected value")
	}
	
	if err := stopServer(server); err != nil {
		t.Fatal("Server shutdown failed")
	}
}

func startServer() *http.Server {
	router := http.NewServeMux() 
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("10.0.0.1"))
	})
	
	server := &http.Server{
		Addr:         ":7000",
		Handler:      router,
	}

	go server.ListenAndServe()
	time.Sleep(2*time.Second)

	return server
}

func stopServer(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3* time.Second)
  	defer cancel()

  	server.SetKeepAlivesEnabled(false)
  	if err := server.Shutdown(ctx); err != nil {
    	return err
  	}
	
	return server.Shutdown(ctx)
}


func setEnvVars() {
	os.Setenv(DNS_UPDATER_GO_DADDY_API_URL, "https://api.godaddy.com/v1/")
	os.Setenv(DNS_UPDATER_GO_DADDY_API_KEY, "SOME_API_KEY")
	os.Setenv(DNS_UPDATER_HOST, "poiuytre.nl")
	os.Setenv(DNS_UPDATER_DOMAIN, "www")
	os.Setenv(DNS_UPDATER_IPFY_URL, "https://api.ipify.org")
}

func unsetEnvVars(varName string) {
	os.Unsetenv(varName)
}

