package main

import (
	"fmt"
	"os"

	"github.com/pborman/getopt/v2"

	"io.rocktest/rocktest/scenario"

	log "github.com/sirupsen/logrus"
)

func banner() {
	version := "1.0.0-go"

	banner := `
________               ______  ________              _____
___  __ \______ __________  /_____  __/_____ __________  /_
__  /_/ /_  __ \_  ___/__  //_/__  /   _  _ \__  ___/_  __/
_  _, _/ / /_/ // /__  _  ,<   _  /    /  __/_(__  ) / /_
/_/ |_|  \____/ \___/  /_/|_|  /_/     \___/ /____/  \__/

                 Test automation that rocks ! (v` + version + ")"

	fmt.Printf("%s\n\n", banner)

}

func initLog(verbose *string) {

	switch *verbose {
	case "0":
		log.SetLevel(log.ErrorLevel)
	case "1":
		log.SetLevel(log.WarnLevel)
	case "2":
		log.SetLevel(log.InfoLevel)
	case "3":
		log.SetLevel(log.DebugLevel)
	case "4":
		log.SetLevel(log.TraceLevel)
	default:
		fmt.Printf("Bad -v value. Allowed values are 0 - 4")
		os.Exit(1)
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})
}

func main() {

	verbose := getopt.StringLong("verbose", 'v', "0", "Log level (0=ERROR -> 4=TRACE")
	getopt.Parse()
	args := getopt.Args()

	initLog(verbose)

	banner()

	scenFile := args[0]

	log.Infof("%s", "Load scenario")

	s := scenario.NewScenario()

	err := s.Run(scenFile)

	if err == nil {
		fmt.Printf(`
========================================
=     Scenario Success ! It Rocks      =
========================================
`)
	} else {
		fmt.Printf(`
=======================================
            Scenario failure             
		
%s	
=======================================
`, err.Error())
	}

}
