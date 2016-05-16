package main

import (
	"fmt"
	"os"

	"github.com/makhomed/nsd-tool/config"
	"github.com/makhomed/nsd-tool/zonelist"
)

const (
	ConfigFileName = "/opt/nsd-tool/conf/nsd-tool.conf"
	usage = `
usage:
	nsd-tool generate zonelist <pattern> </path/to/zonelist.conf>
`
)

func main() {
	conf, err := config.New(ConfigFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't read config '%s' : %v\n\n", ConfigFileName, err)
		os.Exit(2)
	}
	// generate zonelist <pattern> <filename>
	if len(os.Args) == 5 && os.Args[1] == "generate" && os.Args[2] == "zonelist" {
		pattern := os.Args[3]
		filename := os.Args[4]
		err := zonelist.Generate(conf, pattern, filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "generate zonelist : %v\n\n", err)
			os.Exit(2)
		}
	} else {
		fmt.Fprintf(os.Stderr, usage)
		os.Exit(2)
	}
}
