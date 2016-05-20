package ns

import (
	"log"
	"github.com/makhomed/nsd-tool/config"
	"github.com/makhomed/nsd-tool/util"
)

func Check(conf *config.Config) error {
	zones, err := util.ConfigZones(conf)
	if err != nil {
		return err
	}
	for _, zone := range zones {
		configNS, err := util.ConfigNS(conf, zone)
		if err != nil {
			log.Printf("configNS(%s): %v\n", zone, err)
		}
		compareSOA(zone, configNS)
	}
	return nil
}

func compareSOA(zone string, configNS []string) {
	soas := make([]string, 0)
	for _, ns := range configNS {
		soa, err := util.SOA(zone, ns)
		if err != nil {
			log.Printf("SOA(%s): %v\n", zone, err)
			continue
		}
		soas = append(soas, soa)
	}
	if len(soas) == 0 {
		return
	}
	etalon := soas[0]
	for i, soa := range soas {
		if soa != etalon {
			log.Printf("%s has different SOA on %s\n", zone, configNS[i])
		}
	}
}
