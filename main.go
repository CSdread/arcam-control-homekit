package main

import (
	"arcam-controller/arcam"
	"context"
	"fmt"
	"os"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/log"

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

	mylogger := syslog.New(os.Stdout, "ARCAM-HOMEKIT ", syslog.LstdFlags|syslog.Lshortfile)
	log.Debug = &log.Logger{mylogger}

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

	// create homekit resources
	bridge := accessory.NewBridge(accessory.Info{
		Name:         "ARCAM Receiver Bridge",
		SerialNumber: Version,
		Model:        "ARCAM Receiver Bridge",
		Manufacturer: "CSdread",
	})

	tv := accessory.NewTelevision(accessory.Info{
		Name:         "ARCAM Receiver",
		SerialNumber: Version,
		Model:        cfg.Model,
		Manufacturer: "Arcam",
	})
	tv.Television.ConfiguredName.SetValue("Arcam Receiver")
	tv.Television.SleepDiscoveryMode.SetValue(characteristic.SleepDiscoveryModeAlwaysDiscoverable)
	attachInputs(ctx, &arcamClient, tv)
	tv.Television.Active.OnValueUpdate(powerOnCallback(ctx, &arcamClient))
	tv.Television.ActiveIdentifier.OnValueUpdate(inputSelectionCallback(ctx, &arcamClient))
	/*
		dimmer := service.NewLightbulb()
		brightness := characteristic.NewBrightness()
		brightness.Description = "Volume"
		dimmer.AddC(brightness.C)
		tv.Television.AddS(dimmer.S)
		tv.A.AddS(dimmer.S)

		tv.Television.AddS(tv.Speaker.S)

		speakerActive := characteristic.NewActive()
		speakerActive.SetValue(characteristic.ActiveActive)
		tv.Speaker.AddC(speakerActive.C)

		speakerVolControlType := characteristic.NewVolumeControlType()
		speakerVolControlType.SetValue(characteristic.VolumeControlTypeAbsolute)
		tv.Speaker.AddC(speakerVolControlType.C)

		speakerVolume := characteristic.NewVolume()
		tv.Speaker.AddC(speakerVolume.C)

		// Register event handlers for updates originating from homekit
		brightness.OnValueUpdate(brightnessVolCallback(ctx, &arcamClient))
		dimmer.On.OnValueUpdate(muteCallback(ctx, &arcamClient))
		tv.Speaker.Mute.OnValueUpdate(muteCallback(ctx, &arcamClient))
		speakerVolume.OnValueUpdate(speakerVolumeCallback(ctx, &arcamClient))

		// Register event handlers for updates originating from the Receiver
	*/
	arcamClient.RegisterEventHandler(arcam.PowerCommand, func(zone arcam.ZoneNumber, data []byte) error {
		active := characteristic.ActiveInactive
		if arcam.PowerStatus(data[0]) == arcam.PowerStatusActive {
			active = characteristic.ActiveActive
		}
		return tv.Television.Active.SetValue(active)
	})
	arcamClient.RegisterEventHandler(arcam.RequestCurrentSource, func(zone arcam.ZoneNumber, data []byte) error {
		return tv.Television.ActiveIdentifier.SetValue(int(data[0]))
	})
	/*
		arcamClient.RegisterEventHandler(arcam.RequestMuteStatus, func(zone arcam.ZoneNumber, data []byte) error {
			isMuted := arcam.MuteState(data[0]) == arcam.MuteStateMuted
			tv.Speaker.Mute.SetValue(isMuted)
			dimmer.On.SetValue(isMuted)

			return nil
		})

		arcamClient.RegisterEventHandler(arcam.SetRequestVolume, func(zone arcam.ZoneNumber, data []byte) error {
			err := speakerVolume.SetValue(int(data[0]))
			if err != nil {
				return err
			}
			return brightness.SetValue(int(data[0]))
		})
	*/
	arcamClient.RefreshState(ctx)

	s, err := hap.NewServer(hap.NewFsStore("./db"), bridge.A, tv.A)
	if err != nil {
		log.Info.Panic(err)
	}

	s.Pin = cfg.Pin
	s.Addr = fmt.Sprintf("[::]:%d", cfg.ListenPort)

	err = s.ListenAndServe(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
