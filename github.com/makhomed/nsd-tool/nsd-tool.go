package main

import (
	"os"
	"github.com/makhomed/nsd-tool/config"
	"github.com/makhomed/nsd-tool/zonelist"
	"github.com/makhomed/nsd-tool/delegation"
	"log"
	"github.com/makhomed/nsd-tool/ns"
)

const (
	ConfigFileName = "/opt/nsd-tool/conf/nsd-tool.conf"
	usage = `
usage:
	nsd-tool generate zonelist <pattern> </path/to/zonelist.conf>
	nsd-tool check delegation
	nsd-tool check ns
`
)

//go:generate go get github.com/miekg/dns

func main() {
	log.SetFlags(0)
	conf, err := config.New(ConfigFileName)
	if err != nil {
		log.Fatalf("can't read config '%s' : %v\n\n", ConfigFileName, err)
	}
	switch {
	case len(os.Args) == 5 && os.Args[1] == "generate" && os.Args[2] == "zonelist":
		pattern := os.Args[3]
		filename := os.Args[4]
		if err := zonelist.Generate(conf, pattern, filename); err != nil {
			log.Fatalf("generate zonelist: %v\n\n", err)
		}
	case len(os.Args) == 3 && os.Args[1] == "check" && os.Args[2] == "delegation":
		if err := delegation.Check(conf); err != nil {
			log.Fatalf("check delegation: %v\n\n", err)
		}
	case len(os.Args) == 3 && os.Args[1] == "check" && os.Args[2] == "ns":
		if err := ns.Check(conf); err != nil {
			log.Fatalf("check ns: %v\n\n", err)
		}
	default:
		log.Fatalf(usage)
	}
}
