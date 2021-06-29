package main

import (
	"fmt"
	"os"

	"io.rocktest/rocktest/scenario"

	log "github.com/sirupsen/logrus"
)

func banner() {
	version := "2.0.0"

	banner := `
________               ______  ________              _____
___  __ \______ __________  /_____  __/_____ __________  /_
__  /_/ /_  __ \_  ___/__  //_/__  /   _  _ \__  ___/_  __/
_  _, _/ / /_/ // /__  _  ,<   _  /    /  __/_(__  ) / /_
/_/ |_|  \____/ \___/  /_/|_|  /_/     \___/ /____/  \__/

                 Test automation that rocks ! (v` + version + ")"

	fmt.Printf("%s\n\n", banner)

}

func initLog() {
	log.SetLevel(log.TraceLevel)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})
}

func main() {

	initLog()

	banner()

	scenFile := os.Args[1]

	log.Infof("%s", "Load scenario")

	s := scenario.NewScenario()

	s.Run(scenFile)

}
