package soa

import (
	"log"
	"github.com/makhomed/nsd-tool/util"
	"github.com/makhomed/nsd-tool/config"
	"path/filepath"
	"os"
	"sort"
	"fmt"
	"io"
	"io/ioutil"
)

type ZoneInfo struct {
	Zone     string
	Serial   uint32
	Checksum string
}

type ZonesInfo map[string]ZoneInfo

func Check(conf *config.Config) error {
	var errorCount int
	resultZonesInfo := make(ZonesInfo)
	configZonesInfo, err := ConfigZonesInfo(conf)
	if err != nil {
		return err
	}
	memoryZonesInfo, err := ReadMemoryZonesInfo(conf)
	if err != nil {
		return err
	}

	var zones []string
	for zone := range *configZonesInfo {
		zones = append(zones, zone)
	}
	sort.Strings(zones)
	for _, zone := range zones {
		configZoneInfo := (*configZonesInfo)[zone]
		configSerial := configZoneInfo.Serial
		configChecksum := configZoneInfo.Checksum
		if memoryZoneInfo, ok := (*memoryZonesInfo)[zone]; !ok {
			// zoneInfo not exists in memory, just copy it
			resultZonesInfo[zone] = configZoneInfo
		} else {
			memorySerial := memoryZoneInfo.Serial
			memoryChecksum := memoryZoneInfo.Checksum
			if configChecksum == memoryChecksum {
				// zoneInfo not changed, just copy it
				resultZonesInfo[zone] = configZoneInfo
			} else {
				if configSerial > memorySerial {
					// zoneInfo changed, serial incremented
					resultZonesInfo[zone] = configZoneInfo
				} else {
					// save information about serial error
					resultZonesInfo[zone] = memoryZoneInfo
					log.Printf("BAD SERIAL: %s / %d", zone, configSerial)
					errorCount++
				}
			}
		}
	}
	err = WriteMemoryZonesInfo(conf, resultZonesInfo)
	if err != nil {
		return err
	}
	if errorCount > 0 {
		return fmt.Errorf("total %d errors deteted", errorCount)
	}
	return nil
}

const MemoryFileName = "/opt/nsd-tool/var/memory"

func ReadMemoryZonesInfo(conf *config.Config) (*ZonesInfo, error) {
	zonesInfo := make(ZonesInfo)
	file, err := os.Open(MemoryFileName)
	if err != nil {
		// first run, memory file not exists
		return &zonesInfo, nil
	}
	defer file.Close()
	for {
		var zone string
		var serial uint32
		var checksum string
		_, err := fmt.Fscanf(file, "%s\t%d\t%s\n", &zone, &serial, &checksum)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		zonesInfo[zone] = ZoneInfo{Zone: zone, Serial: serial, Checksum: checksum}
	}
	return &zonesInfo, nil
}

func WriteMemoryZonesInfo(conf *config.Config, zonesInfo ZonesInfo) error {
	dir := filepath.Dir(MemoryFileName)
	prefix := filepath.Base(MemoryFileName)
	newMemoryFile, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		return err
	}
	defer func(newMemoryFile string) {
		_, err := os.Stat(newMemoryFile)
		if os.IsNotExist(err) {
			return
		}
		os.Remove(newMemoryFile)
	}(newMemoryFile.Name())

	var zones []string
	for zone := range zonesInfo {
		zones = append(zones, zone)
	}
	sort.Strings(zones)
	for _, zone := range zones {
		serial := zonesInfo[zone].Serial
		checksum := zonesInfo[zone].Checksum
		io.WriteString(newMemoryFile, fmt.Sprintf("%s\t%d\t%s\n", zone, serial, checksum))
	}

	newMemoryFile.Close()
	err = os.Rename(newMemoryFile.Name(), MemoryFileName)
	if err != nil {
		return err
	}
	return nil
}

func ConfigZonesInfo(conf *config.Config) (*ZonesInfo, error) {
	zones, err := util.ConfigZones(conf)
	if err != nil {
		return nil, err
	}
	zonesInfo := make(ZonesInfo)
	for _, zone := range zones {
		serial, err := util.ConfigSerial(conf, zone)
		if err != nil {
			log.Printf("ConfigSerial(%s): %v\n", zone, err)
			continue
		}
		checksum, err := util.ConfigChecksum(conf, zone)
		if err != nil {
			log.Printf("ConfigChecksum(%s): %v\n", zone, err)
			continue
		}
		zonesInfo[zone] = ZoneInfo{Zone: zone, Serial: serial, Checksum: checksum}
	}
	return &zonesInfo, nil
}
