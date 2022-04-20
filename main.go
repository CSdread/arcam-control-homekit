package main

import (
	"arcam-controller/arcam"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/brutella/hap"
	"github.com/brutella/hap/log"
	"github.com/brutella/hap/service"

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

	tv.Television.ConfiguredName.SetValue(fmt.Sprintf("Arcam %s Receiver", cfg.Model))
	tv.Television.SleepDiscoveryMode.SetValue(characteristic.SleepDiscoveryModeAlwaysDiscoverable)

	for _, input := range arcamClient.GetAllInputs() {
		inputSource := service.NewInputSource()

		displayName := arcam.InputDisplayNameMap[input]
		inputId := int(input)

		id := characteristic.NewIdentifier()
		id.SetValue(inputId)
		inputSource.AddC(id.C)

		inputSource.ConfiguredName.SetValue(displayName)
		inputSource.IsConfigured.SetValue(characteristic.IsConfiguredConfigured)
		inputSource.InputSourceType.SetValue(characteristic.InputSourceTypeHdmi)
		inputSource.CurrentVisibilityState.SetValue(characteristic.CurrentVisibilityStateShown)
		inputSource.ConfiguredName.Permissions = []string{characteristic.PermissionRead}

		tv.Television.AddS(inputSource.S)
		tv.A.AddS(inputSource.S)
	}

	source, err := arcamClient.GetSource(ctx)
	if err != nil {
		log.Info.Fatalln("")
	}
	tv.Television.ActiveIdentifier.SetValue(int(source))
	tv.Television.ActiveIdentifier.OnValueUpdate(func(newVal, oldVal int, r *http.Request) {
		arcamClient.SetSource(ctx, arcam.InputSource(newVal))
	})

	tv.Television.AddS(tv.Speaker.S)

	speakerActive := characteristic.NewActive()
	speakerActive.SetValue(characteristic.ActiveActive)
	tv.Speaker.AddC(speakerActive.C)

	speakerVolControlType := characteristic.NewVolumeControlType()
	speakerVolControlType.SetValue(characteristic.VolumeControlTypeAbsolute)
	tv.Speaker.AddC(speakerVolControlType.C)

	vol, err := arcamClient.GetVolume(ctx)

	speakerVolume := characteristic.NewVolume()
	speakerVolume.SetValue(vol)
	speakerVolume.OnValueUpdate(func(newVal, oldVal int, r *http.Request) {
		arcamClient.SetVolume(ctx, newVal)
	})
	tv.Speaker.AddC(speakerVolume.C)

	isMute, err := arcamClient.IsMuted(ctx)
	tv.Speaker.Mute.SetValue(isMute)
	tv.Speaker.Mute.OnValueUpdate(func(oldVal, newVal bool, r *http.Request) {
		arcamClient.ToggleMute(ctx)
	})

	s, err := hap.NewServer(hap.NewFsStore("./db"), bridge.A, tv.A)
	if err != nil {
		log.Info.Panic(err)
	}

	s.Pin = cfg.Pin
	s.Addr = "[::]:33859"

	mylogger := syslog.New(os.Stdout, "ARCAM-HOMEKIT", syslog.LstdFlags|syslog.Lshortfile)
	log.Debug = &log.Logger{mylogger}

	err = s.ListenAndServe(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
