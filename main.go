package main

import (
	"info_exporter/configs"
	"info_exporter/pkg/host"
	"info_exporter/pkg/prom"
	"info_exporter/pkg/switches"
	"info_exporter/pkg/tools"
)

func main() {
	// log init
	tools.LogsInit()

	// config init
	configs.ConfInit()

	// ecdn host
	host.Monitor()

	// ecdn switch
	switches.Monitor()

	prom.Run()

	select {}
}
