package main

import (
	"errors"
	"log"
	"net"
	"os"

	"github.com/dacruz/dns_updater/godaddy"
	"github.com/dacruz/dns_updater/ipify"
)

const (
	DNS_UPDATER_GO_DADDY_API_URL = "DNS_UPDATER_GO_DADDY_API_URL"
	DNS_UPDATER_GO_DADDY_API_KEY = "DNS_UPDATER_GO_DADDY_API_KEY"
	DNS_UPDATER_HOST             = "DNS_UPDATER_HOST"
	DNS_UPDATER_DOMAIN           = "DNS_UPDATER_DOMAIN"
	DNS_UPDATER_IPIFY_URL         = "DNS_UPDATER_IPIFY_URL"
)

type configuration struct {
	GoDaddyAPIUrl string
	GoDaddyAPIKey string
	Host          string
	Domain        string
	IpifyAPIUrl    string
}

func main() {
	log.SetPrefix("dns updater: ")
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conf, err := loadConf()
	if err != nil {
		return err
	}

	currentIpChannel := make(chan net.IP)
	go ipify.FetchCurrentIp(currentIpChannel, conf.IpifyAPIUrl)

	currentDnsValueChannel := make(chan net.IP)
	go godaddy.FetchCurrentRecordValue(currentDnsValueChannel, conf.GoDaddyAPIUrl, conf.Domain, conf.Host, conf.GoDaddyAPIKey)

	currentIp, ok := <-currentIpChannel
	if !ok {
		return errors.New("could not fetch currect ip")
	}

	currentDnsValue, ok := <-currentDnsValueChannel
	if !ok {
		return errors.New("could not fetch currect dns value")
	}

	log.Printf("dns value: %s", currentDnsValue.String())
	log.Printf("current ip: %s", currentIp.String())

	if !currentDnsValue.Equal(currentIp) {
		updatedDnsValue, err := godaddy.UpdateRecordValue(currentIp, conf.GoDaddyAPIUrl, conf.Domain, conf.Host, conf.GoDaddyAPIKey)
		if err != nil {
			return err
		}
		log.Printf("updatedDnsValue: %s", updatedDnsValue.String())
	} else {
		log.Println("no update required")
	}

	return nil
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

	ipifyAPIUrl, ipifyAPIUrlExists := os.LookupEnv(DNS_UPDATER_IPIFY_URL)
	config.IpifyAPIUrl = ipifyAPIUrl

	if !goDaddyAPIUrlExists || !goDaddyAPIKeyExists || !hostExists || !domainExists || !ipifyAPIUrlExists {
		return config, errors.New("environment variable missing. e.g: " +
			"export DNS_UPDATER_GO_DADDY_API_URL=FOO " +
			"export DNS_UPDATER_GO_DADDY_API_KEY=BAR " +
			"export DNS_UPDATER_HOST=BAZ " +
			"export DNS_UPDATER_DOMAIN=OPA " +
			"export DNS_UPDATER_IPIFY_URL=OMA")
	}

	return config, nil
}
