package main

import (
	"flag"
	"log"
	"os"

	dot "github.com/krakend/krakend-config2dot/v3"
	"github.com/luraproject/lura/v3/config"
)

func main() {
	configFile := flag.String("c", "config.json", "path to the krakend config file")
	flag.Parse()

	cfg, err := config.NewParser().Parse(*configFile)
	if err != nil {
		log.Fatal(err)
		return
	}

	dotCfg := dot.ServiceConfig(cfg)
	dotCfg.WriteTo(os.Stdout)
}
