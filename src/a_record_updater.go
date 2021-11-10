package main

import (
	"fmt"
	"log"
	"os"
	"errors"
	"net"
	"net/http"
	"io/ioutil"
)

const (
	DNS_UPDATER_GO_DADDY_API_URL = "DNS_UPDATER_GO_DADDY_API_URL"
	DNS_UPDATER_GO_DADDY_API_KEY = "DNS_UPDATER_GO_DADDY_API_KEY"
	DNS_UPDATER_HOST = "DNS_UPDATER_HOST"
	DNS_UPDATER_DOMAIN = "DNS_UPDATER_DOMAIN"
	DNS_UPDATER_IPFY_URL = "DNS_UPDATER_IPFY_URL"
)

type configuration struct {
    GoDaddyAPIUrl string 
	GoDaddyAPIKey string
	Host   string 
	Domain string
    IpfyAPIUrl string 
}

func main() {
	log.SetPrefix("dns updater: ")
	
	conf, err := loadConf()
	if err != nil {
        log.Fatal(err)
    }

	currentIp, err := fetchCurrentIp(conf.IpfyAPIUrl)
	if err != nil {
        log.Fatal(err)
    }

	fmt.Println(conf)
	fmt.Println(currentIp)

}

func fetchCurrentIp(ipfyUrl string) (net.IP, error){
	resp, err := http.Get(ipfyUrl)
	if err != nil {
        return nil, err
    }
	defer resp.Body.Close()

	// if he dies, he dies...
    bodyBytes, _ := ioutil.ReadAll(resp.Body)
    
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		errorMessage := fmt.Sprintf("request failed: %d, %q", resp.StatusCode, string(bodyBytes))
		return nil, errors.New(errorMessage)
	}

	ip := net.ParseIP(string(bodyBytes))
	if ip == nil {
		return nil, errors.New("invalid IP")
	}

	return ip, nil
}

func loadConf() (configuration, error) {

	config := configuration{}

	goDaddyAPIUrl, goDaddyAPIUrlExists := os.LookupEnv(DNS_UPDATER_GO_DADDY_API_URL)
	config.GoDaddyAPIUrl = goDaddyAPIUrl

	goDaddyAPIKey, goDaddyAPIKeyExists := os.LookupEnv(DNS_UPDATER_GO_DADDY_API_KEY)
	config.GoDaddyAPIKey = goDaddyAPIKey

	host, hostExists := os.LookupEnv(DNS_UPDATER_HOST)
	config.Host = host

	domain, domainExists := os.LookupEnv(DNS_UPDATER_DOMAIN)
	config.Domain = domain

	ipfyAPIUrl, ipfyAPIUrlExists := os.LookupEnv(DNS_UPDATER_IPFY_URL)
	config.IpfyAPIUrl = ipfyAPIUrl

	if !goDaddyAPIUrlExists || !goDaddyAPIKeyExists || !hostExists || !domainExists || !ipfyAPIUrlExists {
		return config, errors.New("environment variable missing. e.g: "+
		"export DNS_UPDATER_GO_DADDY_API_URL=FOO "+
		"export DNS_UPDATER_GO_DADDY_API_KEY=BAR "+
		"export DNS_UPDATER_HOST=BAZ "+
		"export DNS_UPDATER_DOMAIN=OPA "+
		"export DNS_UPDATER_IPFY_URL=OMA")
	}
	
	return config, nil
}

