package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/arriqaaq/flashdb"
	"github.com/arriqaaq/flashdb/cmd"
	"github.com/pelletier/go-toml"
)

var config = flag.String("config", "", "the config file for flashdb")

func main() {
	flag.Parse()

	// Set the config.
	var cfg *flashdb.Config
	if *config == "" {
		log.Println("no config set, using the default config.")
		cfg = flashdb.DefaultConfig()
	} else {
		c, err := newConfigFromFile(*config)
		if err != nil {
			log.Printf("load config err : %+v\n", err)
			return
		}
		cfg = c
	}

	// Listen the server.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGHUP,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	server, err := cmd.NewServer(cfg)
	if err != nil {
		log.Printf("create flashdb server err: %+v\n", err)
		return
	}
	go server.Listen(cfg.Addr)

	<-sig
	server.Stop()
}

func newConfigFromFile(config string) (*flashdb.Config, error) {
	data, err := ioutil.ReadFile(filepath.Clean(config))
	if err != nil {
		return nil, err
	}

	var cfg = new(flashdb.Config)
	err = toml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
