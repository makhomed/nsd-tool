package delegation

import (
	"log"
	"sort"
	"strings"
	"github.com/makhomed/nsd-tool/config"
	"github.com/makhomed/nsd-tool/util"
)

func Check(conf *config.Config) error {
	zones, err := util.ConfigZones(conf)
	if err != nil {
		return err
	}
	err = util.InitNsCache(conf, zones)
	if err != nil {
		return err
	}
	for _, zone := range zones {
		configNS, err := util.ConfigNS(conf, zone)
		if err != nil {
			log.Printf("configNS(%s): %v\n", zone, err)
		}
		delegationNS, err := util.DelegationNS(conf, zone)
		if err != nil {
			log.Printf("delegationNS(%s): %v\n", zone, err)
		}
		compareNS(zone, configNS, delegationNS)
	}
	return nil
}

func compareNS(zone string, configNS []string, delegationNS []string) {
	sort.Strings(configNS)
	sort.Strings(delegationNS)
	config := strings.ToLower(strings.Join(configNS, " "))
	delegation := strings.ToLower(strings.Join(delegationNS, " "))
	if config == delegation {
		return
	}
	switch {
	case util.NotExist(zone):
		log.Printf("%s\tName Error", zone)
	case util.DnsError(zone):
		log.Printf("%s\tServer Failure", zone)
	default:
		log.Printf("%s\tconfigNS: %s\tdelegationNS: %v\n", zone, config, delegation)
	}
}
