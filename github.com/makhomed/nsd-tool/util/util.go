package util

import (
	"io/ioutil"
	"github.com/makhomed/nsd-tool/config"
	"path/filepath"
	"os"
	"github.com/miekg/dns"
	"log"
	"net"
	"fmt"
	"strings"
)

func ConfigZones(conf *config.Config) ([]string, error) {
	files, err := ioutil.ReadDir(conf.ZoneDir)
	if err != nil {
		return nil, err
	}
	dir := make([]string, 0, len(files))
	for _, file := range files {
		dir = append(dir, file.Name())
	}
	return dir, nil
}

func ConfigNS(conf *config.Config, zone string) ([]string, error) {
	file, err := os.Open(filepath.Join(conf.ZoneDir, zone))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	ns := make([]string, 0)
	for line := range dns.ParseZone(file, zone, file.Name()) {
		if line.Error != nil {
			log.Printf("error parsing zone file %s: %v", zone, line.Error)
		} else {
			if line.RR.Header().Rrtype != dns.TypeNS {
				continue
			}
			rr := line.RR.(*dns.NS)
			ns = append(ns, rr.Ns)
		}
	}
	return ns, nil
}

func Tld(zone string) string {
	return strings.Join(dns.SplitDomainName(zone)[1:], ".")
}

var nsCache map[string][]string = make(map[string][]string)

func NsCache(zone string) []string {
	return nsCache[Tld(zone)]
}

func InitNsCache(conf *config.Config, zones []string) error {
	for _, zone := range zones {
		tld := Tld(zone)
		if _, ok := nsCache[tld]; ok {
			continue
		}
		ns, err := NS(conf, tld, nil, false)
		if err != nil {
			return fmt.Errorf("InitNsCache failed: %v", err)
		}
		nsCache[tld] = ns
	}
	return nil
}

var dnsError map[string]bool = make(map[string]bool)

func DnsError(zone string) bool {
	return notExist[zone]
}

var notExist map[string]bool = make(map[string]bool)

func NotExist(zone string) bool {
	return notExist[zone]
}

func NS(conf *config.Config, zone string, servers []string, authority bool) ([]string, error) {
	if servers == nil {
		servers = []string{conf.Resolver}
	}
	client := new(dns.Client)
	message := new(dns.Msg)
	message.SetQuestion(dns.Fqdn(zone), dns.TypeNS)
	message.RecursionDesired = true

	var lastError error
	for _, server := range servers {
		reply, _, err := client.Exchange(message, net.JoinHostPort(server, "53"))
		if err != nil {
			lastError = err
			continue
		}
		ns := make([]string, 0)
		if reply.Rcode == dns.RcodeNameError {
			notExist[zone] = true
			return ns, nil
		}

		if reply.Rcode == dns.RcodeServerFailure {
			dnsError[zone] = true
			return ns, nil
		}

		if reply.Rcode != dns.RcodeSuccess {
			return nil, fmt.Errorf("unexpected dns error %d", reply.Rcode)
		}

		switch authority {
		case false:
			for _, a := range reply.Answer {
				rr := a.(*dns.NS)
				ns = append(ns, rr.Ns)
			}
		case true:
			for _, a := range reply.Ns {
				rr, ok := a.(*dns.NS)
				if ok {
					ns = append(ns, rr.Ns)
				} else {
					// SOA instead of NS, domain exists, but in deletion state
				}
			}
		}
		return ns, nil
	}
	if lastError != nil {
		return nil, lastError
	}
	panic("unreachable code")
}

func DelegationNS(conf *config.Config, zone string) ([]string, error) {
	return NS(conf, zone, NsCache(zone), true)
}
