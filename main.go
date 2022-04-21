package main

import (
	"arcam-controller/arcam"
	"context"
	"fmt"
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
	Model      string
	IP         string
	Port       int
	Pin        string
	ListenPort int
}

func validateModel(Model string) bool {
	return true
}

func main() {
	model := flag.String("model", "", "Receiver model name (eg AVR11)")
	ipAddress := flag.String("ip", "", "IP Address of Receiver")
	port := flag.Int("port", 50001, "Port to communicate with receiver")
	pin := flag.String("pin", "00102003", "Homekit pairing pin")
	listenPort := flag.Int("listen-port", 33859, "Port to listen for homekit")

	flag.Parse()

	validateModel(*model)

	cfg := Config{
		Model:      *model,
		IP:         *ipAddress,
		Port:       *port,
		Pin:        *pin,
		ListenPort: *listenPort,
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
	tv.Television.Active.OnValueUpdate(powerOnCallback(ctx, &arcamClient))

	tv.Television.ConfiguredName.SetValue("Arcam Receiver")
	tv.Television.SleepDiscoveryMode.SetValue(characteristic.SleepDiscoveryModeAlwaysDiscoverable)

	attachInputs(ctx, &arcamClient, tv)

	source, err := arcamClient.GetSource(ctx)
	if err != nil {
		log.Info.Fatalln("")
	}
	tv.Television.ActiveIdentifier.SetValue(int(source))
	tv.Television.ActiveIdentifier.OnValueUpdate(inputSelectionCallback(ctx, &arcamClient))

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
	speakerVolume.OnValueUpdate(speakerVolumeCallback(ctx, &arcamClient))
	tv.Speaker.AddC(speakerVolume.C)

	isMuted, err := arcamClient.IsMuted(ctx)
	tv.Speaker.Mute.SetValue(isMuted)
	tv.Speaker.Mute.OnValueUpdate(muteCallback(ctx, &arcamClient))

	dimmer := service.NewLightbulb()
	dimmer.On.SetValue(!isMuted)
	dimmer.On.OnValueUpdate(muteCallback(ctx, &arcamClient))

	brightness := characteristic.NewBrightness()
	brightness.Description = "Volume"
	brightness.SetValue(vol)
	brightness.OnValueUpdate(brightnessVolCallback(ctx, &arcamClient))
	dimmer.AddC(brightness.C)
	tv.Television.AddS(dimmer.S)
	tv.A.AddS(dimmer.S)

	s, err := hap.NewServer(hap.NewFsStore("./db"), bridge.A, tv.A)
	if err != nil {
		log.Info.Panic(err)
	}

	s.Pin = cfg.Pin
	s.Addr = fmt.Sprintf("[::]:%d", cfg.ListenPort)

	mylogger := syslog.New(os.Stdout, "ARCAM-HOMEKIT ", syslog.LstdFlags|syslog.Lshortfile)
	log.Debug = &log.Logger{mylogger}

	err = s.ListenAndServe(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
