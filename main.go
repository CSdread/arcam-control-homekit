package main

import (
	"arcam-controller/arcam"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/brutella/hap"
	"github.com/brutella/hap/log"

	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"

	"flag"
	syslog "log"
	"os/signal"
	"syscall"
)

type Config struct {
	Model string
	IP    string
	Port  int
	Pin   string
}

func validateModel(Model string) bool {
	return true
}

func main() {
	model := flag.String("model", "", "Receiver model name (eg AVR11)")
	ipAddress := flag.String("ip", "", "IP Address of Receiver")
	port := flag.Int("port", 50001, "Port to communicate with receiver")
	pin := flag.String("pin", "00102003", "Homekit pairing pin")

	flag.Parse()

	validateModel(*model)

	cfg := Config{
		Model: *model,
		IP:    *ipAddress,
		Port:  *port,
		Pin:   *pin,
	}

	arcamClient, err := arcam.NewReceiver(cfg.Model, cfg.IP, cfg.Port)
	if err != nil {
		log.Debug.Fatalln("")
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		signal.Stop(c)
		cancel()
	}()

	err = arcamClient.Connect(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	bridge := accessory.NewBridge(accessory.Info{
		Name:         "ARCAM Receiver Bridge",
		Model:        "ARCAM Receiver Bridge",
		Manufacturer: "CSdread",
	})

	tv := accessory.NewTelevision(accessory.Info{
		Name:         "ARCAM Receiver",
		SerialNumber: "",
		Model:        cfg.Model,
		Manufacturer: "Arcam",
	})

	isOn, err := arcamClient.IsOn(ctx)
	active := characteristic.ActiveInactive
	if isOn {
		active = characteristic.ActiveActive
	}

	tv.Television.Active.SetValue(active)

	tv.Television.Active.OnValueUpdate(func(newVal, oldVal int, r *http.Request) {
		if newVal == characteristic.ActiveActive {
			err := arcamClient.PowerOn(ctx)
			if err != nil {
				log.Info.Fatalln("")
			}
			return
		}

		err := arcamClient.PowerOff(ctx)
		if err != nil {
			log.Info.Fatalln("")
		}
	})

	s, err := hap.NewServer(hap.NewFsStore("./db"), bridge.A, tv.A)
	if err != nil {
		log.Info.Panic(err)
	}

	s.Pin = cfg.Pin

	mylogger := syslog.New(os.Stdout, "BLUB ", syslog.LstdFlags|syslog.Lshortfile)
	log.Debug = &log.Logger{mylogger}

	err = s.ListenAndServe(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
