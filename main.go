package main

import (
	"arcam-controller/arcam"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/brutella/hap"
	"github.com/brutella/hap/log"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"

	"flag"
	"os/signal"
	"syscall"
)

type Config struct {
	Model string
	IP    string
	Port  int
}

func validateModel(Model string) bool {
	return true
}

func main() {
	model := flag.String("model", "", "Receiver model name (eg AVR11)")
	ipAddress := flag.String("ip", "", "")
	port := flag.Int("port", 50000, "")

	flag.Parse()

	validateModel(*model)

	cfg := Config{
		Model: *model,
		IP:    *ipAddress,
		Port:  *port,
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
		SerialNumber: "",
		Model:        "",
		Manufacturer: "",
		Firmware:     "",
	})

	tv := accessory.NewTelevision(accessory.Info{
		Name:         "ARCAM Receiver",
		SerialNumber: "",
		Model:        cfg.Model,
		Manufacturer: "Arcam",
		Firmware:     "",
	})
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

	volControl := characteristic.NewVolume()
	volControl.SetMaxValue(99)
	volControl.SetMinValue(0)
	volControl.SetStepValue(1)
	volControl.OnValueUpdate(func(newVol, oldVol int, r *http.Request) {
		err := arcamClient.SetVolume(ctx, newVol)
		if err != nil {
			log.Info.Fatalln("")
		}
	})

	tv.Speaker.AddC(volControl.C)

	s, err := hap.NewServer(hap.NewFsStore("./db"), bridge.A, tv.A)
	if err != nil {
		log.Info.Panic(err)
	}

	logger := syslog.New(os.Stdout, "ARCAM ", syslog.LstdFlags|syslog.Lshortfile)
	log.Debug = &log.Logger{logger}

	s.ListenAndServe(ctx)
}
