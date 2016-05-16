package config

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	"strconv"
)

type Config struct {
	ZoneDir string
}

func defaultConfig() *Config {
	return &Config{
		"",
	}
}

func New(configFileName string) (*Config, error) {
	conf := defaultConfig()
	configFile, err := os.Open(configFileName)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	scanner := bufio.NewScanner(configFile)
	for scanner.Scan() {
		line := scanner.Text()
		commentPosition := strings.Index(line, "#")
		if commentPosition >= 0 {
			line = line[0:commentPosition]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = strings.Replace(line, "\t", " ", -1)
		split := strings.SplitN(line, " ", 2)
		name, value := split[0], split[1]
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)
		setter := configSetter{}
		setter.setString("zoneDir", &conf.ZoneDir, name, value)
		if err := setter.err(); err != nil {
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return verifyConfig(conf)
}

type configSetter struct {
	parseError error
}

func (setter *configSetter) setString(setName string, setValue *string, name string, value string) {
	if strings.ToLower(setName) == strings.ToLower(name) {
		*setValue = value
	}
}

func (setter *configSetter) setInt(setName string, setValue *int, name string, value string) {
	if strings.ToLower(setName) == strings.ToLower(name) {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			setter.parseError = err
		}
		*setValue = intValue
	}
}

func (setter *configSetter) err() (err error) {
	return setter.parseError
}

func verifyConfig(conf *Config) (*Config, error) {
	stat, err := os.Stat(conf.ZoneDir)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("bad 'zoneDir' value '%s' : %s", conf.ZoneDir, err)
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("bad 'zoneDir' value '%s' : must be directory", conf.ZoneDir)
	}
	return conf, nil
}
