package main

import (
	"context"
	syslog "log"
	"os"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/log"

	"flag"
	"os/signal"
	"syscall"
)

type Config struct {
	Model string
}

func validateModel(Model string) bool {
	return true
}

func main() {
	var model *string = flag.String("model", "", "Receiver model name (eg AVR11)")
	validateModel(*model)

	cfg := Config{
		Model: *model,
	}

	bridge := accessory.NewBridge(accessory.Info{
		Name:         "ARCAM Receiver Bridge",
		SerialNumber: "",
		Model:        cfg.Model,
		Manufacturer: "",
		Firmware:     "",
	})

	s, err := hap.NewServer(hap.NewFsStore("./db"), bridge.A)
	if err != nil {
		log.Info.Panic(err)
	}

	logger := syslog.New(os.Stdout, "ARCAM ", syslog.LstdFlags|syslog.Lshortfile)
	log.Debug = &log.Logger{logger}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		signal.Stop(c)
		cancel()
	}()

	s.ListenAndServe(ctx)
}
