package main

import (
	"arcam-controller/arcam"
	"context"
	"fmt"
	"os"

	"github.com/brutella/hap"
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
		log.Debug.Fatalln("Unable to create new receiver: %s", err)
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
		log.Debug.Fatalf("Could not connect to receiver: %s", err)
	}

	// gather inputs
	sources := []Source{}
	for _, input := range arcamClient.GetInputs() {
		source := Source{name: byte(input), displayName: arcam.InputDisplayNameMap[input]}
		sources = append(sources, source)
	}

	homekitReceiver := NewHomekitReceiver(ctx, cfg.Model, sources)

	// Register homekit callbacks
	homekitReceiver.RegisterPowerCallback(powerCallback(ctx, arcamClient))
	homekitReceiver.RegisterSourceCallback(sourceCallback(ctx, arcamClient))
	homekitReceiver.RegisterMuteCallback(muteCallback(ctx, arcamClient))
	homekitReceiver.RegisterVolumeCallback(volumeCallback(ctx, arcamClient))
	homekitReceiver.RegisterDirectModeCallback(directModeStatusCallback(ctx, arcamClient))

	// Register Arcam callbacks
	arcamClient.RegisterEventHandler(arcam.PowerCommand, powerCommandCallback(homekitReceiver))
	arcamClient.RegisterEventHandler(arcam.RequestCurrentSource, requestCurrentSourceCallback(homekitReceiver))
	arcamClient.RegisterEventHandler(arcam.SetRequestVolume, volumeCommandCallback(homekitReceiver))
	arcamClient.RegisterEventHandler(arcam.RequestMuteStatus, muteCommandCallback(homekitReceiver))
	arcamClient.RegisterEventHandler(arcam.RequestDirectModeStatus, directModeCommandCallback(homekitReceiver))

	// initial state
	arcamClient.RefreshState(ctx)

	// create HAP server
	s, err := hap.NewServer(hap.NewFsStore("./db"), homekitReceiver.Bridge.A, homekitReceiver.Tv.A)
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
